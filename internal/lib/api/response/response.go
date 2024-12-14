package response

type ResponseStatus string

const (
	StatusOk    ResponseStatus = "OK"
	StatusError ResponseStatus = "Error"
)

type Response struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message,omitempty"`
}

func OK() Response {
	return Response{
		Status: StatusOk,
	}
}

func Error(msg string) Response {
	return Response{
		Status:  StatusError,
		Message: msg,
	}
}
