package validation

import (
	"fmt"
	"reflect"
)

type ValueAccessor struct {
	obj reflect.Value
}

func (va *ValueAccessor) GetValue(key string) (any, bool) {
	// 处理结构体
	if va.obj.Kind() == reflect.Struct {
		field := va.obj.FieldByName(key)
		if !field.IsValid() {
			return nil, false
		}
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		return field.Interface(), true
	}
	// 处理 map
	if va.obj.Kind() == reflect.Map {
		mapKey := reflect.ValueOf(key)
		field := va.obj.MapIndex(mapKey)
		if !field.IsValid() {
			return nil, false
		}
		return field.Interface(), true
	}

	return nil, false
}

type VerifyProvider func(fieldName string, fieldValue any, accessor *ValueAccessor) (bool, string)

type Rule struct {
	Field         string
	Value         any
	ResponseField string
	accessor      *ValueAccessor
	Providers     []VerifyProvider
	required      bool
	requiredMsg   string
}

type tempRule struct {
	Rule
}

func (r *Rule) Validate() (bool, string) {
	fName := r.UseField()
	for _, provider := range r.Providers {
		if ok, msg := provider(fName, r.Value, r.accessor); !ok {
			return false, msg
		}
	}
	return true, ""
}

func R(field string) *tempRule {
	return &tempRule{
		Rule: Rule{
			Field: field,
		},
	}

}
func (r *Rule) UseField() string {
	if r.ResponseField != "" {
		return r.ResponseField
	}
	return r.Field
}

func (r *Rule) ConvertField(field string) *Rule {
	r.ResponseField = field
	return r
}

func (r *Rule) Required(msg ...string) *Rule {
	r.required = true

	// 保存用户自定义消息 (有线使用用户传入的, 否则使用默认的)
	if len(msg) > 0 {
		r.requiredMsg = msg[0]
	}

	r.Providers = append(r.Providers, RequiredProvider(msg...))
	return r
}

func (r *Rule) MinLength(min int, msg ...string) *Rule {
	r.Providers = append(r.Providers, MinLengthProvider(min, msg...))
	return r
}

func (r *Rule) MaxLength(max int, msg ...string) *Rule {
	r.Providers = append(r.Providers, MaxLengthProvider(max, msg...))
	return r
}

func (r *Rule) EqualField(field string, msg ...string) *Rule {
	r.Providers = append(r.Providers, EqualFieldProvider(field, msg...))
	return r
}
func (r *Rule) Email(msg ...string) *Rule {
	r.Providers = append(r.Providers, EmailProvider(msg...))
	return r
}

func RV(name string) *tempRule {
	return &tempRule{
		Rule: Rule{
			Field: name,
		},
	}
}

type Rules []*Rule

func (rs Rules) ValidateStruct(obj interface{}) []map[string]string {
	result := make([]map[string]string, 0)
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	accessor := ValueAccessor{
		obj: val,
	}

	for _, rule := range rs {
		rule.accessor = &accessor
		fieldValue := val.FieldByName(rule.Field)

		if fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}

		if !fieldValue.IsValid() {
			if rule.required {
				errMsg := rule.requiredMsg
				if errMsg == "" {
					errMsg = fmt.Sprintf("%s不能为空", rule.UseField())
				}
				result = append(result, map[string]string{
					"field": rule.UseField(),
					"msg":   errMsg,
				})
			}
			continue
		}

		fmt.Println("field: ", rule.Field)
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}

		rule.Value = fieldValue.Interface()
		if ok, msg := rule.Validate(); !ok {
			result = append(result, map[string]string{
				"field": rule.UseField(),
				"msg":   msg,
			})
		}
	}
	return result
}

func (r Rules) ValidateValue(value any) []map[string]string {
	result := make([]map[string]string, 0)
	// 获取 value 的反射值
	val := reflect.ValueOf(value)

	accessor := ValueAccessor{obj: val}
	for _, rule := range r {
		rule.accessor = &accessor

		// 非结构体类型, obj 没有意义, 设为零值
		if val.Kind() == reflect.Map {
			// 检查 Map 键是否为 string 类型 ()
			if val.Type().Key().Kind() != reflect.String {
				panic("map key must be string")
			}

			// 提取 Map 中对应 Field 的值
			key := reflect.ValueOf(rule.Field)
			fieldVal := val.MapIndex(key)

			// 处理键不存在的情况
			if !fieldVal.IsValid() {
				if rule.required {
					errMsg := rule.requiredMsg
					if errMsg == "" {
						errMsg = fmt.Sprintf("%s不能为空", rule.UseField())
					}
					result = append(result, map[string]string{
						"field": rule.UseField(),
						"msg":   errMsg,
					})
				}
				continue
			}
			rule.Value = fieldVal.Interface() // 赋值为 map 中对应的值
		} else {
			// 基础类型 直接赋值
			rule.Value = value
		}
		if ok, msg := rule.Validate(); !ok {
			result = append(result, map[string]string{
				"field": rule.UseField(),
				"msg":   msg,
			})
		}
	}
	return result
}
