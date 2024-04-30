package model

// Login 处理登陆业务,根据手机/密码获取用户名
func Login(phoneNum string, password string) (string, error) {
	user := User{
		Mobile:       phoneNum,
		PasswordHash: password,
	}

	err := GlobalDB.Find(&user).Error
	return user.Name, err
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

// PutUserInfo 修改用户名
func PutUserInfo(userName string, newName string) error {
	return GlobalDB.Model(&User{}).Where("name = ?", userName).Update("name", newName).Error
}

// PutUserAvatar 修改用户头像
func PutUserAvatar(userName string, avatarUrl string) error {
	return GlobalDB.Model(&User{}).Where("name = ?", userName).Update("avatar_url", avatarUrl).Error
}
