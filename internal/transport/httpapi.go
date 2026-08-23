package transport

import (
	"context"
	"encoding/json"
	"example.com/distributed-timeseries-engine/internal/metric_domain"
	"example.com/distributed-timeseries-engine/internal/observability"
	"example.com/distributed-timeseries-engine/internal/query_executor"
	"example.com/distributed-timeseries-engine/internal/query_parser"
	"example.com/distributed-timeseries-engine/internal/query_planner"
	"example.com/distributed-timeseries-engine/internal/repository"
	"example.com/distributed-timeseries-engine/internal/shard_router"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

func StatusForQueryError(err error) int {
	switch err.Error() {
	case query_parser.ErrInvalidQuery.Error(), query_planner.ErrInvalidPlan.Error():
		return http.StatusBadRequest
	case query_executor.ErrQueryRead.Error():
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

type Server struct {
	Repo    *repository.Repository
	Router  *shard_router.Router
	Metrics *observability.Metrics
	Started time.Time
}

func New(repo *repository.Repository) *Server {
	return &Server{Repo: repo, Router: shard_router.New(4), Metrics: &observability.Metrics{Started: time.Now()}, Started: time.Now()}
}
func (s *Server) Routes() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { s.health(w, r) })
	m.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) { s.health(w, r) })
	m.HandleFunc("/metrics", s.Metrics.Handler)
	m.HandleFunc("/api/v1/ingest", s.ingest)
	m.HandleFunc("/api/v1/query", s.query)
	m.HandleFunc("/api/v1/query-range", s.queryRange)
	m.HandleFunc("/api/v1/labels", s.labels)
	m.HandleFunc("/api/v1/label-values", s.labelValues)
	m.HandleFunc("/api/v1/series/cardinality", s.cardinality)
	m.HandleFunc("/api/v1/metrics", s.metrics)
	m.HandleFunc("/api/v1/shards", s.shards)
	m.HandleFunc("/api/v1/replicas", s.replicas)
	m.HandleFunc("/api/v1/tenants/", s.usage)
	return logging(m)
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]any{"status": "ok", "uptime": time.Since(s.Started).String()})
}

type ingestRequest struct {
	Samples []metric_domain.Sample `json:"samples"`
}

func (s *Server) ingest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
		return
	}
	var req ingestRequest
	if e := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<20)).Decode(&req); e != nil {
		atomic.AddUint64(&s.Metrics.Errors, 1)
		writeJSON(w, 400, map[string]string{"error": e.Error()})
		return
	}
	n, e := s.Repo.Ingest(r.Context(), req.Samples)
	atomic.AddUint64(&s.Metrics.Ingested, uint64(n))
	if e != nil {
		atomic.AddUint64(&s.Metrics.Errors, 1)
		writeJSON(w, 400, map[string]any{"accepted": n, "error": e.Error()})
		return
	}
	writeJSON(w, 202, map[string]any{"accepted": n})
}
func (s *Server) query(w http.ResponseWriter, r *http.Request) {
	atomic.AddUint64(&s.Metrics.Queries, 1)
	expr, e := query_parser.Parse(r.URL.Query().Get("query"))
	if e != nil {
		writeJSON(w, StatusForQueryError(e), map[string]string{"error": e.Error()})
		return
	}
	end := time.Now()
	start := end.Add(-expr.Range)
	if v := r.URL.Query().Get("start"); v != "" {
		if n, er := strconv.ParseInt(v, 10, 64); er == nil {
			start = time.Unix(n, 0)
		}
	}
	if v := r.URL.Query().Get("end"); v != "" {
		if n, er := strconv.ParseInt(v, 10, 64); er == nil {
			end = time.Unix(n, 0)
		}
	}
	out, e := s.Repo.Query(r.Context(), expr.Metric, expr.Matchers, start, end)
	if e != nil {
		writeJSON(w, 504, map[string]string{"error": e.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"data": out, "start": start, "end": end})
}
func (s *Server) queryRange(w http.ResponseWriter, r *http.Request) { s.query(w, r) }
func (s *Server) labels(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]any{"data": s.Repo.Labels()})
}
func (s *Server) labelValues(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"data": s.Repo.LabelValues(r.URL.Query().Get("name"))})
}
func (s *Server) cardinality(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]any{"cardinality": s.Repo.Usage()["cardinality"]})
}
func (s *Server) metrics(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]any{"usage": s.Repo.Usage()})
}
func (s *Server) shards(w http.ResponseWriter, _ *http.Request) { writeJSON(w, 200, s.Router.List()) }
func (s *Server) replicas(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]any{"replicas": []string{"local"}, "consistency": "quorum"})
}
func (s *Server) usage(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	tenant := ""
	if len(parts) > 2 {
		tenant = parts[2]
	}
	u := s.Repo.Usage()
	u["tenant"] = tenant
	writeJSON(w, 200, u)
}
func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
