package model

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Success(data any) Response {
	return Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

func Failure(apiErr *APIError) Response {
	if apiErr == nil {
		return Response{
			Code:    CodeInternal,
			Message: "internal server error",
		}
	}

	return Response{
		Code:    apiErr.Code,
		Message: apiErr.Message,
	}
}
