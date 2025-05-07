package zvr

import "github.com/zR-Zr/z-validation/validation"

type RulesBuilder func() validation.Rules

func (rsb RulesBuilder) Add(fileName string, callback func(*validation.Rule)) RulesBuilder {
	rs := rsb()
	rule := validation.R(fileName).ConvertField(fileName)
	callback(rule)
	rs = append(rs, rule)
	return func() validation.Rules {
		return rs
	}
}

func Rs() RulesBuilder {
	rs := &validation.Rules{}
	return func() validation.Rules {
		return *rs
	}
}
