package ownership_checks_test

import (
	"context"
	"io"
	"log"
	"reflect"
	"testing"
	"time"

	"example.com/distributed-timeseries-engine/internal/aggregation"
	"example.com/distributed-timeseries-engine/internal/chunk_codec"
	"example.com/distributed-timeseries-engine/internal/metric_domain"
	"example.com/distributed-timeseries-engine/internal/observability"
)

func TestDteR11RateZ11(t *testing.T) {
	points := []metric_domain.Point{
		{Timestamp: int64(3 * time.Second), Value: 3},
		{Timestamp: int64(time.Second), Value: 1},
		{Timestamp: int64(2 * time.Second), Value: 2},
	}
	want := append([]metric_domain.Point(nil), points...)
	if got := aggregation.Rate(points); got != 1 {
		t.Fatalf("rate=%v", got)
	}
	if !reflect.DeepEqual(points, want) {
		t.Fatalf("rate reordered caller points: got=%v want=%v", points, want)
	}
}

func TestDteR11DownsampleZ11(t *testing.T) {
	points := []chunk_codec.Point{
		{Timestamp: 30, Value: 3, Quality: 1},
		{Timestamp: 10, Value: 1, Quality: 1},
		{Timestamp: 20, Value: 2, Quality: 1},
	}
	want := append([]chunk_codec.Point(nil), points...)
	got := chunk_codec.Downsample(points, 10)
	wantOutput := []chunk_codec.Point{
		{Timestamp: 10, Value: 1, Quality: 1},
		{Timestamp: 20, Value: 2, Quality: 1},
		{Timestamp: 30, Value: 3, Quality: 1},
	}
	if !reflect.DeepEqual(got, wantOutput) {
		t.Fatalf("downsampled points=%v want=%v", got, wantOutput)
	}
	if !reflect.DeepEqual(points, want) {
		t.Fatalf("downsample reordered caller points: got=%v want=%v", points, want)
	}
}

func TestDteR11ResampleZ11(t *testing.T) {
	points := []metric_domain.Point{
		{Timestamp: 30, Value: 3},
		{Timestamp: 10, Value: 1},
		{Timestamp: 20, Value: 2},
	}
	want := append([]metric_domain.Point(nil), points...)
	if got := metric_domain.Resample(points, 0, 40, 10); len(got) != 4 {
		t.Fatalf("buckets=%d", len(got))
	}
	if !reflect.DeepEqual(points, want) {
		t.Fatalf("resample reordered caller points: got=%v want=%v", points, want)
	}
}

func TestDteR11TraceZ11(t *testing.T) {
	tracer := observability.NewTracer(log.New(io.Discard, "", 0))
	_, span := tracer.Start(context.Background(), "query")
	observability.WithAttribute(span, "tenant", "alpha")
	tracer.End(span, nil)

	recent := tracer.Recent()
	recent[0].Attributes["tenant"] = "changed"
	if got := tracer.Recent()[0].Attributes["tenant"]; got != "alpha" {
		t.Fatalf("trace snapshot mutated stored attributes: %q", got)
	}
}
