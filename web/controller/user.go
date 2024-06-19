package controller

import (
	modelMysql "RentHouse/web/model/mysql"
	modelRedis "RentHouse/web/model/redis"
	"RentHouse/web/proto/getArea"
	"RentHouse/web/proto/getCaptcha"
	"RentHouse/web/proto/user"
	"RentHouse/web/third_party"
	"RentHouse/web/utils"
	"context"
	"encoding/json"
	"fmt"
	"image/png"

	"github.com/afocus/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

// TODO: 修改prefixUrl为自己的七牛云存储空间地址
var prefixUrl = "http://sf9x6ixky.hd-bkt.clouddn.com/"

// GetImageCd 获取图片验证码
func GetImageCd(c *gin.Context) {
	uuid := c.Param("uuid")
	consulSrv := utils.InitMicro()
	client := getCaptcha.NewGetCaptchaService("getcaptcha", consulSrv.Client())
	request := &getCaptcha.CallRequest{
		Uuid: uuid,
	}
	resp, err := client.Call(context.Background(), request)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 将得到的数据反序列化
	var img captcha.Image
	err = json.Unmarshal(resp.Img, &img)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 将图片发送给前端
	err = png.Encode(c.Writer, img)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// GetSmscd 获取短信验证码
func GetSmscd(c *gin.Context) {
	phoneNum := c.Param("phonenum")
	imgCode := c.Query("text")
	imgUUID := c.Query("id")
	// 微服务初始化
	consulSrv := utils.InitMicro()
	client := user.NewUserService("user", consulSrv.Client())
	request := &user.SmsRequest{
		PhoneNum: phoneNum,
		ImgCode:  imgCode,
		Uuid:     imgUUID,
	}
	// 调用微服务
	resp, err := client.SendSms(context.Background(), request)
	if err != nil {
		fmt.Println(err)
		return
	}

	utils.ResponseData(c, resp.Errno, nil)
}

// PostRet 发送注册信息
func PostRet(c *gin.Context) {
	// 获取在请求荷载中的数据(随着POST请求发送的数据)
	var regData struct {
		Mobile   string `json:"mobile"`
		Password string `json:"password"`
		SmsCode  string `json:"sms_code"`
	}
	err := c.Bind(&regData)
	if err != nil {
		fmt.Println(err)
		return
	}

	consulSrv := utils.InitMicro()
	client := user.NewUserService("user", consulSrv.Client())
	request := &user.RegisterRequest{
		PhoneNum: regData.Mobile,
		Password: regData.Password,
		SmsCode:  regData.SmsCode,
	}
	resp, err := client.Register(context.Background(), request)
	if err != nil {
		fmt.Println(err)
		return
	}

	utils.ResponseData(c, resp.Errno, nil)
}

// GetArea 获取地区信息
func GetArea(c *gin.Context) {
	// 先从redis获取数据
	redisConn := modelRedis.RedisPool.Get()
	areaData, _ := redis.Bytes(redisConn.Do("GET", "areaData"))
	var areas []modelMysql.Area
	if len(areaData) == 0 { // redis中没有数据
		// 从mysql获取数据
		modelMysql.GlobalDB.Find(&areas)
		// 把数据写入redis
		areasJson, _ := json.Marshal(areas)
		_, err := redisConn.Do("SET", "areaData", areasJson)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else { // redis中有数据
		err := json.Unmarshal(areaData, &areas)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// 微服务初始化
	consulSrv := utils.InitMicro()
	client := getArea.NewGetAreaService("getarea", consulSrv.Client())
	request := &getArea.CallRequest{}
	// 调用微服务
	resp, err := client.Call(context.Background(), request)
	if err != nil {
		fmt.Println(err)
		return
	}

	utils.ResponseData(c, resp.Errno, resp.Data)

}

// PostLogin 登录
func PostLogin(c *gin.Context) {
	// 获取在请求荷载中的数据(随着POST请求发送的数据)
	var loginData struct {
		Mobile   string `json:"mobile"`
		Password string `json:"password"`
	}
	err := c.Bind(&loginData)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 微服务初始化
	consulSrv := utils.InitMicro()
	client := user.NewUserService("user", consulSrv.Client())
	request := &user.LoginRequest{
		Mobile:   loginData.Mobile,
		Password: loginData.Password,
	}
	// 调用微服务
	resp, err := client.Login(context.Background(), request)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 登陆成功后设置session
	session := sessions.Default(c)
	session.Set("userName", resp.Name)
	err = session.Save()
	if err != nil {
		fmt.Println("session save failed: ", err)
		return
	}

	utils.ResponseData(c, resp.Errno, nil)
}

// GetSession 获取session
func GetSession(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("userName")

	utils.ResponseData(c, utils.RECODE_OK, map[string]interface{}{"name": userName})
}

// DeleteSession 退出登录
func DeleteSession(c *gin.Context) {
	session := sessions.Default(c)
	// 删除session
	session.Delete("userName")
	err := session.Save()
	if err != nil {
		fmt.Println(err)
		return
	}

	utils.ResponseData(c, utils.RECODE_OK, nil)
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	// 获取session,得到用户名
	session := sessions.Default(c)
	userName := session.Get("userName")

	// 微服务初始化
	consulSrv := utils.InitMicro()
	client := user.NewUserService("user", consulSrv.Client())
	request := &user.UserInfoRequest{
		Name: userName.(string),
	}
	// 调用微服务
	resp, err := client.GetUserInfo(context.Background(), request)
	if err != nil {
		fmt.Println(err)
		return
	}

	utils.ResponseData(c, resp.Errno, resp.UserData)
}

// PutUserInfo 修改用户名
func PutUserInfo(c *gin.Context) {
	// 获取session,得到用户名
	session := sessions.Default(c)
	userName := session.Get("userName")

	// 获取在请求荷载中的数据(随着POST请求发送的数据)
	var putData struct {
		Name string `json:"name"`
	}
	err := c.Bind(&putData)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 微服务初始化
	consulSrv := utils.InitMicro()
	client := user.NewUserService("user", consulSrv.Client())
	request := &user.PutUserRequest{
		PrevName: userName.(string),
		Name:     putData.Name,
	}
	// 调用微服务
	resp, err := client.PutUserInfo(context.Background(), request)
	if err != nil {
		utils.ResponseData(c, resp.Errno, nil)
		return
	}

	// 更新session中的用户名
	session.Set("userName", putData.Name)
	err = session.Save()
	if err != nil {
		utils.ResponseData(c, utils.RECODE_SESSIONERR, nil)
		return
	}

	utils.ResponseData(c, resp.Errno, map[string]interface{}{"name": putData.Name})
}

// PostAvatar 上传头像
func PostAvatar(c *gin.Context) {
	// 获取session,得到用户名
	session := sessions.Default(c)
	userName := session.Get("userName")

	// 获取头像文件
	file, err := c.FormFile("avatar")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 调用七牛云上传文件
	url, err := third_party.UploadFile(file)
	if err != nil {
		utils.ResponseData(c, utils.RECODE_THIRDERR, nil)
		return
	}

	// 更新mysql中的头像url
	err = modelMysql.PutUserAvatar(userName.(string), url)
	if err != nil {
		utils.ResponseData(c, utils.RECODE_DATAERR, nil)
		return

	}

	url = prefixUrl + url

	utils.ResponseData(c, utils.RECODE_OK, map[string]interface{}{"avatar_url": url})
}

// PutUserAuth 实名认证
func PutUserAuth(c *gin.Context) {
	var authData struct {
		RealName string `json:"real_name"`
		IdCard   string `json:"id_card"`
	}
	err := c.Bind(&authData)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取session,得到用户名
	session := sessions.Default(c)
	userName := session.Get("userName")

	// 微服务初始化
	consulSrv := utils.InitMicro()
	client := user.NewUserService("user", consulSrv.Client())
	request := &user.AuthRequest{
		UserName: userName.(string),
		RealName: authData.RealName,
		IdCard:   authData.IdCard,
	}
	// 调用微服务
	resp, _ := client.PutAuth(context.Background(), request)

	utils.ResponseData(c, resp.Errno, nil)
}

// GetHouses 获取用户发布的房源
func GetHouses(c *gin.Context) {
	// 获取session,得到用户名
	session := sessions.Default(c)
	userName := session.Get("userName")

	// 微服务初始化
	consulSrv := utils.InitMicro()
	client := user.NewUserService("user", consulSrv.Client())
	request := &user.GetHouseRequest{
		UserName: userName.(string),
	}
	// 调用微服务
	resp, err := client.GetHouse(context.Background(), request)
	if err != nil {
		fmt.Println(err)
		return
	}
	utils.ResponseData(c, resp.Errno, resp.HouseData)
}
