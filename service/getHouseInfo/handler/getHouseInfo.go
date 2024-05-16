package handler

import (
	"context"
	"encoding/json"
	mysqlModel "getHouseInfo/model/mysql"
	redisModel "getHouseInfo/model/redis"
	pb "getHouseInfo/proto"
	"getHouseInfo/utils"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type GetHouseInfo struct{}

func (e *GetHouseInfo) Call(ctx context.Context, req *pb.CallRequest, rsp *pb.CallResponse) error {
	var data pb.Data
	//从redis数据库中获取数据
	redisConn := redisModel.RedisPool.Get()
	houseBuffer, _ := redis.Bytes(redisConn.Do("GET", "houseId_"+strconv.Itoa(int(req.HouseId))))
	if len(houseBuffer) == 0 {
		//如果查不到数据,就从数据库获取
		house, orders, orderUsers, user, facs, imgs, err := mysqlModel.GetDetail(req.HouseId)
		if err != nil {
			rsp.Errno = utils.RECODE_DATAERR
			return err
		}

		//获取评论数据
		var comments []*pb.CommentInfo

		for k, v := range orders {
			var comment pb.CommentInfo
			comment.Ctime = v.UpdatedAt.Format("2006:01:02 15:04:05")
			comment.UserName = orderUsers[k].Name
			comment.Comment = v.Comment

			comments = append(comments, &comment)
		}

		//获取家具数据
		var fids []int32
		for _, v := range facs {
			fids = append(fids, int32(v.Id))
		}

		//获取图片
		var imgPaths []string
		imgPaths = append(imgPaths, house.IndexImageUrl)
		for _, v := range imgs {
			imgPaths = append(imgPaths, v.Url)
		}

		//设置返回数据
		var houseInfo pb.HouseInfo
		houseInfo.Acreage = int32(house.Acreage)
		houseInfo.Address = house.Address
		houseInfo.Beds = house.Beds
		houseInfo.Capacity = int32(house.Capacity)
		houseInfo.Comments = comments
		houseInfo.Deposit = int32(house.Deposit)
		houseInfo.Facilities = fids
		houseInfo.Hid = int32(house.ID)
		houseInfo.ImgUrls = imgPaths
		houseInfo.MinDays = int32(house.MinDays)
		houseInfo.MaxDays = int32(house.MaxDays)
		houseInfo.Price = int32(house.Price)
		houseInfo.RoomCount = int32(house.RoomCount)
		houseInfo.Title = house.Title
		houseInfo.Unit = house.Unit
		houseInfo.UserAvatar = user.AvatarUrl
		houseInfo.UserId = int32(house.UserId)
		houseInfo.UserName = user.Name
		data.House = &houseInfo
		data.UserId = int32(user.ID)

		//把查询到的数据存储到redis中
		buffer, err := json.Marshal(data)
		if err != nil {
			return err
		}

		redisConn.Do("set", "houseId_"+strconv.Itoa(int(house.ID)), buffer)
		rsp.Errno = utils.RECODE_OK
		rsp.Data = &data
		return nil
	}

	//如果查到数据,就直接返回
	err := json.Unmarshal(houseBuffer, &data)
	if err != nil {
		rsp.Errno = utils.RECODE_DATAERR
		return err
	}
	//返回数据
	rsp.Errno = utils.RECODE_OK
	rsp.Data = &data
	return nil
}

func (e *GetHouseInfo) GetHouseIndex(ctx context.Context, req *pb.IndexRequest, resp *pb.IndexResponse) error {
	//获取房屋信息
	houseResp, err := mysqlModel.GetIndexHouse()
	if err != nil {
		resp.Errno = utils.RECODE_DBERR
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Data = &pb.GetData{Houses: houseResp}

	return nil
}

func (e *GetHouseInfo) SearchHouse(ctx context.Context, req *pb.SearchRequest, resp *pb.SearchResponse) error {
	//根据传入的参数,查询符合条件的房屋信息
	houseResp, err := mysqlModel.SearchHouse(req.Aid, req.Sd, req.Ed, req.Sk)
	if err != nil {
		resp.Errno = utils.RECODE_DATAERR
		return nil
	}

	resp.Errno = utils.RECODE_OK

	resp.Data = &pb.GetData{Houses: houseResp}
	return nil
}
