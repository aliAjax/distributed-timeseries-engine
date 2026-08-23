package query_planner

type QueryOptions struct {
	Hints map[string]string
}

func (o *QueryOptions) AddHint(name, value string) {
	o.Hints[name] = value
}
