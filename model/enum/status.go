package enum

type UserStatusValue string

type userStatusEnum struct {
	Active   UserStatusValue
	Inactive UserStatusValue
}

var UserStatus = &userStatusEnum{
	Active:   "ACTIVE",
	Inactive: "INACTIVE",
}
