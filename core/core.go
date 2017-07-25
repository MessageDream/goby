package core

import (
	"errors"

	"github.com/MessageDream/goby/model"
	"github.com/MessageDream/goby/module/infrastructure"
)

func WrapIntentError(err error, code int) error {
	return infrastructure.MakeIntentError(err, code)
}

func OwnerCan(uid uint64, name string) (*model.Collaborator, error) {
	col, err := model.FindCollaboratorByAppNameAndUID(name, uid)

	if err != nil || col == nil {
		if err != nil {
			return nil, err
		}
		return nil, ErrAppNotExist
	}

	if col.Role != 1 {
		return nil, WrapIntentError(errors.New("Permission Deny, You are not owner"), INTENT_ERROR_CODE_PERMISSION_DENY)
	}

	return col, nil
}

func OwnerOf(name string) (*model.Collaborator, error) {
	col, err := model.FindOwnerByAppName(name)
	if err != nil || col == nil {
		if err != nil {
			return nil, err
		}
		return nil, ErrAppNotExist
	}
	return col, nil
}

func CollaboratorOf(uid uint64, name string) (*model.Collaborator, error) {
	col, err := model.FindCollaboratorByAppNameAndUID(name, uid)

	if err != nil || col == nil {
		if err != nil {
			return nil, err
		}
		return nil, ErrAppNotExist
	}

	return col, nil
}

const (
	APP_SUFFIX_IOS     = "-ios"
	APP_SUFFIX_ANDROID = "-android"
)

//common
const (
	INTENT_ERROR_CODE_NONE_BUSINESS = iota + 1000
	INTENT_ERROR_CODE_PERMISSION_DENY
	INTENT_ERROR_CODE_DATA_ALREADY_EXIST
	INTENT_ERROR_CODE_DATA_NOT_EXIST
	INTENT_ERROR_CODE_DATA_PARAMS_ERR
)

//User
const (
	INTENT_ERROR_CODE_USER_NAME_NOT_ALLOWNED = iota + 1100

	INTENT_ERROR_CODE_USER_ACTIVATE_EMAIL
	INTENT_ERROR_CODE_USER_ACTIVATE_HAS_ACTIVATED
	INTENT_ERROR_CODE_USER_ACTIVATE_TIME_LIMIT_LENGTH
	INTENT_ERROR_CODE_USER_ACTIVATE_VERIFY_FAILED

	INTENT_ERROR_CODE_USER_USERNAME_OR_PWD_INVALID
	INTENT_ERROR_CODE_USER_FORBIDDEN
)

var (
	ErrUserNameOrPasswordInvalide = WrapIntentError(errors.New("User name or password is invalide"), INTENT_ERROR_CODE_USER_USERNAME_OR_PWD_INVALID)

	ErrUserNotExist                    = WrapIntentError(errors.New("User not exist"), INTENT_ERROR_CODE_DATA_NOT_EXIST)
	ErrNameReserved                    = WrapIntentError(errors.New("User name is reserved"), INTENT_ERROR_CODE_DATA_PARAMS_ERR)
	ErrNameEmpty                       = WrapIntentError(errors.New("User name can't be nil"), INTENT_ERROR_CODE_DATA_PARAMS_ERR)
	ErrUserNameAlreadyExist            = WrapIntentError(errors.New("User name has been used"), INTENT_ERROR_CODE_DATA_ALREADY_EXIST)
	ErrEmailInvalide                   = WrapIntentError(errors.New("Email is invalide"), INTENT_ERROR_CODE_DATA_PARAMS_ERR)
	ErrEmailNotActivated               = WrapIntentError(errors.New("E-mail address has not been activated"), INTENT_ERROR_CODE_USER_ACTIVATE_EMAIL)
	ErrEmailAlreadyExist               = WrapIntentError(errors.New("E-mail address has been used"), INTENT_ERROR_CODE_DATA_ALREADY_EXIST)
	ErrUserAlreadyActivated            = WrapIntentError(errors.New("User has been activated before"), INTENT_ERROR_CODE_USER_ACTIVATE_HAS_ACTIVATED)
	ErrUserActivateTimeLimitCodeLength = WrapIntentError(errors.New("Time limit code length is too short"), INTENT_ERROR_CODE_USER_ACTIVATE_TIME_LIMIT_LENGTH)
	ErrUserActivateVerifyFailed        = WrapIntentError(errors.New("Active code is invalid"), INTENT_ERROR_CODE_USER_ACTIVATE_VERIFY_FAILED)
	ErrUserForbidden                   = WrapIntentError(errors.New("User is forbidden"), INTENT_ERROR_CODE_USER_FORBIDDEN)
)

//Token
const (
	INTENT_ERROR_CODE_TOKEN = iota + 1200
)

//Token
var (
	ErrTokenNotExist     = WrapIntentError(errors.New("Token not exist"), INTENT_ERROR_CODE_DATA_NOT_EXIST)
	ErrTokenAlreadyExist = WrapIntentError(errors.New("Token has existed"), INTENT_ERROR_CODE_DATA_ALREADY_EXIST)
)

//App
const (
	INTENT_ERROR_CODE_APP_PARAMS_ERR = iota + 1300
)

//App
var (
	ErrAppNameEmpty    = WrapIntentError(errors.New("App name is empty"), INTENT_ERROR_CODE_DATA_PARAMS_ERR)
	ErrAppNotExist     = WrapIntentError(errors.New("App not exist"), INTENT_ERROR_CODE_DATA_NOT_EXIST)
	ErrAppAlreadyExist = WrapIntentError(errors.New("App has existed"), INTENT_ERROR_CODE_DATA_ALREADY_EXIST)
)

