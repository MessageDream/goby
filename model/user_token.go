package model

import (
	"time"

	"github.com/MessageDream/goby/model/dto"
)

type UserToken struct {
	ID          uint64 `xorm:"pk autoincr"`
	UID         uint64 `xorm:"unique(s) notnull default(0)"`
	Name        string `xorm:"unique(s) notnull"`
	Token       string `xorm:"notnull"`
	CreatedBy   string
	IsSession   bool `xorm:"notnull default(0)"`
	Description string
	ExpiresAt   time.Time
	CreatedAt   time.Time `xorm:"created"`
	DeletedAt   time.Time `xorm:"deleted"`
}

func (self *UserToken) Update(engine Engine, specFields ...string) error {
	return ModelUpdate(engine, self.ID, self, specFields...)
}

func (self *UserToken) Exist() (bool, error) {
	exi, err := x.Where("id!=?", 0).Get(&UserToken{Name: self.Name})
	return exi, err
}

func (self *UserToken) Create(engine Engine) error {
	return ModelCreate(engine, self)
}

func (self *UserToken) Get() (bool, error) {

	return ModelGet(self.ID, self)
}

func (self *UserToken) GetUnexpired() (bool, error) {
	sess := x.Where("expires_at>=?", time.Now())
	if self.ID <= 0 {
		sess.And("id != ?", self.ID)
	}
	has, err := sess.Get(self)
	if err != nil {
		return false, err
	}
	return has, nil
}

func (self *UserToken) Delete(engine Engine) error {
	engine = EngineGenerate(engine)
	if _, err := engine.Id(self.ID).Delete(new(UserToken)); err != nil {
		return err
	}

	return nil
}

func (self *UserToken) DeleteByNameAndUID(engine Engine) error {
	engine = EngineGenerate(engine)
	if _, err := engine.Where("uid = ?", self.UID).And("name = ?", self.Name).Delete(new(UserToken)); err != nil {
		return err
	}
	return nil
}

func (self *UserToken) DeleteByCreatorAndUID(engine Engine) error {
	engine = EngineGenerate(engine)
	if _, err := engine.Where("uid = ?", self.UID).And("created_by = ?", self.CreatedBy).Delete(new(UserToken)); err != nil {
		return err
	}
	return nil
}

func (self *UserToken) Convert() interface{} {
	return &dto.UserToken{
		ID:           self.ID,
		Name:         self.Token,
		FriendlyName: self.Name,
		Expires:      self.ExpiresAt,
		IsSession:    self.IsSession,
		CreatedBy:    self.CreatedBy,
		CreatedTime:  self.CreatedAt,
		Description:  self.Description,
	}
}
