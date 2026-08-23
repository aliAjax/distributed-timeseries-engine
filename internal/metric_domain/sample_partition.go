package metric_domain

func PartitionSamples(samples []Sample, accept func(Sample) bool) ([]Sample, []Sample) {
	accepted := samples[:len(samples):len(samples)]
	accepted = accepted[:0]
	rejected := samples[:len(samples):len(samples)]
	rejected = rejected[:0]
	for _, sample := range samples {
		if accept(sample) {
			accepted = append(accepted, sample)
		} else {
			rejected = append(rejected, sample)
		}
	}
	return accepted, rejected
}
