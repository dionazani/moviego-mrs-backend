package contextsignup

import (
	"context"
	"time"

	"github.com/dionazani/moviego-mrs-backend/internal/infrastructure/model"
	"github.com/dionazani/moviego-mrs-backend/internal/infrastructure/repository"
	"github.com/google/uuid"
)

// SignUpService defines the business logic contract for sign-up operations.
type SignUpService interface {
	AddNew(ctx context.Context, dto SignUpDTO) (*model.AppPerson, error)
	LoadById(ctx context.Context, id uuid.UUID) (*model.AppPerson, error)
}

type signUpServiceImpl struct {
	appPersonRepository repository.AppPersonRepository
}

// NewSignUpService creates a new instance of SignUpService.
func NewSignUpService(appPersonRepository repository.AppPersonRepository) SignUpService {
	return &signUpServiceImpl{
		appPersonRepository: appPersonRepository,
	}
}

// AddNew maps the SignUpDTO to AppPerson model and calls the repository's Insert function.
func (s *signUpServiceImpl) AddNew(ctx context.Context, dto SignUpDTO) (*model.AppPerson, error) {
	signUpFrom := "WEB" // Defaulting to web sign-up
	loc := time.FixedZone("GMT+7", 7*60*60)
	now := time.Now().In(loc)

	person := &model.AppPerson{
		ID:          uuid.New(),
		Fullname:    dto.Fullname,
		Gender:      dto.Gender,
		Email:       dto.Email,
		MobilePhone: dto.MobilePhone,
		SignUpFrom:  &signUpFrom,
		SignUpAt:    &now,
	}

	if err := s.appPersonRepository.Insert(ctx, person); err != nil {
		return nil, err
	}

	return person, nil
}

// LoadById retrieves a single user's detail using repository's FindById.
func (s *signUpServiceImpl) LoadById(ctx context.Context, id uuid.UUID) (*model.AppPerson, error) {
	return s.appPersonRepository.FindById(ctx, id)
}
