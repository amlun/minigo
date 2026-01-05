package entity

import (
	"context"
	apperrors "minigo/internal/domain/errors"
	"minigo/pkg/utils"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID         int64      `bun:"id,pk,autoincrement" json:"id,string"`
	Name       string     `bun:"name,notnull" json:"name"`
	Phone      string     `bun:"phone,notnull" json:"phone"`
	Password   string     `bun:"password,notnull" json:"-"`
	Status     int16      `bun:"status,notnull,default:0" json:"status"`
	ReferrerID *int64     `bun:"referrer_id" json:"referrer_id,string,omitempty"`
	DeletedAt  *time.Time `bun:"deleted_at,soft_delete,nullzero" json:"deleted_at,omitempty"`
	CreatedAt  time.Time  `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt  time.Time  `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`

	// -- 关系
	Referrer *User `bun:"rel:belongs-to,join:referrer_id=id" json:"referrer,omitempty"`
}

var _ bun.BeforeInsertHook = (*User)(nil)
var _ bun.BeforeUpdateHook = (*User)(nil)

// BeforeInsert - 插入前处理
func (u *User) BeforeInsert(ctx context.Context, query *bun.InsertQuery) error {
	account := query.GetModel().Value().(*User)
	account.BcryptPassword()
	account.Name = strings.ToLower(account.Name)
	return nil
}

// BeforeUpdate - 更新前处理
func (u *User) BeforeUpdate(ctx context.Context, query *bun.UpdateQuery) error {
	account := query.GetModel().Value().(*User)
	account.BcryptPassword()
	account.Name = strings.ToLower(account.Name)
	return nil
}

// BcryptPassword - 对密码进行加密处理
func (u *User) BcryptPassword() {
	// 使用BcryptHash进行加密
	if u.Password != "" {
		u.Password = utils.BcryptHash(u.Password)
	}
}

// CheckPassword - 校验密码是否一致
func (u *User) CheckPassword(password string) error {
	if utils.BcryptCheck(password, u.Password) {
		return nil
	}
	return apperrors.NewAuthError("AUTH_001", "用户名或密码错误")
}
