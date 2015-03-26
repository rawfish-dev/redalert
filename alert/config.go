package alert

import (
	"encoding/json"
	"io/ioutil"
)

type AlertConfig struct {
	Gmail  *GmailConfig  `json:"gmail,omitempty"`
	Slack  *SlackConfig  `json:"slack,omitempty"`
	Twilio *TwilioConfig `json:"twilio,omitempty"`
}

type GmailConfig struct {
	User                  string   `json:"user"`
	Pass                  string   `json:"pass"`
	NotificationAddresses []string `json:"notification_addresses"`
}

type SlackConfig struct {
	WebhookURL string `json:"webhook_url"`
	Channel    string `json:"channel"`
	Username   string `json:"username"`
	IconEmoji  string `json:"icon_emoji"`
}

type TwilioConfig struct {
	AccountSID          string   `json:"account_sid"`
	AuthToken           string   `json:"auth_token"`
	TwilioNumber        string   `json:"twilio_number"`
	NotificationNumbers []string `json:"notification_numbers"`
}

func ReadConfigFile() (*AlertConfig, error) {
	file, err := ioutil.ReadFile("alert/config.json")
	if err != nil {
		return nil, err
	}
	var config AlertConfig
	err = json.Unmarshal(file, &config)
	return &config, err
}
