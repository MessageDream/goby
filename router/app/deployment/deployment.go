package deployment

import (
	"encoding/json"

	"github.com/MessageDream/goby/core/deploymentService"
	"github.com/MessageDream/goby/module/context"
	forms "github.com/MessageDream/goby/module/form"
)

func Post(ctx *context.APIContext, option forms.AppOrDeploymentOption) {
	appName := ctx.Params("appName")
	deploymentName := option.Name

	deployment, err := deploymentService.AddDeployment(ctx.User, appName, deploymentName)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, map[string]interface{}{
		"deployment": deployment,
	})
}

func Patch(ctx *context.APIContext, option forms.AppOrDeploymentOption) {
	appName := ctx.Params("appName")
	deploymentName := ctx.Params("deploymentName")
	newDeploymentName := option.Name
	if err := deploymentService.RenameDeployment(ctx.User.ID, appName, deploymentName, newDeploymentName); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(200)
}

func Delete(ctx *context.APIContext) {
	appName := ctx.Params("appName")
	deploymentName := ctx.Params("deploymentName")
	if err := deploymentService.DeleteDeployment(ctx.User.ID, appName, deploymentName); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(200)
}

func Get(ctx *context.APIContext) {
	name := ctx.Params("appName")
	uid := ctx.User.ID
	depls, err := deploymentService.GetDeployments(uid, name)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"deployments": depls,
	})

}

func MetricsGet(ctx *context.APIContext) {
	appName := ctx.Params("appName")
	deploymentName := ctx.Params("deploymentName")
	uid := ctx.User.ID
	metrics, err := deploymentService.GetDeploymentMetrics(uid, appName, deploymentName)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, map[string]interface{}{
		"metrics": metrics,
	})
}

func HistoryGet(ctx *context.APIContext) {
	appName := ctx.Params("appName")
	deploymentName := ctx.Params("deploymentName")
	pkgs, err := deploymentService.GetDeploymentHistories(ctx.User.ID, appName, deploymentName)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, map[string]interface{}{
		"history": pkgs,
	})
}

func HistoryDelete(ctx *context.APIContext) {
	appName := ctx.Params("appName")
	deploymentName := ctx.Params("deploymentName")
	err := deploymentService.DeleteDeploymentHistory(ctx.User.ID, appName, deploymentName)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(200)
}

func ReleasePost(ctx *context.APIContext) {
	appName := ctx.Params("appName")
	deploymentName := ctx.Params("deploymentName")
	pakageInfoStr := ctx.Query("packageInfo")

	packageInfo := &forms.PackageInfo{}
	if err := json.Unmarshal([]byte(pakageInfoStr), packageInfo); err != nil {
		ctx.Error(err)
		return
	}

	file, header, err := ctx.GetFile("package")
	if err != nil {
		ctx.Error(err)
		return
	}
	defer file.Close()

	fileType := header.Header.Get("Content-Type")
	fileName := header.Filename

	if err := deploymentService.ReleaseDeployment(ctx.User, ctx.Cache, appName, deploymentName, fileType, fileName, file, packageInfo); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(200)
}

func ReleasePatch(ctx *context.APIContext, option forms.UpdatePackageInfo) {
	appName := ctx.Params("appName")
	deploymentName := ctx.Params("deploymentName")
	if err := deploymentService.UpdateDeploymentPackage(ctx.User.ID, appName, deploymentName, option.PackageInfo); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(200)
}

func RobackPost(ctx *context.APIContext) {
	appName := ctx.Params("appName")
	deploymentName := ctx.Params("deploymentName")
	if err := deploymentService.RollbackDeployment(ctx.User, ctx.Cache, appName, deploymentName, ""); err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(200)
}

func RobackForLabelPost(ctx *context.APIContext) {
	appName := ctx.Params("appName")
	deploymentName := ctx.Params("deploymentName")
	label := ctx.Params("label")
	if err := deploymentService.RollbackDeployment(ctx.User, ctx.Cache, appName, deploymentName, label); err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(200)
}

func DestPost(ctx *context.APIContext) {
	appName := ctx.Params("appName")
	sourceDeploymentName := ctx.Params("sourceDeploymentName")
	destDeploymentName := ctx.Params("destDeploymentName")

	if err := deploymentService.PromoteDeployment(ctx.User, ctx.Cache, appName, sourceDeploymentName, destDeploymentName); err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(200)
}
