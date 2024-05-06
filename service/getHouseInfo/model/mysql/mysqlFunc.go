package model

import (
	"strconv"
)

// TODO: 修改prefixUrl为自己的七牛云存储空间地址
var prefixUrl = "http://scpper6fg.hd-bkt.clouddn.com/"

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
