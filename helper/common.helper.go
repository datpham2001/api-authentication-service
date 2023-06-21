package helper

type APIResponse struct {
	Status    string      `json:"status,omitempty"`
	Message   string      `json:"message,omitempty"`
	ErrorCode string      `json:"errorCode,omitempty"`
	Data      interface{} `json:"data,omitempty"`
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
