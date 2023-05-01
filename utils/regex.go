package utils

import "regexp"

/*
	Username must be:
		+ Include - or _
		+ Length must be in [5,12]

*/
func ValidateUsername(username string) bool {
	const regex = "[A-Za-z-_]{5,12}"

	ok, _ := regexp.MatchString(regex, username)
	return ok
}

/*
	Email must be: xxx@gmail.com
		+ Pattern before '@' must be 1 character
		+ Pattern after '@' must be have length in range [2,6]
		+ Pattern after '.' must be have length in range [2,3]
*/
func ValidateEmail(email string) bool {
	const regex = "(\\w+)@(\\w{2,6})\\.(\\w{2,3})"

	ok, _ := regexp.MatchString(regex, email)
	return ok
}

/*
	Password must be:
		+ At least 8 character long
		+ At least 1 lowercase and uppercase letter
		+ At least 2 numbers
		+ At least 1 special character
*/
func ValidatePassword(password string) bool {
	const regex = "^(?=.*\\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[a-zA-Z]).{8,}$"

	ok, _ := regexp.MatchString(regex, password)
	return ok
}
