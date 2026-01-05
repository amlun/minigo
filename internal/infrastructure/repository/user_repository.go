package repository

import (
	"context"

	"minigo/internal/domain/entity"
	"minigo/internal/domain/repository"
	"minigo/internal/infrastructure/dbctx"

	"github.com/uptrace/bun"
)

// BunUserRepository implements UserRepository using Bun ORM
type BunUserRepository struct {
	DB *bun.DB
}

// NewBunUserRepository creates a new BunUserRepository
func NewBunUserRepository(db *bun.DB) repository.UserRepository {
	return &BunUserRepository{DB: db}
}

func (r *BunUserRepository) Create(ctx context.Context, user *entity.User) error {
	db := dbctx.FromCtx(ctx, r.DB)
	_, err := db.NewInsert().Model(user).Exec(ctx)
	return ConvertExecError(err)
}

func (r *BunUserRepository) Update(ctx context.Context, user *entity.User) error {
	db := dbctx.FromCtx(ctx, r.DB)
	result, err := db.NewUpdate().
		Model(user).
		WherePK().
		Exec(ctx)
	return CheckUpdateResult(result, err)
}

func (r *BunUserRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	db := dbctx.FromCtx(ctx, r.DB)
	var user = entity.User{ID: id}
	err := db.NewSelect().
		Model(&user).
		WherePK().
		Scan(ctx)
	if err != nil {
		return nil, ConvertQueryError(err)
	}
	return &user, nil
}

func (r *BunUserRepository) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
	db := dbctx.FromCtx(ctx, r.DB)
	var user entity.User
	err := db.NewSelect().
		Model(&user).
		Where("phone = ?", phone).
		Scan(ctx)
	if err != nil {
		return nil, ConvertQueryError(err)
	}
	return &user, nil
}

func (r *BunUserRepository) GetForUpdate(ctx context.Context, id int64) (*entity.User, error) {
	db := dbctx.FromCtx(ctx, r.DB)
	var user = entity.User{ID: id}
	err := db.NewSelect().
		Model(&user).
		WherePK().
		For("UPDATE").
		Scan(ctx)
	if err != nil {
		return nil, ConvertQueryError(err)
	}
	return &user, nil
}

func (r *BunUserRepository) Delete(ctx context.Context, id int64) error {
	db := dbctx.FromCtx(ctx, r.DB)
	user := &entity.User{ID: id}
	result, err := db.NewDelete().
		Model(user).
		WherePK().
		Exec(ctx)
	return CheckDeleteResult(result, err)

}
