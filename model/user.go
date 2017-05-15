package model

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/Unknwon/com"
	"golang.org/x/crypto/pbkdf2"

	"github.com/MessageDream/goby/model/dto"
	"github.com/MessageDream/goby/module/infrastructure"
	"github.com/MessageDream/goby/module/mailer"
)

type User struct {
	ID        uint64 `xorm:"pk autoincr"`
	Email     string `xorm:"unique notnull"`
	Password  string `xorm:"notnull"`
	Rands     string `xorm:"index"`
	UserName  string
	LowerName string `xorm:"unique notnull"`
	IsAdmin   bool   `xorm:"notnull default(0)"`
	IsActive  bool   `xorm:"notnull default(0)"`
	Salt      string
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}

func (self *User) EncodePasswd() {
	newPasswd := pbkdf2.Key([]byte(self.Password), []byte(self.Salt), 10000, 50, sha256.New)
	self.Password = fmt.Sprintf("%x", newPasswd)
}

func (self *User) ValidatePassword(passwd string) bool {
	newUser := &User{Password: passwd, Salt: self.Salt}
	newUser.EncodePasswd()
	return subtle.ConstantTimeCompare([]byte(self.Password), []byte(newUser.Password)) == 1
}

func (self *User) GenerateRands() {
	self.Rands = infrastructure.GetRandomString(9)
}

func (self *User) GenerateSalt() {
	self.Salt = infrastructure.GetRandomString(10)
}

func (self *User) Tokens() ([]*UserToken, error) {
	tokens := make([]*UserToken, 0, 10)
	err := x.Find(&tokens, &UserToken{UID: self.ID})
	return tokens, err
}

func (self *User) Collaborators() ([]*Collaborator, error) {
	collaborators := make([]*Collaborator, 0, 10)
	err := x.Find(&collaborators, &Collaborator{UID: self.ID})
	return collaborators, err
}

func (self *User) Apps() ([]*App, error) {
	apps := make([]*App, 0, 10)
	err := x.Sql("SELECT b.* FROM collaborator AS a LEFT JOIN app AS b ON a.app_id = b.id WHERE a.uid = ? AND b.deleted_at IS NULL AND a.deleted_at IS NULL", self.ID).Find(&apps)
	return apps, err
}

func (self *User) Update(engine Engine, specFields ...string) error {
	return ModelUpdate(engine, self.ID, self, specFields...)
}

func (self *User) Create(engine Engine) error {

	self.GenerateSalt()
	self.GenerateRands()
	self.EncodePasswd()
	self.LowerName = strings.ToLower(self.UserName)

	return ModelCreate(engine, self)
}

func (self *User) Get() (bool, error) {

	return ModelGet(self.ID, self)
}

func (self *User) Delete(engine Engine) error {
	engine = EngineGenerate(engine)
	if _, err := engine.Id(self.ID).Delete(new(User)); err != nil {
		return err
	}
	return nil
}

func (self *User) Exist() (bool, error) {
	exi, err := x.Where("id!=?", 0).Get(&User{LowerName: strings.ToLower(self.UserName)})
	return exi, err
}

func (self *User) EmailUsable() (bool, error) {
	has, err := x.Get(&User{Email: self.Email})
	if err != nil {
		return false, err
	}
	if has {
		return false, nil
	}
	return true, nil
}

func (self *User) Convert() interface{} {
	return &dto.User{
		ID:    self.Rands,
		Name:  self.UserName,
		Email: self.Email,
	}
}

func FindUserByIDs(ids []uint64) ([]*User, error) {
	users := make([]*User, 0, 10)
	return users, x.In("id", ids).Find(&users)
}

type mailerUser struct {
	user            *User
	activeCodeLives int
}

func (self mailerUser) ID() int64 {
	return int64(self.user.ID)
}

func (self mailerUser) DisplayName() string {
	return self.user.UserName
}

func (self mailerUser) Email() string {
	return self.user.Email
}

func (self mailerUser) GenerateActivateCode() string {
	return self.GenerateEmailActivateCode(self.user.Email)
}
func (self mailerUser) GenerateEmailActivateCode(email string) string {
	code := infrastructure.CreateTimeLimitCode(
		com.ToStr(self.user.ID)+email+self.user.LowerName+self.user.Password+self.user.Rands,
		self.activeCodeLives, nil)

	code += hex.EncodeToString([]byte(self.user.LowerName))
	return code
}

func NewMailerUser(u *User, activeCodeLives int) mailer.User {
	return mailerUser{u, activeCodeLives}
}
