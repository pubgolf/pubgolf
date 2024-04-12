// Package sms contains logic for sending communications via SMS.
package sms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/pubgolf/pubgolf/api/internal/lib/config"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
)

var (
	// AuthCodeFormat allows pre-validating auth codes.
	AuthCodeFormat = regexp.MustCompile(`\d{8}`)

	// MockAuthCode is the default auth code used when a phone number is not on the live request allowlist for a given environment.
	MockAuthCode = "12345678"

	// ErrUpstreamProviderIssue indicates an error response from the SMS provider.
	ErrUpstreamProviderIssue = errors.New("error response from SMS provider")
)

// Client allows making requests to the SMS verification provider's API.
type Client struct {
	hc      *http.Client
	Auth    config.TwilioAuth
	allowed config.PhoneNumSet
}

// New takes Twilio credentials and returns a client which is only allowed to make requests on behalf of the specified set of phone numbers.
func New(auth config.TwilioAuth, allowed config.PhoneNumSet) *Client {
	return &Client{
		hc: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
			Timeout:   10 * time.Second,
		},
		Auth:    auth,
		allowed: allowed,
	}
}

func (c *Client) doVerificationServicePost(ctx context.Context, endpoint string, data url.Values) (*http.Response, error) {
	path := fmt.Sprintf("https://verify.twilio.com/v2/Services/%s/%s", c.Auth.VerificationSID, endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.Auth.AccountSID, c.Auth.AuthToken)

	res, err := c.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("make HTTP request: %w", err)
	}

	return res, nil
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

	data := url.Values{}

	data.Add("To", to.String())
	data.Add("Channel", "sms")

	res, err := c.doVerificationServicePost(ctx, "Verifications", data)
	if err != nil {
		return fmt.Errorf("call Twilio CreateVerification endpoint: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode > 399 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("read Twilio API call response body: %w", err)
		}

		return fmt.Errorf("response from Twilio CreateVerification call %q: %w", body, ErrUpstreamProviderIssue)
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

	data := url.Values{}

	data.Add("To", to.String())
	data.Add("Code", code)

	res, err := c.doVerificationServicePost(ctx, "VerificationCheck", data)
	if err != nil {
		return false, fmt.Errorf("call Twilio VerificationCheck endpoint: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return false, nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("read Twilio API call response body: %w", err)
	}

	if res.StatusCode > 399 {
		return false, fmt.Errorf("response from Twilio VerificationCheck call %q: %w", body, ErrUpstreamProviderIssue)
	}

	r := struct {
		Status string `json:"status"`
	}{}

	err = json.Unmarshal(body, &r)
	if err != nil {
		return false, fmt.Errorf("unmarshal JSON response: %w", err)
	}

	return r.Status == "approved", nil
}
