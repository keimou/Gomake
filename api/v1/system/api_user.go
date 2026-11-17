package system

import (
	"gmkdb/global"

	"gmkdb/models/common/response"
	systemReq "gmkdb/models/system/request"
	systemRes "gmkdb/models/system/response"

	"gmkdb/models/system"
	"gmkdb/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// @Tags Base
// @Summary 用户登录
// @Produce  application/json
// @Param data body systemReq.Login true "用户名, 密码, 验证码"
// @Success 200 {object} response.Response{data=systemRes.LoginResponse,msg=string} "返回包括用户信息,token,过期时间"
// @Router /base/login [post]
func (b *BaseApi) Login(c *gin.Context) {
	var l systemReq.Login
	_ = c.ShouldBindJSON(&l)
	if err := utils.Verify(l, utils.LoginVerify); err != nil {
		// response.FailWithMessage(err.Error(), c)
		response.Fail(c, err)
		return
	}
	if store.Verify(l.CaptchaId, l.Captcha, true) {
		u := &system.SysUser{Username: l.Username, Password: l.Password}
		if user, err := userService.Login(u); err != nil {
			// global.GIA_LOG.Error("登陆失败! 用户名不存在或者密码错误!", zap.Error(err))
			log.Println("登陆失败! 用户名不存在或者密码错误!", err)

			// response.FailWithMessage("用户名不存在或者密码错误", c)
			response.Fail(c, "用户名不存在或者密码错误")
		} else {
			if user.Enable != 1 {
				// global.GIA_LOG.Error("登陆失败! 用户被禁止登录!")
				log.Println("登陆失败! 用户被禁止登录!")

				// response.FailWithMessage("用户被禁止登录", c)
				response.Fail(c, "用户被禁止登录")
				return
			}
			b.TokenNext(c, *user)
		}
	} else {
		// response.FailWithMessage("验证码错误", c)
		response.Fail(c, "验证码错误")
	}
}

// 登录以后签发jwt
func (b *BaseApi) TokenNext(c *gin.Context, user system.SysUser) {
	j := &utils.JWT{SigningKey: []byte(global.GIA_CONFIG.JWT.SigningKey)} // 唯一签名
	claims := j.CreateClaims(systemReq.BaseClaims{
		UUID:        user.UUID,
		ID:          user.ID,
		NickName:    user.NickName,
		Username:    user.Username,
		AuthorityId: user.AuthorityId,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		// global.GIA_LOG.Error("获取token失败!", zap.Error(err))
		log.Println("获取token失败!", err)
		// response.FailWithMessage("获取token失败", c)
		response.Fail(c, "获取token失败")
		return
	}
	if !global.GIA_CONFIG.System.SingleSignOn {
		// response.OkWithDetailed(systemRes.LoginResponse{
		// 	User:      user,
		// 	Token:     token,
		// 	ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		// }, "登录成功", c)
		response.Ok(c, systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功")
		return
	}

	if jwtStr, err := jwtService.GetRedisJWT(user.Username); err == redis.Nil {
		if err := jwtService.SetRedisJWT(token, user.Username); err != nil {
			// global.GIA_LOG.Error("设置登录状态失败!", zap.Error(err))
			log.Println("设置登录状态失败!", err)

			// response.FailWithMessage("设置登录状态失败", c)
			response.Fail(c, "设置登录状态失败")
			return
		}
		// response.OkWithDetailed(systemRes.LoginResponse{
		// 	User:      user,
		// 	Token:     token,
		// 	ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		// }, "登录成功", c)
		response.Ok(c, systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功")
	} else if err != nil {
		// global.GIA_LOG.Error("设置登录状态失败!", zap.Error(err))
		log.Println("设置登录状态失败!", err)

		// response.FailWithMessage("设置登录状态失败", c)
		response.Fail(c, "设置登录状态失败")
	} else {
		var blackJWT system.JwtBlacklist
		blackJWT.Jwt = jwtStr
		if err := jwtService.JsonInBlacklist(blackJWT); err != nil {
			// response.FailWithMessage("jwt作废失败", c)
			response.Fail(c, "jwt作废失败")
			return
		}
		if err := jwtService.SetRedisJWT(token, user.Username); err != nil {
			// response.FailWithMessage("设置登录状态失败", c)
			response.Fail(c, "设置登录状态失败")
			return
		}
		// response.OkWithDetailed(systemRes.LoginResponse{
		// 	User:      user,
		// 	Token:     token,
		// 	ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		// }, "登录成功", c)
		response.Ok(c, systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功")
	}
}
