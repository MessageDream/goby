package model

import (
	"time"
)

type PackageDetail struct {
	Package     `xorm:"extends"`
	PackageDiff `xorm:"extends"`
}

type Package struct {
	ID                 uint64 `xorm:"pk autoincr"`
	DeployID           uint64 `xorm:"index notnull default(0)"`
	DeployVersionID    uint64 `xorm:"index notnull default(0)"`
	Description        string
	Hash               string
	BlobURL            string
	Size               int64 `xorm:"notnull default(0)"`
	ManifestBlobURL    string
	ReleaseMethod      string
	Label              string `xorm:"index notnull"`
	OriginalLabel      string
	OriginalDeployment string
	ReleasedBy         string
	IsMandatory        bool      `xorm:"notnull default(0)"`
	IsDisabled         bool      `xorm:"notnull default(0)"`
	Rollout            uint8     `xorm:"notnull default(100)"`
	CreatedAt          time.Time `xorm:"created"`
	UpdatedAt          time.Time `xorm:"updated"`
	DeletedAt          time.Time `xorm:"deleted"`
}

func (self *Package) Update(engine Engine, specFields ...string) error {
	return ModelUpdate(engine, self.ID, self, specFields...)
}

func (self *Package) Exist() (bool, error) {
	exi, err := x.Where("id!=?", 0).Get(&Package{Hash: self.Hash, DeployID: self.DeployID})
	return exi, err
}

func (self *Package) Create(engine Engine) error {
	return ModelCreate(engine, self)
}

func (self *Package) Get() (bool, error) {

	return ModelGet(self.ID, self)
}

func (self *Package) Delete(engine Engine) error {
	engine = EngineGenerate(engine)
	if _, err := engine.Id(self.ID).Delete(new(Package)); err != nil {
		return err
	}

	return nil
}

func (self *Package) FindPrePackages(number int, isAsc bool) ([]*Package, error) {

	packages := make([]*Package, 0, number)
	sess := x.Where("id < ?", self.ID).And("deploy_version_id = ?", self.DeployVersionID)
	if isAsc {
		sess = sess.Asc("id")
	} else {
		sess = sess.Desc("id")
	}
	if err := sess.Limit(number).Find(&packages); err != nil {
		return nil, err
	}
	return packages, nil

}

func (self *Package) FindPackagesByVersionIDAndReleaseMethods(number int, isAsc bool, methods ...string) ([]*Package, error) {

	packages := make([]*Package, 0, number)
	sess := x.Where("deploy_version_id = ?", self.DeployVersionID).In("release_method", methods)
	if isAsc {
		sess = sess.Asc("id")
	} else {
		sess = sess.Desc("id")
	}
	if err := sess.Limit(number).Find(&packages); err != nil {
		return nil, err
	}
	return packages, nil

}

type PackageDiff struct {
	ID                     uint64 `xorm:"pk autoincr"`
	OriginalPackageHash    string
	DiffAgainstPackageHash string
	DiffBlobURL            string
	DiffSize               int64     `xorm:"notnull default(0)"`
	CreatedAt              time.Time `xorm:"created"`
	UpdatedAt              time.Time `xorm:"updated"`
	DeletedAt              time.Time `xorm:"deleted"`
}

func (self *PackageDiff) Update(engine Engine, specFields ...string) error {
	return ModelUpdate(engine, self.ID, self, specFields...)
}

func (self *PackageDiff) Exist() (bool, error) {
	exi, err := x.Where("id!=?", 0).Get(&PackageDiff{OriginalPackageHash: self.OriginalPackageHash, DiffAgainstPackageHash: self.DiffAgainstPackageHash})
	return exi, err
}

func (self *PackageDiff) Create(engine Engine) error {
	return ModelCreate(engine, self)
}

func (self *PackageDiff) Get() (bool, error) {

	return ModelGet(self.ID, self)
}

func (self *PackageDiff) Delete(engine Engine) error {
	engine = EngineGenerate(engine)
	if _, err := engine.Id(self.ID).Delete(new(PackageDiff)); err != nil {
		return err
	}

	return nil
}

type PackageMetrics struct {
	ID         uint64    `xorm:"pk autoincr"`
	PackageID  uint64    `xorm:"notnull default(0)"`
	Active     uint64    `xorm:"notnull default(0)"`
	Downloaded uint64    `xorm:"notnull default(0)"`
	Failed     uint64    `xorm:"notnull default(0)"`
	Installed  uint64    `xorm:"notnull default(0)"`
	CreatedAt  time.Time `xorm:"created"`
	UpdatedAt  time.Time `xorm:"updated"`
	DeletedAt  time.Time `xorm:"deleted"`
}

func (self *PackageMetrics) Update(engine Engine, specFields ...string) error {
	return ModelUpdate(engine, self.ID, self, specFields...)
}

func (self *PackageMetrics) Exist() (bool, error) {
	exi, err := x.Where("id!=?", 0).Get(&PackageMetrics{PackageID: self.PackageID})
	return exi, err
}

func (self *PackageMetrics) Create(engine Engine) error {
	return ModelCreate(engine, self)
}

func (self *PackageMetrics) Get() (bool, error) {

	return ModelGet(self.ID, self)
}

func (self *PackageMetrics) Delete(engine Engine) error {
	engine = EngineGenerate(engine)
	if _, err := engine.Id(self.ID).Delete(new(PackageMetrics)); err != nil {
		return err
	}

	return nil
}

type PackageMetricsDetail struct {
	*PackageMetrics `xorm:"extends"`
	Label           string
}
