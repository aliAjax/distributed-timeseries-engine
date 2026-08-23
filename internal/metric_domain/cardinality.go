package metric_domain

type Cardinality struct{ Series, Labels, Samples uint64 }

func (c *Cardinality) Add(series int, labels, samples int) {
	if series > 0 {
		c.Series += uint64(series)
	}
	if labels > 0 {
		c.Labels += uint64(labels)
	}
	if samples > 0 {
		c.Samples += uint64(samples)
	}
}
func (c Cardinality) Exceeds(series, labels, samples uint64) bool {
	return c.Series > series || c.Labels > labels || c.Samples > samples
}
func (c Cardinality) Ratio(other Cardinality) float64 {
	total := c.Series + c.Labels + c.Samples
	if total == 0 {
		return 0
	}
	return float64(other.Series+other.Labels+other.Samples) / float64(total)
}
