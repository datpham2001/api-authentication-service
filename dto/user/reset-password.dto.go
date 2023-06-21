package user

type UserResetPasswordDto struct {
	User struct {
		CurrentPassword string `json:"currentPassword" validate:"required"`
		NewPassword     string `json:"newPassword" validate:"required"`
	} `json:"user" validate:"required"`
}
