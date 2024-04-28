package main

import (
	modelMysql "RentHouse/web/model/mysql"
	modelRedis "RentHouse/web/model/redis"
	"RentHouse/web/router"
)

func main() {
	_, err := modelMysql.GormInit()
	if err != nil {
		return
	}

	modelRedis.RedisInit()

	router.Init()
}
