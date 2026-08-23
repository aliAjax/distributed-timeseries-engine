package observability

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

type Metrics struct {
	Ingested uint64
	Queries  uint64
	Errors   uint64
	Started  time.Time
}

func (m *Metrics) Handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, "ts_ingested_samples_total %d\nts_queries_total %d\nts_errors_total %d\nts_uptime_seconds %.0f\n", atomic.LoadUint64(&m.Ingested), atomic.LoadUint64(&m.Queries), atomic.LoadUint64(&m.Errors), time.Since(m.Started).Seconds())
}
