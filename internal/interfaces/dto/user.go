package dto

// UserCreateRequest 管理端创建用户请求
type UserCreateRequest struct {
	Name       string `json:"name" binding:"required,min=2,max=50"`
	Phone      string `json:"phone" binding:"required,len=11"`
	Password   string `json:"password" binding:"required,min=6,max=20"`
	Status     int16  `json:"status" binding:"min=0,max=2"`
	ReferrerID *int64 `json:"referrer_id,string"`
}

// UserBatchCreateRequest 批量创建用户请求
type UserBatchCreateRequest struct {
	Users []UserCreateRequest `json:"users" binding:"required,dive"`
}

type FailedUser struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Error string `json:"error"`
}

type UserBatchCreateResponse struct {
	SuccessCount int          `json:"success_count"`
	FailedCount  int          `json:"failed_count"`
	FailedUsers  []FailedUser `json:"failed_users"`
}

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=50"`
	Phone    string `json:"phone" binding:"required,len=11"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	Keyword string `form:"keyword"`
	Status  *int16 `form:"status"`
	Page    int    `form:"page" binding:"min=1" default:"1"`
	Size    int    `form:"size" binding:"min=1,max=100" default:"20"`
}

// UserUpdateRequest 用户信息更新请求
type UserUpdateRequest struct {
	Name   string `json:"name" binding:"required,min=2,max=50"`
	Phone  string `json:"phone" binding:"required,len=11"`
	Status *int16 `json:"status"`
}

// UserPasswordResetRequest 用户密码重置请求
type UserPasswordResetRequest struct {
	Password string `json:"password" binding:"required,min=6,max=20"`
}
