package controller

import (
	mysqlModel "RentHouse/web/model/mysql"
	"RentHouse/web/proto/orders"
	"RentHouse/web/utils"
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type OrderStu struct {
	EndDate   string `json:"end_date"`
	HouseId   string `json:"house_id"`
	StartDate string `json:"start_date"`
}

// PostOrders 下订单
func PostOrders(ctx *gin.Context) {
	// 获取数据
	var order OrderStu
	err := ctx.Bind(&order)

	// 校验数据
	if err != nil {
		fmt.Println("获取数据错误", err)
		return
	}
	// 获取用户名
	userName := sessions.Default(ctx).Get("userName")

	// 调用微服务
	consulSrv := utils.InitMicro()
	client := orders.NewOrdersService("orders", consulSrv.Client())
	request := &orders.Request{
		StartDate: order.StartDate,
		EndDate:   order.EndDate,
		HouseId:   order.HouseId,
		UserName:  userName.(string),
	}
	resp, err := client.CreateOrder(context.TODO(), request)
	// 返回数据
	ctx.JSON(http.StatusOK, resp)
}

// GetUserOrder 获取订单信息
func GetUserOrder(ctx *gin.Context) {
	// 获取get请求传参
	role := ctx.Query("role")
	// 校验数据
	if role == "" {
		fmt.Println("获取数据失败")
		return
	}

	// 处理数据  服务端
	consulSrv := utils.InitMicro()
	client := orders.NewOrdersService("orders", consulSrv.Client())
	// 获取用户名
	userName := sessions.Default(ctx).Get("userName")
	// 调用微服务
	resp, _ := client.GetOrderInfo(context.TODO(), &orders.GetReq{
		Role:     role,
		UserName: userName.(string),
	})
	//返回数据
	ctx.JSON(http.StatusOK, resp)
}

type StatusStu struct {
	Action string `json:"action"`
	Reason string `json:"reason"`
}

// PutOrders 更新订单状态
func PutOrders(ctx *gin.Context) {
	// 获取数据
	id := ctx.Param("id")
	var statusStu StatusStu
	err := ctx.Bind(&statusStu)

	// 校验数据
	if err != nil || id == "" {
		fmt.Println("获取数据错误", err)
		return
	}

	// 处理数据   更新订单状态
	consulSrv := utils.InitMicro()
	microClient := orders.NewOrdersService("orders", consulSrv.Client())
	// 调用微服务
	resp, _ := microClient.UpdateStatus(context.TODO(), &orders.UpdateReq{
		Action: statusStu.Action,
		Reason: statusStu.Reason,
		Id:     id,
	})

	// 返回数据
	ctx.JSON(http.StatusOK, resp)
}

func PutComment(ctx *gin.Context) {
	// 获取数据
	id := ctx.Param("id")
	type Comment struct {
		OrderId string `json:"order_id"`
		Comment string `json:"comment"`
	}
	var comment Comment
	err := ctx.Bind(&comment)
	if err != nil || id == "" {
		fmt.Println("获取数据错误", err)
		return
	}
	idInt, _ := strconv.Atoi(id)
	err = mysqlModel.PutComment(idInt, comment.Comment)

	if err != nil {
		fmt.Println("更新数据错误", err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"errno":  utils.RECODE_OK,
		"errmsg": utils.RecodeText(utils.RECODE_OK),
	})
}
