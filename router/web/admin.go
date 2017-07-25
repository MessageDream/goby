package web

import (
	"github.com/MessageDream/goby/core/userService"
	"github.com/MessageDream/goby/module/context"
	"github.com/MessageDream/goby/module/infrastructure"
)

const (
	ADMIN_USERS infrastructure.TplName = "admin/users"
)

func UsersGet(ctx *context.HTMLContext) {
	ctx.Data["PageIsAdminUsers"] = true
	ctx.HTML(200, ADMIN_USERS)
}

func UsersQuery(ctx *context.APIContext) {
	pageIndex := ctx.ParamsInt("pageIndex")
	pageCount := ctx.ParamsInt("pageCount")
	email := ctx.Query("email")

	result, err := userService.QueryUsers(ctx.User.ID, pageIndex, pageCount, email)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"users": result,
	})
}
