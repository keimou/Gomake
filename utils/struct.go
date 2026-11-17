package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// 展示输出结构体
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

// Struct/Map[string]any 解析成 json 格式的 string
func ParseJson(i interface{}) string {
	s, _ := json.Marshal(i)
	return string(s)
}

// 合并结构体
func MergeStructs(now, old interface{}, preserveZeroFields ...map[string]bool) {
	vNow := reflect.ValueOf(now).Elem()
	vOld := reflect.ValueOf(old).Elem()

	var pzf map[string]bool
	if len(preserveZeroFields) > 0 {
		pzf = preserveZeroFields[0]
	}

	for i := 0; i < vNow.NumField(); i++ {
		fieldNow := vNow.Field(i)
		fieldOld := vOld.Field(i)

		fieldName := vNow.Type().Field(i).Name
		if fieldNow.IsZero() && !pzf[fieldName] { // isZero(fieldNow)
			fieldNow.Set(fieldOld)
		}

		fmt.Println()
	}
}

//	func isZero(v reflect.Value) bool {
//		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
//	}

//	将结构体slice中对应字段名fieldName组成一个新切片存入&resultPtr中,其中resultPtr切片类型需对应字段类型
//
// 提取结构体切片对于字段的切片
func ExtractField(slice any, fieldName string, resultPtr any) error {
	// 获取 slice 的反射值和类型
	// sliceValue := reflect.ValueOf(slice)
	// if sliceValue.Kind() != reflect.Ptr {
	// 	return fmt.Errorf("slice must be a pointer to a slice")
	// }
	// sliceValue = sliceValue.Elem()
	sliceValue := refValue(reflect.ValueOf(slice))
	if sliceValue.Kind() != reflect.Slice {
		return fmt.Errorf("slicePtr must be a slice")
	}

	// 获取 resultPtr 的反射值和类型
	resultValue := reflect.ValueOf(resultPtr)
	if resultValue.Kind() != reflect.Ptr {
		return fmt.Errorf("resultPtr must be a pointer to a slice")
	}
	resultValue = resultValue.Elem()
	if resultValue.Kind() != reflect.Slice {
		return fmt.Errorf("resultPtr must be a pointer to a slice")
	}

	// 获取切片元素的类型（即结构体类型）
	elemType := refType(sliceValue.Type().Elem())
	// if elemType.Kind() != reflect.Struct {
	// 	return fmt.Errorf("slice elements must be structs")
	// }
	// 根据字段名查找结构体中的字段
	field, found := elemType.FieldByName(fieldName)
	if !found {
		return fmt.Errorf("field %s not found in struct %s", fieldName, elemType.Name())
	}

	// 获取字段的类型
	fieldType := field.Type

	// 检查结果切片的元素类型是否与字段类型匹配
	resultElemType := resultValue.Type().Elem()
	if !fieldType.AssignableTo(resultElemType) && !fieldType.ConvertibleTo(resultElemType) {
		return fmt.Errorf("cannot assign or convert field type %s to result element type %s", fieldType, resultElemType)
	}

	// 遍历原始切片
	for i := 0; i < sliceValue.Len(); i++ {
		// 获取当前结构体元素
		elem := refValue(sliceValue.Index(i))

		// 获取字段值
		fieldValue := elem.FieldByName(fieldName)
		if !fieldValue.IsValid() {
			// 理论上不会发生，因为上面已经检查过字段存在
			continue
		}

		// 将字段值转换为目标切片元素类型
		var convertedValue reflect.Value
		if fieldType.AssignableTo(resultElemType) {
			convertedValue = fieldValue
		} else if fieldType.ConvertibleTo(resultElemType) {
			convertedValue = fieldValue.Convert(resultElemType)
		}

		// 将转换后的值追加到结果切片
		resultValue.Set(reflect.Append(resultValue, convertedValue))
	}

	return nil
}

// 获取当前结构, 忽略并跳过指针
func refValue(elem reflect.Value) reflect.Value {
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	return elem
}
func refType(elem reflect.Type) reflect.Type {
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	return elem
}
