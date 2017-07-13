package model

import (
	"time"

	"github.com/MessageDream/goby/model/dto"
)

type DeploymentDetail struct {
	Deployment        `xorm:"extends"`
	DeploymentVersion `xorm:"extends"`
	Package           `xorm:"extends"`
	PackageMetrics    `xorm:"extends"`
}

func (self *DeploymentDetail) Convert() interface{} {
	var pkg *dto.Package
	if self.Package.ID > 0 {
		pkg = &dto.Package{
			PackageInfo: &dto.PackageInfo{
				AppVersion:  self.AppVersion,
				Description: self.Package.Description,
				IsDisabled:  self.IsDisabled,
				IsMandatory: self.IsMandatory,
				Label:       self.Label,
				PackageHash: self.Hash,
				Rollout:     self.Rollout,
			},
			BlobURL:            self.BlobURL,
			DiffPackageMap:     nil,
			OriginalLabel:      self.OriginalLabel,
			OriginalDeployment: self.OriginalDeployment,
			ReleasedBy:         self.ReleasedBy,
			ReleaseMethod:      self.ReleaseMethod,
			Size:               self.Size,
			UploadTime:         self.Package.CreatedAt.Unix() * 1000,
		}
	}

	var metrics *dto.PackageMetrics
	if self.PackageMetrics.ID > 0 {
		metrics = &dto.PackageMetrics{
			Active:     self.PackageMetrics.Active,
			Downloaded: self.PackageMetrics.Downloaded,
			Failed:     self.PackageMetrics.Failed,
			Installed:  self.PackageMetrics.Installed,
		}
	}
	return &dto.Deployment{
		Key:            self.Key,
		Name:           self.Name,
		Package:        pkg,
		PackageMetrics: metrics,
	}
}

type DeploymentHistoryDetail struct {
	DeploymentVersion `xorm:"extends"`
	Package           `xorm:"extends"`
	PackageMetrics    `xorm:"extends"`
}

func (self *DeploymentHistoryDetail) Convert() interface{} {
	var pkg *dto.Package
	if self.Package.ID > 0 {
		pkg = &dto.Package{
			PackageInfo: &dto.PackageInfo{
				AppVersion:  self.AppVersion,
				Description: self.Package.Description,
				IsDisabled:  self.IsDisabled,
				IsMandatory: self.IsMandatory,
				Label:       self.Label,
				PackageHash: self.Hash,
				Rollout:     self.Rollout,
			},
			BlobURL:            self.BlobURL,
			DiffPackageMap:     nil,
			OriginalLabel:      self.OriginalLabel,
			OriginalDeployment: self.OriginalDeployment,
			ReleasedBy:         self.ReleasedBy,
			ReleaseMethod:      self.ReleaseMethod,
			Size:               self.Size,
			UploadTime:         self.Package.CreatedAt.Unix() * 1000,
		}
	}
	return pkg
}

type Deployment struct {
	ID            uint64 `xorm:"pk autoincr"`
	Name          string `xorm:"notnull default 'Staging'"`
	AppID         uint64 `xorm:"index notnull default(0)"`
	Key           string `xorm:"unique index notnull"`
	LastVersionID uint64 `xorm:"notnull default(0)"`
	LabelCursor   int    `xorm:"notnull default(0)"`
	Description   string
	CreatedAt     time.Time `xorm:"created"`
	UpdatedAt     time.Time `xorm:"updated"`
	DeletedAt     time.Time `xorm:"deleted"`
}

func (self *Deployment) Update(engine Engine, specFields ...string) error {
	return ModelUpdate(engine, self.ID, self, specFields...)
}

func (self *Deployment) Exist() (bool, error) {
	exi, err := x.Where("id!=?", 0).Get(&Deployment{AppID: self.AppID, Name: self.Name})
	return exi, err
}

func (self *Deployment) Create(engine Engine) error {
	return ModelCreate(engine, self)
}

func (self *Deployment) Get() (bool, error) {
	return ModelGet(self.ID, self)
}

