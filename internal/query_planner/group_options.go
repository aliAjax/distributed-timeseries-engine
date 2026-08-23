package query_planner

type GroupOptions struct {
	Limits map[string]int
}

func (o *GroupOptions) Add(group string, limit int) {
	if o.Limits == nil {
		o.Limits = make(map[string]int)
	}
	o.Limits[group] = limit
}
