package web

import (
	"github.com/MessageDream/goby/core"
	"github.com/MessageDream/goby/core/appService"
	"github.com/MessageDream/goby/core/collaboratorService"
	"github.com/MessageDream/goby/core/deploymentService"
	"github.com/MessageDream/goby/module/context"
	"github.com/MessageDream/goby/module/form"
	"github.com/MessageDream/goby/module/infrastructure"
)

const (
	APPS                     infrastructure.TplName = "app/list"
	APP_DETAIL_COLLABORATORS infrastructure.TplName = "app/detail/collaborators"
	APP_DETAIL_DEPLOYMENTS   infrastructure.TplName = "app/detail/deployments"
)

//list
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

//detail of collaborators
func AppCollaboratorsGet(ctx *context.HTMLContext) {
	appName := ctx.Params("appName")
	cols, err := collaboratorService.List(ctx.User.ID, appName)
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}
	owner, err := core.OwnerOf(appName)
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}

	ctx.Data["AppName"] = appName
	ctx.Data["IsCollaboratorPage"] = true
	ctx.Data["Owner"] = owner.Email
	ctx.Data["Collaborators"] = cols
	ctx.HTML(200, APP_DETAIL_COLLABORATORS)
}

//detail of deployments
func AppDeploymentsGet(ctx *context.HTMLContext) {
	appName := ctx.Params("appName")
	deployments, err := deploymentService.GetDeployments(ctx.User.ID, appName)
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}
	owner, err := core.OwnerOf(appName)
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}

	ctx.Data["AppName"] = appName
	ctx.Data["IsDeploymentPage"] = true
	ctx.Data["Owner"] = owner.Email
	ctx.Data["Deployments"] = deployments
	ctx.HTML(200, APP_DETAIL_DEPLOYMENTS)
}

//create app
func AppAdd(ctx *context.HTMLContext, appForm form.AppAddForm) {
	name := appForm.Name
	platform := appForm.Platform
	appName := name + "-" + platform

	_, err := appService.Create(ctx.User, appName)
	if err != nil {
		ctx.Data["NotifyType"] = "error"
		ctx.Data["NotifyMsg"] = err.Error()
		AppsGet(ctx)
		return
	}
	ctx.Redirect("list")
}
