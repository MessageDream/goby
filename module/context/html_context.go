package context

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/go-macaron/cache"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	log "gopkg.in/clog.v1"
	macaron "gopkg.in/macaron.v1"

	"github.com/MessageDream/goby/core/userService"
	"github.com/MessageDream/goby/model"
	forms "github.com/MessageDream/goby/module/form"
	"github.com/MessageDream/goby/module/infrastructure"
)

type HTMLContext struct {
	*Context
	csrf    csrf.CSRF
	Flash   *session.Flash
	Session session.Store
}

func (ctx *HTMLContext) GetErrMsg() string {
	return ctx.Data["ErrorMsg"].(string)
}

func (ctx *HTMLContext) HasError() bool {
	hasErr, ok := ctx.Data["HasError"]
	if !ok {
		return false
	}
	ctx.Flash.ErrorMsg = ctx.Data["ErrorMsg"].(string)
	ctx.Data["Flash"] = ctx.Flash
	return hasErr.(bool)
}

func (ctx *HTMLContext) HTML(status int, name infrastructure.TplName) {
	ctx.Context.HTML(status, string(name))
}

func (ctx *HTMLContext) RenderWithErr(msg string, tpl infrastructure.TplName, form interface{}) {
	if form != nil {
		forms.AssignForm(form, ctx.Data)
	}
	ctx.Flash.ErrorMsg = msg
	ctx.Data["Flash"] = ctx.Flash
	ctx.HTML(200, tpl)
}

func (ctx *HTMLContext) Handle(status int, title string, err error) {
	if err != nil {
		log.Error(4, "%s: %v", title, err)
		if macaron.Env != macaron.PROD {
			ctx.Data["ErrorMsg"] = err
		}
	}

	switch status {
	case 404:
		ctx.Data["Title"] = "Page Not Found"
	case 500:
		ctx.Data["Title"] = "Internal Server Error"
	}
	ctx.HTML(status, infrastructure.TplName(fmt.Sprintf("status/%d", status)))
}

func (ctx *HTMLContext) Auth() {
	if !model.HasEngine {
		return
	}
	user := ctx.Session.Get("u")
	if user == nil {
		if enableReverseProxyAuth {
			webAuthUser := ctx.Req.Header.Get(reverseProxyAuthUser)
			if len(webAuthUser) > 0 {
				u, err := userService.SignInWithUserName(webAuthUser)
				if err != nil {
					ctx.SignError = err
					return
				}
				ctx.User = u
				ctx.IsSigned = true
				ctx.Data["IsSigned"] = ctx.IsSigned
				ctx.Data["SignedUser"] = ctx.User
			}
		}
		return
	}

	if usr, ok := user.(*model.User); ok {
		ctx.Session.Set("u", usr)
		ctx.User = usr
		ctx.IsSigned = true
		ctx.Data["IsSigned"] = ctx.IsSigned
		ctx.Data["SignedUser"] = ctx.User
	}
}

func HTMLContexter() macaron.Handler {
	return func(c *Context, cache cache.Cache, sess session.Store, f *session.Flash, x csrf.CSRF) {

		ctx := &HTMLContext{
			Context: c,
			csrf:    x,
			Flash:   f,
			Session: sess,
		}

		ctx.Auth()

		link := ctx.Req.RequestURI
		i := strings.Index(link, "?")
		if i > -1 {
			link = link[:i]
		}
		ctx.Data["Link"] = link
		ctx.Data["PageStartTime"] = time.Now()
		ctx.Data["ShowRegistrationButton"] = !disableRegistration
		csrfToken := x.GetToken()
		ctx.Data["CsrfToken"] = csrfToken
		ctx.Data["CsrfTokenHtml"] = template.HTML(`<input type="hidden" name="_csrf" value="` + csrfToken + `">`)

		log.Trace("Session ID: %s", sess.ID())
		log.Trace("CSRF Token: %v", csrfToken)

		c.Map(ctx)
	}
}
