package context

import (
	"net/url"
	"path"
	"strings"

	"github.com/MessageDream/goby/model"

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
			if !ctx.IsSigned {
				if ctx.SignError != nil {
					ctx.Error(ctx.SignError)
				} else {
					ctx.Status(403)
				}
			}
		}

		if options.AdminRequire {
			if !ctx.User.IsAdmin {
				ctx.Status(403)
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
				if ctx.SignError != nil {
					ctx.Error(403, ctx.SignError.Error())
					return
				}
				es := url.QueryEscape(ctx.Req.RequestURI)
				ctx.SetCookie("redirect_to", es)
				ctx.Redirect(path.Join(subURL, "/web/auth/signin"))
				return

			} else if registerEmailConfirm && ctx.User.Status == model.USER_STATUS_UN_ACTIVE {
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

		if ctx.User != nil && ctx.User.IsAdmin {
			ctx.Data["IsAdmin"] = true
		}

	}
}
