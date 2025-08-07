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

package stateful

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/generator"
	wasmjs "github.com/panyam/protoc-gen-go-wasmjs/proto/gen/go/wasmjs/v1"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// Config holds configuration for the stateful generator
type Config struct {
	ClientImportPath string
}

// Generator generates stateful TypeScript proxy classes
type Generator struct {
	plugin    *protogen.Plugin
	config    *Config
	services  []*StatefulService
	templates *template.Template
}

// StatefulService represents a service marked with stateful annotations
type StatefulService struct {
	Service          *protogen.Service
	Options          *wasmjs.StatefulOptions
	Methods          []*StatefulMethod
	OutputPath       string
	ClientImportPath string
}

// StatefulMethod represents a method that returns patches
type StatefulMethod struct {
	Method  *protogen.Method
	Options *wasmjs.StatefulMethodOptions
}

// NewGenerator creates a new stateful generator with default config
func NewGenerator(plugin *protogen.Plugin) *Generator {
	return NewGeneratorWithConfig(plugin, &Config{})
}

// NewGeneratorWithConfig creates a new stateful generator with custom config
func NewGeneratorWithConfig(plugin *protogen.Plugin, config *Config) *Generator {
	if config == nil {
		config = &Config{}
	}
	return &Generator{
		plugin: plugin,
		config: config,
	}
}

// Generate generates all stateful proxy files
func (g *Generator) Generate() error {
	// Load templates
	if err := g.loadTemplates(); err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	// Discover stateful services
	if err := g.discoverStatefulServices(); err != nil {
		return fmt.Errorf("failed to discover stateful services: %w", err)
	}

	// Skip generation if no stateful services found
	if len(g.services) == 0 {
		return nil
	}

	// Generate TypeScript files
	if err := g.generateTypeScriptFiles(); err != nil {
		return fmt.Errorf("failed to generate TypeScript files: %w", err)
	}

	return nil
}

// loadTemplates loads the embedded template files
func (g *Generator) loadTemplates() error {
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
		"camelCase": func(s string) string {
			if len(s) == 0 {
				return s
			}
			return strings.ToLower(s[:1]) + s[1:]
		},
	}

	tmpl, err := template.New("stateful").Funcs(funcMap).ParseFS(templateFS, "templates/*.tmpl")
	if err != nil {
		return err
	}

	g.templates = tmpl
	return nil
}

// discoverStatefulServices finds all services with stateful annotations
func (g *Generator) discoverStatefulServices() error {
	for _, file := range g.plugin.Files {
		if !file.Generate {
			continue
		}

		for _, service := range file.Services {
			statefulOpts := g.getStatefulOptions(service)
			if statefulOpts == nil || !statefulOpts.GetEnabled() {
				continue
			}

			// Create stateful service
			outputPath := g.getOutputPath(file, service)
			clientImportPath := g.calculateClientImportPath(outputPath, file)

			statefulService := &StatefulService{
				Service:          service,
				Options:          statefulOpts,
				OutputPath:       outputPath,
				ClientImportPath: clientImportPath,
			}

			// Find stateful methods
			for _, method := range service.Methods {
				methodOpts := g.getStatefulMethodOptions(method)
				if methodOpts != nil {
					statefulService.Methods = append(statefulService.Methods, &StatefulMethod{
						Method:  method,
						Options: methodOpts,
					})
				}
			}

			g.services = append(g.services, statefulService)
		}
	}

	return nil
}

// getStatefulOptions extracts stateful options from service
func (g *Generator) getStatefulOptions(service *protogen.Service) *wasmjs.StatefulOptions {
	opts := service.Desc.Options()
	if opts == nil {
		return nil
	}

	if !proto.HasExtension(opts, wasmjs.E_Stateful) {
		return nil
	}

	ext := proto.GetExtension(opts, wasmjs.E_Stateful)
	if statefulOpts, ok := ext.(*wasmjs.StatefulOptions); ok {
		return statefulOpts
	}

	return nil
}

// getStatefulMethodOptions extracts stateful method options from method
func (g *Generator) getStatefulMethodOptions(method *protogen.Method) *wasmjs.StatefulMethodOptions {
	opts := method.Desc.Options()
	if opts == nil {
		return nil
	}

	if !proto.HasExtension(opts, wasmjs.E_StatefulMethod) {
		return nil
	}

	ext := proto.GetExtension(opts, wasmjs.E_StatefulMethod)
	if methodOpts, ok := ext.(*wasmjs.StatefulMethodOptions); ok {
		return methodOpts
	}

	return nil
}

// getOutputPath generates the output path for a service's stateful files
func (g *Generator) getOutputPath(file *protogen.File, service *protogen.Service) string {
	// Base directory from protogen
	dir := filepath.Dir(file.GeneratedFilenamePrefix)

	// Create stateful subdirectory
	statefulDir := filepath.Join(dir, "stateful")

	// Generate filename: service_name_stateful.ts
	serviceName := strings.ToLower(string(service.Desc.Name()))
	filename := fmt.Sprintf("%s_stateful.ts", serviceName)

	return filepath.Join(statefulDir, filename)
}

// calculateClientImportPath calculates the relative import path from stateful proxy to WASM client
func (g *Generator) calculateClientImportPath(statefulOutputPath string, file *protogen.File) string {
	// If client import path is configured, use it to calculate relative path
	if g.config.ClientImportPath != "" {
		// Use the helper from generator package
		return generator.CalculateRelativePath(statefulOutputPath, g.config.ClientImportPath)
	}

	// Default: calculate relative path to the WASM client
	// This assumes the WASM client is in the same directory as the proto-generated files
	dir := filepath.Dir(file.GeneratedFilenamePrefix)

	// Default WASM client filename (following the pattern from main generator)
	baseName := strings.Replace(filepath.Base(file.Desc.Path()), ".proto", "", 1)
	// Capitalize first letter (simple replacement for strings.Title)
	if len(baseName) > 0 {
		baseName = strings.ToUpper(baseName[:1]) + baseName[1:]
	}
	wasmClientFile := filepath.Join(dir, fmt.Sprintf("%sClient.client.ts", baseName))

	return generator.CalculateRelativePath(statefulOutputPath, wasmClientFile)
}

// generateTypeScriptFiles generates all TypeScript proxy files
func (g *Generator) generateTypeScriptFiles() error {
	// Generate shared patch types (once)
	if err := g.generatePatchTypes(); err != nil {
		return err
	}

	// Generate individual service proxies
	for _, service := range g.services {
		if err := g.generateStatefulProxy(service); err != nil {
			return err
		}
	}

	return nil
}

// generatePatchTypes generates the shared patch type definitions
func (g *Generator) generatePatchTypes() error {
	// Create patches.ts file
	patchFile := g.plugin.NewGeneratedFile("stateful/patches.ts", "")

	data := struct {
		PackageName string
	}{
		PackageName: "stateful",
	}

	if err := g.templates.ExecuteTemplate(patchFile, "patches.ts.tmpl", data); err != nil {
		return fmt.Errorf("failed to execute patches template: %w", err)
	}

	return nil
}

// generateStatefulProxy generates TypeScript code for a stateful proxy
func (g *Generator) generateStatefulProxy(service *StatefulService) error {
	// Create the generated file
	file := g.plugin.NewGeneratedFile(service.OutputPath, "")

	if err := g.templates.ExecuteTemplate(file, "stateful_proxy.ts.tmpl", service); err != nil {
		return fmt.Errorf("failed to execute stateful proxy template for %s: %w", service.Service.GoName, err)
	}

	return nil
}
