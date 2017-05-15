package router

import (
	"os"
	"path"
	"path/filepath"

	"errors"
	"fmt"
	"strings"

	"github.com/Unknwon/com"
	"github.com/go-xorm/xorm"

	log "gopkg.in/clog.v1"
	ini "gopkg.in/ini.v1"
	"gopkg.in/macaron.v1"

	"github.com/MessageDream/goby/core/userService"
	"github.com/MessageDream/goby/model"
	"github.com/MessageDream/goby/module/context"
	forms "github.com/MessageDream/goby/module/form"
	"github.com/MessageDream/goby/module/infrastructure"
	"github.com/MessageDream/goby/module/mailer"
	"github.com/MessageDream/goby/module/setting"
	"github.com/MessageDream/goby/module/storage"
	"github.com/MessageDream/goby/module/template"
)

const (
	INSTALL infrastructure.TplName = "install"
)

var (
	dbTypesWrite      = map[string]string{"MySQL": "mysql", "PostgreSQL": "postgres", "SQLite3": "sqlite3"}
	storageTypesWrite = map[string]string{"本地": "local", "七牛": "qiniu", "阿里云OSS": "oss"}

	dbTypesRead      = map[string]string{"mysql": "MySQL", "postgres": "PostgreSQL", "sqlite3": "SQLite3"}
	storageTypesRead = map[string]string{"local": "本地", "qiniu": "七牛", "oss": "阿里云OSS"}
)

func checkRunMode() {
	if setting.ProdMode {
		macaron.Env = macaron.PROD
		macaron.ColorLog = false
	}
	log.Info("Run Mode: %s", strings.Title(macaron.Env))
}

func initModules() {
	infrastructure.InitInfrastructure(setting.AppName, setting.SecretKey, setting.TimeFormat)

	template.InitTemplate(setting.AppName, setting.AppURL, setting.AppSubURL, setting.Domain, setting.AppVer)

	context.InitContextConfig(setting.Service.EnableReverseProxyAuth,
		setting.InstallLock,
		setting.Service.RegisterEmailConfirm,
		setting.Service.DisableRegistration,
		setting.CacheTokenTimeOut,
		setting.AttachmentMaxSize,
		setting.AppSubURL,
		setting.ReverseProxyAuthUser)

	model.InitDBConfig(&setting.DBConfig)

	storage.InitStorage(setting.Storage.StorageType, setting.Storage.StorageConfig)

	mailer.InitMail(setting.Service.ActiveCodeLives, setting.Service.ResetPwdCodeLives, &macaron.RenderOptions{
		Directory:         path.Join(setting.StaticRootPath, "template/mail"),
		AppendDirectories: []string{path.Join(setting.CustomPath, "template/mail")},
		Funcs:             template.NewFuncMap(),
		Extensions:        []string{".tmpl", ".html"},
	}, setting.MailService)

	// social.NewOauthService()
}

// GlobalInit is for global configuration reload-able.
func GlobalInit() {
	setting.InitConfig()
	log.Trace("Custom path: %s", setting.CustomPath)
	log.Trace("Log path: %s", setting.LogRootPath)
	initModules()

	if setting.InstallLock {
		if err := model.NewEngine(); err != nil {
			log.Fatal(4, "Fail to initialize ORM engine: %v", err)
		}
		model.HasEngine = true
	}
	checkRunMode()
}

func InstallInit(ctx *context.Context) {
	if setting.InstallLock {
		ctx.Handle(404, "Install", errors.New("Installation is prohibited"))
		return
	}

	ctx.Data["Title"] = "安装首页"
	ctx.Data["PageIsInstall"] = true

	dbOpts := []string{"MySQL", "PostgreSQL", "TiDB"}
	if model.EnableSQLite3 {
		dbOpts = append(dbOpts, "SQLite3")
	}
	ctx.Data["DbOptions"] = dbOpts
	ctx.Data["CurDbOption"] = dbTypesRead[setting.DBConfig.Type]

	storageOpts := []string{"本地", "七牛", "阿里云OSS"}
	ctx.Data["StorageOptions"] = storageOpts
	ctx.Data["CurStorageOption"] = storageTypesRead[setting.Storage.StorageType]
}

