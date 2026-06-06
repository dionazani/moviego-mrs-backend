package signin

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	infrastructuredto "github.com/dionazani/moviego-mrs-backend/internal/infrastructure/dto"
	infrastructuremodel "github.com/dionazani/moviego-mrs-backend/internal/infrastructure/model"
	infrastructurerepository "github.com/dionazani/moviego-mrs-backend/internal/infrastructure/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SignInService defines the contract for handling user sign-in.
type SignInService interface {
	SignIn(ctx context.Context, req SignInRequest) (infrastructuredto.Response, error)
}

// maxFailedAttempts is the number of consecutive failures before a temporary lockout is applied.
const maxFailedAttempts = 3

type signInServiceImpl struct {
	appPersonRepository    infrastructurerepository.AppPersonRepository
	appUserRepository      infrastructurerepository.AppUserRepository
	appUserTokenRepository infrastructurerepository.AppUserTokenRepository
	jwtSecret              string
	lockoutDuration        time.Duration
}

// NewSignInService creates a new instance of SignInService.
func NewSignInService(
	appPersonRepository infrastructurerepository.AppPersonRepository,
	appUserRepository infrastructurerepository.AppUserRepository,
	appUserTokenRepository infrastructurerepository.AppUserTokenRepository,
	jwtSecret string,
	lockoutDuration time.Duration,
) SignInService {
	return &signInServiceImpl{
		appPersonRepository:    appPersonRepository,
		appUserRepository:      appUserRepository,
		appUserTokenRepository: appUserTokenRepository,
		jwtSecret:              jwtSecret,
		lockoutDuration:        lockoutDuration,
	}
}

