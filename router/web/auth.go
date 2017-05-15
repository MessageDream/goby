package web

import (
	"fmt"
	"net/url"
	"path"

	"github.com/go-macaron/captcha"
	log "gopkg.in/clog.v1"

	. "github.com/MessageDream/goby/core"
	"github.com/MessageDream/goby/core/tokenService"
	"github.com/MessageDream/goby/core/userService"
	"github.com/MessageDream/goby/model"
	"github.com/MessageDream/goby/module/context"
	forms "github.com/MessageDream/goby/module/form"
	"github.com/MessageDream/goby/module/infrastructure"
	"github.com/MessageDream/goby/module/mailer"
	"github.com/MessageDream/goby/module/setting"
)

const (
	SIGNUP   infrastructure.TplName = "auth/signup"
	SIGNIN   infrastructure.TplName = "auth/signin"
	ACTIVATE infrastructure.TplName = "auth/activate"
)

func SignUpGet(ctx *context.Context) {
	ctx.Data["Title"] = "注册"

	ctx.Data["EnableCaptcha"] = setting.Service.EnableCaptcha

	if setting.Service.DisableRegistration {
		ctx.Data["DisableRegistration"] = true
		ctx.HTML(200, SIGNUP)
		return
	}

	ctx.HTML(200, SIGNUP)
}

func SignUpPost(ctx *context.Context, cpt *captcha.Captcha, form forms.SignUpForm) {

	ctx.Data["Title"] = "注册"

	ctx.Data["EnableCaptcha"] = setting.Service.EnableCaptcha

	if setting.Service.DisableRegistration {
		ctx.Error(403)
		return
	}

	if ctx.HasError() {
		ctx.HTML(200, SIGNUP)
		return
	}

	if setting.Service.EnableCaptcha && !cpt.VerifyReq(ctx.Req) {
		ctx.Data["Err_Captcha"] = true
		ctx.RenderWithErr("验证码未匹配。", SIGNUP, &form)
		return
	}

	if form.Password != form.Retype {
		ctx.Data["Err_Password"] = true
		ctx.RenderWithErr("密码与确认密码未匹配", SIGNUP, &form)
		return
	}

	u, err := userService.Create(form.UserName, form.Password, form.Email, !setting.Service.RegisterEmailConfirm, false)

	if err != nil {
		switch err {
		case ErrUserNameAlreadyExist:
			ctx.Data["Err_UserName"] = true
			ctx.RenderWithErr("用户名已经被占用", SIGNUP, &form)
		case ErrEmailAlreadyExist:
			ctx.Data["Err_Email"] = true
			ctx.RenderWithErr("邮箱已被注册", SIGNUP, &form)
		case ErrNameReserved:
			ctx.Data["Err_UserName"] = true
			ctx.RenderWithErr(fmt.Sprintf("用户名 '%s' 是被保留的。", u.UserName), SIGNUP, &form)
		default:
			if err.(*infrastructure.IntentError).Code == INTENT_ERROR_CODE_USER_NAME_NOT_ALLOWNED {
				ctx.Data["Err_UserName"] = true
				ctx.RenderWithErr("用户名不允许此格式", SIGNUP, &form)
			} else {
				ctx.Handle(500, "CreateUser", err)
			}
		}
		return
	}

	log.Trace("Account created: %s", u.UserName)

	if setting.Service.RegisterEmailConfirm && u.ID > 1 {
		mailer.SendActivateAccountMail(ctx.Context, model.NewMailerUser(u, setting.Service.ActiveCodeLives))
		ctx.Data["IsSendRegisterMail"] = true
		ctx.Data["Email"] = u.Email
		ctx.Data["Hours"] = setting.Service.ActiveCodeLives / 60
		ctx.HTML(200, ACTIVATE)

		if err := ctx.Cache.Put("MailResendLimit_"+u.LowerName, u.LowerName, 180); err != nil {
			log.Error(4, "Set cache(MailResendLimit) fail: %v", err)
		}
		return
	}

	if !setting.Service.RegisterEmailConfirm {
		createToken(u)
	}

	ctx.Redirect(path.Join(setting.AppSubURL, "web/auth/signin"))
}

func SignInGet(ctx *context.Context) {
	ctx.Data["Title"] = "登陆"

	// Check auto-login.
	isSucceed, err := AutoSignIn(ctx)
	if err != nil {
		ctx.Handle(500, "AutoSignIn", err)
		return
	}

	redirectTo := ctx.Query("redirect_to")
	if len(redirectTo) > 0 {
		ctx.SetCookie("redirect_to", redirectTo, 0, setting.AppSubURL)
	} else {
		redirectTo, _ = url.QueryUnescape(ctx.GetCookie("redirect_to"))
		ctx.SetCookie("redirect_to", "", -1, setting.AppSubURL)
	}

	if isSucceed {
		if isValidRedirect(redirectTo) {
			ctx.Redirect(redirectTo)
		} else {
			ctx.Redirect(path.Join(setting.AppSubURL, "/"))
		}
		return
	}

	ctx.HTML(200, SIGNIN)
}

