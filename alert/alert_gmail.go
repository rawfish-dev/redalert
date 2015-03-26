package alert

import (
	"net/smtp"
	"strings"

	"redalert/common"
)

type Gmail struct {
	user                  string
	pass                  string
	notificationAddresses []string
}

func NewGmail() Gmail {
	return Gmail{
		user: config.Gmail.User,
		pass: config.Gmail.Pass,
		notificationAddresses: config.Gmail.NotificationAddresses,
	}
}

func (a Gmail) Name() string {
	return "Gmail"
}

func (a Gmail) Trigger(alertPackage *AlertPackage) error {

	body := "To: " + strings.Join(a.notificationAddresses, ",") +
		"\r\nSubject: " + alertPackage.Message +
		"\r\n\r\n" + alertPackage.Message

	auth := smtp.PlainAuth("", a.user, a.pass, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, a.user,
		a.notificationAddresses, []byte(body))
	if err != nil {
		return err
	}

	alertPackage.AlertLogger.Println(common.White, "Gmail alert successfully triggered.", common.Reset)
	return nil
}
