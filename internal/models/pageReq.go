package models

import "encoding/json"

type PageRequest struct {
	PageSize    json.Number `json:"pageSize"`
	CurrentPage json.Number `json:"currentPage"`
	OrderBy     *string     `json:"orderby,omitempty"`
	Sort        *string     `json:"sort,omitempty"`
}
