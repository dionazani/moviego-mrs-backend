package repository

import (
	"context"
	"github.com/dionazani/moviego-mrs-backend/internal/infrastructure/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FindAllParams holds the query parameters for filtering and pagination.
type FindAllParams struct {
	Fullname string
	Email    string
	Gender   string
	Page     int // 1-based index, defaults to 1 if <= 0
	Limit    int // Items per page, defaults to 10 if <= 0
}

// AppPersonRepository defines the contract for AppPerson database operations.
type AppPersonRepository interface {
	Insert(ctx context.Context, person *model.AppPerson) error
	Update(ctx context.Context, person *model.AppPerson) error
	FindAll(ctx context.Context, params FindAllParams) ([]model.AppPerson, int64, error)
	FindById(ctx context.Context, id uuid.UUID) (*model.AppPerson, error)
	FindByFullname(ctx context.Context, fullname string) ([]model.AppPerson, error)
	FindByEmail(ctx context.Context, email string) (*model.AppPerson, error)
}

type appPersonRepositoryImpl struct {
	db *gorm.DB
}

// NewAppPersonRepository creates a new instance of AppPersonRepository.
func NewAppPersonRepository(db *gorm.DB) AppPersonRepository {
	return &appPersonRepositoryImpl{db: db}
}

// Insert saves a new AppPerson record into the database.
func (r *appPersonRepositoryImpl) Insert(ctx context.Context, person *model.AppPerson) error {
	return r.db.WithContext(ctx).Omit("UpdatedAt").Create(person).Error
}

// Update modifies an existing AppPerson record in the database.
func (r *appPersonRepositoryImpl) Update(ctx context.Context, person *model.AppPerson) error {
	return r.db.WithContext(ctx).Save(person).Error
}

// FindAll retrieves AppPerson records matching the specified filters and pagination options.
// It also returns the total count of records matching the filters.
func (r *appPersonRepositoryImpl) FindAll(ctx context.Context, params FindAllParams) ([]model.AppPerson, int64, error) {
	var people []model.AppPerson
	var total int64

	db := r.db.WithContext(ctx).Model(&model.AppPerson{})

	// Apply WHERE filters if provided
	if params.Fullname != "" {
		db = db.Where("fullname ILIKE ?", "%"+params.Fullname+"%")
	}
	if params.Email != "" {
		db = db.Where("email = ?", params.Email)
	}
	if params.Gender != "" {
		db = db.Where("gender = ?", params.Gender)
	}

	// Count total records matching the filters
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination defaults
	page := params.Page
	if page <= 0 {
		page = 1
	}
	limit := params.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Fetch paginated data
	if err := db.Offset(offset).Limit(limit).Find(&people).Error; err != nil {
		return nil, 0, err
	}

	return people, total, nil
}

// FindById retrieves a single AppPerson record by its UUID.
func (r *appPersonRepositoryImpl) FindById(ctx context.Context, id uuid.UUID) (*model.AppPerson, error) {
	var person model.AppPerson
	if err := r.db.WithContext(ctx).First(&person, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &person, nil
}

// FindByFullname retrieves AppPerson records matching the exact fullname.
func (r *appPersonRepositoryImpl) FindByFullname(ctx context.Context, fullname string) ([]model.AppPerson, error) {
	var people []model.AppPerson
	if err := r.db.WithContext(ctx).Where("fullname = ?", fullname).Find(&people).Error; err != nil {
		return nil, err
	}
	return people, nil
}

// FindByEmail retrieves a single AppPerson record matching the exact email.
func (r *appPersonRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.AppPerson, error) {
	var person model.AppPerson
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&person).Error; err != nil {
		return nil, err
	}
	return &person, nil
}
