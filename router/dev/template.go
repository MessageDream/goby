package dev

import (
	"github.com/MessageDream/goby/model"
	"github.com/MessageDream/goby/module/infrastructure"
	"github.com/MessageDream/goby/module/context"
	"github.com/MessageDream/goby/module/setting"
)

func TemplatePreview(ctx *context.Context) {
	ctx.Data["User"] = model.User{UserName: setting.RunUser}
	ctx.Data["AppName"] = setting.AppName
	ctx.Data["AppVer"] = setting.AppVer
	ctx.Data["AppURL"] = setting.AppURL
	ctx.Data["AppLogo"] = setting.AppLogo
	ctx.Data["Code"] = "2014031910370000009fff6782aadb2162b4a997acb69d4400888e0b9274657374"
	ctx.Data["ActiveCodeLives"] = setting.Service.ActiveCodeLives / 60
	ctx.Data["ResetPwdCodeLives"] = setting.Service.ResetPwdCodeLives / 60
	ctx.Data["CurDbValue"] = ""
	ctx.HTML(200, infrastructure.TplName(ctx.Params("*")))
}
