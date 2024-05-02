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
	user := User{
		Mobile:       phoneNum,
		PasswordHash: password,
	}

	err := GlobalDB.Find(&user).Error
	return user.Name, err
}

// PutUserInfo 修改用户名
func PutUserInfo(userName string, newName string) error {
	return GlobalDB.Model(&User{}).Where("name = ?", userName).Update("name", newName).Error
}
