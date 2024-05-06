package model

import "github.com/gomodule/redigo/redis"

var redisPool *redis.Pool

func RedisInit() {
	redisPool = &redis.Pool{
		MaxIdle:     20,
		MaxActive:   50,
		IdleTimeout: 300,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379")
		},
	}
}
