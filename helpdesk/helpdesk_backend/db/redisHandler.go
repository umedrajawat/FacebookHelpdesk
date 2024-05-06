package db

import (
	"fmt"
	configs "helpdesk_backend/config"
	"helpdesk_backend/logger"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// RedisClient : Main redis client
var RedisClient *redis.Client

// Intialise and setup all Redis connections
func SetupRedisDB() {
	rp := "127.0.0.1"
	rpo := "31019"
	ra := ""
	rdb := configs.REDIS_DB
	// rsdb := configs.REDIS_SCREENER_DB
	rd, _ := strconv.ParseInt(rdb, 10, 64)
	// rsd, _ := strconv.ParseInt(rsdb, 10, 64)
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     rp + ":" + rpo,
		Password: ra,
		DB:       int(rd),
		PoolSize: 1000,
	})
	fmt.Println("RedisClient", RedisClient)
}

func GetRedisKeyTtl(con *redis.Client, key string) time.Duration {
	res := con.TTL(key)
	val := res.Val()
	return val
}

func GetRedisKey(con *redis.Client, key string) string {
	res := con.Get(key)
	val := res.Val()
	return val
}

func GetRedisPatternKeys(con *redis.Client, pattern string) []string {
	res := con.Keys(pattern)
	val := res.Val()
	return val
}

func SetRedisKey(con *redis.Client, key string, val interface{}, ttl time.Duration) {
	fmt.Println("SETTING REDIS KEY ", key)
	err := con.Set(key, val, ttl).Err()

	if err != nil {
		fmt.Println("Error setting key", key, "to", val, err)
	}
}

func GetUserVersionPrefFromRedis(userUUID string) string {
	key := fmt.Sprintf("user_version_pref%s", userUUID)
	fmt.Println("key", key)
	str, err := RedisClient.Get(key).Result()
	if err != nil {
		logger.ZapLogger.Error(err)
	}
	return str
}

func DeleteRedisKey(con *redis.Client, key string) {
	err := con.Del(key).Err()
	if err != nil {
		fmt.Println("Error deleting key", key, err)
	}
}

func GetHashAll(con *redis.Client, key string) map[string]string {
	res := con.HGetAll(key)
	val := res.Val()
	return val
}

func GetHash(con *redis.Client, key string, field string) string {
	res := con.HGet(key, field)
	val := res.Val()
	return val
}
