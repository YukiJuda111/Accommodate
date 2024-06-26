package model

import (
	"fmt"
	pb "orders/proto"
	"strconv"
	"time"
)

// TODO: 修改prefixUrl为自己的七牛云存储空间地址
var prefixUrl = "http://sf9x6ixky.hd-bkt.clouddn.com/"

func InsertOrder(houseId, beginDate, endDate, userName string) (int, error) {
	//获取插入对象
	var order OrderHouse

	//给对象赋值
	hid, _ := strconv.Atoi(houseId)
	order.HouseId = uint(hid)

	//把string类型的时间转换为time类型
	bDate, _ := time.Parse("2006-01-02", beginDate)
	order.BeginDate = bDate

	eDate, _ := time.Parse("2006-01-02", endDate)
	order.EndDate = eDate

	//需要userId
	/*var user User
	GlobalDB.Where("name = ?",userName).Find(&user)*/
	//select id form user where name = userName

	type UserData struct {
		Id int
	}

	var userData UserData
	if err := GlobalDB.Raw("select id from user where name = ?", userName).Scan(&userData).Error; err != nil {
		fmt.Println("获取用户数据错误", err)
		return 0, err
	}

	//获取days
	dur := eDate.Sub(bDate)
	order.Days = int(dur.Hours()) / 24
	order.Status = "WAIT_ACCEPT"

	//房屋的单价和总价
	var house House
	GlobalDB.Where("id = ?", hid).Find(&house).Select("price")
	order.HousePrice = house.Price
	order.Amount = house.Price * order.Days

	order.UserId = uint(userData.Id)
	if err := GlobalDB.Create(&order).Error; err != nil {
		fmt.Println("插入订单失败", err)
		return 0, err
	}
	return int(order.ID), nil
}

func GetOrderInfo(userName, role string) ([]*pb.OrdersData, error) {
	//最终需要的数据
	var orderResp []*pb.OrdersData
	//获取当前用户的所有订单
	var orders []OrderHouse

	type UserData struct {
		Id int
	}
	var userData UserData
	//用原生查询的时候,查询的字段必须跟数据库中的字段保持一直
	GlobalDB.Raw("select id from user where name = ?", userName).Scan(&userData)

	//查询租户的所有的订单
	if role == "custom" {
		if err := GlobalDB.Where("user_id = ?", userData.Id).Find(&orders).Error; err != nil {
			fmt.Println("获取当前用户所有订单失败")
			return nil, err
		}
	} else {
		//查询房东的订单  以房东视角来查看订单
		var houses []House
		GlobalDB.Where("user_id = ?", userData.Id).Find(&houses)

		for _, v := range houses {
			var tempOrders []OrderHouse
			GlobalDB.Where("house_id = ?", v.ID).Find(&tempOrders)

			orders = append(orders, tempOrders...)
		}
	}

	//循环遍历一下orders
	for _, v := range orders {
		var orderTemp pb.OrdersData
		orderTemp.OrderId = int32(v.ID)
		orderTemp.EndDate = v.EndDate.Format("2006-01-02")
		orderTemp.StartDate = v.BeginDate.Format("2006-01-02")
		orderTemp.Ctime = v.CreatedAt.Format("2006-01-02")
		orderTemp.Amount = int32(v.Amount)
		orderTemp.Comment = v.Comment
		orderTemp.Days = int32(v.Days)
		orderTemp.Status = v.Status

		//关联house表
		var house House
		GlobalDB.Where("id = ?", v.HouseId).Find(&house)
		orderTemp.ImgUrl = prefixUrl + house.IndexImageUrl
		orderTemp.Title = house.Title

		orderResp = append(orderResp, &orderTemp)
	}
	return orderResp, nil
}

// UpdateStatus 更新订单状态
func UpdateStatus(action, id, reason string) error {
	db := GlobalDB.Model(new(OrderHouse)).Where("id = ?", id)

	if action == "accept" {
		//表示房东同意订单
		return db.Update("status", "WAIT_COMMENT").Error
	} else {
		//表示房东不同意订单  如果拒单把拒绝的原因写到comment中
		return db.Updates(map[string]interface{}{"status": "REJECTED", "comment": reason}).Error
	}
}