func (self *Deployment) Detail() (*DeploymentDetail, error) {
	deploymentDetail := &DeploymentDetail{}
	exist, err := x.Sql(`SELECT deployment.*, deployment_version.*, package.*, package_metrics.* 
	FROM deployment 
	LEFT JOIN deployment_version 
	ON deployment.last_version_id = deployment_version.id 
	LEFT JOIN package 
	ON deployment_version.package_id = package.id 
	LEFT JOIN package_metrics 
	ON package_metrics.package_id = package.id 
	WHERE deployment.deleted_at IS NULL 
	AND deployment.app_id = ? AND deployment.name = ? 
	LIMIT 1`, self.AppID, self.Name).Get(deploymentDetail)
	if exist {
		return deploymentDetail, err
	}
	return nil, err
}

func (self *Deployment) PackageOfVersion(appVersion string) (*Package, error) {
	pkg := &Package{}

	exist, err := x.Sql(`SELECT package.* 
	FROM package 
	WHERE id 
	IN (SELECT deployment_version.package_id 
	FROM deployment 
	INNER JOIN deployment_version 
	ON deployment.id = deployment_version.deploy_id 
	WHERE deployment.deleted_at IS NULL 
	AND deployment.key = ? 
	AND deployment_version.app_version = ?) 
	LIMIT 1`, self.Key, appVersion).Get(pkg)

	if exist {
		return pkg, err
	}
	return nil, err
}

func (self *Deployment) Histories() ([]*DeploymentHistoryDetail, error) {
	histories := make([]*DeploymentHistoryDetail, 0, 100)
	if err := x.Sql(`SELECT deployment_version.*, package.*, package_metrics.* 
	FROM deployment_history 
	LEFT JOIN package 
	ON deployment_history.package_id = package.id 
	LEFT JOIN deployment_version 
	ON deployment_version.package_id = package.id 
	LEFT JOIN package_metrics 
	ON package_metrics.package_id = package.id 
	WHERE deployment_history.deleted_at IS NULL 
	AND deployment_history.deploy_id = ? 
	LIMIT 1000`, self.ID).Find(&histories); err != nil {
		return nil, err
	}
	return histories, nil
}

func (self *Deployment) PackageMetricsOfLabel(label string) (*PackageMetrics, error) {
	metrics := &PackageMetrics{}
	exist, err := x.Sql(`SELECT package_metrics.* 
	FROM package_metrics 
	WHERE package_id IN 
	(SELECT package.id 
	FROM deployment 
	INNER JOIN package 
	ON deployment.id = package.deploy_id 
	WHERE deployment.deleted_at IS NULL 
	AND deployment.key = ? 
	AND package.label = ?) 
	LIMIT 1`, self.Key, label).Get(metrics)
	if exist {
		return metrics, err
	}
	return nil, err
}

func (self *Deployment) GetPackageMetricsByAppIDAndName() ([]*PackageMetricsDetail, error) {

	packageMetricsDetails := make([]*PackageMetricsDetail, 0, 50)

	if err := x.Sql(`SELECT a.*,b.label 
	FROM package_metrics AS a 
	INNER JOIN
	package AS b
	ON a.package_id = b.id
	WHERE b.deploy_id IN
	 (SELECT id FROM deployment
	 WHERE deployment.deleted_at IS NULL AND deployment.app_id = ? 
	 AND deployment.name = ?) 
	LIMIT 1000`, self.AppID, self.Name).Find(&packageMetricsDetails); err != nil {
		return nil, err
	}
	return packageMetricsDetails, nil
}

func (self *Deployment) Delete(engine Engine) error {
	engine = EngineGenerate(engine)
	if _, err := engine.Id(self.ID).Delete(new(Deployment)); err != nil {
		return err
	}
	return nil
}

func (self *Deployment) DeleteByAppID(engine Engine) error {
	engine = EngineGenerate(engine)
	if _, err := engine.Where("app_id = ?", self.AppID).Delete(new(Deployment)); err != nil {
		return err
	}

	return nil
}

func (self *Deployment) DeleteAllHistory(engine Engine) error {
	engine = EngineGenerate(engine)
	if _, err := engine.Where("deploy_id = ?", self.ID).Delete(new(DeploymentHistory)); err != nil {
		return err
	}

	return nil
}

func (self *Deployment) DeleteAllVersion(engine Engine) error {
	engine = EngineGenerate(engine)
	if _, err := engine.Where("deploy_id = ?", self.ID).Delete(new(DeploymentVersion)); err != nil {
		return err
	}

	return nil
}

