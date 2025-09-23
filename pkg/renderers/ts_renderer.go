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

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/builders"
)

// TSRenderer handles pure template execution for TypeScript files.
// This is a focused renderer that only executes templates - all file management
// is handled by the generator layer for maximum flexibility and control.
type TSRenderer struct {
	// No file management dependencies - renderer is pure template executor
}

// NewTSRenderer creates a new TypeScript renderer.
// This renderer has no dependencies since it only executes templates.
func NewTSRenderer() *TSRenderer {
	return &TSRenderer{}
}

// RenderToFile executes a template with the given data and writes to the provided GeneratedFile.
// This is the core method that all other render methods use - it's a pure template executor.
func (tr *TSRenderer) RenderToFile(file *protogen.GeneratedFile, templateContent string, data interface{}) error {
	if file == nil {
		return fmt.Errorf("GeneratedFile cannot be nil")
	}

	if templateContent == "" {
		return fmt.Errorf("template content cannot be empty")
	}

	// Parse and execute template using shared helpers
	return ExecuteTemplate("typescript", templateContent, data, file)
}

// RenderBrowserServiceManager is no longer needed as BrowserServiceManager
// is now imported from @protoc-gen-go-wasmjs/runtime package
// Removed: RenderBrowserServiceManager() - now imported from runtime package

// RenderServiceClient generates a TypeScript client for a single service using the provided GeneratedFile.
// Uses the same template as RenderClient but for individual services.
func (tr *TSRenderer) RenderServiceClient(file *protogen.GeneratedFile, data *builders.TSTemplateData) error {
	return tr.RenderClient(file, data)
}

// RenderClient generates the TypeScript client using the provided GeneratedFile.
// This is a convenience method that uses the embedded client template.
func (tr *TSRenderer) RenderClient(file *protogen.GeneratedFile, data *builders.TSTemplateData) error {
	if data == nil {
		return nil // No data to render
	}

	// Validate TypeScript template data before rendering
	if err := tr.ValidateTSTemplateData(data); err != nil {
		return fmt.Errorf("invalid client data: %w", err)
	}

	return tr.RenderToFile(file, TSSimpleClientTemplate, data)
}

// RenderInterfaces generates TypeScript interface definitions using the provided GeneratedFile.
// This is a convenience method that uses the embedded interfaces template.
func (tr *TSRenderer) RenderInterfaces(file *protogen.GeneratedFile, data *builders.TSTemplateData) error {
	if data == nil {
		return nil
	}

	// Validate TypeScript template data before rendering
	if err := tr.ValidateTSTemplateData(data); err != nil {
		return fmt.Errorf("invalid interfaces data: %w", err)
	}

	return tr.RenderToFile(file, TSInterfacesTemplate, data)
}

// RenderModels generates TypeScript model class implementations using the provided GeneratedFile.
// This is a convenience method that uses the embedded models template.
func (tr *TSRenderer) RenderModels(file *protogen.GeneratedFile, data *builders.TSTemplateData) error {
	if data == nil {
		return nil
	}

	// Validate TypeScript template data before rendering
	if err := tr.ValidateTSTemplateData(data); err != nil {
		return fmt.Errorf("invalid models data: %w", err)
	}

	return tr.RenderToFile(file, TSModelsTemplate, data)
}

// RenderFactory generates TypeScript factory classes using the provided GeneratedFile.
// This is a convenience method that uses the embedded factory template.
func (tr *TSRenderer) RenderFactory(file *protogen.GeneratedFile, data *builders.TSTemplateData) error {
	if data == nil {
		return nil
	}

	// Validate TypeScript template data before rendering
	if err := tr.ValidateTSTemplateData(data); err != nil {
		return fmt.Errorf("invalid factory data: %w", err)
	}

	return tr.RenderToFile(file, TSFactoryTemplate, data)
}

// RenderSchemas generates TypeScript schema definitions using the provided GeneratedFile.
// This is a convenience method that uses the embedded schemas template.
func (tr *TSRenderer) RenderSchemas(file *protogen.GeneratedFile, data *builders.TSTemplateData) error {
	if data == nil {
		return nil
	}

	// Validate TypeScript template data before rendering
	if err := tr.ValidateTSTemplateData(data); err != nil {
		return fmt.Errorf("invalid schemas data: %w", err)
	}

	return tr.RenderToFile(file, TSSchemaTemplate, data)
}

