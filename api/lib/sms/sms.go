package sms

import "log"

func SendAuthCodeSms(phoneNumber *string, authCode uint32) error {
	log.Printf("Sending text message to %s:\n==========\n%06d\n==========",
		*phoneNumber, authCode)
	return nil
}
