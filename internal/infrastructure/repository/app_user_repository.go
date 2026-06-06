package infrastructurerepository

import (
	"context"

	"github.com/dionazani/moviego-mrs-backend/internal/infrastructure/database"
	"github.com/dionazani/moviego-mrs-backend/internal/infrastructure/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AppUserRepository defines the contract for AppUser database operations.
type AppUserRepository interface {
	Insert(ctx context.Context, user *infrastructuremodel.AppUser) error
	Update(ctx context.Context, user *infrastructuremodel.AppUser) error
	FindByAppPersonId(ctx context.Context, appPersonId uuid.UUID) (*infrastructuremodel.AppUser, error)
}

type appUserRepositoryImpl struct {
	db *gorm.DB
}

// NewAppUserRepository creates a new instance of AppUserRepository.
func NewAppUserRepository(db *gorm.DB) AppUserRepository {
	return &appUserRepositoryImpl{db: db}
}

// Insert saves a new AppUser record into the database.
func (r *appUserRepositoryImpl) Insert(ctx context.Context, user *infrastructuremodel.AppUser) error {
	db := r.db
	if tx := infrastructuredatabase.GetTx(ctx); tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(user).Error
}

// Update modifies an existing AppUser record in the database.
func (r *appUserRepositoryImpl) Update(ctx context.Context, user *infrastructuremodel.AppUser) error {
	db := r.db
	if tx := infrastructuredatabase.GetTx(ctx); tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Save(user).Error
}

// FindByAppPersonId retrieves a single AppUser record matching the appPersonId.
func (r *appUserRepositoryImpl) FindByAppPersonId(ctx context.Context, appPersonId uuid.UUID) (*infrastructuremodel.AppUser, error) {
	var user infrastructuremodel.AppUser
	if err := r.db.WithContext(ctx).Where("app_person_id = ?", appPersonId).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
