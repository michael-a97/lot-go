package app_errors

import "errors"

var ErrRecordNotFound = errors.New("record not found")
var ErrAccountNotFound = errors.New("account not found")
var ErrInvalidPhoneNumberOrPassword = errors.New("invalid phone number or password")
var ErrInvalidPhoneNumberVerificationToken = errors.New("invalid phone number verification token")