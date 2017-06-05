package context

import (
	"errors"
	"strings"

	"github.com/MessageDream/goby/core/tokenService"

	log "gopkg.in/clog.v1"
	"gopkg.in/macaron.v1"

	"github.com/MessageDream/goby/model"
	"github.com/MessageDream/goby/module/infrastructure"
	"github.com/go-macaron/cache"
)

type APIContext struct {
	*Context
}

func (ctx *APIContext) Error(obj interface{}) {
	var message string
	status := 406
	code := status
	if err, ok := obj.(*infrastructure.IntentError); ok {
		message = err.Error()
		code = err.Code
		if err.Status > 0 {
			status = err.Status
		}
	} else if err, ok := obj.(error); ok {
		status = 500
		code = 500
		message = err.Error()
	} else {
		message = obj.(string)
	}

	if status == 500 {
		log.Error(4, "%s", message)
		message = "server error"
	}

	ctx.JSON(status, infrastructure.ApiJsonError{
		Message: message,
		Status:  code,
	})
}

func (ctx *APIContext) JSON(status int, obj interface{}) {

	if converter, ok := obj.(model.JsonObjectConverter); ok {
		ctx.Context.JSON(status, converter.Convert())
	} else {
		ctx.Context.JSON(status, obj)
	}
}

func (ctx *APIContext) Auth() {
	if !model.HasEngine {
		return
	}
	token := ctx.GetToken()
	if len(token) > 0 {
		u, err := checkAccessToken(token, ctx.Cache)
		if err != nil {
			return
		}
		ctx.User = u
		ctx.IsSigned = true
		ctx.IsBasicAuth = true
	}
}

func (ctx *APIContext) GetToken() string {
	tokenSHA := ""
	if IsAPIPath(ctx.Req.URL.Path) {
		auHead := ctx.Req.Header.Get("Authorization")
		if len(auHead) > 0 {
			auths := strings.Fields(auHead)
			if len(auths) == 2 {
				if auths[0] == "Bearer" {
					tokenSHA = auths[1]
				} else if auths[0] == "Basic" {
					_, tokenSHA, _ = infrastructure.BasicAuthDecode(auths[1])
				}
			}
		}

		if len(tokenSHA) == 0 {
			tokenSHA = ctx.Params("access_token")
		}
	}

	return tokenSHA
}

func checkAccessToken(token string, cache cache.Cache) (*model.User, error) {
	if len(token) < 37 {
		return nil, errors.New("token is invalid")
	}

	if !cache.IsExist(token) {
		user, err := tokenService.CheckTokenSession(token)
		if err == nil && user != nil {
			cache.Put(token, user, tokenCacheLiveTimes)
		}

		return user, err
	}

	return (cache.Get(token)).(*model.User), nil

}

func APIContexter() macaron.Handler {
	return func(c *Context) {
		ctx := &APIContext{
			Context: c,
		}
		c.Map(ctx)
	}
}
