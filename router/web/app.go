package web

import (
	"github.com/MessageDream/goby/core/appService"
	"github.com/MessageDream/goby/module/context"
	"github.com/MessageDream/goby/module/infrastructure"
)

const (
	APPS infrastructure.TplName = "app/list"
)

func AppsGet(ctx *context.HTMLContext) {
	apps, err := appService.List(ctx.User)
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}
	ctx.Data["PageIsApps"] = true
	ctx.Data["Apps"] = apps
	ctx.HTML(200, APPS)
}
