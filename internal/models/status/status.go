package status

type Status struct {
	Version   string         `json:"version"`
	Env       string         `json:"env"`
	Component string         `json:"component"`
	ServerID  string         `json:"ServerID"`
	RedisInfo ConnectionInfo `json:"redisInfo"`
	DbInfo    ConnectionInfo `json:"dbInfo"`
	GrpcInfo  interface{}    `json:"grpcInfo"`
}

type ConnectionInfo struct {
	Host        string `json:"host"`
	Database    string `json:"database,omitempty"`
	IsConnected bool   `json:"isConnected"`
	Message     string `json:"message"`
}
