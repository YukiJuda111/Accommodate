package controller

import (
	modelMysql "RentHouse/web/model/mysql"
	modelRedis "RentHouse/web/model/redis"
	"RentHouse/web/proto/getCaptcha"
	"RentHouse/web/proto/user"
	"RentHouse/web/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"image/png"
	"net/http"
)

// GetSession 获取session
func GetSession(c *gin.Context) {
	errno := utils.RECODE_SESSIONERR
	errmsg := utils.RecodeText(errno)
	c.JSON(http.StatusOK, gin.H{
		"errno":  errno,
		"errmsg": errmsg,
	})
}

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

	consulSrv := utils.InitMicro()
	client := user.NewUserService("user", consulSrv.Client())
	request := &user.SmsRequest{
		PhoneNum: phoneNum,
		ImgCode:  imgCode,
		Uuid:     imgUUID,
	}
	resp, err := client.SendSms(context.Background(), request)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"errno":  resp.Errno,
		"errmsg": utils.RecodeText(resp.Errno),
	})
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
	c.JSON(http.StatusOK, gin.H{
		"errno":  resp.Errno,
		"errmsg": utils.RecodeText(resp.Errno),
	})
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
		redisConn := modelRedis.RedisPool.Get()
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

	c.JSON(http.StatusOK, gin.H{
		"errno":  utils.RECODE_OK,
		"errmsg": utils.RecodeText(utils.RECODE_OK),
		"data":   areas,
	})

}
