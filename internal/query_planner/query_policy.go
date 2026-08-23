package query_planner

import "reflect"

type QueryPolicy interface {
	Enabled() bool
}

func PolicyEnabled(policy QueryPolicy) bool {
	if policy == nil {
		return false
	}
	value := reflect.ValueOf(policy)
	if value.Kind() == reflect.Pointer && value.IsNil() {
		return false
	}
	return policy.Enabled()
}
