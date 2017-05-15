package client

import (
	"github.com/MessageDream/goby/core/clientService"
	"github.com/MessageDream/goby/module/context"
	forms "github.com/MessageDream/goby/module/form"
)

func UpdateGet(ctx *context.APIContext) {
	deploymentKey := ctx.Query("deploymentKey")
	appVersion := ctx.Query("appVersion")
	packageHash := ctx.Query("packageHash")
	label := ctx.Query("label")
	// isCompanion := ctx.Query("isCompanion")
	clientUniqueID := ctx.Query("clientUniqueId")

	ret, err := clientService.UpdateCheck(ctx.Cache, deploymentKey, appVersion, label, packageHash, clientUniqueID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"updateInfo": ret,
	})

}

func Download(ctx *context.APIContext, option forms.ReportStatus) {
	deploymentKey := ""
	label := ""

	if option.DeploymentKey != nil {
		deploymentKey = *(option.DeploymentKey)
	}

	if option.Label != nil {
		label = *(option.Label)
	}

	clientService.ReportStatusDownload(deploymentKey, label)
}

func Deploy(ctx *context.APIContext, option forms.ReportStatus) {
	appVersion := ""
	deploymentKey := ""
	label := ""
	status := ""

	if option.AppVersion != nil {
		appVersion = *(option.AppVersion)
	}

	if option.DeploymentKey != nil {
		deploymentKey = *(option.DeploymentKey)
	}

	if option.Label != nil {
		label = *(option.Label)
	}

	if option.Status != nil {
		status = *(option.Status)
	}

	clientService.ReportStatusDeploy(deploymentKey, label, appVersion, status)

}
