package model

func Login(phoneNum string, password string) (string, error) {
	user := User{
		Mobile:       phoneNum,
		PasswordHash: password,
	}

	err := GlobalDB.Find(&user).Error
	return user.Name, err
}
