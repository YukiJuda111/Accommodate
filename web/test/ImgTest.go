package main

import (
	"github.com/afocus/captcha"
	"image/color"
	"image/png"
	"net/http"
)

func main() {
	capt := captcha.New()
	// 可以设置多个字体 或使用cap.AddFont("xx.ttf")追加
	capt.SetFont("../conf/comic.ttf")
	// 设置验证码大小
	capt.SetSize(128, 64)
	// 设置干扰强度
	capt.SetDisturbance(captcha.MEDIUM)
	// 设置前景色 可以多个 随机替换文字颜色 默认黑色
	capt.SetFrontColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	// 设置背景色 可以多个 随机替换背景色 默认白色
	capt.SetBkgColor(color.RGBA{R: 255, A: 255}, color.RGBA{B: 255, A: 255}, color.RGBA{G: 153, A: 255})
	// 生成字体
	http.HandleFunc("/r", func(w http.ResponseWriter, r *http.Request) {
		img, str := capt.Create(6, captcha.ALL)
		err := png.Encode(w, img)
		if err != nil {
			return
		}
		println(str)
	})
	err := http.ListenAndServe(":8085", nil)
	if err != nil {
		return
	}
}
