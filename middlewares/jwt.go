package middlewares

import (
	"fmt"
	"gmkdb/global"
	"gmkdb/models/common/request"
	"gmkdb/models/common/response"
	"gmkdb/models/system"
	"gmkdb/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

var jwtService = services.ServiceGroupApp.SystemSvcGroup.JwtService

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 sbl-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		// token := c.Request.Header.Get("sbl-token")
		token := c.Request.Header.Get(global.GIA_CONFIG.JWT.HeaderKey)

		// // fmt.Println("\n 进入判断")
		// // fmt.Printf("%s: %v\n", global.GIA_CONFIG.JWT.HeaderKey, token)
		// fmt.Println("--- header ---")
		// for k := range c.Request.Header {
		// 	fmt.Printf("%v: %v\n", k, c.GetHeader(k))
		// }
		// fmt.Println(" --- end ---\n ")

		if token == "" {
			fmt.Println("未登录或非法访问")
			response.Redirect(c, gin.H{"reload": true}, "未登录或非法访问")
			c.Abort()
			return
		}
		if jwtService.IsBlacklist(token) {
			fmt.Println("您的帐户异地登陆或令牌失效")
			response.Redirect(c, gin.H{"reload": true}, "您的帐户异地登陆或令牌失效")
			c.Abort()
			return
		}

		j := request.NewJWT()
		// // parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		// fmt.Printf("claimsJWT: %v \n err:%v\n", claims, err)		// claims, err := request.NewJWT().ParseToken(token)
		if err != nil {
			if err == request.ErrTokenExpired {
				fmt.Println("授权已过期")
				response.Redirect(c, gin.H{"reload": true}, "授权已过期")
				c.Abort()
				return
			}
			fmt.Println("解析token失败:" + err.Error())
			response.Redirect(c, gin.H{"reload": true}, err.Error())
			c.Abort()
			return
		}

		// 已登录用户被管理员禁用 需要使该用户的jwt失效 此处比较消耗性能 如果需要 请自行打开
		// 用户被删除的逻辑 需要优化 此处比较消耗性能 如果需要 请自行打开
		// if user, err := userService.FindUserByUuid(claims.UUID.String()); err != nil || user.Enable == 2 {
		// 	_ = jwtService.JsonInBlacklist(system.JwtBlacklist{Jwt: token})
		// 	response.RedirectWithDetailed(gin.H{"reload": true}, err.Error(), c)
		// 	c.Abort()
		// }

		if claims.ExpiresAt.Unix()-time.Now().Unix() < claims.BufferTime {
			// 延长过期时间
			claims.ExpiresAt = jwt.NewNumericDate(time.Unix(time.Now().Unix()+global.GIA_CONFIG.JWT.ExpiresTime, 0))
			// claims.ExpiresAt = jwt.NewNumericDate(time.UnixMilli(time.Now().UnixMilli() + global.GIA_CONFIG.JWT.ExpiresTime)) // 毫秒级
			newToken, _ := j.CreateTokenByOldToken(token, *claims)
			newClaims, _ := j.ParseToken(newToken)
			c.Header("new-token", newToken)
			// c.Header(global.GIA_CONFIG.JWT.HeaderKey, newToken)
			c.Header("new-expires-at", strconv.FormatInt(newClaims.ExpiresAt.Unix(), 10))
			if global.GIA_CONFIG.System.SingleSignOn && global.GIA_CONFIG.System.UseRedis { // 判断是否开启单点登录, 需要redis
				RedisJwtToken, err := jwtService.GetRedisJWT(newClaims.UserName)
				if err != nil {
					global.GIA_LOG.Error("get redis jwt failed", zap.Error(err))
				} else { // 当之前的取成功时才进行拉黑操作
					_ = jwtService.JsonInBlacklist(system.JwtBlacklist{Jwt: RedisJwtToken})
				}
				// 无论如何都要记录当前的活跃状态
				_ = jwtService.SetRedisJWT(newToken, newClaims.UserName)
			}
		}

		// // token := c.Request.Header.Get(global.GIA_CONFIG.JWT.HeaderKey)
		// // c.Request.Response.Header.Set("new-token", "newTokenxxx")
		// c.Header("new-token", "newTokenxxxxx")
		// // kei-token
		// fmt.Printf("jwt-token: %v\n", token)

		c.Set("claims", claims)
		c.Next()
	}
}

// func JWTAuth() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
// 		token := c.Request.Header.Get("x-token")
// 		if token == "" {
// 			response.FailWithDetailed(gin.H{"reload": true}, "未登录或非法访问", c)
// 			c.Abort()
// 			return
// 		}
// 		if jwtService.IsBlacklist(token) {
// 			response.FailWithDetailed(gin.H{"reload": true}, "您的帐户异地登陆或令牌失效", c)
// 			c.Abort()
// 			return
// 		}
// 		j := utils.NewJWT()
// 		// parseToken 解析token包含的信息
// 		claims, err := j.ParseToken(token)
// 		if err != nil {
// 			if err == utils.TokenExpired {
// 				response.FailWithDetailed(gin.H{"reload": true}, "授权已过期", c)
// 				c.Abort()
// 				return
// 			}
// 			response.FailWithDetailed(gin.H{"reload": true}, err.Error(), c)
// 			c.Abort()
// 			return
// 		}

// 		// 已登录用户被管理员禁用 需要使该用户的jwt失效 此处比较消耗性能 如果需要 请自行打开
// 		// 用户被删除的逻辑 需要优化 此处比较消耗性能 如果需要 请自行打开

// 		//if user, err := userService.FindUserByUuid(claims.UUID.String()); err != nil || user.Enable == 2 {
// 		//	_ = jwtService.JsonInBlacklist(system.JwtBlacklist{Jwt: token})
// 		//	response.FailWithDetailed(gin.H{"reload": true}, err.Error(), c)
// 		//	c.Abort()
// 		//}
// 		if claims.ExpiresAt-time.Now().Unix() < claims.BufferTime {
// 			claims.ExpiresAt = time.Now().Unix() + global.GIA_CONFIG.JWT.ExpiresTime
// 			newToken, _ := j.CreateTokenByOldToken(token, *claims)
// 			newClaims, _ := j.ParseToken(newToken)
// 			c.Header("new-token", newToken)
// 			c.Header("new-expires-at", strconv.FormatInt(newClaims.ExpiresAt, 10))
// 			if global.GIA_CONFIG.System.UseMultipoint {
// 				RedisJwtToken, err := jwtService.GetRedisJWT(newClaims.Username)
// 				if err != nil {
// 					global.GIA_LOG.Error("get redis jwt failed", zap.Error(err))
// 				} else { // 当之前的取成功时才进行拉黑操作
// 					_ = jwtService.JsonInBlacklist(system.JwtBlacklist{Jwt: RedisJwtToken})
// 				}
// 				// 无论如何都要记录当前的活跃状态
// 				_ = jwtService.SetRedisJWT(newToken, newClaims.Username)
// 			}
// 		}
// 		c.Set("claims", claims)
// 		c.Next()
// 	}
// }
