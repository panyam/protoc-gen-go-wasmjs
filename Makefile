
all: tool wasm stateful install

tool:
	go build -o ./bin/protoc-gen-go-wasmjs ./cmd/protoc-gen-go-wasmjs

stateful:
	go build -o ./bin/protoc-gen-go-wasmjs-stateful ./cmd/protoc-gen-go-wasmjs-stateful

wasm:
	GOOS=wasip1 GOARCH=wasm go build -o ./bin/protoc-gen-go-wasmjs.wasm ./cmd/protoc-gen-go-wasmjs

install:
	go build -o ${GOBIN}/protoc-gen-go-wasmjs ./cmd/protoc-gen-go-wasmjs
	go build -o ${GOBIN}/protoc-gen-go-wasmjs-stateful ./cmd/protoc-gen-go-wasmjs-stateful

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
