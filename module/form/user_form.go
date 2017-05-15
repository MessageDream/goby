package form

import (
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
)

type InstallForm struct {
	RunUser  string `form:"run_user" binding:"Required"`
	Domain   string `form:"domain" binding:"Required"`
	AppName  string `form:"app_name"`
	AppURL   string `form:"app_url" binding:"Required"`
	HTTPPort string `form:"http_port"`

	DbType   string `form:"db_type" binding:"Required"`
	DbHost   string `form:"db_host"`
	DbUser   string `form:"db_user"`
	DbPasswd string `form:"db_passwd"`
	DbName   string `form:"db_name"`
	SSLMode  string `form:"ssl_mode"`
	DbPath   string `form:"db_path"`

	StorageType        string `form:"storage_type"`
	StorageDownloadURL string `form:"storage_download_url"`
	StoragePath        string `form:"storage_path"`
	StorageAccessKey   string `form:"storage_access_key"`
	StorageSecretKey   string `form:"storage_secret_key"`
	StorageBucketName  string `form:"storage_bucket_name"`
	StorageEndpoint    string `form:"storage_endpoint"`
	StorageZone        string `form:"storage_zone"`
	StoragePrefix      string `form:"storage_prefix"`

	LogRootPath string `form:"log_path"`

	SMTPHost   string `form:"smtp_host"`
	SMTPUser   string `form:"mailer_user"`
	SMTPFrom   string `form:"mailer_from"`
	SMTPPasswd string `form:"mailer_pwd"`

	RegisterConfirm     bool `form:"register_confirm"`
	MailNotify          bool `form:"mail_notify"`
	DisableRegistration bool `form:"disable_registration"`
	EnableCaptcha       bool `form:"enable_captcha"`
	RequireSignInView   bool `form:"require_sign_in_view"`

	AdminName     string `form:"admin_name" binding:"Required;AlphaDashDot;MaxSize(30)"`
	AdminPasswd   string `form:"admin_pwd" binding:"Required;MinSize(6);MaxSize(255)"`
	ConfirmPasswd string `form:"confirm_passwd" binding:"Required;MinSize(6);MaxSize(255)"`
	AdminEmail    string `form:"admin_email" binding:"Required;Email;MaxSize(50)"`
}

func (f *InstallForm) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f)
}

type SignUpForm struct {
	UserName string `form:"user_name" binding:"Required;AlphaDashDot;MaxSize(35)"`
	Email    string `form:"email" binding:"Required;Email;MaxSize(50)"`
	Password string `form:"password" binding:"Required;MinSize(6);MaxSize(255)"`
	Retype   string
}

func (f *SignUpForm) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f)
}

type SignInForm struct {
	UserName string `form:"user_name" binding:"Required;MaxSize(35)"`
	Password string `form:"password" binding:"Required;MinSize(6);MaxSize(255)"`
	Remember bool   `form:"remember"`
}

func (f *SignInForm) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f)
}

type ChangePasswordForm struct {
	OldPassword string `form:"old_password" binding:"Required;MinSize(6);MaxSize(255)"`
	Password    string `form:"password" binding:"Required;MinSize(6);MaxSize(255)"`
	Retype      string `form:"retype"`
}

func (f *ChangePasswordForm) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f)
}
