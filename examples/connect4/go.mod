module github.com/panyam/protoc-gen-go-wasmjs/examples/connect4

go 1.24.6

replace github.com/panyam/protoc-gen-go-wasmjs => ../../

require (
	github.com/panyam/goutils v0.1.9
	github.com/panyam/templar v0.0.20
	google.golang.org/grpc v1.74.2
	google.golang.org/protobuf v1.36.7
)

require (
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
)
