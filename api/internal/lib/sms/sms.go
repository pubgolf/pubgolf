// Package sms contains logic for sending communications via SMS.
package sms

import (
	"fmt"
	"log"
	"regexp"

	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"

	"github.com/pubgolf/pubgolf/api/internal/lib/config"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

var AuthCodeFormat = regexp.MustCompile(`\d{8}`)

var MockAuthCode = "12345678"

type Client struct {
	tc              *twilio.RestClient
	VerificationSID string
	allowed         config.PhoneNumSet
}

func New(c config.TwilioAuth, allowed config.PhoneNumSet) *Client {
	return &Client{
		tc: twilio.NewRestClientWithParams(twilio.ClientParams{
			Username: c.AccountSID,
			Password: c.AuthToken,
		}),
		VerificationSID: c.VerificationSID,
		allowed:         allowed,
	}
}

func (c *Client) SendVerification(to models.PhoneNum) error {
	if !c.allowed.Match(to) {
		log.Printf("Attempted to send SMS to phone number %q, which is not in the allowlist. Use auth code %q for testing purposes.", to, MockAuthCode)

		return nil
	}

	p := verify.CreateVerificationParams{}
	p.SetTo(to.String())
	p.SetChannel("sms")

	_, err := c.tc.VerifyV2.CreateVerification(c.VerificationSID, &p)
	if err != nil {
		return fmt.Errorf("call Twilio CreateVerification endpoint: %w", err)
	}

	return nil
}

func (c *Client) CheckVerification(to models.PhoneNum, code string) (bool, error) {
	if !c.allowed.Match(to) {
		log.Printf("Verifying auth code for phone number %q, which is not in the allowlist. Use auth code %q for testing purposes.", to, MockAuthCode)

		return code == MockAuthCode, nil
	}

	p := verify.CreateVerificationCheckParams{}
	p.SetTo(to.String())
	p.SetCode(code)

	resp, err := c.tc.VerifyV2.CreateVerificationCheck(c.VerificationSID, &p)
	if err != nil {
		return false, fmt.Errorf("call Twilio CreateVerificationCheck endpoint: %w", err)
	}

	return (resp != nil && *resp.Status == "approved"), nil
}
