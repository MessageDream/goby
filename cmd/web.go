package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/go-macaron/binding"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/captcha"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/gzip"
	. "github.com/go-macaron/session"
	"github.com/go-macaron/toolbox"
	log "gopkg.in/clog.v1"
	"gopkg.in/macaron.v1"
	cli "gopkg.in/urfave/cli.v1"

	"github.com/MessageDream/goby/model"
	"github.com/MessageDream/goby/module/context"
	forms "github.com/MessageDream/goby/module/form"
	"github.com/MessageDream/goby/module/setting"
	"github.com/MessageDream/goby/module/template"
	"github.com/MessageDream/goby/router"
	"github.com/MessageDream/goby/router/accessKey"
	"github.com/MessageDream/goby/router/account"
	"github.com/MessageDream/goby/router/app"
	"github.com/MessageDream/goby/router/app/deployment"
	"github.com/MessageDream/goby/router/auth"
	"github.com/MessageDream/goby/router/client"
	"github.com/MessageDream/goby/router/collaborator"
	"github.com/MessageDream/goby/router/dev"
	"github.com/MessageDream/goby/router/web"
)

var CmdWeb = cli.Command{
	Name:        "server",
	Usage:       "Start code push web server",
	Description: ``,
	Action:      runWeb,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "port, p",
			Value: "3000",
			Usage: "Temporary port number to prevent conflict",
		},
		cli.StringFlag{
			Name:  "config, c",
			Value: "custom/conf/app.ini",
			Usage: "Custom configuration file path",
		},
	},
}

func checkVersion() {
	data, err := ioutil.ReadFile(path.Join(setting.StaticRootPath, "template/.VERSION"))
	if err != nil {
		log.Fatal(4, "Fail to read 'template/.VERSION': %v", err)
	}
	if string(data) != setting.AppVer {
		log.Fatal(4, "Binary and template file version does not match, did you forget to recompile?")
	}
}

func newMacaron() *macaron.Macaron {
	m := macaron.New()
	if !setting.DisableRouterLog {
		m.Use(macaron.Logger())
	}
	m.Use(macaron.Recovery())
	m.Use(macaron.Static(
		path.Join(setting.StaticRootPath, "public"),
		macaron.StaticOptions{
			SkipLogging: !setting.DisableRouterLog,
		},
	))

	if setting.EnableGzip {
		m.Use(gzip.Gziper())
	}

	funcMap := template.NewFuncMap()
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory:         path.Join(setting.StaticRootPath, "template"),
		AppendDirectories: []string{path.Join(setting.CustomPath, "template")},
		Funcs:             funcMap,
		IndentJSON:        macaron.Env != macaron.PROD,
		Charset:           "utf-8",
	}))

	m.Use(cache.Cacher(cache.Options{
		Adapter:       setting.CacheAdapter,
		AdapterConfig: setting.CacheConn,
		Interval:      setting.CacheInternal,
	}))

	m.Use(captcha.Captchaer(captcha.Options{
		SubURL: setting.AppSubURL,
	}))

	m.Use(Sessioner(setting.SessionConfig))

	m.Use(csrf.Csrfer(csrf.Options{
		Secret:     setting.SecretKey,
		Cookie:     setting.CSRFCookieName,
		SetCookie:  true,
		Header:     "X-Csrf-Token",
		CookiePath: setting.AppSubURL,
	}))

	m.Use(toolbox.Toolboxer(m, toolbox.Options{
		HealthCheckFuncs: []*toolbox.HealthCheckFuncDesc{
			&toolbox.HealthCheckFuncDesc{
				Desc: "Database connection",
				Func: model.Ping,
			},
		},
	}))
	m.Use(context.Contexter())
	m.Use(context.HTMLContexter())
	m.Use(context.APIContexter())

	return m
}

