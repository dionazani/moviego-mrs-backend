package contextsignup

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/dionazani/moviego-mrs-backend/internal/infrastructure/database"
	"github.com/dionazani/moviego-mrs-backend/internal/infrastructure/dto"
	"github.com/dionazani/moviego-mrs-backend/internal/infrastructure/model"
	"github.com/dionazani/moviego-mrs-backend/internal/infrastructure/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SignUpService defines the business logic contract for sign-up operations.
type SignUpService interface {
	AddNew(ctx context.Context, dto SignUpDTO) (infrastructuredto.Response, error)
	LoadById(ctx context.Context, id uuid.UUID) (infrastructuredto.Response, error)
}

type signUpServiceImpl struct {
	db                     *gorm.DB
	appPersonRepository    infrastructurerepository.AppPersonRepository
	appUserRepository      infrastructurerepository.AppUserRepository
	masterUserRoleRegular  string
}

// NewSignUpService creates a new instance of SignUpService.
func NewSignUpService(
	db *gorm.DB,
	appPersonRepository infrastructurerepository.AppPersonRepository,
	appUserRepository infrastructurerepository.AppUserRepository,
	masterUserRoleRegular string,
) SignUpService {
	return &signUpServiceImpl{
		db:                     db,
		appPersonRepository:    appPersonRepository,
		appUserRepository:      appUserRepository,
		masterUserRoleRegular:  masterUserRoleRegular,
	}
}

// AddNew maps the SignUpDTO to AppPerson and AppUser models and calls the repository's Insert functions inside a transaction.
func (s *signUpServiceImpl) AddNew(ctx context.Context, dto SignUpDTO) (infrastructuredto.Response, error) {
	// 1. Validate Password and Confirmation matching
	if dto.Password != dto.PasswordConfirmation {
		return infrastructuredto.Response{
			Timestamp:       time.Now().Format(time.RFC3339),
			ResponseStatus:  http.StatusBadRequest,
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "Password and password confirmation do not match",
			Data:            nil,
		}, nil
	}

	signUpFrom := "WEB" // Defaulting to web sign-up
	loc := time.FixedZone("GMT+7", 7*60*60)
	now := time.Now().In(loc)

	// 2. Prepare AppPerson model
	person := &infrastructuremodel.AppPerson{
		ID:          uuid.New(),
		Fullname:    dto.Fullname,
		Gender:      dto.Gender,
		Email:       dto.Email,
		MobilePhone: dto.MobilePhone,
		SignUpFrom:  &signUpFrom,
		SignUpAt:    &now,
	}

	// 3. Hash password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return infrastructuredto.Response{}, err
	}

	// 4. Resolve default regular user role UUID from config/env
	defaultRoleID, err := uuid.Parse(s.masterUserRoleRegular)
	if err != nil {
		return infrastructuredto.Response{}, err
	}

	nextChangePasswordDate := now.AddDate(0, 0, 30)

	// 5. Prepare AppUser model
	user := &infrastructuremodel.AppUser{
		ID:                     uuid.New(),
		AppPersonID:            person.ID,
		MstUserRoleID:          defaultRoleID,
		AppPassword:            string(hashedPassword),
		MustChangePassword:     0,
		NextChangePasswordDate: nextChangePasswordDate,
		IsLocked:               0,
		CreatedAt:              now,
	}

	// 6. Execute inserts within GORM transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		ctxTx := infrastructuredatabase.WithTransaction(ctx, tx)

		if err := s.appPersonRepository.Insert(ctxTx, person); err != nil {
			return err
		}

		if err := s.appUserRepository.Insert(ctxTx, user); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return infrastructuredto.Response{}, err
	}

	return infrastructuredto.Response{
		Timestamp:       time.Now().Format(time.RFC3339),
		ResponseStatus:  http.StatusCreated,
		ResponseCode:    http.StatusCreated,
		ResponseMessage: "User signed up successfully",
		Data: map[string]interface{}{
			"appUserId": person.ID,
		},
	}, nil
}

// LoadById retrieves a single user's detail using repository's FindById.
func (s *signUpServiceImpl) LoadById(ctx context.Context, id uuid.UUID) (infrastructuredto.Response, error) {
	person, err := s.appPersonRepository.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return infrastructuredto.Response{
				Timestamp:       time.Now().Format(time.RFC3339),
				ResponseStatus:  http.StatusNotFound,
				ResponseCode:    http.StatusNotFound,
				ResponseMessage: "User not found",
				Data:            nil,
			}, nil
		}
		return infrastructuredto.Response{}, err
	}

	return infrastructuredto.Response{
		Timestamp:       time.Now().Format(time.RFC3339),
		ResponseStatus:  http.StatusOK,
		ResponseCode:    http.StatusOK,
		ResponseMessage: "User retrieved successfully",
		Data:            person,
	}, nil
}