//Collaborator
const (
	INTENT_ERROR_CODE_COLLABORATOR = iota + 1400
)

//Collaborator
var (
	ErrCollaboratorNotExist     = WrapIntentError(errors.New("Collaborator not exist"), INTENT_ERROR_CODE_DATA_NOT_EXIST)
	ErrCollaboratorAlreadyExist = WrapIntentError(errors.New("Collaborator has existed"), INTENT_ERROR_CODE_DATA_ALREADY_EXIST)
)

//Deployment
const (
	INTENT_ERROR_CODE_DEPLOYMENT_PACKAGE_FILE_TYPE_NOT_ALLOWNED = iota + 1500
	INTENT_ERROR_CODE_DEPLOYMENT_PACKAGE_CONTENT_FILE_UNRECOGNIZED
	INTENT_ERROR_CODE_DEPLOYMENT_PACKAGE_CONTENT_FILE_OF_ANDROID
	INTENT_ERROR_CODE_DEPLOYMENT_PACKAGE_CONTENT_FILE_OF_IOS
	INTENT_ERROR_CODE_DEPLOYMENT_PACKAGE_VERSION_NAME_IS_INVALID
	INTENT_ERROR_CODE_DEPLOYMENT_PACKAGE_ROLLOUT_IS_INVALID

	INTENT_ERROR_CODE_DEPLOYMENT_ROLLBACK_NOT_RELEASE_PACKAGE_BEFORE
	INTENT_ERROR_CODE_DEPLOYMENT_ROLLBACK_HAVE_NO_PACKAGE_TO_ROLLBACK
)

//Deployment
var (
	ErrDeploymentNameEmpty              = WrapIntentError(errors.New("Deployment name is empty"), INTENT_ERROR_CODE_DATA_PARAMS_ERR)
	ErrDeploymentNotExist               = WrapIntentError(errors.New("Deployment not exist"), INTENT_ERROR_CODE_DATA_NOT_EXIST)
	ErrDeploymentOfThisNameAlreadyExist = WrapIntentError(errors.New("Deployment of this name has existed"), INTENT_ERROR_CODE_DATA_ALREADY_EXIST)
	ErrDeploymentPackageNotExist        = WrapIntentError(errors.New("Package not exist"), INTENT_ERROR_CODE_DATA_NOT_EXIST)
	ErrDeploymentPackageAlreadyExist    = WrapIntentError(errors.New("The same hash package of this deployment has existed"), INTENT_ERROR_CODE_DATA_ALREADY_EXIST)
	ErrDeploymentPackageRolloutInvalide = WrapIntentError(errors.New("The package rollout must be an aliquot of 10 and from 10 to 100"), INTENT_ERROR_CODE_DEPLOYMENT_PACKAGE_ROLLOUT_IS_INVALID)

	ErrDeploymentPackageFileTypeNotAllowned           = WrapIntentError(errors.New("Package file of this type is not allowned"), INTENT_ERROR_CODE_DEPLOYMENT_PACKAGE_FILE_TYPE_NOT_ALLOWNED)
	ErrDeploymentPackageContentsUnrecognized          = WrapIntentError(errors.New("The file format of this package are unrecognized"), INTENT_ERROR_CODE_DEPLOYMENT_PACKAGE_CONTENT_FILE_UNRECOGNIZED)
	ErrDeploymentPackageContentsShouldBeIOSFormat     = WrapIntentError(errors.New("The file format of this package can not be support by iOS"), INTENT_ERROR_CODE_DEPLOYMENT_PACKAGE_CONTENT_FILE_OF_IOS)
	ErrDeploymentPackageContentsShouldBeAndroidFormat = WrapIntentError(errors.New("The file format of this package can not be support by Android"), INTENT_ERROR_CODE_DEPLOYMENT_PACKAGE_CONTENT_FILE_OF_ANDROID)
	ErrDeploymentPackageVersionNameFormatNotAllowned  = WrapIntentError(errors.New("The app version name is empty or format is not allowned"), INTENT_ERROR_CODE_DEPLOYMENT_PACKAGE_VERSION_NAME_IS_INVALID)
	ErrDeploymentRollbackNotReleasePackageBefore      = WrapIntentError(errors.New("You have not release a package before"), INTENT_ERROR_CODE_DEPLOYMENT_ROLLBACK_NOT_RELEASE_PACKAGE_BEFORE)
	ErrDeploymentRollbackHaveNoPackageToRollback      = WrapIntentError(errors.New("Have no package to rollback"), INTENT_ERROR_CODE_DEPLOYMENT_ROLLBACK_HAVE_NO_PACKAGE_TO_ROLLBACK)
)

//Client
const (
	INTENT_ERROR_CODE_CLIENT_DEPLOYMENT_KEY_OR_APP_VERSION_IS_EMPTY = iota + 1600
	INTENT_ERROR_CODE_CLIENT_DEPLOYMENT_KEY_OR_LABEL_IS_EMPTY
)

//Client
var (
	ErrClientDeploymentKeyOrAppVersionEmpty = WrapIntentError(errors.New("Deployment key or app version is empty"), INTENT_ERROR_CODE_CLIENT_DEPLOYMENT_KEY_OR_APP_VERSION_IS_EMPTY)
	ErrClientDeploymentKeyOrLabelEmpty      = WrapIntentError(errors.New("Deployment key or Label is empty"), INTENT_ERROR_CODE_CLIENT_DEPLOYMENT_KEY_OR_LABEL_IS_EMPTY)
)