func runWeb(cliCtx *cli.Context) {

	if cliCtx != nil && cliCtx.IsSet("config") {
		setting.CustomConf = cliCtx.String("config")
	}

	router.GlobalInit()
	checkVersion()

	m := newMacaron()

	apiReqSignIn := context.APIToggle(&context.ToggleOptions{SignInRequire: true, DisableCSRF: true})
	apiIgnSignIn := context.APIToggle(&context.ToggleOptions{SignInRequire: false})

	reqSignIn := context.Toggle(&context.ToggleOptions{SignInRequire: true})
	ignSignIn := context.Toggle(&context.ToggleOptions{SignInRequire: setting.Service.RequireSignInView})
	reqSignOut := context.Toggle(&context.ToggleOptions{SignOutRequire: true})
	reqAdmin := context.Toggle(&context.ToggleOptions{SignInRequire: true, AdminRequire: true})
	bindIgnErr := binding.BindIgnErr

	apiBind := binding.Bind

	m.Get("/", ignSignIn, router.Home)
	m.Combo("/install", router.InstallInit).Get(router.Install).Post(bindIgnErr(forms.InstallForm{}), router.InstallPost)

	m.Group("/web", func() {
		m.Group("/auth", func() {
			m.Combo("/signup").Get(web.SignUpGet).Post(bindIgnErr(forms.SignUpForm{}), web.SignUpPost)
			m.Combo("/signin").Get(web.SignInGet).Post(bindIgnErr(forms.SignInForm{}), web.SignInPost)
			m.Get("/signout", web.SignOutGet)
			m.Get("/activate", web.ActivateGet)
		}, ignSignIn)

		m.Group("/app", func() {
			m.Combo("/list").Get(web.AppsGet)
			m.Combo("/add").Post(bindIgnErr(forms.AppAddForm{}), web.AppAdd)
			m.Get("/detail/:appName/collaborators", web.AppCollaboratorsGet)
			m.Get("/detail/:appName/deployments", web.AppDeploymentsGet)
		}, reqSignIn)

		m.Group("/access_key", func() {
			m.Combo("/list").Get(web.AccessKeysGet)
		}, reqSignIn)

		m.Group("/admin/api", func() {
			m.Get("/users/:pageIndex/:pageCount", web.UsersQuery)
			m.Post("/users/add", web.UserAddPost)
			m.Patch("/users/:email/status", web.UserPatch)
			m.Patch("/users/:email/role", web.UserPatch)
		}, reqAdmin)

		m.Group("/admin", func() {
			m.Get("/users", web.UsersGet)
		}, reqAdmin)

	})

	//cli
	m.Group("/auth", func() {
		m.Combo("/login").Get(auth.SignInGet)
		m.Combo("/register").Get(auth.SignUpGet)
		m.Post("/logout", auth.SignOutPost)
		m.Get("/link", auth.LinkGet)
	}, reqSignOut)

	m.Get("/README.md", ignSignIn, router.ReadMeGet)

	m.Get("/tokens", ignSignIn, router.TokensGet)

	// m.Delete("/sessions/:machineName", apiReqSignIn, session.Delete)
	m.Delete("/sessions/:machineName", apiReqSignIn, auth.SignOutPost)
	m.Get("/authenticated", apiIgnSignIn, auth.Authenticated)

	m.Get("/updateCheck", apiIgnSignIn, client.UpdateGet)

	m.Group("/reportStatus", func() {
		m.Post("/download", apiBind(forms.ReportStatus{}), client.Download)
		m.Post("/deploy", apiBind(forms.ReportStatus{}), client.Deploy)
	}, apiIgnSignIn)

	m.Group("/accessKeys", func() {
		m.Combo("/").Get(accessKey.Get).Post(accessKey.Add)
		m.Combo("/:name").Delete(accessKey.Delete).Patch(accessKey.Patch)
	}, apiReqSignIn)

	m.Get("/account", apiReqSignIn, account.Info)

	m.Group("/apps", func() {
		m.Combo("/").Get(app.Get).Post(apiBind(forms.AppOrDeploymentOption{}), app.Post)
		m.Group("/:appName", func() {
			m.Combo("/").Delete(app.Delete).Patch(apiBind(forms.AppOrDeploymentOption{}), app.Patch)
			m.Post("/transfer/:email", app.TransferPost)
			m.Group("/collaborators", func() {
				m.Get("/", collaborator.Get)
				m.Combo("/:email").Post(collaborator.Post).Delete(collaborator.Delete)
			})
			m.Group("/deployments", func() {
				m.Combo("/").Get(deployment.Get).Post(apiBind(forms.AppOrDeploymentOption{}), deployment.Post)
				m.Post("/:sourceDeploymentName/promote/:destDeploymentName", deployment.DestPost)
				m.Group("/:deploymentName", func() {
					m.Combo("/").Patch(apiBind(forms.AppOrDeploymentOption{}), deployment.Patch).Delete(deployment.Delete)
					m.Group("/rollback", func() {
						m.Post("/", deployment.RobackPost)
						m.Post("/:label", deployment.RobackForLabelPost)
					})
					m.Combo("/release").Post(deployment.ReleasePost).Patch(apiBind(forms.UpdatePackageInfo{}), deployment.ReleasePatch)
					m.Get("/metrics", deployment.MetricsGet)
					m.Combo("/history").Get(deployment.HistoryGet).Delete(deployment.HistoryDelete)
				})
			})
		})
	}, apiReqSignIn)

	if setting.Storage.StorageType == "local" {
		urlPath := strings.Replace(setting.Storage.StorageConfig.DownloadURL, setting.AppURL, "", -1)
		m.Get(path.Join(urlPath, ":fileName"), func(ctx *macaron.Context) {
			ctx.ServeFile(path.Join(setting.Storage.StorageConfig.LocalStoragePath, ctx.Params("fileName")))
		})
	}

	if macaron.Env == macaron.DEV {
		m.Get("/template/*", dev.TemplatePreview)
	}

	m.NotFound(router.NotFound)

	if cliCtx != nil && cliCtx.IsSet("port") {
		setting.AppURL = strings.Replace(setting.AppURL, setting.HTTPPort, cliCtx.String("port"), 1)
		setting.HTTPPort = cliCtx.String("port")
	}

	var err error
	listenAddr := fmt.Sprintf("%s:%s", setting.HTTPAddr, setting.HTTPPort)
	log.Info("Listen: %v://%s", setting.Protocol, listenAddr)
	switch setting.Protocol {
	case setting.SCHEME_HTTP:
		err = http.ListenAndServe(listenAddr, m)
	case setting.SCHEME_HTTPS:
		err = http.ListenAndServeTLS(listenAddr, setting.CertFile, setting.KeyFile, m)
	default:
		log.Fatal(4, "Invalid protocol: %s", setting.Protocol)
	}

	if err != nil {
		log.Fatal(4, "Fail to start server: %v", err)
	}
}
