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

	// Determine output filename
	// baseName := strings.TrimSuffix(filepath.Base(g.file.Desc.Path()), ".proto")
	filename := filepath.Join(g.config.WasmExportPath, data.ModuleName+WasmGeneratedFilenameExtension)

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

	// Determine output filename - use TSOut if specified, otherwise co-locate with WASM
	var filename string
	if g.config.TSOut != "" {
		// TSOut should be relative to the buf.gen.yaml working directory, not the protoc out directory
		// Calculate relative path from WasmExportPath (out directory) to TSOut
		relativeToOut, err := filepath.Rel(g.config.WasmExportPath, g.config.TSOut)
		if err != nil {
			// Fallback to absolute path if calculation fails
			relativeToOut = g.config.TSOut
		}
		filename = filepath.Join(relativeToOut, data.ModuleName+"Client"+TSGeneratedFilenameExtension)
	} else {
		filename = data.ModuleName + "Client" + TSGeneratedFilenameExtension
	}

	// Create generated file
	generatedFile := g.plugin.NewGeneratedFile(filename, "")

	// Parse and execute template
	tmpl, err := template.New("typescript").Funcs(templateFuncMap).Parse(typescriptTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(generatedFile, data)
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
