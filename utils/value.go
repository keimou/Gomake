package utils

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func NoAny(...any) {}

// Value 类型判定

// 判断变量是否为 空或零值
func IsNil(value any) (ok bool) {
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Ptr:
		// fmt.Println("Ptr:", val.Interface(), val.IsNil())
		if ok = val.IsNil(); !ok {
			ok = IsNil(val.Elem().Interface())
		}
	case reflect.Int, reflect.Int64, reflect.Float64, reflect.Uint, reflect.Uint64:
		// fmt.Println("Number:", val.Interface(), IsZero(val.Interface()))
		ok = IsZero(val.Interface())
	case reflect.String:
		// fmt.Println("String:", val.Interface(), IsNull(val.Interface()))
		ok = IsNull(val.Interface())
	case reflect.Bool:
		// fmt.Printf("Bool: %v\n", val.Interface())
		ok = !val.Interface().(bool)
	case reflect.Struct:
		switch v := val.Interface().(type) {
		case time.Time:
			ok = v.IsZero()
		default:
			// fmt.Printf("Struct: %v\n", val.Interface())
			ok = true // 未完成, 其余结构体均判定为空
		}
	// case reflect.:
	default:
		// fmt.Printf("default: %v\n", val.Interface())
		ok = true // 未完成, 没有录入元素均临时判定为nil
	}
	return
}
func IsNil_X(str any) (ok bool) {
	switch v := str.(type) {
	case string:
		ok = len(strings.TrimSpace(v)) == 0
	case int:
		ok = v == 0
	case int64:
		ok = v == 0
	case float64:
		ok = v == 0
	case bool:
		ok = v
	case *string:
		if ok = v == nil; !ok {
			ok = IsNil(*v)
		}
	case *int:
		if ok = v == nil; !ok {
			ok = IsNil(*v)
		}
	case *int64:
		if ok = v == nil; !ok {
			ok = IsNil(*v)
		}
	case *float64:
		if ok = v == nil; !ok {
			ok = IsNil(*v)
		}
	case *bool:
		if ok = v == nil; !ok {
			ok = IsNil(*v)
		}
	case time.Time:
		fmt.Println("time", v.IsZero())
		ok = v.IsZero()
	case *time.Time:
		// if true {
		// 	v = nil
		// }
		if ok = v == nil; !ok {
			fmt.Println("*time.Time:", v, v.UnixMilli(), v.IsZero())
			ok = IsNil(*v)
		}
		// fmt.Println("*time.Time:", v)
		// ok = v.IsZero()
		// ok = true
	default:
		// ok = len(fmt.Sprintf("%v", v)) == 0
		fmt.Println("isNil:", v)
		ok = true // 未完成, 没有录入元素均临时判定为nil
	}
	return
}

// 判断字符串是否为空 (过滤空格,全是空格也为空)
func IsNull(str any) (ok bool) {
	if str, ok = RemPtr(str); ok {
		return
	}

	switch s := str.(type) {
	case string:
		ok = len(strings.TrimSpace(s)) == 0
	default:
		ok = len(fmt.Sprintf("%v", s)) == 0
	}
	return
}

// 摘掉指针, 取出指向的变量. 若nil则ok为真, 否则用v取出变量
func RemPtr(value any) (v any, ok bool) {
	if val := reflect.ValueOf(value); val.Kind() == reflect.Ptr {
		if ok = val.IsNil(); !ok {
			v, ok = RemPtr(val.Elem().Interface())
		}
	} else {
		v = val.Interface()
	}
	return
}

