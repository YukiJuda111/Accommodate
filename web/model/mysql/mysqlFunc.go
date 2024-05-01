package model

// PutUserInfo 修改用户名
func PutUserInfo(userName string, newName string) error {
	return GlobalDB.Model(&User{}).Where("name = ?", userName).Update("name", newName).Error
}

// PutUserAvatar 修改用户头像
func PutUserAvatar(userName string, avatarUrl string) error {
	return GlobalDB.Model(&User{}).Where("name = ?", userName).Update("avatar_url", avatarUrl).Error
}
