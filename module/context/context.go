package context

import (
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	log "gopkg.in/clog.v1"

	"github.com/go-macaron/cache"
	"gopkg.in/macaron.v1"

	"github.com/MessageDream/goby/model"
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

type Context struct {
	*macaron.Context
	Cache       cache.Cache
	User        *model.User
	SignError   error
	IsSigned    bool
	IsBasicAuth bool
}

func (ctx *Context) Query(name string) string {
	ctx.Req.ParseForm()
	return ctx.Req.Form.Get(name)
}

func (ctx *Context) ServeFile(file string, names ...string) {
	var name string
	if len(names) > 0 {
		name = names[0]
	} else {
		name = path.Base(file)
	}
	ctx.Resp.Header().Set("Content-Description", "File Transfer")
	ctx.Resp.Header().Set("Content-Type", "application/octet-stream")
	ctx.Resp.Header().Set("Content-Disposition", "attachment; filename="+name)
	ctx.Resp.Header().Set("Content-Transfer-Encoding", "binary")
	ctx.Resp.Header().Set("Expires", "0")
	ctx.Resp.Header().Set("Cache-Control", "must-revalidate")
	ctx.Resp.Header().Set("Pragma", "public")
	http.ServeFile(ctx.Resp, ctx.Req.Request, file)
}

func (ctx *Context) ServeContent(name string, r io.ReadSeeker, params ...interface{}) {
	modtime := time.Now()
	for _, p := range params {
		switch v := p.(type) {
		case time.Time:
			modtime = v
		}
	}
	ctx.Resp.Header().Set("Content-Description", "File Transfer")
	ctx.Resp.Header().Set("Content-Type", "application/octet-stream")
	ctx.Resp.Header().Set("Content-Disposition", "attachment; filename="+name)
	ctx.Resp.Header().Set("Content-Transfer-Encoding", "binary")
	ctx.Resp.Header().Set("Expires", "0")
	ctx.Resp.Header().Set("Cache-Control", "must-revalidate")
	ctx.Resp.Header().Set("Pragma", "public")
	http.ServeContent(ctx.Resp, ctx.Req.Request, name, modtime, r)
}

func (ctx *Context) Handle(status int, title string, err error) {
	if err != nil {
		log.Error(4, "%s: %v", title, err)
		if macaron.Env != macaron.PROD {

		}
	}

	switch status {
	case 404:
		ctx.Status(404)
	case 500:
		ctx.Status(status)
	}

}

func (ctx *Context) Auth() {

}

func Contexter() macaron.Handler {
	return func(c *macaron.Context, cache cache.Cache) {

		ctx := &Context{
			Context: c,
			Cache:   cache,
		}

		if ctx.Req.Method == "POST" && strings.Contains(ctx.Req.Header.Get("Content-Type"), "multipart/form-data") {
			if err := ctx.Req.ParseMultipartForm(attachmentFileMaxSize << 20); err != nil && !strings.Contains(err.Error(), "EOF") { // 32MB max size
				ctx.Handle(500, "ParseMultipartForm", err)
				return
			}
		}
		c.Map(ctx)
	}
}
