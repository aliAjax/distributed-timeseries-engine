#!/usr/bin/env bash
set -euo pipefail
PORT="${TS_SMOKE_PORT:-18080}"
DIR="$(mktemp -d)"
TS_HTTP_ADDR=":$PORT" TS_DATA_DIR="$DIR" ./bin/timeseriesd >"$DIR/server.log" 2>&1 &
PID=$!
trap 'kill "$PID" 2>/dev/null || true; rm -rf "$DIR"' EXIT
ready=0
for i in $(seq 1 30); do
  if curl -fsS "http://127.0.0.1:$PORT/healthz" >/dev/null 2>&1; then ready=1; break; fi
  sleep .2
done
[ "$ready" -eq 1 ]
curl -fsS -X POST "http://127.0.0.1:$PORT/api/v1/ingest" -H 'content-type: application/json' -d '{"samples":[{"tenant":"demo","metric":"smoke_metric","labels":[{"Name":"host","Value":"a"}],"timestamp":"2026-08-21T00:00:00Z","value":1.5,"quality":1}]}' >/dev/null
curl -fsS --get "http://127.0.0.1:$PORT/api/v1/query" --data-urlencode 'query=smoke_metric{host="a"}' --data-urlencode 'start=0' --data-urlencode 'end=4102444800' >/dev/null
curl -fsS "http://127.0.0.1:$PORT/api/v1/series/cardinality" >/dev/null
curl -fsS "http://127.0.0.1:$PORT/api/v1/tenants/demo/usage" >/dev/null
echo smoke-ok
