package query_planner

import "reflect"

type QueryValidator interface {
	Validate(Plan) error
}

type ValidationChain struct {
	validators []QueryValidator
}

func (c *ValidationChain) Add(validator QueryValidator) {
	if validator == nil {
		return
	}
	value := reflect.ValueOf(validator)
	if value.Kind() == reflect.Pointer && value.IsNil() {
		return
	}
	c.validators = append(c.validators, validator)
}

func (c ValidationChain) Len() int { return len(c.validators) }
