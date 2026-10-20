.PHONY: run bench test clean

run:
	go run cmd/api/main.go

bench:
	./scripts/benchmark.sh
