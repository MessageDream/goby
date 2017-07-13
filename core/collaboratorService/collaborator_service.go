package collaboratorService

import (
	"errors"

	. "github.com/MessageDream/goby/core"
	"github.com/MessageDream/goby/model"
	"github.com/MessageDream/goby/module/infrastructure"
)

func List(uid uint64, appName string) ([]*model.Collaborator, error) {
	col, err := CollaboratorOf(uid, appName)
	if err != nil {
		return nil, err
	}

	return model.FindCollaboratorsByAppIDs([]uint64{col.AppID})

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
