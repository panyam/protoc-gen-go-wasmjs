
all: old test wasm install

# Default generators (split by language) - using layered architecture
default: default-go default-ts

default-go:
	go build -o ./bin/protoc-gen-go-wasmjs-go ./cmd/protoc-gen-go-wasmjs-go

default-ts:
	go build -o ./bin/protoc-gen-go-wasmjs-ts ./cmd/protoc-gen-go-wasmjs-ts

# Old monolithic generator (for backward compatibility)
old:
	go build -o ./bin/protoc-gen-go-wasmjs-old ./cmd/protoc-gen-go-wasmjs-old

wasm:
	GOOS=wasip1 GOARCH=wasm go build -o ./bin/protoc-gen-go-wasmjs-go.wasm ./cmd/protoc-gen-go-wasmjs-go
	GOOS=wasip1 GOARCH=wasm go build -o ./bin/protoc-gen-go-wasmjs-ts.wasm ./cmd/protoc-gen-go-wasmjs-ts
	GOOS=wasip1 GOARCH=wasm go build -o ./bin/protoc-gen-go-wasmjs-old.wasm ./cmd/protoc-gen-go-wasmjs-old

install:
	go build -o ${GOBIN}/protoc-gen-go-wasmjs-go ./cmd/protoc-gen-go-wasmjs-go
	go build -o ${GOBIN}/protoc-gen-go-wasmjs-ts ./cmd/protoc-gen-go-wasmjs-ts

# Install old monolithic generator for backward compatibility
install-old:
	go build -o ${GOBIN}/protoc-gen-go-wasmjs-old ./cmd/protoc-gen-go-wasmjs-old

clean:
	rm -rf ./bin/*
	rm -rf ./examples/*/web/gen/*

# Test default split generators
test: default
	go test -v ./...
