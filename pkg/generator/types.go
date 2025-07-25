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
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

// TemplateData holds all data needed for template generation
type TemplateData struct {
	// Core metadata
	PackageName string
	SourcePath  string
	GoPackage   string
	Services    []ServiceData
	Config      *Config

	// Customization fields
	JSNamespace  string
	ModuleName   string
	APIStructure string // namespaced|flat|service_based

	// Import management
	Imports     []ImportInfo // Unique package imports with aliases
	PackageMap  map[string]string // Maps full package path to alias

	// TypeScript import management
	TSRelativeImportPath string // Relative path from TS export to TS import

	// Build info
	GeneratedImports []string
	BuildScript      string
}

// ImportInfo represents a Go package import with alias
type ImportInfo struct {
	Path  string // Full import path
	Alias string // Package alias (e.g., "libraryv1")
}

// ServiceData represents a gRPC service for template generation
type ServiceData struct {
	Name         string       // Service name (e.g., "LibraryService")
	GoType       string       // Go type for the service interface (qualified, e.g., "libraryv1.LibraryServiceServer")
	JSName       string       // JavaScript name (e.g., "library")
	Methods      []MethodData // All methods in the service
	PackagePath  string       // Go package path for imports
	PackageAlias string       // Package alias for qualified type names
}

// MethodData represents a gRPC method for template generation
type MethodData struct {
	Name           string // Original protobuf method name (e.g., "FindBooks")
	JSName         string // JavaScript method name (e.g., "findBooks")
	GoFuncName     string // Go function name for WASM wrapper (e.g., "libraryFindBooks")
	ShouldGenerate bool   // Whether to generate this method based on filters
	RequestType    string // Fully qualified Go request type
	ResponseType   string // Fully qualified Go response type
	Comment        string // Method comment from protobuf
}

// FileGenerator handles generation for a single proto file
type FileGenerator struct {
	file   *protogen.File
	plugin *protogen.Plugin
	config *Config
}

// NewFileGenerator creates a new file generator
func NewFileGenerator(file *protogen.File, plugin *protogen.Plugin, config *Config) *FileGenerator {
	return &FileGenerator{
		file:   file,
		plugin: plugin,
		config: config,
	}
}

// Generate generates WASM wrapper and TypeScript client for the proto file
func (g *FileGenerator) Generate() error {
	// Validate configuration
	if err := g.config.Validate(); err != nil {
		return err
	}

	// Skip files with no services
	if len(g.file.Services) == 0 {
		return nil
	}

	// Build template data
	templateData, err := g.buildTemplateData()
	if err != nil {
		return err
	}

	// Generate WASM wrapper
	if err := g.generateWasmWrapper(templateData); err != nil {
		return err
	}

	// Generate TypeScript client
	if err := g.generateTypeScriptClient(templateData); err != nil {
		return err
	}

	// Generate build script if requested
	if g.config.GenerateBuildScript {
		if err := g.generateBuildScript(templateData); err != nil {
			return err
		}
	}

	return nil
}

// buildTemplateData constructs the template data from the proto file
func (g *FileGenerator) buildTemplateData() (*TemplateData, error) {
	packageName := string(g.file.Desc.Package())

	// Collect unique package imports
	imports, packageMap := g.collectUniqueImports()

	// Build service data
	var services []ServiceData
	for _, service := range g.file.Services {
		serviceData := g.buildServiceData(service, packageMap)
		if serviceData != nil {
			services = append(services, *serviceData)
		}
	}

	// Skip if no services to generate
	if len(services) == 0 {
		return nil, nil
	}

	return &TemplateData{
		PackageName:          packageName,
		SourcePath:           g.file.Desc.Path(),
		GoPackage:            string(g.file.GoPackageName),
		Services:             services,
		Config:               g.config,
		JSNamespace:          g.config.GetDefaultJSNamespace(packageName),
		ModuleName:           g.config.GetDefaultModuleName(packageName),
		APIStructure:         g.config.JSStructure,
		Imports:              imports,
		PackageMap:           packageMap,
		TSRelativeImportPath: g.config.GetRelativeTSImportPathForProto(g.file.Desc.Path()),
	}, nil
}

