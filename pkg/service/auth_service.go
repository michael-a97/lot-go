package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	dto "lot/api/dto/auth"
	"lot/config"
	"lot/pkg/entity"
	"lot/pkg/repository"
	"time"
)

const (
	AccessTokenTTL  = time.Minute * 5
	RefreshTokenTTL = time.Hour * 7 * 24
)

type AuthService interface {
	SignIn(request dto.LoginRequest) (*dto.AuthenticationResponse, error)
	RefreshToken(request dto.TokenRefreshRequest) (*dto.AuthenticationResponse, error)
	VerifyPhoneNumberAuthenticationToken(token string) (bool, error)
}

type authService struct {
	userRepository  repository.UserRepository
	authRepository  repository.AuthRepository
	smsTokenVerfier SmsTokenVerifier
}

func (a authService) SignIn(request dto.LoginRequest) (*dto.AuthenticationResponse, error) {
	user := a.userRepository.FindByPhoneNumber(request.PhoneNumber)
	if user == nil {
		return nil, errors.New("invalid phone number or password")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(request.Password),
	); err != nil {
		return nil, errors.New("invalid phone number or password")
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

	a.authRepository.SaveRefreshToken(hashAndEncodeToHex(refreshToken), user.ID)

	response := dto.AuthenticationResponse{
		ID:           user.ID,
		Role:         user.Role.Name,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
	}

	return &response, nil
}

func (a authService) RefreshToken(request dto.TokenRefreshRequest) (*dto.AuthenticationResponse, error) {
	token, err := jwt.Parse(request.RefreshToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Config("secret")), nil
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

	ok, err := a.authRepository.IsValidRefreshToken(hashAndEncodeToHex(request.RefreshToken), id)

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

	a.authRepository.SaveRefreshToken(hashAndEncodeToHex(refreshToken), user.ID)

	response := dto.AuthenticationResponse{
		ID:           user.ID,
		Role:         user.Role.Name,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
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

func hashAndEncodeToHex(token string) string {
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
	claims["issued_at"] = time.Now()
	claims["expires_at"] = expiryTime
	signedToken, err := signer.SignedString([]byte(config.Config("secret")))
	return signedToken, err
}
