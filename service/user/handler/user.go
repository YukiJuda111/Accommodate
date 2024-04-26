package handler

import (
	"context"

	"user/model"
	"user/proto"
	"user/third_party"
	"user/utils"
)

type User struct{}

func (e *User) SendSms(ctx context.Context, req *user.SmsRequest, rsp *user.SmsResponse) error {
	// 校验图片验证码
	if !model.CheckImgCode(req.Uuid, req.ImgCode) {
		rsp.Errno = utils.RECODE_CAPTCHAERR
		return nil
	}

	smsCode, err := third_party.SendSms(req.PhoneNum)
	if err != nil {
		rsp.Errno = utils.RECODE_SMSERR
		return nil
	}

	err = model.SaveSmsCode(req.PhoneNum, smsCode)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		return nil
	}

	rsp.Errno = utils.RECODE_OK
	return nil
}
