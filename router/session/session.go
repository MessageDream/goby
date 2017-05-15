package session

import (
	"github.com/MessageDream/goby/core/tokenService"
	"github.com/MessageDream/goby/module/context"
)

func Delete(ctx *context.APIContext) {
	createdBy := ctx.Params("machineName")
	uid := ctx.User.ID

	if err := tokenService.Delete(uid, createdBy); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(200)
}
