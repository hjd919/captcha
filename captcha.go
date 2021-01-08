package captcha

// https://gitee.com/longfei6671/gocaptcha

import (
	"fmt"
	"net/http"

	"github.com/lifei6671/gocaptcha"
)

const (
	dx = 150
	dy = 50
)

// example
// func gincaptcha() {
// 	// 设置图片验证码字体目录
// 	captcha.SetFontPath("../captcha/fonts")
// 	r := gin.Default()
// 	var cache captcha.Cache = &captchaCache{}
// 	r.GET("/ping", func(c *gin.Context) {
// 		captcha.New("register-uid", cache).Output(c.Writer, 150, 50)
// 	})
// 	r.Run() // listen and serve on 0.0.0.0:8080
// }

type RandTextType int

const (
	RandTextTypeNum RandTextType = iota
)

func SetFontPath(fontPath string) {
	err := gocaptcha.ReadFonts(fontPath, ".ttf")
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 缓存接口，需要实现以下方法
type Cache interface {
	// 设置验证码
	Set(key string, randText string, expire int)
	// 获取验证码
	Get(key string) (randText string)
	// 删除验证码
	Del(key string)
}

// 图片验证码类  缓存 校验 key
type Captcha struct {
	key   string // 项目-业务-uuid：xxx-register-uuid
	cache Cache  // 缓存
}

// 创建图片验证码对象
func New(key string, cache Cache) *Captcha {
	return &Captcha{
		key,
		cache,
	}
}

// 输出带验证码图片
func (c *Captcha) Output(w http.ResponseWriter, width, height int) {
	// 生成验证码
	randText := c.createRandText()
	// 把验证码存到缓存
	c.saveRandTextCache(randText)
	// 生成图片验证码
	captchaImage := gocaptcha.NewCaptchaImage(width, height, gocaptcha.RandLightColor())
	err := captchaImage.
		DrawText(randText).                           // 图片验证码
		DrawNoise(gocaptcha.CaptchaComplexLower).     //画随机噪点
		DrawTextNoise(gocaptcha.CaptchaComplexLower). //画随机文字噪点
		DrawBorder(gocaptcha.ColorToRGB(0x17A7A7A)).  //画边框
		// DrawSineLine().                               //画上三条随机直线
		Error

	if err != nil {
		fmt.Println(err)
	}
	_ = captchaImage.SaveImage(w, gocaptcha.ImageFormatJpeg)
	//将验证码保持到输出流种，可以是文件或HTTP流等
	captchaImage.SaveImage(w, gocaptcha.ImageFormatJpeg)
}

// 校验缓存验证码 key 用户输入的验证码 是否正确
func (c *Captcha) CheckRandText(code string) (isRight bool) {
	randText := c.cache.Get(c.key)
	if randText == "" {
		return
	}
	if randText != code {
		return
	}
	// 已使用，删除
	c.cache.Del(c.key)
	return true
}

// 把验证码存到缓存
func (c *Captcha) saveRandTextCache(randText string) {
	c.cache.Set(c.key, randText, 60) // 缓存60秒
	return
}

// 生成验证码，4位数字 TODO
func (c *Captcha) createRandText() (randText string) {
	return gocaptcha.RandText(4)
}
