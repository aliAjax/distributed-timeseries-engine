package query_planner

type QueryOptions struct {
	Hints map[string]string
}

func (o *QueryOptions) AddHint(name, value string) {
	if o.Hints == nil {
		o.Hints = make(map[string]string)
	}
	o.Hints[name] = value
}
