package mailer

import (
	"fmt"

	log "gopkg.in/clog.v1"
	"gopkg.in/gomail.v2"
	"gopkg.in/macaron.v1"

	"github.com/MessageDream/goby/module/infrastructure"
)

const (
	MAIL_AUTH_ACTIVATE        infrastructure.TplName = "auth/activate"
	MAIL_AUTH_ACTIVATE_EMAIL  infrastructure.TplName = "auth/activate_email"
	MAIL_AUTH_RESET_PASSWORD  infrastructure.TplName = "auth/reset_passwd"
	MAIL_AUTH_REGISTER_NOTIFY infrastructure.TplName = "auth/register_notify"
)

var (
	activeCodeTimeLives   int
	resetPwdCodeTimeLives int
	mailRender            MailRender
)

func init() {
	activeCodeTimeLives = 180
	resetPwdCodeTimeLives = 180
}

func InitMail(activeCodeLives, resetPwdCodeLives int, renderOption *macaron.RenderOptions, mailer *Mailer) {
	activeCodeTimeLives = activeCodeLives
	resetPwdCodeTimeLives = resetPwdCodeLives

	ts := macaron.NewTemplateSet()
	ts.Set(macaron.DEFAULT_TPL_SET_NAME, renderOption)

	mailRender = &macaron.TplRender{
		TemplateSet: ts,
		Opt:         renderOption,
	}

	newContext(mailer)
}

type MailRender interface {
	HTMLString(string, interface{}, ...macaron.HTMLOptions) (string, error)
}

func SendTestMail(email string) error {
	return gomail.Send(&Sender{}, NewMessage([]string{email}, "goby Test Email!", "goby Test Email!").Message)
}

type User interface {
	ID() int64
	DisplayName() string
	Email() string
	GenerateActivateCode() string
	GenerateEmailActivateCode(string) string
}

func SendUserMail(c *macaron.Context, u User, tpl infrastructure.TplName, code, subject, info string) {
	data := map[string]interface{}{
		"Username":          u.DisplayName(),
		"ActiveCodeLives":   activeCodeTimeLives / 60,
		"ResetPwdCodeLives": resetPwdCodeTimeLives / 60,
		"Code":              code,
	}
	body, err := mailRender.HTMLString(string(tpl), data)
	if err != nil {
		log.Error(3, "HTMLString: %v", err)
		return
	}

	msg := NewMessage([]string{u.Email()}, subject, body)
	msg.Info = fmt.Sprintf("UID: %d, %s", u.ID(), info)

	SendAsync(msg)
}

func SendActivateAccountMail(c *macaron.Context, u User) {
	SendUserMail(c, u, MAIL_AUTH_ACTIVATE, u.GenerateActivateCode(), "请激活您的帐户", "activate account")
}

func SendResetPasswordMail(c *macaron.Context, u User) {
	SendUserMail(c, u, MAIL_AUTH_RESET_PASSWORD, u.GenerateActivateCode(), "重置您的密码", "reset password")
}

func SendActivateEmailMail(c *macaron.Context, u User, email string) {
	data := map[string]interface{}{
		"Username":        u.DisplayName(),
		"ActiveCodeLives": activeCodeTimeLives / 60,
		"Code":            u.GenerateEmailActivateCode(email),
		"Email":           email,
	}
	body, err := mailRender.HTMLString(string(MAIL_AUTH_ACTIVATE_EMAIL), data)
	if err != nil {
		log.Error(3, "HTMLString: %v", err)
		return
	}

	msg := NewMessage([]string{email}, "请验证您的邮箱地址", body)
	msg.Info = fmt.Sprintf("UID: %d, activate email", u.ID())

	SendAsync(msg)
}

func SendRegisterNotifyMail(c *macaron.Context, u User) {
	data := map[string]interface{}{
		"Username": u.DisplayName(),
	}
	body, err := mailRender.HTMLString(string(MAIL_AUTH_REGISTER_NOTIFY), data)
	if err != nil {
		log.Error(3, "HTMLString: %v", err)
		return
	}

	msg := NewMessage([]string{u.Email()}, "欢迎使用", body)
	msg.Info = fmt.Sprintf("UID: %d, registration notify", u.ID())

	SendAsync(msg)
}
