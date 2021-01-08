package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hjd919/captcha"
)

type captchaCache struct{}

// 设置验证码
func (c *captchaCache) Set(key string, randText string, expire int) {
}

// 获取验证码
func (c *captchaCache) Get(key string) (randText string) {
	return
}

// 删除验证码
func (c *captchaCache) Del(key string) {
}

func gincaptcha() {
	// 设置图片验证码字体目录
	captcha.SetFontPath("../captcha/fonts")
	r := gin.Default()
	var cache captcha.Cache = &captchaCache{}
	r.GET("/ping", func(c *gin.Context) {
		captcha.New("register-uid", cache).Output(c.Writer, 150, 50)
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func main() {
	// ExcelExport()
	// timerfc()

	gincaptcha()
}
