package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func main() {
	// 连接数据库
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis error: ", err)
		return
	}
	defer conn.Close()

	// 操作数据库
	reply, err := redis.String(conn.Do("SET", "name", "Tom"))
	if err != nil {
		fmt.Println("redis set failed: ", err)
		return
	}
	fmt.Println(reply)
}
