package collaborator

import (
	"github.com/MessageDream/goby/core/collaboratorService"
	"github.com/MessageDream/goby/module/context"
)

func Get(ctx *context.APIContext) {
	appName := ctx.Params("appName")
	cols, err := collaboratorService.List(ctx.User.ID, appName)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"collaborators": cols,
	})
}

func Post(ctx *context.APIContext) {
	appName := ctx.Params("appName")
	email := ctx.Params("email")
	if err := collaboratorService.Add(ctx.User, appName, email); err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(200)
}

func Delete(ctx *context.APIContext) {
	appName := ctx.Params("appName")
	email := ctx.Params("email")
	if err := collaboratorService.Delete(ctx.User, appName, email); err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(200)
}
