package handler

import (
	"context"
	"fmt"
	mysqlModel "orders/model/mysql"
	pb "orders/proto"
	"orders/utils"
	"strconv"
)

type Orders struct{}

func (e *Orders) CreateOrder(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	//获取到相关数据,插入到数据库
	orderId, err := mysqlModel.InsertOrder(req.HouseId, req.StartDate, req.EndDate, req.UserName)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	var orderData pb.OrderData
	orderData.OrderId = strconv.Itoa(orderId)

	rsp.Data = &orderData

	return nil
}

func (e *Orders) GetOrderInfo(ctx context.Context, req *pb.GetReq, resp *pb.GetResp) error {
	//要根据传入数据获取订单信息   mysql
	respData, err := mysqlModel.GetOrderInfo(req.UserName, req.Role)
	if err != nil {
		resp.Errno = utils.RECODE_DATAERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	var getData pb.GetData
	getData.Orders = respData
	resp.Data = &getData

	return nil
}

func (e *Orders) UpdateStatus(ctx context.Context, req *pb.UpdateReq, resp *pb.UpdateResp) error {
	//根据传入数据,更新订单状态
	err := mysqlModel.UpdateStatus(req.Action, req.Id, req.Reason)
	if err != nil {
		fmt.Println("更新订单装填错误", err)
		resp.Errno = utils.RECODE_DATAERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	return nil
}
