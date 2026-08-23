package main

import (
	"context"
	"encoding/json"
	"example.com/distributed-timeseries-engine/internal/config"
	"example.com/distributed-timeseries-engine/internal/metric_domain"
	"example.com/distributed-timeseries-engine/internal/repository"
	"example.com/distributed-timeseries-engine/internal/transport"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.FromEnv()
	if e := os.MkdirAll(cfg.DataDir, 0750); e != nil {
		log.Fatal(e)
	}
	repo, e := repository.Open(cfg.DataDir, cfg.MaxSamples)
	if e != nil {
		log.Fatal(e)
	}
	defer repo.Close()
	seed(repo)
	srv := transport.New(repo)
	httpSrv := &http.Server{Addr: cfg.HTTPAddr, Handler: srv.Routes(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 20 * time.Second, WriteTimeout: 20 * time.Second, IdleTimeout: 60 * time.Second}
	go func() {
		log.Printf("timeseriesd listening on %s", cfg.HTTPAddr)
		if e := httpSrv.ListenAndServe(); e != nil && e != http.ErrServerClosed {
			log.Fatal(e)
		}
	}()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	httpSrv.Shutdown(ctx)
}
func seed(r *repository.Repository) {
	m := metric_domain.Metric{Tenant: "demo", Name: "device_temperature", Unit: "celsius", Labels: metric_domain.LabelSet{{Name: "device", Value: "seed"}}, CreatedAt: time.Now()}
	_ = r.Register(m)
	_ = json.Valid
}
