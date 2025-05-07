package validation

import (
	"fmt"
	"reflect"
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
