package system

import (
	"fmt"
	"gmkdb/global"
	"gmkdb/models/common/response"
	systemRes "gmkdb/models/system/response"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

type BaseApi struct{}

// Captcha
// @Tags Base
// @Summary 生成验证码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=systemRes.SysCaptchaResponse,msg=string} "生成验证码,返回包括随机数id,base64,验证码长度"
// @Router /base/captcha [post]
func (b *BaseApi) Captcha(c *gin.Context) {
	// 字符,公式,验证码配置
	// 生成默认数字的driver
	driver := base64Captcha.NewDriverDigit(global.GIA_CONFIG.Captcha.ImgHeight, global.GIA_CONFIG.Captcha.ImgWidth, global.GIA_CONFIG.Captcha.KeyLong, 0.7, 80)
	// cp := base64Captcha.NewCaptcha(driver, store.UseWithCtx(c))   // v8下使用redis
	cp := base64Captcha.NewCaptcha(driver, store)
	if id, b64s, answer, err := cp.Generate(); err != nil {
		// global.GIA_LOG.Error("验证码获取失败!", zap.Error(err))
		log.Println("验证码获取失败!", err.Error())

		// response.FailWithMessage("验证码获取失败", c)
		response.Fail(c, "验证码获取失败")
	} else {
		// response.OkWithDetailed(systemRes.SysCaptchaResponse{
		// 	CaptchaId:     id,
		// 	PicPath:       b64s,
		// 	CaptchaLength: global.GIA_CONFIG.Captcha.KeyLong,
		// }, "验证码获取成功", c)
		fmt.Println("验证码获取成功!\n", "id:", id, "b64s:", b64s, "answer:", answer)
		response.Ok(c, systemRes.SysCaptchaResponse{
			CaptchaId:     id,
			PicPath:       b64s,
			CaptchaLength: global.GIA_CONFIG.Captcha.KeyLong,
		}, "验证码获取成功")
	}
}
