package alert

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"redalert/common"
)

type Twilio struct {
	accountSid   string
	authToken    string
	phoneNumbers []string
	twilioNumber string
}

func NewTwilio() Twilio {
	return Twilio{
		accountSid:   config.Twilio.AccountSID,
		authToken:    config.Twilio.AuthToken,
		phoneNumbers: config.Twilio.NotificationNumbers,
		twilioNumber: config.Twilio.TwilioNumber,
	}
}

func (a Twilio) Name() string {
	return "Twilio"
}

func (a Twilio) Trigger(alertPackage *AlertPackage) (err error) {

	msg := alertPackage.Message
	for _, num := range a.phoneNumbers {
		err = SendSMS(a.accountSid, a.authToken, num, a.twilioNumber, msg)
		if err != nil {
			return
		}
	}
	alertPackage.AlertLogger.Println(common.White, "Twilio alert successfully triggered.", common.Reset)
	return nil

}

func SendSMS(accountSID string, authToken string, to string, from string, body string) error {

	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSID + "/Messages.json"

	v := url.Values{}
	v.Set("To", to)
	v.Set("From", from)
	v.Set("Body", body)
	rb := *strings.NewReader(v.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &rb)
	req.SetBasicAuth(accountSID, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return errors.New("Invalid Twilio status code")
	}
	return err

}
