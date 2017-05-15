package appService

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-xorm/xorm"

	. "github.com/MessageDream/goby/core"
	"github.com/MessageDream/goby/model"
	"github.com/MessageDream/goby/model/dto"
	"github.com/MessageDream/goby/module/infrastructure"
	. "gopkg.in/ahmetb/go-linq.v3"
)

func List(user *model.User) ([]*dto.App, error) {
	apps, err := user.Apps()
	if err != nil {
		return nil, err
	}

	if len(apps) == 0 {
		return nil, nil
	}

	ids := make([]uint64, 0, 10)

	From(apps).Select(func(app interface{}) interface{} {
		return app.(*model.App).ID
	}).Distinct().ToSlice(&ids)

	cols, err := model.FindCollaboratorsByAppIDs(ids)
	if err != nil {
		return nil, err
	}

	deploys, err := model.FindDeploymentsByAppIDs(ids)
	if err != nil {
		return nil, err
	}

	uids := make([]uint64, 0, 10)
	From(cols).Select(func(cl interface{}) interface{} {
		return cl.(*model.Collaborator).UID
	}).ToSlice(&uids)

	results := make([]*dto.App, 0, 10)
	From(apps).Select(func(app interface{}) interface{} {
		ap := app.(*model.App)
		name := ap.Name
		appID := ap.ID

		canOwner := false
		colsMap := map[string]*dto.Collaborator{}
		From(cols).Where(func(col interface{}) bool {
			return col.(*model.Collaborator).AppID == appID
		}).Select(func(col interface{}) interface{} {
			co := col.(*model.Collaborator)
			isOwner := user.ID == co.UID
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

		var depls []string
		From(deploys).Where(func(dyp interface{}) bool {
			return dyp.(*model.Deployment).AppID == appID
		}).Select(func(dp interface{}) interface{} {
			return dp.(*model.Deployment).Name
		}).ToSlice(&depls)

		return &dto.App{
			Collaborators: colsMap,
			Deployments:   depls,
			Name:          name,
			IsOwner:       canOwner,
		}
	}).ToSlice(&results)

	return results, nil
}

func Create(user *model.User, appName string) (*dto.App, error) {

	if len(appName) == 0 {
		return nil, ErrAppNameEmpty

	}
	if !strings.HasSuffix(appName, APP_SUFFIX_IOS) && !strings.HasSuffix(appName, APP_SUFFIX_ANDROID) {
		return nil, WrapIntentError(fmt.Errorf("appName have to point -android or -ios suffix! eg. %s-android | %s-ios", appName, appName), INTENT_ERROR_CODE_DATA_PARAMS_ERR)
	}

	app := &model.App{UID: user.ID, Name: appName}
	if exist, err := app.Exist(); err != nil || exist {
		if exist {
			return nil, ErrAppAlreadyExist
		}
		return nil, err
	}

	result, err := model.Transaction(func(sess *xorm.Session) (interface{}, error) {

		if _, err := sess.Insert(app); err != nil {
			return nil, err
		}

		deployProduct := &model.Deployment{
			AppID: app.ID,
			Name:  "Production",
			Key:   infrastructure.GetRandomString(28) + user.Rands,
		}

		if err := deployProduct.Create(sess); err != nil {
			return nil, err
		}

		deployStaging := &model.Deployment{
			AppID: app.ID,
			Name:  "Staging",
			Key:   infrastructure.GetRandomString(28) + user.Rands,
		}

		if err := deployStaging.Create(sess); err != nil {
			return nil, err
		}

		col := &model.Collaborator{
			AppID: app.ID,
			UID:   user.ID,
			Email: user.Email,
			Role:  1,
		}
		if err := col.Create(sess); err != nil {
			return nil, err
		}

		return &dto.App{
			Collaborators: map[string]*dto.Collaborator{user.Email: &dto.Collaborator{IsCurrentAccount: true, Permission: "Owner"}},
			Deployments:   []string{"Production", "Staging"},
			Name:          app.Name,
			IsOwner:       true,
		}, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*dto.App), nil

}

func Delete(uid uint64, appName string) error {
	col, err := OwnerCan(uid, appName)
	if err != nil {
		return err
	}

	_, err = model.Transaction(func(sess *xorm.Session) (interface{}, error) {
		if err := (&model.App{ID: col.AppID}).Delete(sess); err != nil {
			return nil, err
		}

		if err := col.DeleteByAppID(sess); err != nil {
			return nil, err
		}

		if err := (&model.Deployment{AppID: col.AppID}).DeleteByAppID(sess); err != nil {
			return nil, err
		}
		return nil, nil
	})

	return err
}

func ChangeName(uid uint64, appName, newName string) error {

	if len(newName) == 0 || len(appName) == 0 {
		return ErrAppNameEmpty
	}
	if !strings.HasSuffix(newName, APP_SUFFIX_IOS) && !strings.HasSuffix(newName, APP_SUFFIX_ANDROID) {
		return WrapIntentError(fmt.Errorf("appName have to point -android or -ios suffix! eg. %s-android | %s-ios", newName, newName), INTENT_ERROR_CODE_DATA_PARAMS_ERR)
	}

	_, err := OwnerCan(uid, appName)
	if err != nil {
		return err
	}

	app := &model.App{Name: appName}

	if exist, err := app.Get(); err != nil || !exist {
		if !exist {
			return ErrAppNotExist
		}
		return err
	}

	app.Name = newName
	return app.Update(nil)
}

func Transfer(uid uint64, appName, fromEmail, toEmail string) error {

	if len(appName) == 0 {
		return ErrAppNameEmpty
	}

	if !infrastructure.VerifyEmail(toEmail) {
		return ErrEmailInvalide
	}

	if toEmail == fromEmail {
		return WrapIntentError(errors.New("You can't transfer to yourself!"), INTENT_ERROR_CODE_DATA_PARAMS_ERR)

	}
	fromCol, err := OwnerCan(uid, appName)
	if err != nil {
		return err
	}
	toUser := model.User{Email: toEmail}
	if exist, err := toUser.Get(); err != nil || !exist {
		if !exist {
			return ErrUserNotExist
		}
		return err
	}

	app := &model.App{Name: appName}
	if exist, err := app.Get(); err != nil || !exist {
		if !exist {
			return ErrAppNotExist
		}
		return err
	}
	_, err = model.Transaction(func(sess *xorm.Session) (interface{}, error) {
		app.UID = toUser.ID
		if err := app.Update(sess, "uid"); err != nil {
			return nil, err
		}

		if err := fromCol.DeleteByAppIDAndUID(sess); err != nil {
			return nil, err
		}

		toCol := &model.Collaborator{
			AppID: app.ID,
			UID:   toUser.ID,
		}

		if err := toCol.DeleteByAppIDAndUID(sess); err != nil {
			return nil, err
		}

		newCol := &model.Collaborator{
			AppID: app.ID,
			UID:   toUser.ID,
			Role:  1,
			Email: toUser.Email,
		}
		if err := newCol.Create(sess); err != nil {
			return nil, err
		}

		return nil, nil
	})
	return err
}
