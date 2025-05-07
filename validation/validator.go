package validation

type ValidationError struct {
	msg     string
	details []map[string]string
}

func (ve *ValidationError) Error() string {
	return ve.msg
}

func (ve *ValidationError) Details() []map[string]string {
	return ve.details
}

func ValidateStruct(obj interface{}, rules *Rules) error {

	result := rules.ValidateStruct(obj)
	if len(result) > 0 {
		return &ValidationError{
			msg:     "validation failed",
			details: result,
		}
	}
	return nil
}

func ValidateValue(value any, rule *Rule) error {
	rule.Value = value
	ok, msg := rule.Validate()
	if !ok {
		return &ValidationError{
			msg:     msg,
			details: []map[string]string{{"field": rule.Field, "msg": msg}},
		}
	}

	return nil
}

func ValidateMap(values map[string]any, rules *Rules) error {
	result := rules.ValidateValue(values)
	if len(result) > 0 {
		return &ValidationError{
			msg:     "validation failed",
			details: result,
		}
	}
	return nil
}
