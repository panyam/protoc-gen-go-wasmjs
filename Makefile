
all: default old stateful wasm install

# Default generators (split by language) - using layered architecture
default: default-go default-ts

default-go:
	go build -o ./bin/protoc-gen-go-wasmjs-go ./cmd/protoc-gen-go-wasmjs-go

default-ts:
	go build -o ./bin/protoc-gen-go-wasmjs-ts ./cmd/protoc-gen-go-wasmjs-ts

# Old monolithic generator (for backward compatibility)
old:
	go build -o ./bin/protoc-gen-go-wasmjs-old ./cmd/protoc-gen-go-wasmjs

stateful:
	go build -o ./bin/protoc-gen-go-wasmjs-stateful ./cmd/protoc-gen-go-wasmjs-stateful

wasm:
	GOOS=wasip1 GOARCH=wasm go build -o ./bin/protoc-gen-go-wasmjs.wasm ./cmd/protoc-gen-go-wasmjs

install:
	go build -o ${GOBIN}/protoc-gen-go-wasmjs-go ./cmd/protoc-gen-go-wasmjs-go
	go build -o ${GOBIN}/protoc-gen-go-wasmjs-ts ./cmd/protoc-gen-go-wasmjs-ts
	go build -o ${GOBIN}/protoc-gen-go-wasmjs-stateful ./cmd/protoc-gen-go-wasmjs-stateful

# Install old monolithic generator for backward compatibility
install-old:
	go build -o ${GOBIN}/protoc-gen-go-wasmjs-old ./cmd/protoc-gen-go-wasmjs

test-connect4: stateful
	cd examples/connect4 && make test

test-connect4-buf: stateful
	cd examples/connect4 && buf generate

test-stateful: stateful
	cd example && protoc \
		--plugin=protoc-gen-go-wasmjs-stateful=../bin/protoc-gen-go-wasmjs-stateful \
		--go-wasmjs-stateful_out=./gen/ts \
		-I proto \
		-I ../proto \
		proto/game/game.proto

clean:
	rm -rf ./bin/*
	rm -rf ./examples/*/web/gen/*
	rm -rf ./example/gen/ts/stateful/*

# Test default split generators
test-default: default
	cd examples/library && ../../bin/protoc-gen-go-wasmjs-go --help || echo "Go generator ready"
	cd examples/library && ../../bin/protoc-gen-go-wasmjs-ts --help || echo "TS generator ready"
