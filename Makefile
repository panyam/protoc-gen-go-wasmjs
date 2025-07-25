
tool:
	go build -o ./bin/protoc-gen-go-wasmjs ./cmd/protoc-gen-go-wasmjs

wasm:
	GOOS=wasip1 GOARCH=wasm go build -o ./bin/protoc-gen-go-wasmjs.wasm ./cmd/protoc-gen-go-wasmjs
