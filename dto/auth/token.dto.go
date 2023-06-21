package auth

type RefreshTokenRequestDto struct {
	RefreshToken string `form:"refreshToken" binding:"required"`
}
