package handlers

import (
	"minigo/internal/application/service"
	"minigo/internal/domain/entity"
	"minigo/internal/interfaces/dto"
	"minigo/internal/interfaces/middleware"
	resp "minigo/internal/interfaces/response"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	authService *service.AuthService
	userService *service.UserService
}

func NewAuthHandler(authService *service.AuthService, userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

// Login implements POST /api/auth/login
// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var (
		err error
		req dto.ShopLoginRequest
		ctx = c.Request.Context()
	)

	// 绑定查询参数
	if !middleware.ValidateAndBindJSON(c, &req) {
		return
	}

	// 调用服务层登录逻辑
	token, err := h.authService.Login(ctx, req.Phone, req.Password)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	resp.Ok(c, token)
}

// Register implements POST /api/auth/register
// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var (
		err error
		ctx = c.Request.Context()
		req dto.UserRegisterRequest
	)

	// 绑定查询参数
	if !middleware.ValidateAndBindJSON(c, &req) {
		return
	}

	// 构建参数
	params := service.CreateUserParams{
		Name:     req.Name,
		Phone:    req.Phone,
		Password: req.Password,
	}

	// 创建用户
	if _, err = h.userService.CreateUser(ctx, params); err != nil {
		middleware.HandleError(c, err)
		return
	}

	resp.Ok(c, nil)
}

// GetMe implements GET /api/auth/me
// GetMe 获取当前登录用户信息
func (h *AuthHandler) GetMe(c *gin.Context) {
	// 从context获取userID
	var (
		err    error
		user   *entity.User
		ctx    = c.Request.Context()
		userID = middleware.GetUserIDFromContext(c)
	)

	// 调用服务获取当前用户信息
	if user, err = h.userService.GetUserByID(ctx, userID); err != nil {
		middleware.HandleError(c, err)
		return
	}

	resp.Ok(c, user)
}

// ChangePassword implements PUT /api/auth/password
// ChangePassword 变更当前用户密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {

	var (
		req    dto.PasswordChangeRequest
		ctx    = c.Request.Context()
		userID = middleware.GetUserIDFromContext(c)
	)

	// 绑定并验证请求参数
	if !middleware.ValidateAndBindJSON(c, &req) {
		return
	}

	// 调用服务层变更密码逻辑
	if err := h.userService.ChangePassword(ctx, userID, req.OldPassword, req.NewPassword); err != nil {
		middleware.HandleError(c, err)
		return
	}

	resp.Ok(c, nil)
}

// UpdateProfile implements PUT /api/auth/profile
// UpdateProfile 修改当前用户资料
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var (
		req    dto.UpdateProfileRequest
		ctx    = c.Request.Context()
		userID = middleware.GetUserIDFromContext(c)
	)

	// 绑定并验证请求参数
	if !middleware.ValidateAndBindJSON(c, &req) {
		return
	}

	// 构建更新参数
	params := service.UpdateUserParams{
		Name:  req.Name,
		Phone: req.Phone,
	}

	// 调用服务层更新用户信息逻辑
	if err := h.userService.UpdateUser(ctx, userID, params); err != nil {
		middleware.HandleError(c, err)
		return
	}

	resp.Ok(c, nil)
}
