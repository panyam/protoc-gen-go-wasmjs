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
	"fmt"
	"log"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/builders"
)

// GoRenderer handles pure template execution for Go files.
// This is a focused renderer that only executes templates - all file management
// is handled by the generator layer for maximum flexibility and control.
type GoRenderer struct {
	// No file management dependencies - renderer is pure template executor
}

// NewGoRenderer creates a new Go renderer.
// This renderer has no dependencies since it only executes templates.
func NewGoRenderer() *GoRenderer {
	return &GoRenderer{}
}

// RenderToString executes a template with the given data and returns the result as a string.
// This is the core method that all other render methods use - it's a pure template executor.
func (gr *GoRenderer) RenderToString(templateContent string, data interface{}) (string, error) {
	if templateContent == "" {
		return "", fmt.Errorf("template content cannot be empty")
	}

	var buf strings.Builder
	if err := ExecuteTemplate("go", templateContent, data, &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// RenderWasmWrapper generates the Go WASM wrapper and returns it as a string.
// This is a convenience method that uses the embedded WASM template.
func (gr *GoRenderer) RenderWasmWrapper(data *builders.GoTemplateData) (string, error) {
	if data == nil {
		return "", nil // No data to render
	}

	// Validate Go template data before rendering
	if err := gr.ValidateGoTemplateData(data); err != nil {
		return "", fmt.Errorf("invalid WASM wrapper data: %w", err)
	}

	return gr.RenderToString(GoWasmTemplate, data)
}

// RenderWasmWrapperDirect renders WASM wrapper directly to GeneratedFile using old generator pattern.
// This matches exactly how the old generator works.
func (gr *GoRenderer) RenderWasmWrapperDirect(file *protogen.GeneratedFile, data *builders.GoTemplateData) error {
	if file == nil {
		return fmt.Errorf("GeneratedFile cannot be nil")
	}
	if data == nil {
		return nil // No data to render
	}

	// Validate Go template data before rendering
	if err := gr.ValidateGoTemplateData(data); err != nil {
		return fmt.Errorf("invalid WASM wrapper data: %w", err)
	}

	// Execute template and fail early on any errors
	if err := ExecuteTemplateToFile("wasm", GoWasmTemplate, data, file); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	// Note: We cannot call file.Content() here as it might interfere with protogen's
	// internal state. The protogen framework will handle the content properly.
	// The fact that ExecuteTemplateToFile succeeded means content was written.
	log.Printf("WASM: Template rendered successfully")
	return nil
}

// RenderMainExampleDirect renders main example directly to GeneratedFile using old generator pattern.
func (gr *GoRenderer) RenderMainExampleDirect(file *protogen.GeneratedFile, data *builders.GoTemplateData) error {
	if file == nil {
		return fmt.Errorf("GeneratedFile cannot be nil")
	}
	if data == nil {
		return nil // No data to render
	}

	// Validate Go template data before rendering
	if err := gr.ValidateGoTemplateData(data); err != nil {
		return fmt.Errorf("invalid main example data: %w", err)
	}

	// Execute template and fail early on any errors
	if err := ExecuteTemplateToFile("main", GoMainTemplate, data, file); err != nil {
		return fmt.Errorf("main template execution failed: %w", err)
	}

	// Note: We cannot call file.Content() here as it might interfere with protogen's
	// internal state. The protogen framework will handle the content properly.
	log.Printf("MAIN: Template rendered successfully")
	return nil
}

// RenderBuildScriptDirect renders build script directly to GeneratedFile using old generator pattern.
func (gr *GoRenderer) RenderBuildScriptDirect(file *protogen.GeneratedFile, data *builders.GoTemplateData) error {
	if file == nil {
		return fmt.Errorf("GeneratedFile cannot be nil")
	}
	if data == nil {
		return nil // No data to render
	}

	// Validate Go template data before rendering
	if err := gr.ValidateGoTemplateData(data); err != nil {
		return fmt.Errorf("invalid build script data: %w", err)
	}

	// Execute template and fail early on any errors
	if err := ExecuteTemplateToFile("build", GoBuildScriptTemplate, data, file); err != nil {
		return fmt.Errorf("build script template execution failed: %w", err)
	}

	// Note: We cannot call file.Content() here as it might interfere with protogen's
	// internal state. The protogen framework will handle the content properly.
	log.Printf("BUILD: Template rendered successfully")
	return nil
}

// RenderMainExample generates an example main.go file and returns it as a string.
// This is a convenience method that uses the embedded main template.
func (gr *GoRenderer) RenderMainExample(data *builders.GoTemplateData) (string, error) {
	if data == nil {
		return "", nil
	}

	// Validate Go template data before rendering
	if err := gr.ValidateGoTemplateData(data); err != nil {
		return "", fmt.Errorf("invalid main example data: %w", err)
	}

	return gr.RenderToString(GoMainTemplate, data)
}

// RenderBuildScript generates a build script and returns it as a string.
// This is a convenience method that uses the embedded build script template.
func (gr *GoRenderer) RenderBuildScript(data *builders.GoTemplateData) (string, error) {
	if data == nil {
		return "", nil
	}

	// Validate Go template data before rendering
	if err := gr.ValidateGoTemplateData(data); err != nil {
		return "", fmt.Errorf("invalid build script data: %w", err)
	}

	return gr.RenderToString(GoBuildScriptTemplate, data)
}

// ValidateGoTemplateData performs validation on Go template data before rendering.
// This catches common issues that would cause template execution failures.
func (gr *GoRenderer) ValidateGoTemplateData(data *builders.GoTemplateData) error {
	if data == nil {
		return nil // Nil data is valid (means no generation needed)
	}

	// Basic validation
	if data.PackageName == "" {
		return fmt.Errorf("GoTemplateData.PackageName cannot be empty")
	}

	if data.ModuleName == "" {
		return fmt.Errorf("GoTemplateData.ModuleName cannot be empty")
	}

	// Validate services have methods
	for i, service := range data.Services {
		if len(service.Methods) == 0 {
			return fmt.Errorf("service %s at index %d has no methods", service.Name, i)
		}

		if service.GoType == "" {
			return fmt.Errorf("service %s has empty GoType", service.Name)
		}
	}

	// Validate browser services
	for i, service := range data.BrowserServices {
		if !service.IsBrowserProvided {
			return fmt.Errorf("browser service %s at index %d has IsBrowserProvided=false", service.Name, i)
		}
	}

	return nil
}
