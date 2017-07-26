package web

import (
	"strconv"

	"github.com/MessageDream/goby/core/userService"
	"github.com/MessageDream/goby/model"
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

func UserAddPost(ctx *context.APIContext) {
	name := ctx.Query("name")
	email := ctx.Query("email")
	role := ctx.QueryInt("role")
	pwd := "123456"

	_, err := userService.Create(name, pwd, email, model.USER_STATUS_NORMAL, role == 1)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(200)
}

func UserPatch(ctx *context.APIContext) {
	email := ctx.Params("email")
	status := ctx.Query("status")
	if len(status) > 0 {
		if sta, err := strconv.Atoi(status); err == nil {
			if err := userService.UpdateStatus(email, sta); err != nil {
				ctx.Error(err)
				return
			}
		} else {
			ctx.Error(err)
			return
		}
	}

	role := ctx.Query("role")
	if len(role) > 0 {
		if rol, err := strconv.Atoi(role); err == nil {
			if err := userService.UpdateRole(email, rol); err != nil {
				ctx.Error(err)
				return
			}
		} else {
			ctx.Error(err)
			return
		}
	}
}
