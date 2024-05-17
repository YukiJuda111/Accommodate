package model

func RegisterUser(phoneNum string, password string) error {

	user := User{
		Name:         phoneNum,
		PasswordHash: password,
		Mobile:       phoneNum,
	}

	return GlobalDB.Create(&user).Error
}

// GetUserInfo 获取用户信息
func GetUserInfo(userName string) (User, error) {
	user := User{
		Name: userName,
	}
	// 查询用户信息,userName是唯一的,查First加快效率
	err := GlobalDB.First(&user).Error
	return user, err
}

// Login 处理登陆业务,根据手机/密码获取用户名
func Login(phoneNum string, password string) (string, error) {
	var user User
	err := GlobalDB.Where("mobile = ? and password_hash = ?", phoneNum, password).First(&user).Error
	return user.Name, err
}

// PutUserInfo 修改用户名
func PutUserInfo(userName string, newName string) error {
	return GlobalDB.Model(&User{}).Where("name = ?", userName).Update("name", newName).Error
}

// GetHouses 获取用户发布的房源
func GetHouses(userName string) ([]House, User, []Area, error) {
	//查询当前用户对象
	var user User
	err := GlobalDB.Where("name = ?", userName).First(&user).Error
	if err != nil {
		return nil, User{}, nil, err
	}

	//查询当前用户对应的房屋
	var houses []House
	err = GlobalDB.Where("user_id = ?", user.ID).Find(&houses).Error
	if err != nil {
		return nil, User{}, nil, err
	}

	//获取各个房源对应的地址
	var areas []Area
	for _, house := range houses {
		var area Area
		err = GlobalDB.Where("id = ?", house.AreaId).First(&area).Error
		areas = append(areas, area)
	}
	return houses, user, areas, nil
}

// SaveAuth 保存用户实名认证信息
func SaveAuth(userName string, realName string, idCard string) error {
	return GlobalDB.Model(&User{}).Where("name = ?", userName).Updates(User{RealName: realName, IdCard: idCard}).Error
}