// RenderDeserializer generates TypeScript deserializer classes using the provided GeneratedFile.
// This is a convenience method that uses the embedded deserializer template.
func (tr *TSRenderer) RenderDeserializer(file *protogen.GeneratedFile, data *builders.TSTemplateData) error {
	if data == nil {
		return nil
	}

	// Validate TypeScript template data before rendering
	if err := tr.ValidateTSTemplateData(data); err != nil {
		return fmt.Errorf("invalid deserializer data: %w", err)
	}

	return tr.RenderToFile(file, TSDeserializerTemplate, data)
}

// RenderBundle generates TypeScript bundle class using the provided GeneratedFile.
// This renders the shared WASM bundle that manages all service clients for a module.
func (tr *TSRenderer) RenderBundle(file *protogen.GeneratedFile, data *builders.TSTemplateData) error {
	if data == nil {
		return nil
	}

	// Validate bundle data with bundle-specific validation (allows empty methods)
	if err := tr.ValidateBundleTemplateData(data); err != nil {
		return fmt.Errorf("invalid bundle data: %w", err)
	}

	return tr.RenderToFile(file, TSBundleTemplate, data)
}

// ValidateBundleTemplateData validates TSTemplateData specifically for bundle rendering.
// Bundle validation is less strict since bundles don't use method data.
func (tr *TSRenderer) ValidateBundleTemplateData(data *builders.TSTemplateData) error {
	if data == nil {
		return nil
	}
	
	// Basic validation (same as regular validation)
	if data.PackageName == "" {
		return fmt.Errorf("TSTemplateData.PackageName cannot be empty")
	}
	
	if data.PackagePath == "" {
		return fmt.Errorf("TSTemplateData.PackagePath cannot be empty")
	}
	
	// For bundles, we don't validate that services have methods
	// Services are just used for import/property generation
	for _, service := range data.Services {
		if service.Name == "" {
			return fmt.Errorf("service has empty Name")
		}
		
		if service.JSName == "" {
			return fmt.Errorf("service %s has empty JSName", service.Name)
		}
	}
	
	return nil
}

// RenderBrowserServices generates TypeScript browser service interfaces using the provided GeneratedFile.
// This renders interfaces that browser implementations must implement.
func (tr *TSRenderer) RenderBrowserServices(file *protogen.GeneratedFile, data *builders.TSTemplateData) error {
	if data == nil {
		return nil
	}

	// Validate TypeScript template data before rendering
	if err := tr.ValidateTSTemplateData(data); err != nil {
		return fmt.Errorf("invalid browser services data: %w", err)
	}

	return tr.RenderToFile(file, TSBrowserServiceTemplate, data)
}

// ValidateTSTemplateData performs validation on TypeScript template data before rendering.
// This catches common issues that would cause template execution failures.
func (tr *TSRenderer) ValidateTSTemplateData(data *builders.TSTemplateData) error {
	if data == nil {
		return nil // Nil data is valid (means no generation needed)
	}

	// Basic validation
	if data.PackageName == "" {
		return fmt.Errorf("TSTemplateData.PackageName cannot be empty")
	}

	if data.PackagePath == "" {
		return fmt.Errorf("TSTemplateData.PackagePath cannot be empty")
	}

	// Validate services have proper names (methods are optional)
	for _, service := range data.Services {
		if service.Name == "" {
			return fmt.Errorf("service has empty Name")
		}

		if service.JSName == "" {
			return fmt.Errorf("service %s has empty JSName", service.Name)
		}
		
		// Note: We now allow services with 0 methods (e.g., browser services with no RPCs)
		// The service client will be generated but have no callable methods
	}

	// Validate messages have names
	for i, message := range data.Messages {
		if message.Name == "" {
			return fmt.Errorf("message at index %d has empty name", i)
		}
	}

	// Validate enums have names and values
	for i, enum := range data.Enums {
		if enum.Name == "" {
			return fmt.Errorf("enum at index %d has empty name", i)
		}
		if len(enum.Values) == 0 {
			return fmt.Errorf("enum %s has no values", enum.Name)
		}
	}

	return nil
}
