package controllers

import (
	"botgpt/internal/config"
	"botgpt/internal/enum"
	"botgpt/internal/models/status"
	"botgpt/internal/utils"
	"botgpt/internal/utils/response"
	"botgpt/pkg/mysql"
	"botgpt/pkg/redis"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type StatusController struct {
}

func NewStatusController() *StatusController {
	return &StatusController{}
}

func (h StatusController) Status(c *gin.Context) {

	rdb := redis.GetSingleRdb()

	_, err := rdb.Ping(c).Result()
	var isConnected = true
	var errors = "ok"
	if err != nil {
		errors = err.Error()
		isConnected = false
	}
	redisInfo := status.ConnectionInfo{
		Host:        rdb.Options().Addr,
		Database:    strconv.Itoa(rdb.Options().DB),
		IsConnected: isConnected,
		Message:     errors,
	}

	_, err = mysql.GetMysqlDB()
	if err != nil {
		errors = err.Error()
		isConnected = false
	}

	dbInfo := status.ConnectionInfo{
		Host:        fmt.Sprintf("%s:%s", config.GetConfig().GetString("mysql.host"), config.GetConfig().GetString("mysql.port")),
		Database:    config.GetConfig().GetString("mysql.database"),
		IsConnected: isConnected,
		Message:     errors,
	}
	statusInfo := &status.Status{
		Version:   config.Version,
		Env:       config.Env,
		Component: config.Service,
		ServerID:  config.ServerID,
		RedisInfo: redisInfo,
		DbInfo:    dbInfo,
		GrpcInfo:  nil,
	}

	result := response.Make(true, enum.SUCCESS, "success", statusInfo)

	utils.SendResponse(http.StatusOK, result, c)
}