// @router /install [get]
func Install(ctx *context.Context) {
	form := forms.InstallForm{}

	// Get and assign values to install form.
	if len(form.DbHost) == 0 {
		form.DbHost = setting.DBConfig.Host
	}
	if len(form.DbUser) == 0 {
		form.DbUser = setting.DBConfig.User
	}
	if len(form.DbPasswd) == 0 {
		form.DbPasswd = setting.DBConfig.Passwd
	}

	if len(form.DbName) == 0 {
		form.DbName = setting.DBConfig.Name
	}
	if len(form.DbPath) == 0 {
		form.DbPath = setting.DBConfig.Path
	}

	if len(form.RunUser) == 0 {
		form.RunUser = setting.RunUser
	}
	if len(form.Domain) == 0 {
		form.Domain = setting.Domain
	}
	if len(form.AppURL) == 0 {
		form.AppURL = setting.AppURL
	}

	if len(form.StoragePath) == 0 {
		form.StoragePath = setting.Storage.StorageConfig.LocalStoragePath
	}

	forms.AssignForm(form, ctx.Data)
	ctx.HTML(200, INSTALL)
}

func InstallPost(ctx *context.Context, form forms.InstallForm) {
	if setting.InstallLock {
		ctx.Handle(404, "InstallPost", errors.New("Installation is prohibited"))
		return
	}

	ctx.Data["CurDbOption"] = form.DbType
	if ctx.HasError() {
		ctx.HTML(200, INSTALL)
		return
	}

	// if _, err := exec.LookPath("node"); err != nil {
	// 	ctx.RenderWithErr(fmt.Sprintf("测试node错误%v", err), INSTALL, &form)
	// 	return
	// }

	setting.DBConfig.Type = dbTypesWrite[form.DbType]
	setting.DBConfig.Host = form.DbHost
	setting.DBConfig.User = form.DbUser
	setting.DBConfig.Passwd = form.DbPasswd
	setting.DBConfig.Name = form.DbName
	setting.DBConfig.SSLMode = form.SSLMode
	setting.DBConfig.Path = form.DbPath

	log.Info("-------%v", setting.DBConfig)

	if setting.DBConfig.Type == "sqlite3" && len(setting.DBConfig.Path) == 0 {
		ctx.Data["Err_DbPath"] = true
		ctx.RenderWithErr("SQLite 数据库文件路径不能为空", INSTALL, &form)
		return
	}

	var x *xorm.Engine
	if err := model.NewTestEngine(x); err != nil {
		if strings.Contains(err.Error(), `Unknown database type: sqlite3`) {
			ctx.RenderWithErr("您所使用的发行版不支持 SQLite3，请下载官方构建版，而不是 gobuild 版本", INSTALL, &form)
		} else {
			ctx.RenderWithErr(fmt.Sprintf("数据库设置不正确：%v", err), INSTALL, &form)
		}
		return
	}

	// Check run user.
	curUser := os.Getenv("USER")
	if len(curUser) == 0 {
		curUser = os.Getenv("USERNAME")
	}
	// Does not check run user when the install lock is off.
	if form.RunUser != curUser {
		ctx.RenderWithErr(fmt.Sprintf("运行系统用户非当前用户：%s -> %s", form.RunUser, curUser), INSTALL, &form)
		return
	}

	// Check admin password.
	if form.AdminPasswd != form.ConfirmPasswd {
		ctx.RenderWithErr("密码与确认密码未匹配。", INSTALL, form)
		return
	}

	if form.AppURL[len(form.AppURL)-1] != '/' {
		form.AppURL += "/"
	}

	// Save settings.
	cfg := ini.Empty()
	if com.IsFile(setting.CustomConf) {
		// Keeps custom settings if there is already something.
		if err := cfg.Append(setting.CustomConf); err != nil {
			log.Error(4, "Fail to load custom conf '%s': %v", setting.CustomConf, err)
		}
	}

	//database
	cfg.Section("database").Key("DB_TYPE").SetValue(setting.DBConfig.Type)
	cfg.Section("database").Key("HOST").SetValue(setting.DBConfig.Host)
	cfg.Section("database").Key("NAME").SetValue(setting.DBConfig.Name)
	cfg.Section("database").Key("USER").SetValue(setting.DBConfig.User)
	cfg.Section("database").Key("PASSWD").SetValue(setting.DBConfig.Passwd)
	cfg.Section("database").Key("SSL_MODE").SetValue(setting.DBConfig.SSLMode)
	cfg.Section("database").Key("PATH").SetValue(setting.DBConfig.Path)

	//storage
	cfg.Section("storage").Key("STORAGE_TYPE").SetValue(storageTypesWrite[form.StorageType])
	cfg.Section("storage").Key("STORAGE_PATH").SetValue(form.StoragePath)
	cfg.Section("storage").Key("ACCESS_KEY").SetValue(form.StorageAccessKey)
	cfg.Section("storage").Key("SECRET_KEY").SetValue(form.StorageSecretKey)
	cfg.Section("storage").Key("BUCKET").SetValue(form.StorageBucketName)
	cfg.Section("storage").Key("QN_ZONE").SetValue(form.StorageZone)
	cfg.Section("storage").Key("OSS_ENDPOINT").SetValue(form.StorageEndpoint)
	cfg.Section("storage").Key("PREFIX").SetValue(form.StoragePrefix)
	cfg.Section("storage").Key("DOWNLOAD_URL").SetValue(form.StorageDownloadURL)

	//app
	cfg.Section("").Key("APP_NAME").SetValue(form.AppName)
	cfg.Section("").Key("RUN_USER").SetValue(form.RunUser)
	cfg.Section("server").Key("DOMAIN").SetValue(form.Domain)
	cfg.Section("server").Key("HTTP_PORT").SetValue(form.HTTPPort)
	cfg.Section("server").Key("ROOT_URL").SetValue(form.AppURL)

	//mail
	if len(strings.TrimSpace(form.SMTPHost)) > 0 {
		cfg.Section("mailer").Key("ENABLED").SetValue("true")
		cfg.Section("mailer").Key("HOST").SetValue(form.SMTPHost)
		cfg.Section("mailer").Key("FROM").SetValue(form.SMTPFrom)
		cfg.Section("mailer").Key("USER").SetValue(form.SMTPUser)
		cfg.Section("mailer").Key("PASSWD").SetValue(form.SMTPPasswd)
	} else {
		cfg.Section("mailer").Key("ENABLED").SetValue("false")
	}

	cfg.Section("service").Key("REGISTER_EMAIL_CONFIRM").SetValue(com.ToStr(form.RegisterConfirm))
	cfg.Section("service").Key("ENABLE_NOTIFY_MAIL").SetValue(com.ToStr(form.MailNotify))
	cfg.Section("service").Key("DISABLE_REGISTRATION").SetValue(com.ToStr(form.DisableRegistration))
	cfg.Section("service").Key("ENABLE_CAPTCHA").SetValue(com.ToStr(form.EnableCaptcha))
	cfg.Section("service").Key("REQUIRE_SIGNIN_VIEW").SetValue(com.ToStr(form.RequireSignInView))

	cfg.Section("").Key("RUN_MODE").SetValue("prod")

	cfg.Section("session").Key("PROVIDER").SetValue("file")

	cfg.Section("log").Key("MODE").SetValue("file")
	cfg.Section("log").Key("LEVEL").SetValue("Info")
	cfg.Section("log").Key("ROOT_PATH").SetValue(form.LogRootPath)

	cfg.Section("security").Key("INSTALL_LOCK").SetValue("true")
	secretKey := infrastructure.GetRandomString(15)
	cfg.Section("security").Key("SECRET_KEY").SetValue(secretKey)

	os.MkdirAll(filepath.Dir(setting.CustomConf), os.ModePerm)
	if err := cfg.SaveTo(setting.CustomConf); err != nil {
		ctx.RenderWithErr(fmt.Sprintf("应用配置保存失败：%v", err), INSTALL, &form)
		return
	}

	GlobalInit()

	// Create admin account.
	if _, err := userService.Create(form.AdminName, form.AdminPasswd, form.AdminEmail, true, true); err != nil {
		if err != nil {
			setting.InstallLock = false
			ctx.RenderWithErr(fmt.Sprintf("管理员帐户设置不正确:%v", err), INSTALL, &form)
			return
		}
		log.Info("Admin account already exist")
	}

	log.Info("First-time run install finished!")
	ctx.Flash.Success("安装成功")
	ctx.Redirect(path.Join(setting.AppSubURL, "/"))
}
