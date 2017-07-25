package router

import (
	"path"

	"github.com/MessageDream/goby/model"
	"github.com/MessageDream/goby/module/context"
	"github.com/MessageDream/goby/module/infrastructure"
	"github.com/MessageDream/goby/module/setting"
)

const (
	HOME     infrastructure.TplName = "home"
	ACTIVATE infrastructure.TplName = "auth/activate"
)

func Home(ctx *context.HTMLContext) {
	if ctx.IsSigned {
		if ctx.User.Status == model.USER_STATUS_UN_ACTIVE && setting.Service.RegisterEmailConfirm {
			ctx.Data["Title"] = "激活您的账户"
			ctx.HTML(200, ACTIVATE)
			return
		}

		ctx.Data["PageIsHome"] = true
		ctx.HTML(200, HOME)
		return
	}

	//Check auto-login.
	uname := ctx.GetCookie(setting.CookieUserName)
	if len(uname) != 0 {
		ctx.Redirect(path.Join(setting.AppSubURL, "web/signin"))
		return
	}

	ctx.Data["PageIsHome"] = true
	ctx.HTML(200, HOME)
}

func ReadMeGet(ctx *context.APIContext) {

}

func TokensGet(ctx *context.APIContext) {

}

func NotFound(ctx *context.HTMLContext) {
	ctx.Data["Title"] = "Page Not Found"
	ctx.Handle(404, "home.NotFound", nil)
}
