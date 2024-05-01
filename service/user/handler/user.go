package handler

import (
	"context"
	modelMysql "user/model/mysql"
	modelRedis "user/model/redis"
	"user/proto"
	"user/third_party"
	"user/utils"
)

// TODO: 修改prefixUrl为自己的七牛云存储空间地址
var prefixUrl = "http://scpper6fg.hd-bkt.clouddn.com/"

type User struct{}

func (e *User) SendSms(ctx context.Context, req *user.SmsRequest, rsp *user.SmsResponse) error {
	// 校验图片验证码
	if !modelRedis.CheckImgCode(req.Uuid, req.ImgCode) {
		rsp.Errno = utils.RECODE_CAPTCHAERR
		return nil
	}

	smsCode, err := third_party.SendSms(req.PhoneNum)
	if err != nil {
		rsp.Errno = utils.RECODE_SMSERR
		return nil
	}

	err = modelRedis.SaveSmsCode(req.PhoneNum, smsCode)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		return nil
	}

	rsp.Errno = utils.RECODE_OK
	return nil
}

func (e *User) Register(ctx context.Context, req *user.RegisterRequest, rsp *user.RegisterResponse) error {
	// 校验短信验证码
	if !modelRedis.CheckSmsCode(req.PhoneNum, req.SmsCode) {
		rsp.Errno = utils.RECODE_DBERR
		return nil
	}
	// 注册用户
	err := modelMysql.RegisterUser(req.PhoneNum, req.Password)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		return err
	}

	rsp.Errno = utils.RECODE_OK
	return nil
}

func (e *User) GetUserInfo(ctx context.Context, req *user.UserInfoRequest, rsp *user.UserInfoResponse) error {
	// 从mysql获取用户信息
	userInfo, err := modelMysql.GetUserInfo(req.Name)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		return err
	}

	userInfo.AvatarUrl = prefixUrl + userInfo.AvatarUrl

	var data user.UserData
	data.UserId = int32(userInfo.ID)
	data.Name = userInfo.Name
	data.Mobile = userInfo.Mobile
	data.RealName = userInfo.RealName
	data.IdCard = userInfo.IdCard
	data.AvatarUrl = userInfo.AvatarUrl

	rsp.Errno = utils.RECODE_OK
	rsp.UserData = &data
	return nil

}

func (e *User) Login(ctx context.Context, req *user.LoginRequest, rsp *user.LoginResponse) error {

	userName, err := modelMysql.Login(req.Mobile, req.Password)
	if err != nil {
		rsp.Errno = utils.RECODE_DATAERR
		return nil
	}

	rsp.Errno = utils.RECODE_OK
	rsp.Name = userName
	return nil
}
