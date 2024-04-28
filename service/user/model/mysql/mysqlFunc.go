package model

func RegisterUser(phoneNum string, password string) error {

	user := User{
		Name:         phoneNum,
		PasswordHash: password,
		Mobile:       phoneNum,
	}

	return GlobalDB.Create(&user).Error
}
