package deploymentService

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/go-macaron/cache"
	"github.com/go-xorm/xorm"
	. "gopkg.in/ahmetb/go-linq.v3"
	log "gopkg.in/clog.v1"

	. "github.com/MessageDream/goby/core"
	"github.com/MessageDream/goby/model"
	"github.com/MessageDream/goby/model/dto"
	forms "github.com/MessageDream/goby/module/form"
	"github.com/MessageDream/goby/module/infrastructure"
	"github.com/MessageDream/goby/module/setting"
)

func GetDeployments(uid uint64, appName string) ([]*dto.Deployment, error) {
	col, err := CollaboratorOf(uid, appName)
	if err != nil {
		return nil, err
	}
	details, err := model.FindDeploymentDetailsByAppID(col.AppID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.Deployment, 0, 10)

	From(details).Select(func(detail interface{}) interface{} {
		return detail.(*model.DeploymentDetail).Convert()
	}).ToSlice(&result)

	return result, nil
}

func GetDeploymentMetrics(uid uint64, appName, deploymentName string) (map[string]*dto.PackageMetrics, error) {
	col, err := CollaboratorOf(uid, appName)
	if err != nil {
		return nil, err
	}

	metricsMap := map[string]*dto.PackageMetrics{}
	metrics, err := (&model.Deployment{AppID: col.AppID, Name: deploymentName}).GetPackageMetricsByAppIDAndName()
	if err != nil {
		return nil, err
	}

	From(metrics).Select(func(metr interface{}) interface{} {
		metric := metr.(*model.PackageMetricsDetail)
		return KeyValue{
			Key: metric.Label,
			Value: &dto.PackageMetrics{
				Active:     metric.Active,
				Downloaded: metric.Downloaded,
				Failed:     metric.Failed,
				Installed:  metric.Installed,
			},
		}
	}).ToMap(&metricsMap)

	return metricsMap, nil
}

func GetDeploymentHistories(uid uint64, appName, deploymentName string) ([]*dto.Package, error) {
	col, err := CollaboratorOf(uid, appName)
	if err != nil {
		return nil, err
	}

	deployment := &model.Deployment{
		AppID: col.AppID,
		Name:  deploymentName,
	}

	if exist, err := deployment.Get(); err != nil || !exist {
		if !exist {
			return nil, ErrDeploymentNotExist
		}
		return nil, err
	}

	histories, err := deployment.Histories()

	if err != nil {
		return nil, err
	}

	result := make([]*dto.Package, 0, 10)

	From(histories).Select(func(history interface{}) interface{} {
		return history.(*model.DeploymentHistoryDetail).Convert()
	}).ToSlice(&result)

	return result, nil
}

func DeleteDeploymentHistory(uid uint64, appName, deploymentName string) error {
	col, err := OwnerCan(uid, appName)
	if err != nil {
		return err
	}

	deployment := &model.Deployment{
		AppID: col.AppID,
		Name:  deploymentName,
	}

	if exist, err := deployment.Get(); err != nil || !exist {
		if !exist {
			return ErrDeploymentNotExist
		}
		return err
	}

	_, err = model.Transaction(func(sess *xorm.Session) (interface{}, error) {

		if err := deployment.DeleteAllHistory(sess); err != nil {
			return nil, err
		}

		if err := deployment.DeleteAllVersion(sess); err != nil {
			return nil, err
		}

		if err := deployment.DeleteAllPackage(sess); err != nil {
			return nil, err
		}

		deployment.LabelCursor = 0
		deployment.LastVersionID = 0

		if err := deployment.Update(nil); err != nil {
			return nil, err
		}
		return nil, nil
	})

	return err
}

