package account

import (
	"github.com/MessageDream/goby/module/context"
)

func Info(ctx *context.APIContext) {
	ctx.JSON(200, map[string]interface{}{
		"account": ctx.User.Convert(),
	})
}
