package sms

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// SendAuthCodeSms handles populating the auth message template with the provided `authCode` and sends it to
// `phoneNumber` via SMS. The `phoneNumber` argument is expected to be in valid E.164 format:
// https://www.twilio.com/docs/glossary/what-e164
func SendAuthCodeSms(phoneNumber *string, authCode uint32) error {
	smsContent := fmt.Sprintf("Your pubgolf.co auth code is: %d", authCode)
	logAction := "Logging"
	var err error = nil

	if os.Getenv("TWILIO_FROM_NUM") != "" && os.Getenv("TWILIO_ACCOUNT_SID") != "" && os.Getenv("TWILIO_AUTH_TOKEN") !=
		"" {
		logAction = "Sending"
		msgReader := formatMessage(phoneNumber, smsContent)
		err = sendRequest(os.Getenv("TWILIO_ACCOUNT_SID"), os.Getenv("TWILIO_AUTH_TOKEN"), msgReader)
	}

	log.Printf("%s text message to %s:\n==========\n%s\n==========", logAction, *phoneNumber, smsContent)

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
	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)
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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("twilio server responded with status %d, could not read body: %s", resp.StatusCode, err)
		}

		return fmt.Errorf("twilio server responded with status %d and body %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
