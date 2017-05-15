package context

import (
	log "gopkg.in/clog.v1"
	"gopkg.in/macaron.v1"

	"github.com/MessageDream/goby/model"
	"github.com/MessageDream/goby/module/infrastructure"
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

func APIContexter() macaron.Handler {
	return func(c *Context) {
		ctx := &APIContext{
			Context: c,
		}
		c.Map(ctx)
	}
}
