package infrastructurerepository

import (
	"context"

	"github.com/dionazani/moviego-mrs-backend/internal/infrastructure/database"
	"github.com/dionazani/moviego-mrs-backend/internal/infrastructure/model"
	"gorm.io/gorm"
)

// AppUserTokenRepository defines the contract for AppUserToken database operations.
type AppUserTokenRepository interface {
	Insert(ctx context.Context, token *infrastructuremodel.AppUserToken) error
	Upsert(ctx context.Context, token *infrastructuremodel.AppUserToken) error
}

type appUserTokenRepositoryImpl struct {
	db *gorm.DB
}

// NewAppUserTokenRepository creates a new instance of AppUserTokenRepository.
func NewAppUserTokenRepository(db *gorm.DB) AppUserTokenRepository {
	return &appUserTokenRepositoryImpl{db: db}
}

// Insert saves a new AppUserToken record into the database.
func (r *appUserTokenRepositoryImpl) Insert(ctx context.Context, token *infrastructuremodel.AppUserToken) error {
	db := r.db
	if tx := infrastructuredatabase.GetTx(ctx); tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(token).Error
}

// Upsert inserts or updates an AppUserToken record based on app_user_id + token_type uniqueness.
// If a record with the same (app_user_id, token_type) already exists, it will update the token and expiry.
func (r *appUserTokenRepositoryImpl) Upsert(ctx context.Context, token *infrastructuremodel.AppUserToken) error {
	db := r.db
	if tx := infrastructuredatabase.GetTx(ctx); tx != nil {
		db = tx
	}
	return db.WithContext(ctx).
		Where(infrastructuremodel.AppUserToken{AppUserID: token.AppUserID, TokenType: token.TokenType}).
		Assign(infrastructuremodel.AppUserToken{TokenUser: token.TokenUser, ExpireAt: token.ExpireAt}).
		FirstOrCreate(token).Error
}
