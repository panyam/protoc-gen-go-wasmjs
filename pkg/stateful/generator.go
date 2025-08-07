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

	wasmjs "github.com/panyam/protoc-gen-go-wasmjs/proto/gen/go/wasmjs/v1"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// Generator generates stateful TypeScript proxy classes
type Generator struct {
	plugin    *protogen.Plugin
	services  []*StatefulService
	templates *template.Template
}

// StatefulService represents a service marked with stateful annotations
type StatefulService struct {
	Service    *protogen.Service
	Options    *wasmjs.StatefulOptions
	Methods    []*StatefulMethod
	OutputPath string
}

// StatefulMethod represents a method that returns patches
type StatefulMethod struct {
	Method  *protogen.Method
	Options *wasmjs.StatefulMethodOptions
}

// NewGenerator creates a new stateful generator
func NewGenerator(plugin *protogen.Plugin) *Generator {
	return &Generator{
		plugin: plugin,
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
			statefulService := &StatefulService{
				Service:    service,
				Options:    statefulOpts,
				OutputPath: g.getOutputPath(file, service),
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