// 判断'数字型字符串'是否为零 (非数字,非有效数,为0和极小值 都返回真)
// 判断是否为零值
func IsZero(num any) (ok bool) {
	// if val := reflect.ValueOf(num); val.Kind() == reflect.Ptr {
	// 	if ok = val.IsNil(); !ok {
	// 		ok = IsZero(val.Elem().Interface())
	// 	}
	// 	return
	// }
	if num, ok = RemPtr(num); ok {
		return
	}

	var str string
	switch v := num.(type) {
	case string:
		str = v
	case int:
		return v == 0
	case int64:
		return v == 0
	case uint:
		return v == 0
	case uint64:
		return v == 0
	// case float64:
	// 	return math.IsNaN(v) || v == 0 || math.Abs(v) < 1e-13
	default:
		str = fmt.Sprintf("%v", v)
	}
	str = strings.TrimSpace(str) // 过滤字符串前后多余空格
	n, err := strconv.ParseFloat(str, 64)
	// 1.转换有错误 2.为NaN 3.为0 4.极小值(趋于0)
	return err != nil || math.IsNaN(n) || n == 0 || math.Abs(n) < 1e-13 // 8 9 10 16
}

// // 非零 -> 由 !IsZero(str) 替代
// func IsNonZero(str string) bool {
// 	num, _ := strconv.ParseFloat(str, 64)
// 	return !math.IsNaN(num) && num != 0
// }

// 判断字符串是否为数字 (str), 默认float; (str, "int")->是否为整形
func IsNumber(value string, types ...string) bool {
	var err error
	var kind string
	if len(types) > 0 {
		kind = types[0]
	}
	switch kind {
	case "int":
		_, err = strconv.ParseInt(value, 10, 64)
	case "float":
		fallthrough
	default:
		if v, e := strconv.ParseFloat(value, 64); e != nil {
			err = e
		} else {
			if math.IsNaN(v) {
				err = errors.New("not a number")
			}
		}
	}
	// fmt.Printf("%v: %v\n", num, err)
	return err == nil
}

// IsInf reports whether f is illegal float.
// 判断是否为极小数(or无穷小)
func IsInf(f float64) bool {
	const SmallestNormal = 2.2250738585072014e-16
	// SmallestNormal := math.SmallestNonzeroFloat64 // 4.94066e-324
	// SmallestNormal := math.SmallestNonzeroFloat32 // 1.4013e-45
	if math.Abs(f) < SmallestNormal && f != 0 { // && f != 0
		return true
	} else {
		return false
	}
}

// --Struct--

func IsPtr(val any) bool {
	t := reflect.TypeOf(val)
	if t == nil {
		return false
	}
	return t.Kind() == reflect.Ptr
}

// 判定 1.model为结构体 2.包含指定字段名 (&User, "UserName")
func HasField(model any, fieldName string) bool {
	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val.Kind() == reflect.Struct && val.FieldByName(fieldName).IsValid()
}

// 判断是否为结构体 或结构体指针 IsStructOrPointerToStruct
func IsStruct(val any) (ok bool) {
	t := reflect.TypeOf(val)
	if t != nil {
		switch t.Kind() {
		case reflect.Struct:
			ok = true
		case reflect.Ptr:
			// 判断指针指向的是否是结构体
			if t.Elem().Kind() == reflect.Struct {
				ok = true
			}
		}
	}
	return
}

// 判断是否为 1.指向结构体的'指针' 2.结构体类型的'切片'或'数组'
func IsStructs(val any) (ok bool) {
	var ptr bool
	t := reflect.TypeOf(val)
	if t == nil {
		return
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		ptr = true
	}
	switch t.Kind() {
	case reflect.Struct:
		// 判断结构体是否由指针指向
		if ptr {
			ok = true
		}
	case reflect.Slice, reflect.Array:
		// 判断指针指向的是否是结构体
		if t.Elem().Kind() == reflect.Struct {
			ok = true
		}
	}
	return
}

// 检查 `data` 是否是一个指针类型且指向 [结构体]
func IsStructx(data any) error {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("sql data err: data must be a pointer to a struct")
	}
	return nil
}

