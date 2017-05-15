package web

import (
	"github.com/MessageDream/goby/core/appService"
	"github.com/MessageDream/goby/module/infrastructure"
	"github.com/MessageDream/goby/module/context"
)

const (
	APPS infrastructure.TplName = "app/list"
)

func AppsGet(ctx *context.Context) {
	apps, err := appService.List(ctx.User)
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}
	ctx.Data["PageIsApps"] = true
	ctx.Data["Apps"] = apps
	ctx.HTML(200, APPS)
}
