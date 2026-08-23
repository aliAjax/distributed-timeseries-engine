package query_planner

type GroupOptions struct {
	Limits map[string]int
}

func (o *GroupOptions) Add(group string, limit int) {
	o.Limits[group] = limit
}
