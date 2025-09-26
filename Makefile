.PHONY: wire
# generate wire
.PHONY: wire
wire:
	wire ./cmd/wire.go

.PHONY: run
run:
	go run ./cmd/main.go ./cmd/wire_gen.go

.PHONY: build
build:
	GOOS=windows GOARCH=amd64 go build -o ./bin/flow-bridge-mcp.exe ./cmd

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o ./bin/flow-bridge-mcp ./cmd
