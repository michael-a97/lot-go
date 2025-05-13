package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"lot/api/dto"
	"lot/config"
	"lot/internal/entity"
	app_errors "lot/internal/errors"
	"lot/internal/repository"
	"lot/internal/utilities"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTokenTTL  = time.Minute * 5
	RefreshTokenTTL = time.Hour * 7 * 24
)

type AuthService interface {
	SignIn(request dto.LoginRequest) (*dto.AuthenticationResponse, error)
	RefreshToken(request dto.TokenRefreshRequest) (*dto.AuthenticationResponse, error)
	VerifyPhoneNumberAuthenticationToken(token string) (bool, error)
	ResetPassword(request dto.PasswordResetRequest) error
	ChangePassword(request dto.ChangePasswordRequest, user entity.User) error
	GetUserFromAccessToken(accesToken string) (*entity.User, error)
}

type authService struct {
	userRepository  repository.UserRepository
	authRepository  repository.AuthRepository
	smsTokenVerfier SmsTokenVerifier
}

func (a authService) SignIn(request dto.LoginRequest) (*dto.AuthenticationResponse, error) {
	user, err := a.userRepository.FindByPhoneNumber(request.PhoneNumber)
	if errors.Is(err, app_errors.ErrAccountNotFound) {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(request.Password),
	); err != nil {
		return nil, app_errors.ErrInvalidPhoneNumberOrPassword
	}

	accessToken, err := GetSignedToken(*user, time.Now().Add(AccessTokenTTL))
	if err != nil {
		return nil, err
	}

	refreshToken, err := GetSignedToken(*user, time.Now().Add(RefreshTokenTTL))
	if err != nil {
		return nil, err
	}

	if err := a.authRepository.RevokeAllRefreshTokensForUser(user.ID); err != nil {
		return nil, err
	}

	if err := a.authRepository.SaveRefreshToken(HashAndEncodeToHex(refreshToken), user.ID); err != nil {
		return nil, err
	}

	response := dto.AuthenticationResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         toDto(*user),
	}

	return &response, nil
}

func (a authService) RefreshToken(request dto.TokenRefreshRequest) (*dto.AuthenticationResponse, error) {
	token, err := jwt.Parse(request.RefreshToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		secret, err := config.Config("secret")
		if err != nil {
			return nil, err
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	id := uint(claims["id"].(float64))

	user, err := a.userRepository.FindById(id)

	if user == nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	if err != nil {
		return nil, err
	}

	ok, err := a.authRepository.IsValidRefreshToken(HashAndEncodeToHex(request.RefreshToken), id)

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf("invalid refresh token")
	}

	accessToken, err := GetSignedToken(*user, time.Now().Add(AccessTokenTTL))
	if err != nil {
		return nil, err
	}

	refreshToken, err := GetSignedToken(*user, time.Now().Add(RefreshTokenTTL))
	if err != nil {
		return nil, err
	}
	if err := a.authRepository.RevokeAllRefreshTokensForUser(user.ID); err != nil {
		return nil, err
	}

	if err := a.authRepository.SaveRefreshToken(HashAndEncodeToHex(refreshToken), user.ID); err != nil {
		return nil, err
	}

	response := dto.AuthenticationResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         toDto(*user),
	}
	return &response, nil
}

func (a authService) VerifyPhoneNumberAuthenticationToken(token string) (bool, error) {
	valid, err := a.smsTokenVerfier.IsTokenValid(token)
	if err != nil {
		return false, err
	}
	return valid, nil
}

func (a authService) ResetPassword(request dto.PasswordResetRequest) error {
	user, err := a.userRepository.FindByPhoneNumber(request.PhoneNumber)
	if err != nil {
		return err
	}

	hashedPassword, err := utilities.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	_, err = a.userRepository.Save(*user)

	if err != nil {
		return err
	}

	return nil

}

func (a authService) ChangePassword(request dto.ChangePasswordRequest, user entity.User) error {

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(request.OldPassword),
	); err != nil {
		return app_errors.ErrInvalidPassword
	}

	hashedPassword, err := utilities.HashPassword(request.NewPassword)

	if err != nil {
		return err
	}

	user.Password = hashedPassword
	_, err = a.userRepository.Update(user)

	return err
}

func HashAndEncodeToHex(token string) string {
	hasher := sha256.New()
	hasher.Write([]byte(token))
	return hex.EncodeToString(hasher.Sum(nil))
}

func NewAuthService(
	authRepository repository.AuthRepository,
	userRepository repository.UserRepository,
	smsTokenVerifier SmsTokenVerifier,
) AuthService {
	return &authService{
		userRepository:  userRepository,
		authRepository:  authRepository,
		smsTokenVerfier: smsTokenVerifier,
	}
}

func GetSignedToken(user entity.User, expiryTime time.Time) (string, error) {
	signer := jwt.New(jwt.SigningMethodHS256)
	claims := signer.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["role"] = user.Role.Name
	claims["issued_at"] = time.Now().Unix()
	claims["exp"] = expiryTime.Unix()
	secret, err := config.Config("secret")
	if err != nil {
		return "", err
	}
	signedToken, err := signer.SignedString([]byte(secret))
	return signedToken, err
}

func (a authService) GetUserFromAccessToken(accesToken string) (*entity.User, error) {
	token, err := jwt.Parse(accesToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		secret, err := config.Config("secret")
		if err != nil {
			return nil, err
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	id := uint(claims["id"].(float64))

	return a.userRepository.FindById(id)

}
