APP=timeseriesd
GO=go
.PHONY: fmt vet test race build run smoke
fmt:; $(GO)fmt -w $$(find . -name '*.go' -not -path './vendor/*')
vet:; $(GO) vet ./...
test:; $(GO) test ./...
race:; $(GO) test -race ./...
build:; $(GO) build -o bin/$(APP) ./cmd/timeseriesd
run:; mkdir -p data; TS_DATA_DIR=./data $(GO) run ./cmd/timeseriesd
smoke: build; ./scripts/smoke.sh
