package clientService

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/go-macaron/cache"

	. "github.com/MessageDream/goby/core"
	"github.com/MessageDream/goby/model"
	"github.com/MessageDream/goby/model/dto"
	"github.com/MessageDream/goby/module/infrastructure"
	"github.com/MessageDream/goby/module/setting"
	"github.com/MessageDream/goby/module/storage"
)

func UpdateCheck(cache cache.Cache, deploymentKey, appVersion, label, packageHash, deviceId string) (*dto.UpdateCheck, error) {
	if len(deploymentKey) == 0 || len(appVersion) == 0 {
		return nil, ErrClientDeploymentKeyOrAppVersionEmpty
	}

	checkResult := &dto.UpdateCheck{
		PackageInfo: &dto.PackageInfo{
			Description: "",
			IsMandatory: false,
			AppVersion:  appVersion,
			PackageHash: "",
			Label:       "",
		},
		DownloadURL:            "",
		IsAvailable:            false,
		PackageSize:            0,
		ShouldRunBinaryVersion: false,
		UpdateAppVersion:       false,
	}

	cacheKey := infrastructure.EncodeSHA256(deploymentKey + appVersion)

	deployment := &model.Deployment{
		Key: deploymentKey,
	}
	var pkg *model.Package
	if !cache.IsExist(cacheKey) {
		pk, err := deployment.PackageOfVersion(appVersion)
		if err != nil || pk == nil {
			if pkg == nil {
				return checkResult, nil
			}
			return nil, err
		}

		cache.Put(cacheKey, pk, setting.CacheClientCheckUpdateTimeOut)

		pkg = pk
	} else {
		pkg = cache.Get(cacheKey).(*model.Package)
	}

	if pkg.Hash == packageHash || !grayCalculate(pkg.Rollout, deviceId) {
		return checkResult, nil
	}

	checkResult.Description = pkg.Description
	checkResult.IsMandatory = pkg.IsMandatory
	checkResult.PackageHash = pkg.Hash
	checkResult.Label = pkg.Label

	checkResult.IsAvailable = true

	checkResult.DownloadURL = storage.GenerateDownloadURL(pkg.BlobURL)
	checkResult.PackageSize = pkg.Size

	pkgDiff := &model.PackageDiff{
		PackageID:              pkg.ID,
		DiffAgainstPackageHash: packageHash,
	}

	if exist, err := pkgDiff.Get(); err != nil || !exist {
		return checkResult, err
	}

	checkResult.DownloadURL = storage.GenerateDownloadURL(pkgDiff.DiffBlobURL)
	checkResult.PackageSize = pkgDiff.DiffSize

	return checkResult, nil

}

func ReportStatusDownload(deploymentKey, label string) error {
	if len(deploymentKey) == 0 || len(label) == 0 {
		return ErrClientDeploymentKeyOrLabelEmpty
	}
	deployment := &model.Deployment{
		Key: deploymentKey,
	}

	metrics, err := deployment.PackageMetricsOfLabel(label)

	if err != nil || metrics == nil {
		if metrics == nil {
			return ErrDeploymentPackageNotExist
		}
		return err
	}

	metrics.Downloaded += 1

	return metrics.Update(nil)
}

func ReportStatusDeploy(deploymentKey, label, appVersion, status string) error {
	if len(deploymentKey) == 0 || len(appVersion) == 0 || len(label) == 0 {
		return ErrClientDeploymentKeyOrAppVersionEmpty
	}
	deployment := &model.Deployment{
		Key: deploymentKey,
	}

	metrics, err := deployment.PackageMetricsOfLabel(label)

	if err != nil || metrics == nil {
		if metrics == nil {
			return ErrDeploymentPackageNotExist
		}
		return err
	}

	metrics.Installed += 1
	switch status {
	case "DeploymentSucceeded":
		metrics.Active += 1
		break
	case "DeploymentFailed":
		metrics.Failed += 1
		break
	default:
		return nil
	}

	return metrics.Update(nil)
}

func grayCalculate(rollout uint8, deviceId string) bool {
	if rollout >= 100 {
		return true
	}
	hexStr := infrastructure.EncodeMd5(deviceId)
	bi := big.NewInt(0)
	if _, ok := bi.SetString(hexStr, 16); !ok {
		return false
	}
	numberStr := fmt.Sprintf("%v", bi)

	numArr := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	persent := rollout / 10
	pool := numArr[:persent]

	last := numberStr[len(numberStr)-1:]

	num, _ := strconv.Atoi(last)

	for _, v := range pool {
		if num == v {
			return true
		}
	}

	return false
}