func SignInPost(ctx *context.Context, form forms.SignInForm) {
	ctx.Data["Title"] = "登陆"

	if ctx.HasError() {
		ctx.HTML(200, SIGNIN)
		return
	}

	u, err := userService.Signin(form.UserName, form.Password)
	if err != nil {
		if err == ErrUserNameOrPasswordInvalide || err == ErrUserNotExist {
			ctx.RenderWithErr("用户名或密码不正确。", SIGNIN, &form)
		} else {
			ctx.Handle(500, "userService.Signin", err)
		}
		return
	}

	if form.Remember {
		days := 86400 * setting.LogInRememberDays
		ctx.SetCookie(setting.CookieUserName, u.UserName, days, setting.AppSubURL, "", setting.CookieSecure, true)
		ctx.SetSuperSecureCookie(u.Rands+u.Password, setting.CookieRememberName, u.UserName, days, setting.AppSubURL, "", setting.CookieSecure, true)
	}

	ctx.Session.Set("u", u)
	ctx.Session.Set("uname", u.UserName)

	// Clear whatever CSRF has right now, force to generate a new one
	ctx.SetCookie(setting.CSRFCookieName, "", -1, setting.AppSubURL)

	redirectTo, _ := url.QueryUnescape(ctx.GetCookie("redirect_to"))
	ctx.SetCookie("redirect_to", "", -1, setting.AppSubURL)
	if isValidRedirect(redirectTo) {
		ctx.Redirect(redirectTo)
		return
	}

	ctx.Redirect(path.Join(setting.AppSubURL, "/"))
}

func SignOutGet(ctx *context.Context) {
	ctx.Session.Delete("u")
	ctx.Session.Delete("uname")
	ctx.Session.Delete("socialId")
	ctx.Session.Delete("socialName")
	ctx.Session.Delete("socialEmail")
	ctx.SetCookie(setting.CookieUserName, "", -1, setting.AppSubURL)
	ctx.SetCookie(setting.CookieRememberName, "", -1, setting.AppSubURL)
	ctx.SetCookie(setting.CSRFCookieName, "", -1, setting.AppSubURL)
	ctx.Redirect(path.Join(setting.AppSubURL, "/"))
}

func ActivateGet(ctx *context.Context) {
	code := ctx.Query("code")
	if len(code) == 0 {
		ctx.Data["IsActivatePage"] = true
		if ctx.User.IsActive {
			ctx.Error(404)
			return
		}
		// Resend confirmation email.
		if setting.Service.RegisterEmailConfirm {
			if ctx.Cache.IsExist("MailResendLimit_" + ctx.User.LowerName) {
				ctx.Data["ResendLimited"] = true
			} else {
				ctx.Data["Hours"] = setting.Service.ActiveCodeLives / 60
				mailer.SendActivateAccountMail(ctx.Context, model.NewMailerUser(ctx.User, setting.Service.ActiveCodeLives))

				if err := ctx.Cache.Put("MailResendLimit_"+ctx.User.LowerName, ctx.User.LowerName, 180); err != nil {
					log.Error(4, "Set cache(MailResendLimit) fail: %v", err)
				}
			}
		} else {
			ctx.Data["ServiceNotEnabled"] = true
		}
		ctx.HTML(200, ACTIVATE)
		return
	}

	user, err := userService.Active(code)

	if err != nil {
		switch err {
		case ErrUserAlreadyActivated:
		case ErrUserNotExist:
			ctx.Error(404)
		case ErrUserActivateVerifyFailed:
			ctx.Data["IsActivateFailed"] = true
			ctx.HTML(200, ACTIVATE)
		default:
			ctx.Handle(500, "userService.Active", err)
		}

		return
	}

	log.Trace("User activated: %s", user.UserName)

	ctx.Session.Set("u", user)
	ctx.Session.Set("uname", user.UserName)

	createToken(user)

	ctx.Redirect(path.Join(setting.AppSubURL, "/"))
}

func AutoSignIn(ctx *context.Context) (bool, error) {
	if !model.HasEngine {
		return false, nil
	}

	uname := ctx.GetCookie(setting.CookieUserName)
	if len(uname) == 0 {
		return false, nil
	}

	isSucceed := false
	defer func() {
		if !isSucceed {
			log.Trace("auto-login cookie cleared: %s", uname)
			ctx.SetCookie(setting.CookieUserName, "", -1, setting.AppSubURL)
			ctx.SetCookie(setting.CookieRememberName, "", -1, setting.AppSubURL)
		}
	}()

	u, err := userService.GetByUserName(uname)
	if err != nil {
		if err != ErrUserNotExist {
			return false, fmt.Errorf("GetUserByName: %v", err)
		}
		return false, nil
	}

	if val, ok := ctx.GetSuperSecureCookie(u.Rands+u.Password, setting.CookieRememberName); !ok || val != u.UserName {
		return false, nil
	}

	isSucceed = true
	ctx.Session.Set("u", u)
	ctx.Session.Set("uname", u.UserName)
	ctx.SetCookie(setting.CSRFCookieName, "", -1, setting.AppSubURL)
	return true, nil
}

func isValidRedirect(url string) bool {
	return len(url) >= 2 && url[0] == '/' && url[1] != '/'
}

func createToken(user *model.User) {
	creator := setting.AppName + "-" + user.UserName

	if _, err := tokenService.Create(user, creator, creator+"-auto", setting.AppName+"-auto", 0); err != nil {
		log.Warn("tokenService.Create:%v", err)
	}
}
