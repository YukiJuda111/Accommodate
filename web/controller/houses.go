package controller

import (
	modelMysql "RentHouse/web/model/mysql"
	"RentHouse/web/proto/postHouses"
	"RentHouse/web/third_party"
	"RentHouse/web/utils"
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
)

// PostHouses 发布房源
func PostHouses(c *gin.Context) {
	//获取前端传递过来的数据
	type HouseInfo struct {
		Title     string   `json:"title"`
		Acreage   string   `json:"acreage"`
		Address   string   `json:"address"`
		AreaId    string   `json:"area_id"`
		Beds      string   `json:"beds"`
		Capacity  string   `json:"capacity"`
		Deposit   string   `json:"deposit"`
		Facility  []string `json:"facility"`
		MaxDays   string   `json:"max_days"`
		MinDays   string   `json:"min_days"`
		Price     string   `json:"price"`
		RoomCount string   `json:"room_count"`
		Unit      string   `json:"unit"`
	}
	var houseInfo HouseInfo
	err := c.Bind(&houseInfo)
	if err != nil {
		utils.ResponseData(c, utils.RECODE_NODATA, nil)
		return
	}

	// 获取session
	session := sessions.Default(c)
	userName := session.Get("userName")

	// 微服务初始化
	consulSrv := utils.InitMicro()
	client := postHouses.NewPostHousesService("posthouses", consulSrv.Client())
	request := &postHouses.CallRequest{
		Name:      userName.(string),
		Acreage:   houseInfo.Acreage,
		Address:   houseInfo.Address,
		AreaId:    houseInfo.AreaId,
		Beds:      houseInfo.Beds,
		Capacity:  houseInfo.Capacity,
		Deposit:   houseInfo.Deposit,
		Facility:  houseInfo.Facility,
		MaxDays:   houseInfo.MaxDays,
		MinDays:   houseInfo.MinDays,
		Price:     houseInfo.Price,
		RoomCount: houseInfo.RoomCount,
		Title:     houseInfo.Title,
		Unit:      houseInfo.Unit,
	}
	// 调用微服务
	resp, err := client.Call(context.Background(), request)
	if err != nil {
		utils.ResponseData(c, resp.Errno, nil)
		return
	}
	utils.ResponseData(c, resp.Errno, resp.Data)
}

// PostHousesImage 上传房屋图片
func PostHousesImage(c *gin.Context) {
	// 得到houseId
	houseId, _ := strconv.Atoi(c.Param("id"))
	// 获取房屋图片文件
	file, err := c.FormFile("house_image")
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

	// 更新mysql中的房屋图片url
	err = modelMysql.PostHouseImage(uint(houseId), url)
	if err != nil {
		utils.ResponseData(c, utils.RECODE_DATAERR, nil)
		return

	}

	url = prefixUrl + url

	utils.ResponseData(c, utils.RECODE_OK, map[string]interface{}{"url": url})
}
