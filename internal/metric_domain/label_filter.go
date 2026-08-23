package metric_domain

func SelectLabels(labels LabelSet, names map[string]struct{}) LabelSet {
	out := make(LabelSet, 0, len(labels))
	for _, label := range labels {
		if _, ok := names[label.Name]; ok {
			out = append(out, label)
		}
	}
	return out
}
