package model

import (
	"time"
)

type Collaborator struct {
	ID        uint64    `xorm:"pk autoincr"`
	Role      int8      `xorm:"notnull default(0)"`
	UID       uint64    `xorm:"index notnull default(0)"`
	AppID     uint64    `xorm:"index notnull default(0)"`
	Email     string    `xorm:"index notnull"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}

func (self *Collaborator) App() (*App, error) {
	app := &App{ID: self.AppID}
	has, err := x.Get(app)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return app, nil
}

func (self *Collaborator) Update(engine Engine, specFields ...string) error {
	return ModelUpdate(engine, self.ID, self, specFields...)
}

func (self *Collaborator) Exist() (bool, error) {
	exi, err := x.Where("id!=?", 0).Get(&Collaborator{AppID: self.AppID, UID: self.UID})
	return exi, err
}

func (self *Collaborator) Create(engine Engine) error {
	return ModelCreate(engine, self)
}

func (self *Collaborator) Get() (bool, error) {
	return ModelGet(self.ID, self)
}

func (self *Collaborator) Delete(engine Engine) error {
	engine = EngineGenerate(engine)
	if _, err := engine.Id(self.ID).Delete(new(Collaborator)); err != nil {
		return err
	}
	return nil
}

func (self *Collaborator) DeleteByAppID(engine Engine) error {
	engine = EngineGenerate(engine)
	if _, err := engine.Where("app_id = ?", self.AppID).Delete(new(Collaborator)); err != nil {
		return err
	}

	return nil
}

func (self *Collaborator) DeleteByAppIDAndUID(engine Engine) error {
	engine = EngineGenerate(engine)
	if _, err := engine.Where("app_id = ?", self.AppID).And("uid = ?", self.UID).Delete(new(Collaborator)); err != nil {
		return err
	}

	return nil
}

func FindCollaboratorsByAppIDs(appIds []uint64) ([]*Collaborator, error) {
	cols := make([]*Collaborator, 0, 10)
	return cols, x.In("app_id", appIds).Find(&cols)
}

func FindOwnerByAppName(appName string) (*Collaborator, error) {
	col := &Collaborator{}

	exsit, err := x.Sql(`SELECT b.* FROM app AS a 
	INNER JOIN collaborator AS b 
	ON a.id = b.app_id
	WHERE a.name = ? AND b.Role = 1 AND a.deleted_at IS NULL AND b.deleted_at IS NULL 
    LIMIT 1`, appName).Get(col)

	if err != nil {
		return nil, err
	}
	if !exsit {
		return nil, nil
	}
	return col, nil
}

func FindCollaboratorByAppNameAndUID(appName string, uid uint64) (*Collaborator, error) {
	col := &Collaborator{}

	exsit, err := x.Sql(`SELECT b.* FROM app AS a 
	INNER JOIN collaborator AS b 
	ON a.id = b.app_id
	WHERE a.name = ? AND b.uid = ? AND a.deleted_at IS NULL AND b.deleted_at IS NULL 
    LIMIT 1`, appName, uid).Get(col)

	if err != nil {
		return nil, err
	}
	if !exsit {
		return nil, nil
	}
	return col, nil
}
