module github.com/panyam/protoc-gen-go-wasmjs/examples/browser-callbacks

go 1.23.5

toolchain go1.24.6

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.2
	github.com/panyam/protoc-gen-go-wasmjs v0.0.0-00010101000000-000000000000
	google.golang.org/genproto/googleapis/api v0.0.0-20250818200422-3122310a409c
	google.golang.org/grpc v1.75.0
	google.golang.org/protobuf v1.36.7
)

require (
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250818200422-3122310a409c // indirect
)

// Use local module during development
replace github.com/panyam/protoc-gen-go-wasmjs => ../..
