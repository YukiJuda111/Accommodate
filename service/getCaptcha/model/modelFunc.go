package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// SaveImg 存储Img到Redis
func SaveImg(code, uuid string) error {
	// 连接数据库
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis error: ", err)
		return err
	}
	defer conn.Close()

	// 操作数据库
	_, err = conn.Do("SETEX", uuid, 300, code)
	if err != nil {
		fmt.Println("redis set failed: ", err)
		return err
	}
	return err
}
