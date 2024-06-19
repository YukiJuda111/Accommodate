package model

import (
	"fmt"
	pb "getHouseInfo/proto"
	"strconv"
	"time"
)

// TODO: 修改prefixUrl为自己的七牛云存储空间地址
var prefixUrl = "http://sf9x6ixky.hd-bkt.clouddn.com/"

func InsertHouse(name string, house House, fids []string) (uint, error) {
	//根据用户名获取用户Id
	var user User
	user.Name = name
	err := GlobalDB.Where("name = ?", name).Find(&user).Error
	if err != nil {
		return 0, err
	}
	//关联当前用户
	house.UserId = uint(user.ID)
	//查询所有家具放到house里面
	for _, fid := range fids {
		id, _ := strconv.Atoi(fid)
		var fac Facility
		err = GlobalDB.Where("id = ?", id).Find(&fac).Error
		if err != nil {
			return 0, err
		}
		house.Facilities = append(house.Facilities, &fac)
	}
	return house.ID, GlobalDB.Create(&house).Error
}

func GetDetail(houseId int32) (House, []OrderHouse, []User, User, []Facility, []HouseImage, error) {
	var house House
	var orders []OrderHouse
	var orderUsers []User
	var user User
	var facs []Facility
	var imgs []HouseImage

	house.ID = uint(houseId)
	err := GlobalDB.Where(&house).First(&house).Error
	if err != nil {
		return house, orders, orderUsers, user, facs, imgs, err
	}

	//获取评论信息
	err = GlobalDB.Where("house_id = ?", houseId).Find(&orders).Error
	if err != nil {
		return house, orders, orderUsers, user, facs, imgs, err
	}

	//获取评论用户信息
	for _, v := range orders {
		var orderUser User
		err = GlobalDB.Where("id = ?", v.UserId).Find(&orderUser).Error
		if err != nil {
			return house, orders, orderUsers, user, facs, imgs, err
		}
		orderUsers = append(orderUsers, orderUser)
	}

	//获取房东信息
	err = GlobalDB.Where("id = ?", house.UserId).Find(&user).Error
	if err != nil {
		return house, orders, orderUsers, user, facs, imgs, err
	}

	//获取家具信息
	var facids []uint
	// 查找房屋下的所有家具id
	err = GlobalDB.Raw("select facility_id from house_facilities where house_id = ?", houseId).Scan(&facids).Error
	if err != nil {
		return house, orders, orderUsers, user, facs, imgs, err
	}

	err = GlobalDB.Where("id in (?)", facids).Find(&facs).Error
	if err != nil {
		return house, orders, orderUsers, user, facs, imgs, err
	}

	err = GlobalDB.Where("id in (?)", facids).Find(&facs).Error
	if err != nil {
		return house, orders, orderUsers, user, facs, imgs, err
	}

	//获取房屋图片
	err = GlobalDB.Where("house_id = ?", houseId).Find(&imgs).Error
	if err != nil {
		return house, orders, orderUsers, user, facs, imgs, err
	}

	house.IndexImageUrl = prefixUrl + house.IndexImageUrl
	user.AvatarUrl = prefixUrl + user.AvatarUrl
	for _, v := range imgs {
		v.Url = prefixUrl + v.Url
	}

	return house, orders, orderUsers, user, facs, imgs, nil
}

// GetIndexHouse 获取首页展示的房屋信息
func GetIndexHouse() ([]*pb.Houses, error) {

	var housesResp []*pb.Houses

	var houses []House
	if err := GlobalDB.Limit(5).Find(&houses).Error; err != nil {
		fmt.Println("获取房屋信息失败", err)
		return nil, err
	}

	for _, house := range houses {
		var houseTemp pb.Houses
		houseTemp.Address = house.Address
		//根据房屋信息获取地域信息
		var area Area
		var user User

		GlobalDB.Where("id= ?", house.AreaId).Find(&area)
		GlobalDB.Where("id= ?", house.UserId).Find(&user)

		houseTemp.AreaName = area.Name
		houseTemp.Ctime = house.CreatedAt.Format("2006-01-02 15:04:05")
		houseTemp.HouseId = int32(house.ID)
		houseTemp.ImgUrl = prefixUrl + house.IndexImageUrl
		houseTemp.OrderCount = int32(house.OrderCount)
		houseTemp.Price = int32(house.Price)
		houseTemp.RoomCount = int32(house.RoomCount)
		houseTemp.Title = house.Title
		houseTemp.UserAvatar = prefixUrl + user.AvatarUrl

		housesResp = append(housesResp, &houseTemp)
	}

	return housesResp, nil
}

// SearchHouse 获取搜索的房屋信息
func SearchHouse(areaId, sd, ed, sk string) ([]*pb.Houses, error) {
	var houseInfos []House

	//   minDays  <  (结束时间  -  开始时间) <  max_days
	//计算一个差值  先把string类型转为time类型
	sdTime, _ := time.Parse("2006-01-02", sd)
	edTime, _ := time.Parse("2006-01-02", ed)
	dur := edTime.Sub(sdTime)

	err := GlobalDB.Where("area_id = ?", areaId).
		Where("min_days < ?", dur.Hours()/24).
		Where("max_days > ?", dur.Hours()/24).
		Order("created_at desc").Find(&houseInfos).Error
	if err != nil {
		fmt.Println("搜索房屋失败", err)
		return nil, err
	}

	//获取[]*house.Houses
	var housesResp []*pb.Houses

	for _, v := range houseInfos {
		var houseTemp pb.Houses
		houseTemp.Address = v.Address
		//根据房屋信息获取地域信息
		var area Area
		var user User

		GlobalDB.Where("id = ?", v.AreaId).Find(&area)
		GlobalDB.Where("id = ?", v.UserId).Find(&user)

		houseTemp.AreaName = area.Name
		houseTemp.Ctime = v.CreatedAt.Format("2006-01-02 15:04:05")
		houseTemp.HouseId = int32(v.ID)
		houseTemp.ImgUrl = prefixUrl + v.IndexImageUrl
		houseTemp.OrderCount = int32(v.OrderCount)
		houseTemp.Price = int32(v.Price)
		houseTemp.RoomCount = int32(v.RoomCount)
		houseTemp.Title = v.Title
		houseTemp.UserAvatar = prefixUrl + user.AvatarUrl

		housesResp = append(housesResp, &houseTemp)

	}

	return housesResp, nil
}
