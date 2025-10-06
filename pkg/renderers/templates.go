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

package renderers

import (
	_ "embed"
)

// Go WASM templates
//
//go:embed templates/wasm_converters.go.tmpl
var GoConvertersTemplate string

//go:embed templates/wasm_exports.go.tmpl
var GoExportsTemplate string

//go:embed templates/wasm_browser_clients.go.tmpl
var GoBrowserClientsTemplate string

//go:embed templates/wasm.go.tmpl
var GoWasmTemplate string

//go:embed templates/main.go.tmpl
var GoMainTemplate string

//go:embed templates/build.sh.tmpl
var GoBuildScriptTemplate string

// TypeScript templates
//
//go:embed templates/client_simple.ts.tmpl
var TSSimpleClientTemplate string

// Removed: TSBrowserServiceManagerTemplate (now imported from @protoc-gen-go-wasmjs/runtime)

//go:embed templates/interfaces.ts.tmpl
var TSInterfacesTemplate string

//go:embed templates/models.ts.tmpl
var TSModelsTemplate string

//go:embed templates/factory.ts.tmpl
var TSFactoryTemplate string

//go:embed templates/schemas.ts.tmpl
var TSSchemaTemplate string

//go:embed templates/deserializer.ts.tmpl
var TSDeserializerTemplate string

//go:embed templates/bundle.ts.tmpl
var TSBundleTemplate string

//go:embed templates/browser_service.ts.tmpl
var TSBrowserServiceTemplate string

// Removed: TSDeserializerSchemasTemplate (now imported from @protoc-gen-go-wasmjs/runtime)

// Removed: TSClientTemplate (unused dead code)
