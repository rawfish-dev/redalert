package alert

import (
	"fmt"
	"log"
	"os"
)

type AlertType string

const (
	STDERR AlertType = "standardError"
	SLACK  AlertType = "slack"
	GMAIL  AlertType = "gmail"
	TWILIO AlertType = "twilio"
)

type AlertPackage struct {
	Message     string
	AlertLogger *log.Logger
}

type Alert interface {
	Trigger(*AlertPackage) error
	Name() string
}

var config *AlertConfig
var alertLogger *log.Logger
var RegisteredAlerts map[AlertType]Alert

func init() {
	var err error
	config, err = ReadConfigFile()
	if err != nil {
		panic(fmt.Sprintf("Missing or invalid config %v", err))
	}

	alertLogger = log.New(os.Stdout, "Alert ", log.Ldate|log.Ltime)

	RegisterConfiguredAlerts()
}

func RegisterConfiguredAlerts() {
	RegisteredAlerts = make(map[AlertType]Alert)

	RegisteredAlerts[STDERR] = NewStandardError()

	if config.Slack == nil || config.Slack.WebhookURL == "" {
		alertLogger.Println("Slack is not configured")
	} else {
		RegisteredAlerts[SLACK] = NewSlackWebhook()
		alertLogger.Println("Slack alert registered")
	}

	if config.Gmail == nil || config.Gmail.User == "" || config.Gmail.Pass == "" || len(config.Gmail.NotificationAddresses) == 0 {
		alertLogger.Println("Gmail is not configured")
	} else {
		RegisteredAlerts[GMAIL] = NewGmail()
		alertLogger.Println("Gmail alert registered")
	}

	if config.Twilio == nil || config.Twilio.AccountSID == "" || config.Twilio.AuthToken == "" || len(config.Twilio.NotificationNumbers) == 0 || config.Twilio.TwilioNumber == "" {
		alertLogger.Println("Twilio is not configured")
	} else {
		RegisteredAlerts[TWILIO] = NewTwilio()
		alertLogger.Println("Twilio alert registered")
	}
}

func Debug(alertPackage *AlertPackage) {
	standardError := RegisteredAlerts[STDERR]
	standardError.Trigger(alertPackage)
}

func Warn(alertPackage *AlertPackage) {
	if slack, exists := RegisteredAlerts[SLACK]; exists {
		Debug(alertPackage)
		slack.Trigger(alertPackage)
	}
}

func Error(alertPackage *AlertPackage) {
	if gmail, exists := RegisteredAlerts[GMAIL]; exists {
		Warn(alertPackage)
		gmail.Trigger(alertPackage)
	}
}

func Critical(alertPackage *AlertPackage) {
	if twilio, exists := RegisteredAlerts[TWILIO]; exists {
		Error(alertPackage)
		twilio.Trigger(alertPackage)
	}
}
