package form

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Unknwon/com"
	"github.com/go-macaron/binding"
)

type Form interface {
	binding.Validator
}

func init() {
	binding.SetNameMapper(com.ToSnakeCase)
}

func AssignForm(form interface{}, data map[string]interface{}) {
	typ := reflect.TypeOf(form)
	val := reflect.ValueOf(form)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldName := field.Tag.Get("form")
		if fieldName == "-" {
			continue
		}

		data[fieldName] = val.Field(i).Interface()
	}
}

func getRuleBody(field reflect.StructField, prefix string) string {
	for _, rule := range strings.Split(field.Tag.Get("binding"), ";") {
		if strings.HasPrefix(rule, prefix) {
			return rule[len(prefix) : len(rule)-1]
		}
	}
	return ""
}

func GetSize(field reflect.StructField) string {
	return getRuleBody(field, "Size(")
}

func GetMinSize(field reflect.StructField) string {
	return getRuleBody(field, "MinSize(")
}

func GetMaxSize(field reflect.StructField) string {
	return getRuleBody(field, "MaxSize(")
}

func GetInclude(field reflect.StructField) string {
	return getRuleBody(field, "Include(")
}

func validate(errs binding.Errors, data map[string]interface{}, f Form) binding.Errors {
	if errs.Len() == 0 {
		return errs
	}

	data["HasError"] = true
	AssignForm(f, data)

	data["HasError"] = true
	AssignForm(f, data)

	typ := reflect.TypeOf(f)
	val := reflect.ValueOf(f)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldName := field.Tag.Get("form")
		if fieldName == "-" {
			continue
		}

		if errs[0].FieldNames[0] == field.Name {
			data["Err_"+field.Name] = true

			trName := field.Tag.Get("locale")
			if len(trName) == 0 {
				trName = "form." + field.Name
			}

			switch errs[0].Classification {
			case binding.ERR_REQUIRED:
				data["ErrorMsg"] = trName + "不能为空。"
			case binding.ERR_ALPHA_DASH:
				data["ErrorMsg"] = trName + "必须为英文字母、阿拉伯数字或横线（-_）"
			case binding.ERR_ALPHA_DASH_DOT:
				data["ErrorMsg"] = trName + "必须为英文字母、阿拉伯数字、横线（-_）或点。"
			case binding.ERR_SIZE:
				data["ErrorMsg"] = trName + fmt.Sprintf("长度必须为 %s。", GetSize(field))
			case binding.ERR_MIN_SIZE:
				data["ErrorMsg"] = trName + fmt.Sprintf("长度最小为 %s 个字符", GetMinSize(field))
			case binding.ERR_MAX_SIZE:
				data["ErrorMsg"] = trName + fmt.Sprintf("长度最大为 %s 个字符", GetMaxSize(field))
			case binding.ERR_EMAIL:
				data["ErrorMsg"] = trName + "不是一个有效的邮箱地址。"
			case binding.ERR_URL:
				data["ErrorMsg"] = trName + "不是一个有效的 URL。"
			case binding.ERR_INCLUDE:
				data["ErrorMsg"] = trName + fmt.Sprintf("必须包含子字符串 '%s'。", GetInclude(field))
			default:
				data["ErrorMsg"] = "未知错误：" + " " + errs[0].Classification
			}
			return errs
		}
	}
	return errs
}
