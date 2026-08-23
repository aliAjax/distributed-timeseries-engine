package query_planner

type QueryPolicy interface {
	Enabled() bool
}

func PolicyEnabled(policy QueryPolicy) bool {
	if policy == nil {
		return false
	}
	return policy.Enabled()
}
