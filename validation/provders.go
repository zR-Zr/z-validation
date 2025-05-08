package validation

import (
	"fmt"
	"net"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func RequiredProvider(msg ...string) VerifyProvider {
	return func(fName string, fVal any, accessor *ValueAccessor) (bool, string) {
		if isEmpty(fVal) {
			if len(msg) > 0 {
				return false, msg[0]
			}
			return false, fmt.Sprintf("%s 不能为空", fName)
		}
		return true, ""
	}
}

func MinLengthProvider(min int, msg ...string) VerifyProvider {
	return func(fName string, fVal any, accessor *ValueAccessor) (bool, string) {
		if len(fVal.(string)) < min {
			if len(msg) > 0 {
				return false, msg[0]
			}
			return false, fmt.Sprintf("%s 最小长度为 %d", fName, min)
		}
		return true, ""
	}
}

func MaxLengthProvider(max int, msg ...string) VerifyProvider {
	return func(fName string, fVal any, accessor *ValueAccessor) (bool, string) {
		if len(fVal.(string)) > max {
			if len(msg) > 0 {
				return false, msg[0]
			}
			return false, fmt.Sprintf("%s 最大长度为 %d", fName, max)
		}
		return true, ""
	}
}

func EqualFieldProvider(field string, msg ...string) VerifyProvider {
	return func(fName string, fVal any, accessor *ValueAccessor) (bool, string) {

		otherVal, exists := accessor.GetValue(field)
		if !exists {
			panic("[z-validation] EqualFieldProvider 未找到字段 : " + field)
		}

		if fVal.(string) != otherVal.(string) {
			if len(msg) > 0 {
				return false, msg[0]
			}

			return false, fName + "必须等于 " + field
		}

		return true, ""

	}
}

func EmailProvider(msg ...string) VerifyProvider {
	return func(fieldName string, fieldValue any, accessor *ValueAccessor) (bool, string) {
		if !strings.Contains(fieldValue.(string), "@") {
			if len(msg) > 0 {
				return false, msg[0]
			}

			return false, fieldName + "必须是一个有效的邮箱地址"
		}

		return true, ""
	}
}

func MinProvider(min float64, msg ...string) VerifyProvider {
	return func(fName string, fVal any, accessor *ValueAccessor) (bool, string) {
		var num float64
		switch v := fVal.(type) {
		case int:
			num = float64(v)
		case int32:
			num = float64(v)
		case int64:
			num = float64(v)
		case float32:
			num = float64(v)
		case float64:
			num = v
		default:
			return false, fmt.Sprintf("%s 必须是数值类型", fName)
		}
		if num < min {
			if len(msg) > 0 {
				return false, msg[0]
			}
			return false, fmt.Sprintf("%s 不能小于 %v", fName, min)
		}
		return true, ""
	}
}

func MaxProvider(max float64, msg ...string) VerifyProvider {
	return func(fName string, fVal any, accessor *ValueAccessor) (bool, string) {
		var num float64
		switch v := fVal.(type) {
		case int:
			num = float64(v)
		case int32:
			num = float64(v)
		case int64:
			num = float64(v)
		case float32:
			num = float64(v)
		case float64:
			num = v
		default:
			return false, fmt.Sprintf("%s 必须是数值类型", fName)
		}
		if num > max {
			if len(msg) > 0 {
				return false, msg[0]
			}
			return false, fmt.Sprintf("%s 不能大于 %v", fName, max)
		}
		return true, ""
	}
}

func RegexpProvider(pattern string, msg ...string) VerifyProvider {
	reg := regexp.MustCompile(pattern)
	return func(fName string, fVal any, accessor *ValueAccessor) (bool, string) {
		str, ok := fVal.(string)
		if !ok {
			return false, fmt.Sprintf("%s 必须是字符串类型", fName)
		}
		if !reg.MatchString(str) {
			if len(msg) > 0 {
				return false, msg[0]
			}
			return false, fmt.Sprintf("%s 格式不符合要求", fName)
		}
		return true, ""
	}
}

// 枚举校验：检查值是否在指定枚举列表中
func EnumProvider(enums []string, msg ...string) VerifyProvider {
	return func(fName string, fVal any, accessor *ValueAccessor) (bool, string) {
		str, ok := fVal.(string)
		if !ok {
			return false, fmt.Sprintf("%s 必须是字符串类型", fName)
		}
		for _, e := range enums {
			if str == e {
				return true, ""
			}
		}
		if len(msg) > 0 {
			return false, msg[0]
		}
		return false, fmt.Sprintf("%s 必须是以下值之一: %v", fName, enums)
	}
}

// 日期格式校验：检查字符串是否符合指定日期格式
func DateFormatProvider(layout string, msg ...string) VerifyProvider {
	return func(fName string, fVal any, accessor *ValueAccessor) (bool, string) {
		str, ok := fVal.(string)
		if !ok {
			return false, fmt.Sprintf("%s 必须是字符串类型", fName)
		}
		_, err := time.Parse(layout, str)
		if err != nil {
			if len(msg) > 0 {
				return false, msg[0]
			}
			return false, fmt.Sprintf("%s 格式不符合要求，正确格式应为: %s", fName, layout)
		}
		return true, ""
	}
}

// IP地址校验：检查是否为有效的IPv4或IPv6地址
func IPProvider(msg ...string) VerifyProvider {
	return func(fName string, fVal any, accessor *ValueAccessor) (bool, string) {
		str, ok := fVal.(string)
		if !ok {
			return false, fmt.Sprintf("%s 必须是字符串类型", fName)
		}
		ip := net.ParseIP(str)
		if ip == nil {
			if len(msg) > 0 {
				return false, msg[0]
			}
			return false, fmt.Sprintf("%s 不是有效的IP地址", fName)
		}
		return true, ""
	}
}

// URL校验：检查是否为有效的URL格式
func URLProvider(msg ...string) VerifyProvider {
	return func(fName string, fVal any, accessor *ValueAccessor) (bool, string) {
		str, ok := fVal.(string)
		if !ok {
			return false, fmt.Sprintf("%s 必须是字符串类型", fName)
		}
		_, err := url.ParseRequestURI(str)
		if err != nil {
			if len(msg) > 0 {
				return false, msg[0]
			}
			return false, fmt.Sprintf("%s 不是有效的URL地址", fName)
		}
		return true, ""
	}
}

// 辅助函数: 判断值是否为空
func isEmpty(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		return rv.Len() == 0

	case reflect.Ptr, reflect.Interface:
		return rv.IsNil() || isEmpty(rv.Elem().Interface())
	default:
		return false
	}
}