func (self *Deployment) DeleteAllPackage(engine Engine) error {
	engine = EngineGenerate(engine)
	engine.Exec("UPDATE package_diff SET deleted_at = ? WHERE package_id IN (SELECT id FROM package WHERE deploy_id = ? AND deleted_at IS NULL)", time.Now(), self.ID)
	engine.Exec("UPDATE package_metrics SET deleted_at = ? WHERE package_id IN (SELECT id FROM package WHERE deploy_id = ? AND deleted_at IS NULL)", time.Now(), self.ID)
	if _, err := engine.Where("deploy_id = ?", self.ID).Delete(new(Package)); err != nil {
		return err
	}

	return nil
}

func FindDeploymentsByAppIDs(appIds []uint64) ([]*Deployment, error) {
	deploys := make([]*Deployment, 0, 10)
	err := x.In("app_id", appIds).Find(&deploys)
	return deploys, err
}

func FindDeploymentDetailsByAppID(appID uint64) ([]*DeploymentDetail, error) {

	deploymentDetails := make([]*DeploymentDetail, 0, 25)
	if err := x.Sql(`SELECT deployment.*, deployment_version.*, package.*, package_metrics.* 
	 FROM deployment
	 LEFT JOIN deployment_version 
	 ON deployment.last_version_id = deployment_version.id 
	 LEFT JOIN package 
	 ON deployment_version.package_id = package.id 
	 LEFT JOIN package_metrics 
	 ON package_metrics.package_id = package.id 
	 WHERE deployment.deleted_at IS NULL AND deployment.app_id = ? 
	 LIMIT  50`, appID).Find(&deploymentDetails); err != nil {
		return nil, err
	}
	return deploymentDetails, nil
}

type DeploymentHistory struct {
	ID        uint64    `xorm:"pk autoincr"`
	DeployID  uint64    `xorm:"notnull default(0)"`
	PackageID uint64    `xorm:"notnull default(0)"`
	CreatedAt time.Time `xorm:"created"`
	DeletedAt time.Time `xorm:"deleted"`
}

func (self *DeploymentHistory) Update(engine Engine, specFields ...string) error {
	return ModelUpdate(engine, self.ID, self, specFields...)
}

func (self *DeploymentHistory) Exist() (bool, error) {
	exi, err := x.Where("id!=?", 0).Get(&DeploymentHistory{PackageID: self.PackageID, DeployID: self.DeployID})
	return exi, err
}

func (self *DeploymentHistory) Create(engine Engine) error {
	return ModelCreate(engine, self)
}

func (self *DeploymentHistory) Get() (bool, error) {

	return ModelGet(self.ID, self)
}

func (self *DeploymentHistory) Delete(engine Engine) error {
	engine = EngineGenerate(engine)
	if _, err := engine.Id(self.ID).Delete(new(DeploymentHistory)); err != nil {
		return err
	}

	return nil
}

type DeploymentVersion struct {
	ID         uint64 `xorm:"pk autoincr"`
	DeployID   uint64 `xorm:"index notnull default(0)"`
	PackageID  uint64 `xorm:"notnull default(0)"`
	AppVersion string
	CreatedAt  time.Time `xorm:"created"`
	UpdatedAt  time.Time `xorm:"updated"`
	DeletedAt  time.Time `xorm:"deleted"`
}

func (self *DeploymentVersion) Update(engine Engine, specFields ...string) error {
	return ModelUpdate(engine, self.ID, self, specFields...)
}

func (self *DeploymentVersion) Exist() (bool, error) {
	exi, err := x.Where("id!=?", 0).Get(&DeploymentVersion{AppVersion: self.AppVersion, DeployID: self.DeployID})
	return exi, err
}

func (self *DeploymentVersion) Create(engine Engine) error {
	return ModelCreate(engine, self)
}

func (self *DeploymentVersion) Get() (bool, error) {

	return ModelGet(self.ID, self)
}

func (self *DeploymentVersion) Delete(engine Engine) error {
	engine = EngineGenerate(engine)
	if _, err := engine.Id(self.ID).Delete(new(DeploymentVersion)); err != nil {
		return err
	}

	return nil
}
