package enum

type ErrorCodeEnumValue string

type (
	errorCodeInvalidEnum struct {
		ParseData ErrorCodeEnumValue
		Email     ErrorCodeEnumValue
		Username  ErrorCodeEnumValue
		Password  ErrorCodeEnumValue
	}

	errorCodeRequiredEnum struct {
		Email    ErrorCodeEnumValue
		Username ErrorCodeEnumValue
		Password ErrorCodeEnumValue
	}

	errorCodeExistedEnum struct {
		UsernameOrEmail ErrorCodeEnumValue
	}

	errorCodePackageEnum struct {
		Bcrypt ErrorCodeEnumValue
	}

	errorCodeDatabaseOperation struct {
		Create ErrorCodeEnumValue
		Read   ErrorCodeEnumValue
		Update ErrorCodeEnumValue
		Delete ErrorCodeEnumValue
	}

	errorCodeNotExisted struct {
		User ErrorCodeEnumValue
	}
)

var (
	ErrorCodeInvalid = &errorCodeInvalidEnum{
		ParseData: "INVALID_PARSE_DATA",
		Email:     "INVALID_EMAIL_FORMAT",
		Username:  "INVALID_USERNAME_FORMAT",
		Password:  "INVALID_PASSWORD_FORMAT",
	}

	ErrorCodeRequired = &errorCodeRequiredEnum{
		Email:    "REQUIRED_EMAIL",
		Username: "REQUIRED_USERNAME",
		Password: "REQUIRED_PASSWORD",
	}

	ErrorCodeExisted = &errorCodeExistedEnum{
		UsernameOrEmail: "EXISTED_USERNAME_OR_EMAIL",
	}

	ErrorCodePackage = &errorCodePackageEnum{
		Bcrypt: "BCRYPT_ERROR",
	}

	ErrorCodeDatabaseOperation = &errorCodeDatabaseOperation{
		Create: "CREATE_RECORD_ERROR",
		Read:   "READ_RECORD_ERROR",
		Update: "UPDATE_RECORD_ERROR",
		Delete: "DELETE_RECORD_ERROR",
	}

	ErrorCodeNotExisted = &errorCodeNotExisted{
		User: "NOT_EXISTED_USER",
	}
)
