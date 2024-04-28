package handler

import (
	"context"
	modelMysql "user/model/mysql"
	modelRedis "user/model/redis"
	"user/proto"
	"user/third_party"
	"user/utils"
)

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
