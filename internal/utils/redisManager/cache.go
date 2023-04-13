package redisManager

import (
	"botgpt/pkg/redis"
	"context"
	"encoding/json"
)

// var expireTimeDuration = 24 * time.Hour
const expiredTime = 0

func GetAndCache(redisKey string, getFunc func() (interface{}, error)) (interface{}, error) {
	// 檢查 Redis 上是否存在資料

	rdb := redis.GetSingleRdb()
	ctx := context.Background()
	data, err := rdb.Get(ctx, redisKey).Result()

	if err == nil {
		// 如果 Redis 上存在該資料，則回傳該資料
		var result interface{}
		_ = json.Unmarshal([]byte(data), &result)

		return result, nil
	}

	result, err := getFunc()
	if err != nil {
		return nil, err
	}

	jsonData, _ := json.Marshal(&result)

	err = rdb.Set(ctx, redisKey, jsonData, expiredTime).Err()

	if err != nil {
		return nil, err
	}

	return result, nil
}
