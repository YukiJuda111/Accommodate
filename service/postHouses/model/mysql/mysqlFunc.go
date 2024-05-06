package model

import "strconv"

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