func ReleaseDeployment(user *model.User,
	cache cache.Cache,
	appName,
	deploymentName,
	fileType,
	fileName string,
	packageFile io.Reader,
	packageInfo *forms.PackageInfo) error {

	appVersion := ""
	isMandatory := false
	isDisabled := false
	description := ""
	var rollout uint8 = 100

	if packageInfo.AppVersion == nil {
		return ErrDeploymentPackageVersionNameFormatNotAllowned
	}

	appVersion = *(packageInfo.AppVersion)

	if packageInfo.IsMandatory != nil {
		isMandatory = *(packageInfo.IsMandatory)
	}

	if packageInfo.IsDisabled != nil {
		isDisabled = *(packageInfo.IsDisabled)
	}

	if packageInfo.Description != nil {
		description = *(packageInfo.Description)
	}

	if packageInfo.Rollout != nil {
		rollout = *(packageInfo.Rollout)
		if rollout > 100 || rollout%10 != 0 {
			return ErrDeploymentPackageRolloutInvalide
		}
	}

	if !strings.Contains(setting.AttachmentAllowedTypes, fileType) {
		return ErrDeploymentPackageFileTypeNotAllowned
	}

	col, err := CollaboratorOf(user.ID, appName)
	if err != nil {
		return err
	}

	deployment := &model.Deployment{
		AppID: col.AppID,
		Name:  deploymentName,
	}

	if exist, err := deployment.Get(); err != nil || !exist {
		if !exist {
			return ErrDeploymentNotExist
		}
		return err
	}

	filePath, err := infrastructure.SaveFileToTemp(fileName, packageFile)
	if err != nil {
		return err
	}

	zipHash, err := infrastructure.FileSHA256(filePath)
	if err != nil {
		return err
	}

	fileInfo, err := os.Lstat(filePath)
	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()

	unzipTo := strings.Replace(filePath, ".zip", "_unzip", -1)

	if err := infrastructure.DeCompress(filePath, unzipTo); err != nil {
		return err
	}

	defer func() {
		os.Remove(filePath)
		os.RemoveAll(unzipTo)
	}()

	typ, err := checkPackageFileType(unzipTo)

	if err != nil {
		return err
	}

	if typ == PACKAGE_TYPE_IOS {
		if !strings.HasSuffix(appName, APP_SUFFIX_IOS) {
			return ErrDeploymentPackageContentsShouldBeIOSFormat
		}
	} else {
		if !strings.HasSuffix(appName, APP_SUFFIX_ANDROID) {
			return ErrDeploymentPackageContentsShouldBeAndroidFormat
		}
	}

	packageHash, _, manifestHash, manifestFilePath, _, err := rearrangingPackage(unzipTo)
	if err != nil {
		return err
	}

	manifesURLPath, err := uploadPackage(manifestHash, manifestFilePath)
	if err != nil {
		return err
	}
	zipURLPath, err := uploadPackage(zipHash, filePath)
	if err != nil {
		return nil
	}

	pkg, err := createPackage(deployment,
		isMandatory,
		isDisabled,
		rollout,
		fileSize,
		appVersion,
		packageHash,
		manifesURLPath,
		zipURLPath,
		"Upload",
		user.Email,
		description,
		"",
		"")

	if err != nil {
		return err
	}

	go func() {
		if _, err := createDiffPackage(pkg); err != nil {
			log.Error(4, "Create diff error for pkg %p:%v", pkg, err)
		}

		os.RemoveAll(path.Join(manifestFilePath, "../"))

		clearCache(cache, deployment.Key, appVersion)

	}()

	return nil
}

