package response

import (
	res "botgpt/internal/enum"
	"botgpt/internal/models"
)

func OK() models.ResponseResult {
	result := models.ResponseResult{
		//RequestId: GetRequestId(),
		Success: true,
		Code:    string(res.SUCCESS),
		Message: "success",
		Content: nil,
	}

	return result
}

func OKHasContent(t interface{}) models.ResponseResult {
	result := models.ResponseResult{
		//RequestId: GetRequestId(),
		Success: true,
		Code:    string(res.SUCCESS),
		Message: "",
		Content: t,
	}

	return result
}

func Failure(msg string, code res.ResponseCode) models.ResponseResult {
	result := models.ResponseResult{
		//RequestId: GetRequestId(),
		Success: false,
		Code:    string(code),
		Message: msg,
		Content: nil,
	}
	return result
}

func Make(success bool, code res.ResponseCode, msg string, t interface{}) models.ResponseResult {

	result := models.ResponseResult{
		//RequestId: GetRequestId(),
		Success: success,
		Code:    string(code),
		Message: msg,
		Content: t,
	}

	return result
}

func MakePage(success bool, code res.ResponseCode, msg string, t interface{}, page *models.PageResponse) models.ResponsePageResult {

	result := models.ResponsePageResult{
		//RequestId: GetRequestId(),
		Success: success,
		Code:    string(code),
		Message: msg,
		Page:    page,
		Content: t,
	}

	return result
}
