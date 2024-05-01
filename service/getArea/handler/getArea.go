package handler

import (
	"context"
	"encoding/json"
	"fmt"
	modelMysql "getArea/model/mysql"
	modelRedis "getArea/model/redis"
	pb "getArea/proto"
	"getArea/utils"
	"github.com/gomodule/redigo/redis"
)

type GetArea struct{}

func (e *GetArea) Call(ctx context.Context, req *pb.CallRequest, rsp *pb.CallResponse) error {
	// 先从redis获取数据
	redisConn := modelRedis.RedisPool.Get()
	areaData, _ := redis.Bytes(redisConn.Do("GET", "areaData"))
	var areas []modelMysql.Area
	if len(areaData) == 0 { // redis中没有数据
		// 从mysql获取数据
		modelMysql.GlobalDB.Find(&areas)
		// 把数据写入redis
		areasJson, _ := json.Marshal(areas)
		_, err := redisConn.Do("SET", "areaData", areasJson)
		if err != nil {
			fmt.Println(err)
			rsp.Errno = utils.RECODE_DATAERR
			rsp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
			rsp.Data = nil
			return err
		}
	} else { // redis中有数据
		err := json.Unmarshal(areaData, &areas)
		if err != nil {
			fmt.Println(err)
			rsp.Errno = utils.RECODE_DATAERR
			rsp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
			rsp.Data = nil
			return err
		}
	}
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	for _, v := range areas {
		temp := pb.Area{
			Aid:   int32(v.Id),
			Aname: v.Name,
		}
		rsp.Data = append(rsp.Data, &temp)
	}
	return nil
}
