package enum

type UserRoleValue string

type userRole struct {
	User   UserRoleValue
	Admin  UserRoleValue
	Author UserRoleValue
}

var UserRole = &userRole{
	User:   "USER",
	Admin:  "ADMIN",
	Author: "AUTHOR",
}

type ProviderNameValue string

type providerName struct {
	Google ProviderNameValue
	Github ProviderNameValue
}

var ProviderName = &providerName{
	Google: "GOOGLE",
	Github: "GITHUB",
}
