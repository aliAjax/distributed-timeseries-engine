package query_planner

import (
	"testing"
)

type optionalPolicy struct {
	enabled bool
}

func (p *optionalPolicy) Enabled() bool { return p.enabled }

type optionalValidator struct{}

func (*optionalValidator) Validate(Plan) error { return nil }

func TestDteR08HintsZ08(t *testing.T) {
	var options QueryOptions
	options.AddHint("tier", "warm")
	if options.Hints["tier"] != "warm" {
		t.Fatalf("hints=%v", options.Hints)
	}
}

func TestDteR08PolicyZ08(t *testing.T) {
	var concrete *optionalPolicy
	var policy QueryPolicy = concrete
	if PolicyEnabled(policy) {
		t.Fatal("typed-nil policy reported enabled")
	}
}

func TestDteR08GroupsZ08(t *testing.T) {
	var options GroupOptions
	options.Add("host", 20)
	if options.Limits["host"] != 20 {
		t.Fatalf("limits=%v", options.Limits)
	}
}

func TestDteR08ValidatorsZ08(t *testing.T) {
	var concrete *optionalValidator
	var validator QueryValidator = concrete
	var chain ValidationChain
	chain.Add(validator)
	if chain.Len() != 0 {
		t.Fatalf("typed-nil validator retained: len=%d", chain.Len())
	}
}