func UpdateDeploymentPackage(uid uint64, appName, deploymentName string, packageInfo *forms.PackageInfo) error {

	col, err := CollaboratorOf(uid, appName)
	if err != nil {
		return err
	}

	deployment := &model.Deployment{
		AppID: col.AppID,
		Name:  deploymentName,
	}

	detail, err := deployment.Detail()
	if err != nil {
		return err
	}

	_, err = model.Transaction(func(sess *xorm.Session) (interface{}, error) {
		pkg := &detail.Package
		deployVersion := &detail.DeploymentVersion

		if pkg != nil {
			needUpdate := false
			if packageInfo.IsMandatory != nil {
				needUpdate = true
				pkg.IsMandatory = *(packageInfo.IsMandatory)
			}

			if packageInfo.IsDisabled != nil {
				needUpdate = true
				pkg.IsDisabled = *(packageInfo.IsDisabled)
			}

			if packageInfo.Description != nil {
				needUpdate = true
				pkg.Description = *(packageInfo.Description)
			}

			if packageInfo.Label != nil {
				needUpdate = true
				pkg.Label = *(packageInfo.Label)
			}

			if packageInfo.Rollout != nil {
				needUpdate = true
				rollout := *(packageInfo.Rollout)
				if rollout > 100 || rollout%10 != 0 {
					return nil, ErrDeploymentPackageRolloutInvalide
				}
				pkg.Rollout = rollout
			}
			if needUpdate {
				if err := pkg.Update(sess, "is_mandatory", "is_disabled", "description", "label"); err != nil {
					return nil, err
				}
			}
		}

		if deployVersion != nil {
			if packageInfo.AppVersion != nil {
				deployVersion.AppVersion = *(packageInfo.AppVersion)
				if err := deployVersion.Update(sess); err != nil {
					return nil, err
				}
			}
		}

		return nil, nil
	})

	return err
}

func AddDeployment(user *model.User, appName, deploymentName string) (*dto.Deployment, error) {
	col, err := OwnerCan(user.ID, appName)
	if err != nil {
		return nil, err
	}

	newDeployment := &model.Deployment{
		AppID: col.AppID,
		Name:  deploymentName,
	}

	if exist, err := newDeployment.Get(); err != nil || exist {
		if exist {
			return nil, ErrDeploymentOfThisNameAlreadyExist
		}
		return nil, err
	}

	newDeployment.Key = infrastructure.GetRandomString(28) + user.Rands

	if err := newDeployment.Create(nil); err != nil {
		return nil, err
	}

	return &dto.Deployment{
		Name: newDeployment.Name,
		Key:  newDeployment.Key,
	}, nil

}

func RenameDeployment(uid uint64, appName, deploymentName, newDeploymentName string) error {
	col, err := OwnerCan(uid, appName)
	if err != nil {
		return err
	}

	deployment := &model.Deployment{
		AppID: col.AppID,
		Name:  newDeploymentName,
	}

	if exist, err := deployment.Get(); err != nil || exist {
		if exist {
			return ErrDeploymentOfThisNameAlreadyExist
		}
		return err
	}

	deployment.Name = deploymentName

	if exist, err := deployment.Get(); err != nil || !exist {
		if !exist {
			return ErrDeploymentNotExist
		}
		return err
	}

	deployment.Name = newDeploymentName

	if err := deployment.Update(nil, "name"); err != nil {
		return err
	}
	return nil
}

func DeleteDeployment(uid uint64, appName, deploymentName string) error {
	col, err := OwnerCan(uid, appName)
	if err != nil {
		return err
	}

	deployment := &model.Deployment{
		AppID: col.AppID,
		Name:  deploymentName,
	}

	if exist, err := deployment.Get(); err != nil || !exist {
		if !exist {
			return ErrDeploymentNotExist
		}
		return err
	}

	if err := deployment.Delete(nil); err != nil {
		return err
	}

	return nil
}

