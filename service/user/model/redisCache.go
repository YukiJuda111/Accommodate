package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

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

func CheckImgCode(uuid, code string) bool {
	conn := redisPool.Get()
	defer conn.Close()
	imgCode, err := redis.String(conn.Do("get", uuid))
	if err != nil {
		fmt.Println("redis get uuid failed, err:", err)
		return false
	}
	return imgCode == code
}

func SaveSmsCode(phoneNum, code string) error {
	conn := redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SETEX", phoneNum+"_code", 300, code)
	if err != nil {
		fmt.Println("redis setex phone failed, err:", err)
	}
	return err
}
