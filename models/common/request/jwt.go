package request

import (
	"errors"
	"fmt"
	"gmkdb/global"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Custom claims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	// jwt.StandardClaims
	jwt.RegisteredClaims
}

type BaseClaims struct {
	UserID      uint
	UserName    string
	AliasName   string
	AuthorityId uint
}

// ------

type JWT struct {
	SigningKey []byte
}

var (
	ErrTokenExpired     = errors.New("token is expired")
	ErrTokenNotValidYet = errors.New("token not active yet")
	ErrTokenMalformed   = errors.New("that's not even a token")
	ErrTokenInvalid     = errors.New("couldn't handle this token")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.GIA_CONFIG.JWT.SigningKey),
	}
}

func (j *JWT) CreateClaims(baseClaims BaseClaims, ts ...string) CustomClaims {
	// claims := reqt.CustomClaims{
	// 	BaseClaims: baseClaims,
	// 	BufferTime: global.GIA_CONFIG.JWT.BufferTime, // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
	// 	StandardClaims: jwt.StandardClaims{
	// 		NotBefore: time.Now().Unix() - 1000,                              // 签名生效时间
	// 		ExpiresAt: time.Now().Unix() + global.GIA_CONFIG.JWT.ExpiresTime, // 过期时间 7天  配置文件
	// 		Issuer:    global.GIA_CONFIG.JWT.Issuer,                          // 签名的发行者
	// 	},
	// }
	ExpiresTime := global.GIA_CONFIG.JWT.ExpiresTime
	if len(ts) > 0 {
		// 1M, 1w, 1d -> unix() 秒级
		fmt.Printf("ExpiresTime: %v\n", ts[0]) // string -> ts_int64
	}

	claims := CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: global.GIA_CONFIG.JWT.BufferTime, // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Add(-120 * time.Second)),                       // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ExpiresTime) * time.Second)), // 过期时间 7天  配置文件
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                               // 签发时间
			Issuer:    global.GIA_CONFIG.JWT.Issuer,                                                 // 签名的发行者
			// ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(global.GIA_CONFIG.JWT.ExpiresTime) * time.Second)), // 过期时间 7天  配置文件
		},
	}
	return claims
}

// 创建一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// CreateTokenByOldToken 旧token 换新token 使用归并回源避免并发问题
func (j *JWT) CreateTokenByOldToken(oldToken string, claims CustomClaims) (string, error) {
	v, err, _ := global.JWT_Concurrency_Control.Do("JWT:"+oldToken, func() (interface{}, error) {
		return j.CreateToken(claims)
	})
	return v.(string), err
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, ErrTokenInvalid

	} else {
		return nil, ErrTokenInvalid
	}
}
