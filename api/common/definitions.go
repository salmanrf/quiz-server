package common

type ApiResponse[resType interface{}] struct {
	Message string `json:"message"`
	Data resType `json:"data"`
}