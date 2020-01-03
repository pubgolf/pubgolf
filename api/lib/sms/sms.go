package sms

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// twilioResponse is a struct representing the identifiers we get back from Twilio on a successful messaging call, or
// the error codes we get back on a failed one.
type twilioResponse struct {
	// Indicates whether an API call was actually made, or just simulated (depends on the presence of env vars).
	LiveCall bool `json:"live_call"`

	// Entity IDs on successful response.
	Sid        string `json:"sid,omitempty"`
	AccountSid string `json:"account_sid,omitempty"`

	// Nonzero if empty.
	Code int32 `json:"code"`

	// Detailed error info.
	Message string `json:"message,omitempty"`
}

// SendAuthCodeSms handles populating the auth message template with the provided `authCode` and sends it to
// `phoneNumber` via SMS. The `phoneNumber` argument is expected to be in valid E.164 format:
// https://www.twilio.com/docs/glossary/what-e164
func SendAuthCodeSms(logCtx *log.Entry, phoneNumber *string, authCode uint32) (err error) {
	smsContent := fmt.Sprintf("Your pubgolf.co auth code is: %d", authCode)
	responseInfo := twilioResponse{}

	if os.Getenv("TWILIO_FROM_NUM") != "" && os.Getenv("TWILIO_ACCOUNT_SID") != "" && os.Getenv("TWILIO_AUTH_TOKEN") !=
		"" {
		// API keys are present, so make live calls.
		msgReader := formatMessage(phoneNumber, smsContent)
		responseInfo, err = sendRequest(os.Getenv("TWILIO_ACCOUNT_SID"), os.Getenv("TWILIO_AUTH_TOKEN"), msgReader)
		responseInfo.LiveCall = true

		if err == nil {
			logCtx.WithField("twilio_api_response", responseInfo).Info("Twilio Success")
		} else {
			logCtx.WithField("twilio_api_response", responseInfo).Error("Twilio Failure")
		}
	} else {
		// Simulated call, so output debugging log.
		if strings.HasSuffix(os.Getenv("PUBGOLF_ENV"), "dev") {
			log.Debugf("\nSending text message to %s:\n==========\n%s\n==========", *phoneNumber, smsContent)
		}

		logCtx.WithField("twilio_api_response", responseInfo).Info("Twilio Simulation")
	}

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

func sendRequest(accountSid string, authToken string, msgDataReader strings.Reader) (parsedBody twilioResponse,
	err error) {
	client := &http.Client{}
	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)
	req, err := http.NewRequest("POST", urlStr, &msgDataReader)
	if err != nil {
		return parsedBody, err
	}

	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return parsedBody, err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&parsedBody)
	if err != nil {
		return parsedBody, fmt.Errorf("twilio server responded with status %d, could not read body: %s", resp.StatusCode,
			err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parsedBody, fmt.Errorf("twilio server responded with status %d and error message: \"%s\" (code: %d)",
			resp.StatusCode, parsedBody.Message, parsedBody.Code)
	}

	return parsedBody, nil
}
