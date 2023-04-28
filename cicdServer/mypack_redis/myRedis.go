package myRedis

import (
	dlogger "cicdServer/log/zlog"
	"context"
	"github.com/redis/go-redis/v9"
	"os"
)

var ctx = context.Background()

func InitRedisConn() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "******", // no password set
		DB:       0,        // use default DB
		PoolSize: 8,
	})
	err := rdb.Ping(ctx).Err()
	if err != nil {
		dlogger.Error(err.Error())
		os.Exit(1)
	}
	return rdb
}

func RpopList(c *redis.Client, keyNamne string) (result string) {
	val, err := c.RPop(ctx, keyNamne).Result()
	if err != nil {
		dlogger.Info("redis has no task " + err.Error())
		result = ""
	} else {
		result = val
	}
	return
}

func LpushList(c *redis.Client, keyNamne string, val string) (res int) {
	err := c.LPush(ctx, keyNamne, val)
	if err != nil {
		dlogger.Info("get task from redis " + err.String())
		res = 0
		return
	}
	res = 1
	return
}

//func Lrange(c *redis.Client, keyNamne string) (result []string) {
//	val, err := c.LRange(ctx, keyNamne, -1, 0).Result()
//	if err != nil {
//		fmt.Println(err)
//		result = []string{}
//	} else {
//		result = val
//	}
//	return
//}
