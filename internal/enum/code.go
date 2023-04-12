package enum

type ResponseCode string

// Response Code
const (
	SUCCESS                     ResponseCode = "0000"
	FALIURE                     ResponseCode = "9999"
	UNKNOW_ERROR                ResponseCode = "9998"
	UNAUTHORIZE                 ResponseCode = "9997"
	OVER_LIMIT                  ResponseCode = "9996"
	DUPLICATE_IDCARD            ResponseCode = "9995"
	VERIFY_CODE_ERROR           ResponseCode = "9994"
	LOGIN_FAILED_TOO_MANY_TIMES ResponseCode = "9993"
	INVALID_ACCOUNT_PASSWORD    ResponseCode = "9992"
	APIAUTHORIZEFAILURE         ResponseCode = "9991"
	APPIDNOTFOUND               ResponseCode = "9990"
	APIVALIDATEERROR            ResponseCode = "9989"
	InsufficientPoint           ResponseCode = "9988"
	ACCOUNT_DISABLE             ResponseCode = "9987"
)
