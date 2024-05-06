package handler

import (
	"context"
	mysqlModel "postHouses/model/mysql"
	pb "postHouses/proto"
	"postHouses/utils"
	"strconv"
)

type PostHouses struct{}

func (e *PostHouses) Call(ctx context.Context, req *pb.CallRequest, rsp *pb.CallResponse) error {
	var house mysqlModel.House
	acrea, _ := strconv.Atoi(req.Acreage)
	house.Acreage = acrea
	house.Address = req.Address
	areaId, _ := strconv.Atoi(req.AreaId)
	house.AreaId = uint(areaId)
	house.Beds = req.Beds
	capicity, _ := strconv.Atoi(req.Capacity)
	house.Capacity = capicity
	dep, _ := strconv.Atoi(req.Deposit)
	house.Deposit = dep
	maxDays, _ := strconv.Atoi(req.MaxDays)
	house.MaxDays = maxDays
	minDays, _ := strconv.Atoi(req.MinDays)
	house.MinDays = minDays
	price, _ := strconv.Atoi(req.Price)
	house.Price = price
	house.Title = req.Title
	house.Unit = req.Unit

	roomCount, _ := strconv.Atoi(req.RoomCount)
	house.RoomCount = roomCount
	// 插入房屋信息
	houseId, err := mysqlModel.InsertHouse(req.Name, house, req.Facility)
	if err != nil {
		return err
	}

	// 返回数据
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	rsp.Data = map[string]int32{"house_id": int32(houseId)}

	return nil
}
