package dto

// ShopLoginRequest represents login payload for shop context.
type ShopLoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// PasswordChangeRequest represents shop password change payload.
type PasswordChangeRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=20"`
}

// UpdateProfileRequest represents shop profile update payload.
type UpdateProfileRequest struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required,min=11,max=11"`
}
