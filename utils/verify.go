package utils

import (
	"fmt"
	"reflect"
)

var (
	IdVerify               = Rules{"ID": []string{NotEmpty()}}
	ApiVerify              = Rules{"Path": {NotEmpty()}, "Description": {NotEmpty()}, "ApiGroup": {NotEmpty()}, "Method": {NotEmpty()}}
	MenuVerify             = Rules{"Path": {NotEmpty()}, "ParentId": {NotEmpty()}, "Name": {NotEmpty()}, "Component": {NotEmpty()}, "Sort": {Ge("0")}}
	MenuMetaVerify         = Rules{"Title": {NotEmpty()}}
	LoginVerify            = Rules{"CaptchaId": {NotEmpty()}, "Captcha": {NotEmpty()}, "Username": {NotEmpty()}, "Password": {NotEmpty()}}
	RegisterVerify         = Rules{"Username": {NotEmpty()}, "NickName": {NotEmpty()}, "Password": {NotEmpty()}, "AuthorityId": {NotEmpty()}}
	PageInfoVerify         = Rules{"Page": {NotEmpty()}, "PageSize": {NotEmpty()}}
	CustomerVerify         = Rules{"CustomerName": {NotEmpty()}, "CustomerPhoneData": {NotEmpty()}}
	AutoCodeVerify         = Rules{"Abbreviation": {NotEmpty()}, "StructName": {NotEmpty()}, "PackageName": {NotEmpty()}, "Fields": {NotEmpty()}}
	AutoPackageVerify      = Rules{"PackageName": {NotEmpty()}}
	AuthorityVerify        = Rules{"AuthorityId": {NotEmpty()}, "AuthorityName": {NotEmpty()}}
	AuthorityIdVerify      = Rules{"AuthorityId": {NotEmpty()}}
	OldAuthorityVerify     = Rules{"OldAuthorityId": {NotEmpty()}}
	ChangePasswordVerify   = Rules{"Password": {NotEmpty()}, "NewPassword": {NotEmpty()}}
	SetUserAuthorityVerify = Rules{"AuthorityId": {NotEmpty()}}
)

// 通用函数，根据传入的多个字符串生成对应的 Rules
func VerifyStr(keys ...string) Rules {
	rules := Rules{}
	for _, key := range keys {
		rules[key] = []string{NotEmpty()}
	}
	return rules
}

// 根据传入的结构体生成对应结构体下所有字段的 Rules
func VerifyStruct(inputStruct any) Rules {
	rules := Rules{}

	// 使用反射获取结构体的类型
	t := reflect.TypeOf(inputStruct)
	if t.Kind() == reflect.Ptr {
		t = t.Elem() // 如果是指针，获取其指向的元素
	}

	// 确保是结构体类型
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			// fmt.Printf("%v: %v\n", field.Name, field.Type.Kind())
			if field.Type.Kind() == reflect.String {
				rules[field.Name] = []string{NotEmpty()}
			}
		}
	}
	return rules
}

// 通用函数，根据传入的字段生成对应的 Rules (未完成)
func VerifyRules(keys ...any) Rules {
	rules := Rules{}
	for _, key := range keys {
		fmt.Printf("key: %v\n", key)
		t := reflect.TypeOf(key)
		fmt.Printf("t.Kind(): %v - %v\n", t.Kind(), t)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		fmt.Printf("t: %v\n", t)
		if t.Kind() == reflect.Struct {
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				rules[field.Name] = []string{NotEmpty()}
			}
		}
	}
	return rules
}
