package metric_domain

func PartitionSamples(samples []Sample, accept func(Sample) bool) ([]Sample, []Sample) {
	accepted := make([]Sample, 0, len(samples))
	rejected := make([]Sample, 0, len(samples))
	for _, sample := range samples {
		if accept(sample) {
			accepted = append(accepted, sample)
		} else {
			rejected = append(rejected, sample)
		}
	}
	return accepted, rejected
}
