package dto

type RefreshTokenRequest struct {
	RefreshToken string `form:"refreshToken" binding:"required"`
}
