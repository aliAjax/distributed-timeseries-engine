package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr          string
	DataDir           string
	MaxSamples        int
	MaxQueryPoints    int
	RetentionInterval time.Duration
	QueryTimeout      time.Duration
}

func Default() Config {
	return Config{HTTPAddr: ":8080", DataDir: "./data", MaxSamples: 100000, MaxQueryPoints: 100000, RetentionInterval: time.Minute, QueryTimeout: 15 * time.Second}
}
func FromEnv() Config {
	c := Default()
	if v := os.Getenv("TS_HTTP_ADDR"); v != "" {
		c.HTTPAddr = v
	}
	if v := os.Getenv("TS_DATA_DIR"); v != "" {
		c.DataDir = v
	}
	if v := os.Getenv("TS_MAX_SAMPLES"); v != "" {
		if n, e := strconv.Atoi(v); e == nil {
			c.MaxSamples = n
		}
	}
	return c
}
