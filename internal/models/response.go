package models

type ResponseResult struct {
	RequestId string      `json:"requestId"`
	Success   bool        `json:"success"`
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Content   interface{} `json:"content"`
}

type ResponsePageResult struct {
	RequestId string        `json:"requestId"`
	Success   bool          `json:"success"`
	Code      string        `json:"code"`
	Message   string        `json:"message"`
	Page      *PageResponse `json:"pages"`
	Content   interface{}   `json:"content"`
}