// 检查 `datas` 是否是一个指针类型且 指向 [结构体], [结构体.切片 or .数组]
func IsStructsx(datas any) error {
	v := reflect.ValueOf(datas)
	if v.Kind() != reflect.Ptr ||
		(v.Elem().Kind() != reflect.Struct &&
			!(v.Elem().Kind() == reflect.Slice || v.Elem().Kind() == reflect.Array) ||
			(v.Elem().Kind() == reflect.Slice || v.Elem().Kind() == reflect.Array) &&
				v.Type().Elem().Elem().Kind() != reflect.Struct) {
		return errors.New("sql data err: data must be a pointer to a dbtables")
	}
	return nil
}

func GetFieldAdd(v reflect.Value) (ptr any) {
	// fmt.Printf("v: %v\n", v.Kind())
	fieldAddr := v.Addr() // .Elem() // .Addr()
	fmt.Printf("fieldAddr: %v\n", fieldAddr)

	switch v.Kind() {
	case reflect.Int64:
		if p, ok := fieldAddr.Interface().(*int64); ok {
			fmt.Printf("p: %v\n", p)
			ptr = p
		}
	case reflect.String:
		if p, ok := fieldAddr.Interface().(*string); ok {
			fmt.Printf("p: %v\n", p)
			ptr = p
		}
	default:
		// ptr = nil
	}

	return
}

func GetFieldSt(value any) (args string, dest []any) {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		for i := 0; i < val.NumField(); i++ {
			if val.Field(i).Kind() == reflect.Struct {
				va2 := val.Field(i)
				for j := 0; j < va2.NumField(); j++ {
					// if va2.Type().Field(j).Tag.Get("db") != "-" {
					// 	if len(args) > 0 {
					// 		args = fmt.Sprintf("%s, %s", args, va2.Type().Field(j).Tag.Get("db"))
					// 	} else {
					// 		args = va2.Type().Field(j).Tag.Get("db")
					// 	}
					// 	dest = append(dest, va2.Field(j).Addr())
					// }
					if field, ok := filterField(va2.Type().Field(j).Tag.Get("db")); ok {
						addStr(&args, field)
						dest = append(dest, va2.Field(j).Addr())
					}

				}
			} else {
				if field, ok := filterField(val.Type().Field(i).Tag.Get("db")); ok {
					addStr(&args, field)
					dest = append(dest, val.Field(i).Addr())
					// val.Field(i)
				}
			}
		}
	}

	return
}

// filter
func filterField(str string) (failed string, ok bool) {

	if str != "-" {
		failed = str
		ok = true
	}

	return
}
func addStr(args *string, str string) {
	if len(*args) > 0 {
		*args = fmt.Sprintf("%s, %s", *args, str)
	} else {
		*args = str
	}
}

// 过滤字段
// 1. db: "-" 2. db:"" 3. db:"field"
// 4. db:"field, omitempty" 5. db:"field,default=0" 6. db:"field,default=0, omitempty"
// 7. db:"field,default=0, omitempty, comment='test'" 8. db:"field,default=0, omitempty, comment='test', index"
// 9. db:"field,default=0, omitempty, comment='test', index:idx_name"
// 10. db:"field,default=0, omitempty, comment='test', index:idx_name, unique"
// 11. db:"field,default=0, omitempty, comment='test', index:idx_name, unique:idx_name"
// 12. db:"field,default=0, omitempty, comment='test', index:idx_name, unique:idx_name, primary"
// 13. db:"field,default=0, omitempty, comment='test', index:idx_name, unique:idx_name, primary:pk_name"
// 14. db:"field,default=0, omitempty, comment='test', index:idx_name, unique:idx_name, primary:pk_name, auto_increment"
// 15. db:"field,default=0, omitempty, comment='test', index:idx_name, unique:idx_name, primary:pk_name, auto_increment:1"
// 16. db:"field,default=0, omitempty, comment='test', index:idx_name, unique:idx_name, primary:pk_name, auto_increment:1, not null"
// 17. db:"field,default=0, omitempty, comment='test', index:idx_name, unique:idx_name, primary:pk_name, auto_increment:1, not null:1"
