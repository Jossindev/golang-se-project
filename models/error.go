package models

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ErrorResponse(code int, message string) Error {
	return Error{
		Code:    code,
		Message: message,
	}
}
