package web

import (
	"github.com/MessageDream/goby/core/tokenService"
	"github.com/MessageDream/goby/module/context"
	"github.com/MessageDream/goby/module/infrastructure"
)

const (
	ACCESSKEYS infrastructure.TplName = "accessKey/list"
)

func AccessKeysGet(ctx *context.HTMLContext) {
	keys, err := tokenService.List(ctx.User)
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}
	ctx.Data["PageIsAccessKeys"] = true
	ctx.Data["Keys"] = keys
	ctx.HTML(200, ACCESSKEYS)
}
