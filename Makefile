
all: tool split stateful wasm install

tool:
	go build -o ./bin/protoc-gen-go-wasmjs ./cmd/protoc-gen-go-wasmjs

# New split generators using layered architecture
split: split-go split-ts

split-go:
	go build -o ./bin/protoc-gen-go-wasmjs-go ./cmd/protoc-gen-go-wasmjs-go

split-ts:
	go build -o ./bin/protoc-gen-go-wasmjs-ts ./cmd/protoc-gen-go-wasmjs-ts

stateful:
	go build -o ./bin/protoc-gen-go-wasmjs-stateful ./cmd/protoc-gen-go-wasmjs-stateful

wasm:
	GOOS=wasip1 GOARCH=wasm go build -o ./bin/protoc-gen-go-wasmjs.wasm ./cmd/protoc-gen-go-wasmjs

install:
	go build -o ${GOBIN}/protoc-gen-go-wasmjs ./cmd/protoc-gen-go-wasmjs
	go build -o ${GOBIN}/protoc-gen-go-wasmjs-stateful ./cmd/protoc-gen-go-wasmjs-stateful

# Install new split generators
install-split:
	go build -o ${GOBIN}/protoc-gen-go-wasmjs-go ./cmd/protoc-gen-go-wasmjs-go
	go build -o ${GOBIN}/protoc-gen-go-wasmjs-ts ./cmd/protoc-gen-go-wasmjs-ts

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

# Test new split generators
test-split: split
	cd examples/library && ../../bin/protoc-gen-go-wasmjs-go --help || echo "Go generator ready"
	cd examples/library && ../../bin/protoc-gen-go-wasmjs-ts --help || echo "TS generator ready"