// SignIn authenticates the user and returns JWT access and refresh tokens.
func (s *signInServiceImpl) SignIn(ctx context.Context, req SignInRequest) (infrastructuredto.Response, error) {
	// 1. Find person by email. If not found, fallback to mobile phone.
	person, err := s.appPersonRepository.FindByEmail(ctx, req.Username)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return infrastructuredto.Response{}, err
		}
		person, err = s.appPersonRepository.FindByMobilePhone(ctx, req.Username)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return s.unauthorizedResponse(), nil
			}
			return infrastructuredto.Response{}, err
		}
	}

	// 2. Fetch the AppUser record (contains the hashed password).
	appUserModel, err := s.appUserRepository.FindByAppPersonId(ctx, person.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.unauthorizedResponse(), nil
		}
		return infrastructuredto.Response{}, err
	}

	// 2.1 Check if the user account is permanently locked.
	if appUserModel.IsLocked == 1 {
		return s.lockedResponse(), nil
	}

	// 2.2 Evaluate the current lockout state.
	//   - lockoutActive  : lockout_until is set and the lockout period is still ongoing.
	//   - lockoutExpired : lockout_until is set but the lockout period has already passed.
	now := time.Now()
	lockoutActive := appUserModel.LockoutUntil != nil && now.Before(*appUserModel.LockoutUntil)
	lockoutExpired := appUserModel.LockoutUntil != nil && !now.Before(*appUserModel.LockoutUntil)

	if lockoutActive {
		// Account is still blocked — reject immediately without checking the password.
		remainingSeconds := int(appUserModel.LockoutUntil.Sub(now).Seconds())
		return s.tooManyAttemptsResponse(remainingSeconds), nil
	}

	// 3. Verify the provided password against the stored hash.
	if err := bcrypt.CompareHashAndPassword([]byte(appUserModel.AppPassword), []byte(req.Password)); err != nil {
		if lockoutExpired {
			// 3.1 The lockout period has expired but the password is still wrong.
			// Reset the counter to 1 (fresh start) and clear the expired lockout_until.
			appUserModel.FailedAttemptCount = 1
			appUserModel.LockoutUntil = nil
		} else {
			// 3.2 Normal failed attempt — increment the counter.
			appUserModel.FailedAttemptCount++

			// 3.3 If max failed attempts reached, set lockout_until = now + USER_LOGIN_LOCKOUT_DURATION.
			if appUserModel.FailedAttemptCount >= maxFailedAttempts {
				lockoutUntil := now.Add(s.lockoutDuration)
				appUserModel.LockoutUntil = &lockoutUntil
			}
		}

		// 3.4 Persist the updated failed_attempt_count and lockout_until to the database.
		if updateErr := s.appUserRepository.Update(ctx, appUserModel); updateErr != nil {
			return infrastructuredto.Response{}, updateErr
		}

		return s.unauthorizedResponse(), nil
	}

	// Reset failed attempts and lockout upon successful sign-in if they were set.
	if appUserModel.FailedAttemptCount > 0 || appUserModel.LockoutUntil != nil {
		appUserModel.FailedAttemptCount = 0
		appUserModel.LockoutUntil = nil
		if err := s.appUserRepository.Update(ctx, appUserModel); err != nil {
			return infrastructuredto.Response{}, err
		}
	}

	// 4. Generate Access Token (JWT, expires in 60 minutes).
	accessExpiry := time.Now().Add(60 * time.Minute)
	accessClaims := jwt.MapClaims{
		"sub": appUserModel.ID.String(),
		"exp": accessExpiry.Unix(),
		"iat": time.Now().Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return infrastructuredto.Response{}, err
	}

	// 5. Generate Refresh Token (JWT, expires in 7 days).
	refreshExpiry := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := jwt.MapClaims{
		"sub": appUserModel.ID.String(),
		"exp": refreshExpiry.Unix(),
		"iat": time.Now().Unix(),
		"typ": "refresh",
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return infrastructuredto.Response{}, err
	}

	// 6. Upsert Refresh Token to the database.
	// Uses Upsert so that re-login simply updates the existing token instead of failing.
	appUserTokenModel := &infrastructuremodel.AppUserToken{
		AppUserID: appUserModel.ID,
		TokenType: "refresh",
		TokenUser: refreshTokenString,
		ExpireAt:  refreshExpiry,
	}
	if err := s.appUserTokenRepository.Upsert(ctx, appUserTokenModel); err != nil {
		return infrastructuredto.Response{}, err
	}

	// 7. Return successful response.
	return infrastructuredto.Response{
		Timestamp:       time.Now().Format(time.RFC3339),
		ResponseStatus:  http.StatusOK,
		ResponseCode:    http.StatusOK,
		ResponseMessage: "User signed in successfully",
		Data: TokenData{
			Token:        accessTokenString,
			RefreshToken: refreshTokenString,
		},
	}, nil
}

// unauthorizedResponse returns a standardized 401 Unauthorized response.
func (s *signInServiceImpl) unauthorizedResponse() infrastructuredto.Response {
	return infrastructuredto.Response{
		Timestamp:       time.Now().Format(time.RFC3339),
		ResponseStatus:  http.StatusUnauthorized,
		ResponseCode:    http.StatusUnauthorized,
		ResponseMessage: "Authentication credentials were not provided",
		Data:            nil,
	}
}

// lockedResponse returns a standardized 403 Forbidden response when the account is permanently locked.
func (s *signInServiceImpl) lockedResponse() infrastructuredto.Response {
	return infrastructuredto.Response{
		Timestamp:       time.Now().Format(time.RFC3339),
		ResponseStatus:  http.StatusForbidden,
		ResponseCode:    http.StatusForbidden,
		ResponseMessage: "User account is locked",
		Data:            nil,
	}
}

// tooManyAttemptsResponse returns a standardized 429 Too Many Requests response
// when the account is temporarily locked out due to too many failed login attempts.
func (s *signInServiceImpl) tooManyAttemptsResponse(remainingSeconds int) infrastructuredto.Response {
	return infrastructuredto.Response{
		Timestamp:       time.Now().Format(time.RFC3339),
		ResponseStatus:  http.StatusTooManyRequests,
		ResponseCode:    http.StatusTooManyRequests,
		ResponseMessage: fmt.Sprintf("Too many failed login attempts. Please try again in %d seconds.", remainingSeconds),
		Data:            nil,
	}
}
