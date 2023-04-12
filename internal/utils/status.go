package utils

var StatusInfo *Status

func SetStatusInfo(status *Status) {
	StatusInfo = status
}

type Status struct {
	Version   string      `json:"version"`
	Env       string      `json:"env"`
	Component string      `json:"component"`
	ServerID  string      `json:"ServerID"`
	RedisInfo interface{} `json:"redisInfo"`
	DbInfo    interface{} `json:"dbInfo"`
	GrpcInfo  interface{} `json:"grpcInfo"`
}
