package model

// PutUserAvatar 修改用户头像
func PutUserAvatar(userName string, avatarUrl string) error {
	return GlobalDB.Model(&User{}).Where("name = ?", userName).Update("avatar_url", avatarUrl).Error
}

// PostHouseImage 上传房屋图片
func PostHouseImage(houseId uint, url string) error {
	var house House

	//查询对应 房屋信息
	err := GlobalDB.First(&house, houseId).Error
	if err != nil {
		return err
	}

	if house.IndexImageUrl == "" {
		err = GlobalDB.Model(&house).Update("index_image_url", url).Error
		if err != nil {
			return err
		}
	} else {
		var houseImage HouseImage
		houseImage.HouseId = houseId
		houseImage.Url = url
		return GlobalDB.Create(&houseImage).Error
	}
	return nil
}

func PutComment(orderId int, comment string) error {
	return GlobalDB.Model(&OrderHouse{}).Where("id = ?", orderId).Update("comment", comment).Error
}

func FindHouseIdByOrderId(orderId int) (int, error) {
	var order OrderHouse
	err := GlobalDB.First(&order, orderId).Error
	ret := (int)(order.HouseId)
	return ret, err
}
