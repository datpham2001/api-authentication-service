package enum

var (
	TRUE  = true
	FALSE = false
)

type UserStatusValue string

type userStatusEnum struct {
	Active   UserStatusValue
	Inactive UserStatusValue
}

var UserStatus = &userStatusEnum{
	Active:   "ACTIVE",
	Inactive: "INACTIVE",
}
