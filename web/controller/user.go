package controller

import (
	getCaptcha "RentHouse/web/proto"
	"RentHouse/web/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/gin-gonic/gin"
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
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

	consulReg := consul.NewRegistry()
	consulSev := micro.NewService(
		micro.Registry(consulReg),
	)
	client := getCaptcha.NewGetCaptchaService("getcaptcha", consulSev.Client())
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

	fmt.Println(uuid)
}
