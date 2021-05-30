package starGo

/*
redis包对Redis的连接池进行了部分封装
对于未封装的方法有两种处理方式：
1、自行添加至代码并添加相应的注释说明
2、调用GetConnection方法，然后自己实现逻辑
*/

import (
	"fmt"

	"github.com/go-redis/redis"
)

type redisConfig struct {
	Addr string
	Pwd  string
	db   int
}

type Redis struct {
	client *redis.Client
	conf   *redisConfig
}

func NewRedis(addr, pwd string, db int, poolSize int) *Redis {
	redisClient := &Redis{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pwd,
			DB:       db,
			PoolSize: poolSize,
		}),
		conf: &redisConfig{
			Addr: addr,
			Pwd:  pwd,
			db:   db,
		},
	}

	// 验证redis链接成功
	pong, err := redisClient.client.Ping().Result()
	if err != nil {
		ErrorLog("redis连接失败，结果:%v 错误信息:%v", pong, err)
		panic(fmt.Errorf("redis连接失败,错误信息:%v", err))
	} else {
		InfoLog("redis连接成功")
	}

	redisCfg = redisClient
	return redisClient
}

func (r *Redis) GetConnection() *redis.Client {
	return r.client
}

func (r *Redis) HGet(key, field string) string {
	result, err := r.client.HGet(key, field).Result()
	if err != nil {
		return ""
	}

	return result
}

func (r *Redis) HDel(key string, field ...string) int64 {
	result, err := r.client.HDel(key, field...).Result()
	if err != nil {
		return 0
	}

	return result
}

func (r *Redis) HExists(key, field string) bool {
	result, err := r.client.HExists(key, field).Result()
	if err != nil {
		return false
	}

	return result
}

func (r *Redis) HGetAll(key string) map[string]string {
	result, err := r.client.HGetAll(key).Result()
	if err != nil {
		return nil
	}

	return result
}

func (r *Redis) HIncrBy(key, field string, incr int64) int64 {
	result, err := r.client.HIncrBy(key, field, incr).Result()
	if err != nil {
		return 0
	}

	return result
}

func (r *Redis) HIncrByFloat(key, field string, incr float64) float64 {
	result, err := r.client.HIncrByFloat(key, field, incr).Result()
	if err != nil {
		return 0
	}

	return result
}

func (r *Redis) HKeys(key string) []string {
	result, err := r.client.HKeys(key).Result()
	if err != nil {
		return nil
	}

	return result
}

func (r *Redis) HLen(key string) int64 {
	result, err := r.client.HLen(key).Result()
	if err != nil {
		return 0
	}

	return result
}

func (r *Redis) HMGet(key string, field ...string) []interface{} {
	result, err := r.client.HMGet(key, field...).Result()
	if err != nil {
		return nil
	}

	return result
}

func (r *Redis) HMSet(key string, field map[string]interface{}) string {
	result, err := r.client.HMSet(key, field).Result()
	if err != nil {
		return ""
	}

	return result
}

func (r *Redis) HSet(key, field string, value interface{}) bool {
	result, err := r.client.HSet(key, field, value).Result()
	if err != nil {
		return false
	}

	return result
}

func (r *Redis) HSetNX(key string, field string, value interface{}) bool {
	result, err := r.client.HSetNX(key, field, value).Result()
	if err != nil {
		return false
	}

	return result
}

func (r *Redis) HVals(key string) []string {
	result, err := r.client.HVals(key).Result()
	if err != nil {
		return nil
	}

	return result
}
