package server

import (
	"crypto/rand"
	"math/big"
	"regexp"
)

func generateAuthCode() (uint32, error) {
	randNum, err := rand.Int(rand.Reader, big.NewInt(899999))
	if err != nil {
		return 0, err
	}
	randCode := randNum.Int64() + 100000
	return uint32(randCode), nil
}

func isEmpty(arg *string) bool {
	return *arg == ""
}

func invalidPhoneNumberFormat(phoneNumber *string) bool {
	isValid, err := regexp.MatchString("^\\+[1-9]\\d{1,14}$", *phoneNumber)
	return !isValid || err != nil
}

func invalidAuthCodeFormat(authCode uint32) bool {
	return authCode < 100000 || authCode > 999999
}
