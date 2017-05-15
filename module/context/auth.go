package context

import (
	"errors"
	"net/url"
	"path"

	"github.com/go-macaron/cache"
	"github.com/go-macaron/csrf"
	"gopkg.in/macaron.v1"

	"github.com/MessageDream/goby/core/tokenService"
	"github.com/MessageDream/goby/core/userService"
	"github.com/MessageDream/goby/model"
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

var (
	enableReverseProxyAuth                     bool
	reverseProxyAuthUser, subURL               string
	appInstallLock                             bool
	registerEmailConfirm                       bool
	disableRegistration                        bool
	attachmentFileMaxSize, tokenCacheLiveTimes int64
)

func init() {
	enableReverseProxyAuth = false
	reverseProxyAuthUser = ""
	subURL = ""
	appInstallLock = false
	registerEmailConfirm = false
	disableRegistration = false
	attachmentFileMaxSize = 4
	tokenCacheLiveTimes = 1000 * 60 * 5
}

func InitContextConfig(reverseProxyAuthAble, installLock, registerEmailConfirmAble, registrationDisabled bool, tokenCacheTimeOut, attachmentMaxSize int64, subUrl, reverseProxyAuthUserName string) {
	enableReverseProxyAuth = reverseProxyAuthAble
	appInstallLock = installLock
	registerEmailConfirm = registerEmailConfirmAble
	reverseProxyAuthUser = reverseProxyAuthUserName
	disableRegistration = registrationDisabled
	attachmentFileMaxSize = attachmentMaxSize
	tokenCacheLiveTimes = tokenCacheTimeOut
	subURL = subUrl
}

func checkAccessToken(token string, cache cache.Cache) (*model.User, error) {
	if len(token) < 37 {
		return nil, errors.New("token is invalid")
	}

	if !cache.IsExist(token) {
		user, err := tokenService.CheckTokenSession(token)
		cache.Put(token, user, tokenCacheLiveTimes)
		return user, err
	}

	return (cache.Get(token)).(*model.User), nil

}

func signedInUser(ctx *Context) (*model.User, bool) {
	if !model.HasEngine {
		return nil, false
	}
	token := ctx.GetToken()
	if len(token) > 0 {
		u, err := checkAccessToken(token, ctx.Cache)
		if err != nil {
			return nil, false
		}
		return u, true
	}

	user := ctx.Session.Get("u")
	if user == nil {
		if enableReverseProxyAuth {
			webAuthUser := ctx.Req.Header.Get(reverseProxyAuthUser)
			if len(webAuthUser) > 0 {
				u, err := userService.GetByUserName(webAuthUser)
				if err != nil {
					return nil, false
				}
				return u, false
			}
		}
		return nil, false
	}

	if usr, ok := user.(*model.User); ok {
		ctx.Session.Set("u", usr)
		return usr, false
	}

	return nil, false
}

func Toggle(options *ToggleOptions) macaron.Handler {
	return func(ctx *Context) {
		if !appInstallLock {
			ctx.Redirect(path.Join(subURL, "install"))
			return
		}

		if options.SignOutRequire && ctx.IsSigned && ctx.Req.RequestURI != "/" {
			ctx.Redirect(path.Join(subURL, "/"))
			return
		}

		if !options.SignOutRequire && !options.DisableCSRF && ctx.Req.Method == "POST" && !IsAPIPath(ctx.Req.URL.Path) {
			csrf.Validate(ctx.Context, ctx.csrf)
			if ctx.Written() {
				return
			}
		}

		if options.SignInRequire {
			if !ctx.IsSigned {
				if IsAPIPath(ctx.Req.URL.Path) {
					ctx.Error(401, "authorize failed")
					return
				}
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
