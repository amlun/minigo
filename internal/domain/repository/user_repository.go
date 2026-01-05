package repository

import (
	"context"
	"minigo/internal/domain/entity"
)

type UserRepository interface {
	// Create persists a new user.
	Create(ctx context.Context, user *entity.User) error

	// Update updates user basic fields.
	Update(ctx context.Context, user *entity.User) error

	// GetByID returns user by id.
	GetByID(ctx context.Context, id int64) (*entity.User, error)

	// GetByPhone returns user by phone.
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)

	// GetForUpdate 加悲观锁读取用户（需在事务上下文中使用）
	GetForUpdate(ctx context.Context, id int64) (*entity.User, error)

	// Delete removes user by id.
	Delete(ctx context.Context, id int64) error
}
