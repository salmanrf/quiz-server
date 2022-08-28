package common

type FindRes[ResType interface{}] struct {
	Items []ResType `json:"items"`
	PageSize int `json:"page_size"`
	PageNumber int `json:"page_number"`
	TotalItems int `json:"total_items"`
}