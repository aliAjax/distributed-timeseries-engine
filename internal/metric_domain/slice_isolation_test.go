package metric_domain

import (
	"reflect"
	"testing"
	"time"
)

func TestDteR04QualityZ04(t *testing.T) {
	original := []Point{{Timestamp: 1, Value: 1, Quality: QualityUnknown}, {Timestamp: 2, Value: 2, Quality: QualityGood}}
	want := append([]Point(nil), original...)
	selected := SelectQuality(original, QualityGood)
	selected[0].Value = 99
	if !reflect.DeepEqual(original, want) {
		t.Fatalf("quality filter mutated source: got=%v want=%v", original, want)
	}
}

func TestDteR04WindowZ04(t *testing.T) {
	original := []Point{{Timestamp: 1, Value: 1}, {Timestamp: 2, Value: 2}, {Timestamp: 3, Value: 3}}
	want := append([]Point(nil), original...)
	window := PointsInWindow(original, 2, 4)
	window[0].Value = 77
	if !reflect.DeepEqual(original, want) {
		panic("window extraction corrupted the source point slice")
	}
}

func TestDteR04LabelsZ04(t *testing.T) {
	original := LabelSet{{Name: "host", Value: "a"}, {Name: "zone", Value: "east"}}
	want := original.Clone()
	selected := SelectLabels(original, map[string]struct{}{"zone": {}})
	selected[0].Value = "west"
	if !reflect.DeepEqual(original, want) {
		t.Fatalf("label filter mutated source: got=%v want=%v", original, want)
	}
}

func TestDteR04PartitionZ04(t *testing.T) {
	now := time.Unix(100, 0)
	original := []Sample{{Tenant: "a", Timestamp: now}, {Tenant: "b", Timestamp: now}, {Tenant: "c", Timestamp: now}}
	want := append([]Sample(nil), original...)
	accepted, rejected := PartitionSamples(original, func(sample Sample) bool { return sample.Tenant != "b" })
	accepted[0].Tenant = "changed"
	if len(rejected) != 1 || rejected[0].Tenant != "b" || !reflect.DeepEqual(original, want) {
		t.Fatalf("partition aliases source: accepted=%v rejected=%v source=%v", accepted, rejected, original)
	}
}
