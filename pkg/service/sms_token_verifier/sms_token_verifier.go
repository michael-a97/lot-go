package service

import (
	"context"

	"firebase.google.com/go/v4/auth"
)

type SmsTokenVerifier interface {
	IsTokenValid(token string) (bool, error)
}

type firebaseSmsTokenVerifier struct {
	firebaseAuthClient *auth.Client
}

func (s firebaseSmsTokenVerifier) IsTokenValid(token string) (bool, error) {
	verifiedToken, err :=
		s.firebaseAuthClient.VerifyIDTokenAndCheckRevoked(
			context.Background(),
			token,
		)

	if err != nil {
		return false, err
	}

	return verifiedToken != nil, nil
}

func NewFirebaseSmsTokenVerifier(authClient *auth.Client) SmsTokenVerifier {
	return firebaseSmsTokenVerifier{
		firebaseAuthClient: authClient,
	}
}