// collectUniqueImports collects unique package imports from all services and messages
func (g *FileGenerator) collectUniqueImports() ([]ImportInfo, map[string]string) {
	packagePaths := make(map[string]bool)
	packageMap := make(map[string]string)
	var imports []ImportInfo

	// Collect all unique package paths from services
	for _, service := range g.file.Services {
		if !g.config.ShouldGenerateService(string(service.Desc.Name())) {
			continue
		}

		packagePath := string(g.file.GoImportPath)
		if !packagePaths[packagePath] {
			packagePaths[packagePath] = true
			alias := g.generatePackageAlias(packagePath)
			imports = append(imports, ImportInfo{
				Path:  packagePath,
				Alias: alias,
			})
			packageMap[packagePath] = alias
		}
	}

	return imports, packageMap
}

// generatePackageAlias creates a package alias from a Go import path
func (g *FileGenerator) generatePackageAlias(packagePath string) string {
	// Extract package name from the path (e.g., "library/v1" -> "libraryv1")
	parts := strings.Split(packagePath, "/")
	if len(parts) >= 2 {
		// Take last two parts and combine them
		pkg := parts[len(parts)-2]
		version := parts[len(parts)-1]
		// Remove dots and slashes to create valid Go identifier
		alias := strings.ReplaceAll(pkg+version, ".", "")
		alias = strings.ReplaceAll(alias, "/", "")
		return alias
	}
	// Fallback to last part only
	last := parts[len(parts)-1]
	return strings.ReplaceAll(last, ".", "")
}

// buildServiceData constructs service data for a protobuf service
func (g *FileGenerator) buildServiceData(service *protogen.Service, packageMap map[string]string) *ServiceData {
	serviceName := string(service.Desc.Name())

	// Check if this service should be generated
	if !g.config.ShouldGenerateService(serviceName) {
		return nil
	}

	// Get package alias for this service
	packagePath := string(g.file.GoImportPath)
	packageAlias := packageMap[packagePath]

	// Build method data
	var methods []MethodData
	for _, method := range service.Methods {
		methodData := g.buildMethodData(method, serviceName, packageAlias)
		if methodData != nil {
			methods = append(methods, *methodData)
		}
	}

	// Skip service if no methods to generate
	if len(methods) == 0 {
		return nil
	}

	return &ServiceData{
		Name:         serviceName,
		GoType:       packageAlias + "." + string(service.GoName) + "Server", // Qualified type name
		JSName:       strings.ToLower(serviceName[:1]) + serviceName[1:],     // camelCase
		Methods:      methods,
		PackagePath:  packagePath,
		PackageAlias: packageAlias,
	}
}

// buildMethodData constructs method data for a protobuf method
func (g *FileGenerator) buildMethodData(method *protogen.Method, serviceName, packageAlias string) *MethodData {
	methodName := string(method.Desc.Name())

	// Check if this method should be generated
	if !g.config.ShouldGenerateMethod(methodName) {
		return nil
	}

	// Only support unary methods for now
	if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
		return nil
	}

	jsName := g.config.GetMethodJSName(methodName)
	goFuncName := strings.ToLower(serviceName[:1]) + serviceName[1:] + methodName

	return &MethodData{
		Name:           methodName,
		JSName:         jsName,
		GoFuncName:     goFuncName,
		ShouldGenerate: true,
		RequestType:    g.getQualifiedTypeName(method.Input, packageAlias),
		ResponseType:   g.getQualifiedTypeName(method.Output, packageAlias),
		Comment:        strings.TrimSpace(string(method.Comments.Leading)),
	}
}

// getQualifiedTypeName returns the fully qualified Go type name for a message
func (g *FileGenerator) getQualifiedTypeName(message *protogen.Message, packageAlias string) string {
	// Return the qualified type name with package alias
	return "*" + packageAlias + "." + string(message.GoIdent.GoName)
}
