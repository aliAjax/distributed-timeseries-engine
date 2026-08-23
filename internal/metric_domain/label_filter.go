package metric_domain

func SelectLabels(labels LabelSet, names map[string]struct{}) LabelSet {
	out := labels[:len(labels):len(labels)]
	out = out[:0]
	for _, label := range labels {
		if _, ok := names[label.Name]; ok {
			out = append(out, label)
		}
	}
	return out
}
