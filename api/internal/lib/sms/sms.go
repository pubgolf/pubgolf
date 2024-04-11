// Package sms contains logic for sending communications via SMS.
package sms

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/pubgolf/pubgolf/api/internal/lib/config"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
)

// AuthCodeFormat allows pre-validating auth codes.
var AuthCodeFormat = regexp.MustCompile(`\d{8}`)

// MockAuthCode is the default auth code used when a phone number is not on the live request allowlist for a given environment.
var MockAuthCode = "12345678"

// Client allows making requests to the SMS verification provider's API.
type Client struct {
	tc              *twilio.RestClient
	VerificationSID string
	allowed         config.PhoneNumSet
}

// New takes Twilio credentials and returns a client which is only allowed to make requests on behalf of the specified set of phone numbers.
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

// SendVerification sends an SMS verification message to the given number.
func (c *Client) SendVerification(ctx context.Context, to models.PhoneNum) error {
	defer telemetry.FnSpan(&ctx)()
	span := trace.SpanFromContext(ctx)

	if !c.allowed.Match(to) {
		log.Printf("Attempted to send SMS to phone number %q, which is not in the allowlist. Use auth code %q for testing purposes.", to, MockAuthCode)
		span.SetAttributes(attribute.Bool("sms.send_live_request", false))

		return nil
	}

	span.SetAttributes(attribute.Bool("sms.send_live_request", true))

	p := verify.CreateVerificationParams{}
	p.SetTo(to.String())
	p.SetChannel("sms")

	_, err := c.tc.VerifyV2.CreateVerification(c.VerificationSID, &p)
	if err != nil {
		return fmt.Errorf("call Twilio CreateVerification endpoint: %w", err)
	}

	return nil
}

// CheckVerification validates the auth code matches the last verification sent to the given number.
func (c *Client) CheckVerification(ctx context.Context, to models.PhoneNum, code string) (bool, error) {
	defer telemetry.FnSpan(&ctx)()
	span := trace.SpanFromContext(ctx)

	if !c.allowed.Match(to) {
		log.Printf("Verifying auth code for phone number %q, which is not in the allowlist. Use auth code %q for testing purposes.", to, MockAuthCode)

		span.SetAttributes(attribute.Bool("sms.send_live_request", false))

		return code == MockAuthCode, nil
	}

	span.SetAttributes(attribute.Bool("sms.send_live_request", true))

	p := verify.CreateVerificationCheckParams{}
	p.SetTo(to.String())
	p.SetCode(code)

	resp, err := c.tc.VerifyV2.CreateVerificationCheck(c.VerificationSID, &p)
	if err != nil {
		return false, fmt.Errorf("call Twilio CreateVerificationCheck endpoint: %w", err)
	}

	return (resp != nil && *resp.Status == "approved"), nil
}
