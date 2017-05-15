package collaboratorService

import (
	"errors"

	. "gopkg.in/ahmetb/go-linq.v3"

	. "github.com/MessageDream/goby/core"
	"github.com/MessageDream/goby/model"
	"github.com/MessageDream/goby/model/dto"
	"github.com/MessageDream/goby/module/infrastructure"
)

func List(uid uint64, appName string) (map[string]*dto.Collaborator, error) {
	col, err := CollaboratorOf(uid, appName)
	if err != nil {
		return nil, err
	}

	cols, err := model.FindCollaboratorsByAppIDs([]uint64{col.AppID})
	if err != nil {
		return nil, err
	}

	canOwner := false
	colsMap := map[string]*dto.Collaborator{}
	From(cols).Select(func(col interface{}) interface{} {
		co := col.(*model.Collaborator)
		isOwner := uid == co.UID
		if isOwner && co.Role == 1 {
			canOwner = true
		}
		permission := "Member"
		if co.Role == 1 {
			permission = "Owner"
		}
		return KeyValue{
			Key: co.Email,
			Value: &dto.Collaborator{
				IsCurrentAccount: isOwner,
				Permission:       permission,
			},
		}

	}).ToMap(&colsMap)

	return colsMap, nil

}

func Add(user *model.User, appName, toEmail string) error {

	if len(appName) == 0 {
		return ErrAppNameEmpty
	}

	if !infrastructure.VerifyEmail(toEmail) {
		return ErrEmailInvalide
	}

	if toEmail == user.Email {
		return WrapIntentError(errors.New("You can't add yourself!"), INTENT_ERROR_CODE_DATA_PARAMS_ERR)
	}

	col, err := OwnerCan(user.ID, appName)
	if err != nil {
		return err
	}

	account := &model.User{
		Email: toEmail,
	}

	if exist, err := account.Get(); err != nil || !exist {
		if !exist {
			return ErrUserNotExist
		}
		return err
	}

	newCol := &model.Collaborator{
		AppID: col.AppID,
		UID:   account.ID,
		Email: account.Email,
		Role:  0,
	}

	return newCol.Create(nil)
}

func Delete(user *model.User, appName, toEmail string) error {
	if len(appName) == 0 {
		return ErrAppNameEmpty
	}

	if !infrastructure.VerifyEmail(toEmail) {
		return ErrEmailInvalide
	}

	if toEmail == user.Email {
		return WrapIntentError(errors.New("You can't delete yourself!"), INTENT_ERROR_CODE_DATA_PARAMS_ERR)
	}

	col, err := OwnerCan(user.ID, appName)
	if err != nil {
		return err
	}

	account := &model.User{
		Email: toEmail,
	}

	if exist, err := account.Get(); err != nil || !exist {
		if !exist {
			return ErrUserNotExist
		}
		return err
	}

	delCol := &model.Collaborator{
		AppID: col.AppID,
		UID:   account.ID,
	}

	return delCol.DeleteByAppIDAndUID(nil)
}
