package utils

type Response struct {
	Data  any `json:"data"`
	Error any `json:"error"`
}

func SuccessResponse(data any) Response {
	return Response{
		Data: data,
	}
}

func ErrorResponse(err any) Response {
	return Response{
		Error: err,
	}
}
