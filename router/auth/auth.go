package auth

import (
	"path"
	"strings"

	"github.com/MessageDream/goby/core/tokenService"
	"github.com/MessageDream/goby/module/context"
	"github.com/MessageDream/goby/module/setting"
)

func SignInGet(ctx *context.Context) {
	ctx.Redirect("/web/auth/signin")
}

func SignUpGet(ctx *context.Context) {
	ctx.Redirect("/web/auth/signup")
}

func SignOutGet(ctx *context.Context) {
	ctx.Session.Delete("uid")
	ctx.Session.Delete("uname")
	// ctx.Session.Delete("socialId")
	// ctx.Session.Delete("socialName")
	// ctx.Session.Delete("socialEmail")
	ctx.SetCookie(setting.CookieUserName, "", -1, setting.AppSubURL)
	ctx.SetCookie(setting.CookieRememberName, "", -1, setting.AppSubURL)
	ctx.SetCookie(setting.CSRFCookieName, "", -1, setting.AppSubURL)
	ctx.Redirect(path.Join(setting.AppSubURL, "/"))
}

func LinkGet(ctx *context.Context) {
	ctx.Redirect(path.Join(setting.AppSubURL, "auth/login"))
}

func Authenticated(ctx *context.APIContext) {

	err := "authorize failed"

	auHead := ctx.Req.Header.Get("Authorization")
	if len(auHead) <= 0 {
		ctx.Context.Error(401, err)
		return
	}

	auths := strings.Fields(auHead)
	if len(auths) < 2 || auths[0] != "Bearer" {
		ctx.Context.Error(401, err)
		return
	}
	if er := tokenService.AuthWithToken(auths[1], true); er != nil {
		ctx.Context.Error(401, err)
		return
	}

	ctx.Status(200)
}

func SignOutPost(ctx *context.APIContext) {
	token := ctx.GetToken()
	if len(token) > 0 {
		if err := tokenService.AuthWithToken(token, false); err != nil {
			ctx.Error(err)
			return
		}
		ctx.Cache.Delete(token)
	}
	ctx.Status(200)
}
