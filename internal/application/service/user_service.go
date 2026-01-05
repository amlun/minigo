package service

import (
	"context"
	"minigo/internal/domain/entity"
	"minigo/internal/domain/repository"
	"minigo/internal/infrastructure/id"
	"minigo/internal/infrastructure/tx"

	"github.com/shopspring/decimal"
)

type UserService struct {
	userRepo  repository.UserRepository
	txManager *tx.Manager
}

// NewUserService 创建用户服务实例
func NewUserService(
	userRepo repository.UserRepository,
	txManager *tx.Manager,
) *UserService {
	return &UserService{
		userRepo:  userRepo,
		txManager: txManager,
	}
}

// CreateUserParams 管理端创建用户参数
type CreateUserParams struct {
	Name       string
	Phone      string
	Password   string
	Status     int16
	ReferrerID *int64 // 邀请人ID
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, params CreateUserParams) (*entity.User, error) {
	// 定义变量
	var (
		err  error
		user *entity.User
	)
	// 构造用户实体
	user = &entity.User{
		ID:         id.NextID(),
		Name:       params.Name,
		Phone:      params.Phone,
		Password:   params.Password,
		Status:     params.Status,
		ReferrerID: params.ReferrerID,
	}
	// 在事务中执行
	if err = s.txManager.InTx(ctx, func(txCtx context.Context) error {
		// 先检查

		// 再添加
		return s.userRepo.Create(txCtx, user)
	}); err != nil {
		return nil, err
	}

	// 返回结果
	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// UpdateUserParams 更新用户参数
type UpdateUserParams struct {
	Name          string
	Phone         string
	Status        *int16
	IsAdmin       *bool
	Remark        *string
	OrderMinPrice decimal.Decimal
	OrderMaxPrice decimal.Decimal
	OrderCount    int
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(ctx context.Context, id int64, params UpdateUserParams) error {
	// 定义变量
	var (
		err  error
		user *entity.User
	)

	// 在事务中保存更新
	if err = s.txManager.InTx(ctx, func(txCtx context.Context) error {
		// 获取用户实体（加锁）
		user, err = s.userRepo.GetForUpdate(txCtx, id)
		if err != nil {
			return err
		}
		if user == nil {
			return ErrUserNotFound
		}

		// 更新用户基本信息
		user.Name = params.Name
		user.Phone = params.Phone
		if params.Status != nil {
			user.Status = *params.Status
		}
		// 更新用户实体
		return s.userRepo.Update(txCtx, user)
	}); err != nil {
		return err
	}

	// 返回结果
	return nil
}

func (s *UserService) ChangePassword(ctx context.Context, id int64, oldPassword, newPassword string) error {
	var (
		err  error
		user *entity.User
	)

	// 在事务中执行
	if err = s.txManager.InTx(ctx, func(txCtx context.Context) error {
		// 获取用户实体
		if user, err = s.userRepo.GetForUpdate(txCtx, id); err != nil {
			return ErrUserNotFound
		}
		// 校验旧密码
		if err = user.CheckPassword(oldPassword); err != nil {
			return err
		}
		// 使用实体的领域方法设置密码，然后通过通用 Update 方法持久化
		user.Password = newPassword
		return s.userRepo.Update(txCtx, user)
	}); err != nil {
		return err
	}

	return nil
}
