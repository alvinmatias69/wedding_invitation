.PHONY: all
all:
	go build -o winv cmd/main.go

.PHONY: run
run:
	go run ./cmd
