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

package generator

import (
	_ "embed"
)

// Embedded template files
//go:embed templates/wasm.go.tmpl
var wasmTemplate string

//go:embed templates/client_simple.ts.tmpl
var typescriptTemplate string

//go:embed templates/browser_service_manager.ts.tmpl
var browserServiceManagerTemplate string

//go:embed templates/build.sh.tmpl
var buildScriptTemplate string

//go:embed templates/main.go.tmpl
var mainExampleTemplate string

// New TypeScript generation templates
//go:embed templates/interfaces.ts.tmpl
var interfacesTemplate string

//go:embed templates/models.ts.tmpl
var modelsTemplate string

//go:embed templates/factory.ts.tmpl
var factoryTemplate string

//go:embed templates/schemas.ts.tmpl
var schemasTemplate string

//go:embed templates/deserializer.ts.tmpl
var deserializerTemplate string

//go:embed templates/deserializer_schemas.ts.tmpl  
var deserializerSchemasTemplate string