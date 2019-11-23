package sms

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func SendAuthCodeSms(phoneNumber *string, authCode uint32) error {
	smsContent := fmt.Sprintf("Your pubgolf.co auth code is: %d", authCode)
	logAction := "Logging"
	var err error = nil

	if os.Getenv("TWILIO_FROM_NUM") != "" && os.Getenv("TWILIO_ACCOUNT_SID") !=
		"" && os.Getenv("TWILIO_AUTH_TOKEN") != "" {
		logAction = "Sending"
		msgReader := formatMessage(phoneNumber, smsContent)
		err = sendRequest(os.Getenv("TWILIO_ACCOUNT_SID"),
			os.Getenv("TWILIO_AUTH_TOKEN"), msgReader)
	}

	log.Printf("%s text message to %s:\n==========\n%s\n==========",
		logAction, *phoneNumber, smsContent)

	return err
}

func formatMessage(phoneNumber *string, smsContent string) strings.Reader {
	msgData := url.Values{}
	msgData.Set("To", *phoneNumber)
	msgData.Set("From", os.Getenv("TWILIO_FROM_NUM"))
	msgData.Set("Body", smsContent)
	msgDataReader := *strings.NewReader(msgData.Encode())
	return msgDataReader
}

func sendRequest(accountSid string, authToken string,
	msgDataReader strings.Reader) error {
	client := &http.Client{}
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid +
		"/Messages.json"
	req, err := http.NewRequest("POST", urlStr, &msgDataReader)
	if err != nil {
		return err
	}

	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		return err
	} else {
		return fmt.Errorf("twilio server responded with status %d", resp.StatusCode)
	}
	return nil
}
