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
