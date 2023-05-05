package enum

const (
	EnUS  = "en-US"
	JaJP  = "ja-JP"
	CmnCN = "cmn-CN"
	ZhTw  = "zh-TW"
)

var LangMap = map[string]string{
	"Japanese": JaJP,
	"Chinese":  CmnCN,
	"English":  EnUS,
}
