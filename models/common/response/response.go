package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"gmkdb/utils"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

type Response struct {
	Code  status `json:"code"`
	Data  any    `json:"data"`
	Msg   string `json:"msg"`
	Times string `json:"ts"`

	Count int `json:"count"`
}

type status int

const (
	ERROR    status = http.StatusBadRequest      // 7
	FAILURE  status = http.StatusBadRequest      // 3
	REDIRECT status = http.StatusMultipleChoices // ?
	SUCCESS  status = http.StatusOK              // 1
)

// Respond Result
// code 处理状态, data 返回值map[string]any{}, msg 提示信息. count 数据条数?
func Respond(c *gin.Context, code status, data any, count int, msg string) {
	response := Response{
		Code:  code,
		Data:  data,
		Msg:   code.message(msg),
		Count: count,
		Times: time.Now().Format(time.RFC3339),
	}
	if count > 0 { // 转map快速去除该字段 map[string]interface{} // 字段顺序丢失, 无序 不能保证标签顺序逻辑
		// // response.Count = count
		// c.JSON(int(code), response)
		// return
	}

	// http.StatusOK
	c.JSON(int(code), response)
}
func Ok(c *gin.Context, args ...any) {
	SUCCESS.send(c, args)
}

func Fail(c *gin.Context, args ...any) {
	ERROR.send(c, args)
}

func Redirect(c *gin.Context, args ...any) { // ?
	REDIRECT.send(c, args)
}

func (code status) send(c *gin.Context, args []any) {
	p := code.parseParam(args)
	// Respond(c, code, p.Data, p.Msg)
	Respond(c, code, p.Data, p.Count, p.Msg)
}

// 传入多个同类参数时, msg和data只取最后一个
func (code status) parseParam(args []any) (r Response) {
	// Analysis Parameters Parse
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			r.Msg += v // message
		case error:
			r.Msg += v.Error() // message
		case int:
			r.Count = v
		case int64:
			r.Count = int(v)
		default:
			r.Data = v // data
		}
	}
	return
}

// default message
func (code status) message(msg string) string {
	// if len(msg) > 0 {
	// 	return msg[0]
	// }
	// return http.StatusText(int(code)) // 默认使用 HTTP 状态消息
	if strings.TrimSpace(msg) != "" {
		return msg
	}
	return http.StatusText(int(code)) // 默认使用 HTTP 状态消息
}

// 从 r 内获取关于 &data 的数据
func (r *Response) Unmarshal(data any) (err error) {
	val := reflect.ValueOf(data)
	// fmt.Println(utils.PrettyPrint(data))

	if r.Data != nil && val.Kind() == reflect.Ptr {
		v := val.Elem()
		switch v.Kind() {
		case reflect.Struct:
			// fmt.Printf("Struct: %v\n", v.Type())
			err = json.Unmarshal([]byte(utils.ParseJson(r.Data)), &data)
		case reflect.Map:
			// fmt.Printf("Map: %v\n", v.Type()) // gin.H
			err = mapstructure.Decode(r.Data, data)
		case reflect.Slice:
			// fmt.Printf("Slice: %v\n", v.Type()) // []gin.H
			// 将 data 转换为反射值
			dataValue := reflect.ValueOf(data).Elem() // 获取指针指向的值
			if da, ok := dataValue.Interface().([]gin.H); ok {
				if inSlice, ok := r.Data.([]any); ok {
					for _, item := range inSlice {
						var d gin.H
						err = mapstructure.Decode(item, &d)
						da = append(da, d)
					}
					// // 关键步骤：将修改后的 da 赋值回 data // data = da
					dataValue.Set(reflect.ValueOf(da).Convert(dataValue.Type()))
				} else {
					fmt.Println("识别的数据 - 不是 Slice - ")
				}
			} else {
				fmt.Printf("传入的值 不ok: %v\n", ok)
			}

		default:
			err = errors.ErrUnsupported
			fmt.Printf(" * default: %v\n", utils.PrettyPrint(data))
		}
	} else {
		err = errors.New("Response Unmarshal: data is null (or) not a pointer")
	}
	return
}

// -------

// type Response struct {
// 	Code int         `json:"code"`
// 	Data interface{} `json:"data"`
// 	Msg  string      `json:"msg"`
// }

// const (
// 	ERROR   = 7
// 	SUCCESS = 0
// )

// func Result(code int, data interface{}, msg string, c *gin.Context) {
// 	// 开始时间
// 	c.JSON(http.StatusOK, Response{
// 		code,
// 		data,
// 		msg,
// 	})
// }

// func Ok(c *gin.Context) {
// 	Result(SUCCESS, map[string]interface{}{}, "操作成功", c)
// }

// func OkWithMessage(message string, c *gin.Context) {
// 	Result(SUCCESS, map[string]interface{}{}, message, c)
// }

// func OkWithData(data interface{}, c *gin.Context) {
// 	Result(SUCCESS, data, "查询成功", c)
// }

// func OkWithDetailed(data interface{}, message string, c *gin.Context) {
// 	Result(SUCCESS, data, message, c)
// }

// func Fail(c *gin.Context) {
// 	Result(ERROR, map[string]interface{}{}, "操作失败", c)
// }

// func FailWithMessage(message string, c *gin.Context) {
// 	Result(ERROR, map[string]interface{}{}, message, c)
// }

// func FailWithDetailed(data interface{}, message string, c *gin.Context) {
// 	Result(ERROR, data, message, c)
// }
