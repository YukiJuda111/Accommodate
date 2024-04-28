package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// CheckImgCode saves the img code
func CheckImgCode(uuid, code string) bool {
	conn := RedisPool.Get()
	defer conn.Close()
	imgCode, err := redis.String(conn.Do("get", uuid))
	if err != nil {
		fmt.Println("redis get uuid failed, err:", err)
		return false
	}
	return imgCode == code
}

// SaveSmsCode saves the sms code
func SaveSmsCode(phoneNum, code string) error {
	conn := RedisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SETEX", phoneNum+"_code", 300, code)
	if err != nil {
		fmt.Println("redis setex phone failed, err:", err)
	}
	return err
}

// CheckSmsCode checks the sms code
func CheckSmsCode(phoneNum, code string) bool {
	conn := RedisPool.Get()
	defer conn.Close()
	smsCode, err := redis.String(conn.Do("GET", phoneNum+"_code"))
	if err != nil {
		fmt.Println("redis get phone failed, err:", err)
		return false
	}
	return smsCode == code
}
