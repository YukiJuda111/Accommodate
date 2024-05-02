package model

// PutUserAvatar 修改用户头像
func PutUserAvatar(userName string, avatarUrl string) error {
	return GlobalDB.Model(&User{}).Where("name = ?", userName).Update("avatar_url", avatarUrl).Error
}
