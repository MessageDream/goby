package collaborator

import (
	"github.com/MessageDream/goby/core/collaboratorService"
	"github.com/MessageDream/goby/model"
	"github.com/MessageDream/goby/model/dto"
	"github.com/MessageDream/goby/module/context"
	. "gopkg.in/ahmetb/go-linq.v3"
)

func Get(ctx *context.APIContext) {
	appName := ctx.Params("appName")
	uid := ctx.User.ID
	cols, err := collaboratorService.List(uid, appName)
	if err != nil {
		ctx.Error(err)
		return
	}

	canOwner := false
	colsMap := map[string]*dto.Collaborator{}
	From(cols).Select(func(col interface{}) interface{} {
		co := col.(*model.Collaborator)
		isOwner := uid == co.UID
		if isOwner && co.Role == 1 {
			canOwner = true
		}
		permission := "Member"
		if co.Role == 1 {
			permission = "Owner"
		}
		return KeyValue{
			Key: co.Email,
			Value: &dto.Collaborator{
				IsCurrentAccount: isOwner,
				Permission:       permission,
			},
		}

	}).ToMap(&colsMap)

	ctx.JSON(200, map[string]interface{}{
		"collaborators": colsMap,
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
