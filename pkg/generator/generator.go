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
	"path/filepath"
	"strings"
	"text/template"
)

// Template helper functions
var templateFuncMap = template.FuncMap{
	"title": func(s string) string {
		if len(s) == 0 {
			return s
		}
		return strings.ToUpper(s[:1]) + s[1:]
	},
	"replaceAll": func(old, new, s string) string {
		return strings.ReplaceAll(s, old, new)
	},
}

const (
	WasmGeneratedFilenameExtension = ".wasm.go"
	TSGeneratedFilenameExtension   = ".client.ts"
	BuildScriptFilename            = "build.sh"
	MainExampleFilename            = "main.go.example"
)

// generateWasmWrapper generates the Go WASM wrapper file
func (g *FileGenerator) generateWasmWrapper(data *TemplateData) error {
	if data == nil {
		return nil // No services to generate
	}

	// Determine output filename using package structure
	packagePath := strings.ReplaceAll(data.PackageName, ".", "/")
	baseName := strings.ReplaceAll(data.PackageName, ".", "_")
	filename := filepath.Join(packagePath, baseName+WasmGeneratedFilenameExtension)

	// Create generated file
	generatedFile := g.plugin.NewGeneratedFile(filename, "")

	// Parse and execute template
	tmpl, err := template.New("wasm").Funcs(templateFuncMap).Parse(wasmTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(generatedFile, data)
}

// generateTypeScriptClient generates the TypeScript client file
func (g *FileGenerator) generateTypeScriptClient(data *TemplateData) error {
	if data == nil {
		return nil // No services to generate
	}

	// TypeScript client co-locates with other artifacts in the protoc out directory
	filename := data.ModuleName + "Client" + TSGeneratedFilenameExtension

	// Create generated file
	generatedFile := g.plugin.NewGeneratedFile(filename, "")

	// Parse and execute template
	tmpl, err := template.New("typescript").Funcs(templateFuncMap).Parse(typescriptTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(generatedFile, data)
}

// generateBrowserServiceManager generates the shared BrowserServiceManager TypeScript file
func (g *FileGenerator) generateBrowserServiceManager() error {
	// Always generate in the root of output directory
	filename := "browserServiceManager.ts"

	// Create generated file
	generatedFile := g.plugin.NewGeneratedFile(filename, "")

	// Parse and execute template
	tmpl, err := template.New("browserServiceManager").Funcs(templateFuncMap).Parse(browserServiceManagerTemplate)
	if err != nil {
		return err
	}

	// No data needed for this template
	return tmpl.Execute(generatedFile, nil)
}

// generateBuildScript generates a build script for compiling the WASM wrapper
func (g *FileGenerator) generateBuildScript(data *TemplateData) error {
	if data == nil {
		return nil // No services to generate
	}

	// Create build script in WASM export directory
	filename := filepath.Join(g.config.WasmExportPath, BuildScriptFilename)

	// Create generated file
	generatedFile := g.plugin.NewGeneratedFile(filename, "")

	// Parse and execute template
	tmpl, err := template.New("buildscript").Funcs(templateFuncMap).Parse(buildScriptTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(generatedFile, data)
}

// generateMainExample generates an example main.go file that users can copy and customize
func (g *FileGenerator) generateMainExample(data *TemplateData) error {
	if data == nil {
		return nil // No services to generate
	}

	// Create main example in WASM export directory
	filename := filepath.Join(g.config.WasmExportPath, MainExampleFilename)

	// Create generated file
	generatedFile := g.plugin.NewGeneratedFile(filename, "")

	// Parse and execute template
	tmpl, err := template.New("mainexample").Funcs(templateFuncMap).Parse(mainExampleTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(generatedFile, data)
}
