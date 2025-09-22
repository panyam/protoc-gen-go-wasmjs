
all: default test wasm install

# Default generators (split by language) - using layered architecture
default: default-go default-ts

default-go:
	go build -o ./bin/protoc-gen-go-wasmjs-go ./cmd/protoc-gen-go-wasmjs-go

default-ts:
	go build -o ./bin/protoc-gen-go-wasmjs-ts ./cmd/protoc-gen-go-wasmjs-ts

wasm:
	GOOS=wasip1 GOARCH=wasm go build -o ./bin/protoc-gen-go-wasmjs-go.wasm ./cmd/protoc-gen-go-wasmjs-go
	GOOS=wasip1 GOARCH=wasm go build -o ./bin/protoc-gen-go-wasmjs-ts.wasm ./cmd/protoc-gen-go-wasmjs-ts

install:
	go build -o ${GOBIN}/protoc-gen-go-wasmjs-go ./cmd/protoc-gen-go-wasmjs-go
	go build -o ${GOBIN}/protoc-gen-go-wasmjs-ts ./cmd/protoc-gen-go-wasmjs-ts

clean:
	rm -rf ./bin/*
	rm -rf ./examples/*/web/gen/*

# Test default split generators
test: default
	go test -v ./...
