package repository

import (
	"botgpt/internal/config"
	"sync"

	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
)

var testAddr string
var rdb *redis.Client
var redisOnce sync.Once

func GetSingleRdb() *redis.Client {
	redisOnce.Do(func() {
		rdb = CreateRdbConnection()
	})
	return rdb
}

func CreateRdbConnection() *redis.Client {

	c := config.GetConfig()
	host := c.GetString("redis.host")
	port := c.GetString("redis.port")
	password := c.GetString("redis.password")
	db := c.GetInt("redis.db")

	addr := host + ":" + port
	if gin.Mode() == "test" {
		addr = testAddr
	}
	rdbSingle := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
	})
	return rdbSingle
}

func SetTestAddr(addr string) {
	testAddr = addr
}
