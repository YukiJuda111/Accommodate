package controller

import (
	modelMysql "RentHouse/web/model/mysql"
	"RentHouse/web/proto/getHouseInfo"
	"RentHouse/web/proto/postHouses"
	"RentHouse/web/third_party"
	"RentHouse/web/utils"
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
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

// GetHouseInfo 获取房屋信息
func GetHouseInfo(c *gin.Context) {
	houseId, _ := strconv.Atoi(c.Param("id"))

	// 初始化微服务
	consulSrv := utils.InitMicro()
	client := getHouseInfo.NewGetHouseInfoService("gethouseinfo", consulSrv.Client())
	request := &getHouseInfo.CallRequest{
		HouseId: int32(houseId),
	}
	// 调用微服务
	resp, err := client.Call(context.Background(), request)
	if err != nil {
		utils.ResponseData(c, resp.Errno, nil)
		return
	}
	utils.ResponseData(c, resp.Errno, resp.Data)
}

// GetIndex 获取首页房屋信息
func GetIndex(c *gin.Context) {
	// 初始化微服务
	consulSrv := utils.InitMicro()
	client := getHouseInfo.NewGetHouseInfoService("gethouseinfo", consulSrv.Client())
	request := &getHouseInfo.IndexRequest{}
	// 调用微服务
	resp, err := client.GetHouseIndex(context.Background(), request)
	if err != nil {
		utils.ResponseData(c, resp.Errno, nil)
		return
	}
	utils.ResponseData(c, resp.Errno, resp.Data)
}

// SearchHouses 获取房源信息
func SearchHouses(c *gin.Context) {
	//获取数据
	// areaId
	aid := c.Query("aid")
	// start day
	sd := c.Query("sd")
	//end day
	ed := c.Query("ed")
	// 排序方式
	sk := c.Query("sk")
	// page  第几页
	// 校验数据
	if aid == "" || sd == "" || ed == "" || sk == "" {
		fmt.Println("传入数据不完整")
		return
	}

	// 初始化微服务
	consulSrv := utils.InitMicro()
	client := getHouseInfo.NewGetHouseInfoService("gethouseinfo", consulSrv.Client())
	request := &getHouseInfo.SearchRequest{
		Aid: aid,
		Sd:  sd,
		Ed:  ed,
		Sk:  sk,
	}
	// 调用微服务
	resp, err := client.SearchHouse(context.Background(), request)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, resp)

}
