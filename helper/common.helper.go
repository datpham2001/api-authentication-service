package helper

type APIResponse struct {
	Code      int
	Status    string
	Message   string
	ErrorCode string
	Data      interface{}
}

type apiStatusEnum struct {
	Ok           string
	Unauthorized string
	Invalid      string
	Error        string
	Notfound     string
}

var APIStatus = &apiStatusEnum{
	"Ok",
	"Unauthorized",
	"Invalid",
	"Error",
	"Notfound",
}
