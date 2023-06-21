package enum

type UserRoleValue string

type userRole struct {
	User   UserRoleValue
	Author UserRoleValue
}

var UserRole = &userRole{
	User:   "USER",
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
