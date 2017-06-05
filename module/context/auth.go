package context

import (
	"net/url"
	"path"
	"strings"

	"9fbank.com/goby/module/infrastructure"

	"github.com/go-macaron/csrf"
	"gopkg.in/macaron.v1"
)

type ToggleOptions struct {
	SignInRequire  bool
	SignOutRequire bool
	AdminRequire   bool
	DisableCSRF    bool
}

const (
	TOKEN_TYPE_AUTH = iota
	TOKEN_TYPE_ACCESS
)

func IsAPIPath(url string) bool {
	return !strings.HasPrefix(url, "/web/")
}

func defaultHandle(options *ToggleOptions, ctx *Context) bool {
	if !appInstallLock {
		ctx.Redirect(path.Join(subURL, "install"))
		return false
	}

	if options.SignOutRequire && ctx.IsSigned && ctx.Req.RequestURI != "/" {
		ctx.Redirect(path.Join(subURL, "/"))
		return false
	}

	return true
}

func APIToggle(options *ToggleOptions) macaron.Handler {
	return func(ctx *APIContext) {
		if !defaultHandle(options, ctx.Context) {
			return
		}

		if options.SignInRequire {
			ctx.Auth()
			if !ctx.IsSigned {
				ctx.Status(401)
				return

			} else if !ctx.User.IsActive && registerEmailConfirm {
				ctx.JSON(200, infrastructure.ApiJsonError{
					Message: "You should active your account",
					Status:  401,
				})
				return
			}
		}
	}
}

func Toggle(options *ToggleOptions) macaron.Handler {
	return func(ctx *HTMLContext) {
		if !defaultHandle(options, ctx.Context) {
			return
		}

		if !options.SignOutRequire && !options.DisableCSRF && ctx.Req.Method == "POST" {
			csrf.Validate(ctx.Context.Context, ctx.csrf)
			if ctx.Written() {
				return
			}
		}

		if options.SignInRequire {

			if !ctx.IsSigned {
				ctx.SetCookie("redirect_to", "/"+url.QueryEscape(ctx.Req.RequestURI))
				ctx.Redirect(path.Join(subURL, "/web/auth/signin"))
				return

			} else if !ctx.User.IsActive && registerEmailConfirm {
				ctx.Data["Title"] = "激活您的账户"
				ctx.HTML(200, "auth/activate")
				return
			}
		}

		if options.AdminRequire {
			if !ctx.User.IsAdmin {
				ctx.Error(403)
				return
			}
			ctx.Data["PageIsAdmin"] = true
		}
	}
}
