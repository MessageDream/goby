package model

import (
	"time"
)

type App struct {
	ID        uint64    `xorm:"pk autoincr"`
	Name      string    `xorm:"unique index notnull"`
	UID       uint64    `xorm:"notnull default(0)"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}

func (self *App) Deployments() ([]*Deployment, error) {
	deploys := make([]*Deployment, 0, 10)
	err := x.Find(&deploys, &Deployment{AppID: self.ID})
	return deploys, err
}

func (self *App) Collaborators() ([]*Collaborator, error) {
	collaborators := make([]*Collaborator, 0, 10)
	err := x.Find(&collaborators, &Collaborator{AppID: self.ID})
	return collaborators, err
}

func (self *App) Update(engine Engine, specFields ...string) error {
	return ModelUpdate(engine, self.ID, self, specFields...)
}

func (self *App) Exist() (bool, error) {
	exi, err := x.Where("id!=?", 0).Get(&App{Name: self.Name, UID: self.UID})
	return exi, err
}

func (self *App) Create(engine Engine) error {
	return ModelCreate(engine, self)
}

func (self *App) Get() (bool, error) {
	return ModelGet(self.ID, self)
}

func (self *App) Delete(engine Engine) error {
	engine = EngineGenerate(engine)
	if _, err := engine.Id(self.ID).Delete(new(App)); err != nil {
		return err
	}
	return nil
}

func FindAppsByIDs(ids []uint64) ([]*App, error) {
	apps := make([]*App, 0, 10)
	err := x.In("id", ids).Find(&apps)
	return apps, err
}
