package app

import (
	"github.com/MessageDream/goby/core/appService"
	"github.com/MessageDream/goby/module/context"
	forms "github.com/MessageDream/goby/module/form"
)

func Get(ctx *context.APIContext) {
	apps, err := appService.List(ctx.User)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"apps": apps,
	})
}

func Post(ctx *context.APIContext, option forms.AppOrDeploymentOption) {
	name := option.Name

	app, err := appService.Create(ctx.User, name)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, app)
}

func Delete(ctx *context.APIContext) {
	uid := ctx.User.ID
	name := ctx.Params("appName")

	err := appService.Delete(uid, name)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(200)
}

func Patch(ctx *context.APIContext, option forms.AppOrDeploymentOption) {
	oriName := ctx.Params("appName")
	newName := option.Name
	uid := ctx.User.ID

	if err := appService.ChangeName(uid, oriName, newName); err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(200)
}

func TransferPost(ctx *context.APIContext) {
	name := ctx.Params("appName")
	email := ctx.Params("email")

	if err := appService.Transfer(ctx.User.ID, name, ctx.User.Email, email); err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(200)
}