func RollbackDeployment(user *model.User, cache cache.Cache, appName, deploymentName, targetLabel string) error {
	col, err := CollaboratorOf(user.ID, appName)
	if err != nil {
		return err
	}

	deployment := &model.Deployment{
		AppID: col.AppID,
		Name:  deploymentName,
	}

	detail, err := deployment.Detail()
	if err != nil || detail == nil {
		if detail == nil {
			return ErrDeploymentNotExist
		}
		return err
	}
	deployment = &detail.Deployment
	version := detail.DeploymentVersion
	pkg := detail.Package
	if version.ID == 0 || pkg.ID == 0 {
		return ErrDeploymentRollbackNotReleasePackageBefore
	}

	rollbackPkgs := make([]*model.Package, 0, 2)
	if len(targetLabel) != 0 {
		rollbackPkg := &model.Package{
			DeployVersionID: version.ID,
			Label:           targetLabel,
		}

		if exist, err := rollbackPkg.Get(); err != nil || !exist {
			if !exist {
				return ErrDeploymentPackageNotExist
			}
			return err
		}
		rollbackPkgs = append(rollbackPkgs, rollbackPkg)
	}

	if len(rollbackPkgs) == 0 {
		rollbackPkg := &model.Package{
			DeployVersionID: version.ID,
		}

		pkgs, err := rollbackPkg.FindPackagesByVersionIDAndReleaseMethods(2, false, "Upload", "Promote")
		if err != nil {
			return err
		}

		rollbackPkgs = pkgs

	}

	result := From(rollbackPkgs).FirstWith(func(pk interface{}) bool {
		return pk.(*model.Package).Hash != pkg.Hash
	})

	if result == nil {
		return ErrDeploymentRollbackHaveNoPackageToRollback
	}

	rollbackPkg := result.(*model.Package)

	_, err = createPackage(deployment,
		rollbackPkg.IsMandatory,
		rollbackPkg.IsDisabled,
		rollbackPkg.Rollout,
		rollbackPkg.Size,
		version.AppVersion,
		rollbackPkg.Hash,
		rollbackPkg.ManifestBlobURL,
		rollbackPkg.BlobURL,
		"Rollback",
		user.Email,
		rollbackPkg.Description,
		rollbackPkg.Label,
		deployment.Name)

	if err != nil {
		return err
	}
	clearCache(cache, deployment.Key, version.AppVersion)
	return nil
}

func PromoteDeployment(user *model.User, cache cache.Cache, appName, sourceDeploymentName, destDeploymentName string) error {
	col, err := CollaboratorOf(user.ID, appName)
	if err != nil {
		return err
	}

	sourceDeployment := &model.Deployment{
		AppID: col.AppID,
		Name:  sourceDeploymentName,
	}

	sourceDetail, err := sourceDeployment.Detail()

	if err != nil || sourceDetail == nil {
		if sourceDetail == nil {
			return WrapIntentError(fmt.Errorf("%s does not exist", sourceDeploymentName), INTENT_ERROR_CODE_DATA_NOT_EXIST)
		}

		return err
	}

	sourceDeployment = &sourceDetail.Deployment
	sourcePkg := sourceDetail.Package
	sourceVersion := sourceDetail.DeploymentVersion

	if sourcePkg.ID == 0 || sourceVersion.ID == 0 {
		return WrapIntentError(fmt.Errorf("%s has no version or package data to promote", sourceDeploymentName), INTENT_ERROR_CODE_DATA_NOT_EXIST)
	}

	destDeployment := &model.Deployment{
		AppID: col.AppID,
		Name:  destDeploymentName,
	}

	if exist, err := destDeployment.Get(); err != nil || !exist {
		if !exist {
			return WrapIntentError(fmt.Errorf("%s does not exist", destDeploymentName), INTENT_ERROR_CODE_DATA_NOT_EXIST)
		}

		return err
	}

	_, err = createPackage(destDeployment,
		sourcePkg.IsMandatory,
		sourcePkg.IsDisabled,
		sourcePkg.Rollout,
		sourcePkg.Size,
		sourceVersion.AppVersion,
		sourcePkg.Hash,
		sourcePkg.ManifestBlobURL,
		sourcePkg.BlobURL,
		"Promote",
		user.Email,
		sourcePkg.Description,
		sourcePkg.Label,
		sourceDeployment.Name)

	if err != nil {
		clearCache(cache, destDeployment.Key, sourceVersion.AppVersion)
	}
	return err

}
