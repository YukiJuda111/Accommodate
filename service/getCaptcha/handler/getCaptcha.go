package handler

import (
	"context"
	"encoding/json"
	"getCaptcha/model"
	"github.com/afocus/captcha"
	"go-micro.dev/v4/logger"
	"image/color"

	pb "getCaptcha/proto"
)

type GetCaptcha struct{}

func (e *GetCaptcha) Call(ctx context.Context, req *pb.CallRequest, rsp *pb.CallResponse) error {
	logger.Infof("Received GetCaptcha.Call request: %v", req)
	capt := captcha.New()
	// 可以设置多个字体 或使用cap.AddFont("xx.ttf")追加
	capt.SetFont("service/getCaptcha/conf/comic.ttf")
	// 设置验证码大小
	capt.SetSize(128, 64)
	// 设置干扰强度
	capt.SetDisturbance(captcha.MEDIUM)
	// 设置前景色 可以多个 随机替换文字颜色 默认黑色
	capt.SetFrontColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	// 设置背景色 可以多个 随机替换背景色 默认白色
	capt.SetBkgColor(color.RGBA{R: 255, A: 255}, color.RGBA{B: 255, A: 255}, color.RGBA{G: 153, A: 255})

	img, str := capt.Create(6, captcha.ALL)
	err := model.SaveImg(str, req.Uuid)
	if err != nil {
		logger.Errorf("Error saving image: %v", err)
		return err
	}
	// 将图片json序列化
	imgBuf, err := json.Marshal(img)
	if err != nil {
		logger.Errorf("Error serializing image: %v", err)
		return err
	}
	rsp.Img = imgBuf
	return nil
}
