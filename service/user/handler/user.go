package handler

import (
	"context"
	modelMysql "user/model/mysql"
	modelRedis "user/model/redis"
	pb "user/proto"
	"user/third_party"
	"user/utils"
)

// TODO: 修改prefixUrl为自己的七牛云存储空间地址
var prefixUrl = "http://scpper6fg.hd-bkt.clouddn.com/"

type User struct{}

func (e *User) SendSms(ctx context.Context, req *pb.SmsRequest, rsp *pb.SmsResponse) error {
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

func (e *User) Register(ctx context.Context, req *pb.RegisterRequest, rsp *pb.RegisterResponse) error {
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

func (e *User) GetUserInfo(ctx context.Context, req *pb.UserInfoRequest, rsp *pb.UserInfoResponse) error {
	// 从mysql获取用户信息
	userInfo, err := modelMysql.GetUserInfo(req.Name)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		return err
	}

	userInfo.AvatarUrl = prefixUrl + userInfo.AvatarUrl

	var data pb.UserData
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

func (e *User) Login(ctx context.Context, req *pb.LoginRequest, rsp *pb.LoginResponse) error {

	userName, err := modelMysql.Login(req.Mobile, req.Password)
	if err != nil {
		rsp.Errno = utils.RECODE_DATAERR
		return nil
	}
	rsp.Errno = utils.RECODE_OK
	rsp.Name = userName
	return nil
}

func (e *User) PutUserInfo(ctx context.Context, req *pb.PutUserRequest, rsp *pb.PutUserResponse) error {
	// 修改用户名
	err := modelMysql.PutUserInfo(req.PrevName, req.Name)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		return err
	}
	rsp.Errno = utils.RECODE_OK
	return nil
}

func (e *User) PutAuth(ctx context.Context, req *pb.AuthRequest, rsp *pb.AuthResponse) error {
	err := modelMysql.SaveAuth(req.UserName, req.RealName, req.IdCard)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		return err
	}
	rsp.Errno = utils.RECODE_OK
	return nil
}

func (e *User) GetHouse(ctx context.Context, req *pb.GetHouseRequest, rsp *pb.GetHouseResponse) error {
	// 从mysql获取用户发布的房源
	houses, user, areas, err := modelMysql.GetHouses(req.UserName)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		return err
	}

	//根据返回数据，构造返回值
	var houseInfos pb.Houses
	for k, v := range houses {
		var house pb.HouseInfo

		house.Address = v.Address
		house.HouseId = int32(v.ID)
		house.Ctime = v.CreatedAt.Format("2006-01-02 15:04:05")
		house.ImgUrl = prefixUrl + v.IndexImageUrl
		house.OrderCount = int32(v.OrderCount)
		house.Price = int32(v.Price)
		house.RoomCount = int32(v.RoomCount)
		house.Title = v.Title
		house.UserAvatar = prefixUrl + user.AvatarUrl
		house.AreaName = areas[k].Name

		houseInfos.Houses = append(houseInfos.Houses, &house)
	}

	//返回数据
	rsp.Errno = utils.RECODE_OK
	rsp.HouseData = &houseInfos
	return nil
}
