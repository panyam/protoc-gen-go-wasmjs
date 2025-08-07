// Copyright 2025 Sri Panyam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"log"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/stateful"
)

var (
	version = flag.Bool("version", false, "print the version and exit")
)

func main() {
	flag.Parse()
	if *version {
		fmt.Printf("protoc-gen-go-wasmjs-stateful %s\n", getVersion())
		return
	}

	var flags flag.FlagSet

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		// Create stateful generator
		statefulGen := stateful.NewGenerator(gen)

		// Generate stateful proxy files
		if err := statefulGen.Generate(); err != nil {
			return fmt.Errorf("failed to generate stateful proxies: %w", err)
		}

		return nil
	})
}

func getVersion() string {
	// TODO: This should be set at build time via ldflags
	return "v0.1.0-dev"
}

func init() {
	log.SetFlags(0)
	log.SetPrefix("protoc-gen-go-wasmjs-stateful: ")
}
