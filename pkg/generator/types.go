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
	TSImports            []TSImportInfo // TypeScript imports with proper extensions

	// Build info
	GeneratedImports []string
	BuildScript      string
}

// ImportInfo represents a Go package import with alias
type ImportInfo struct {
	Path  string // Full import path
	Alias string // Package alias (e.g., "libraryv1")
}

// TSImportInfo represents a TypeScript import for the client
type TSImportInfo struct {
	ProtoFile string   // Original proto file (e.g., "games.proto")
	ImportPath string  // Relative import path with proper extension
	Types     []string // List of types to import from this file
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
	RequestTSType  string // TypeScript request type name
	ResponseTSType string // TypeScript response type name
	Comment        string // Method comment from protobuf
}

// FileGenerator handles generation for a single proto file
type FileGenerator struct {
	file         *protogen.File
	plugin       *protogen.Plugin
	config       *Config
	packageFiles []*protogen.File // All files in the same package
}

// NewFileGenerator creates a new file generator
func NewFileGenerator(file *protogen.File, plugin *protogen.Plugin, config *Config) *FileGenerator {
	return &FileGenerator{
		file:         file,
		plugin:       plugin,
		config:       config,
		packageFiles: []*protogen.File{file}, // Default to just this file
	}
}

// SetPackageFiles sets all files that belong to the same package
func (g *FileGenerator) SetPackageFiles(files []*protogen.File) {
	g.packageFiles = files
}

// Generate generates WASM wrapper and TypeScript client for the proto file
func (g *FileGenerator) Generate() error {
	// Validate configuration
	if err := g.config.Validate(); err != nil {
		return err
	}

	// Skip packages with no services
	hasServices := false
	for _, file := range g.packageFiles {
		if len(file.Services) > 0 {
			hasServices = true
			break
		}
	}
	if !hasServices {
		return nil
	}

	// Build template data
	templateData, err := g.buildTemplateData()
	if err != nil {
		return err
	}

	// Generate WASM wrapper if enabled
	if g.config.GenerateWasm {
		if err := g.generateWasmWrapper(templateData); err != nil {
			return err
		}
	}

	// Generate TypeScript client if enabled
	if g.config.GenerateTypeScript {
		if err := g.generateTypeScriptClient(templateData); err != nil {
			return err
		}
	}

	// Generate build script and main example only when generating WASM
	if g.config.GenerateWasm {
		// Generate build script if requested
		if g.config.GenerateBuildScript {
			if err := g.generateBuildScript(templateData); err != nil {
				return err
			}
		}

		// Generate example main.go file
		if err := g.generateMainExample(templateData); err != nil {
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

	// Build service data from all files in the package
	var services []ServiceData
	for _, file := range g.packageFiles {
		for _, service := range file.Services {
			serviceData := g.buildServiceDataForFile(file, service, packageMap)
			if serviceData != nil {
				services = append(services, *serviceData)
			}
		}
	}

	// Skip if no services to generate
	if len(services) == 0 {
		return nil, nil
	}

	// Build TypeScript imports from all methods in all services
	tsImports := g.buildTypeScriptImports(services)

	return &TemplateData{
		PackageName:          packageName,
		SourcePath:           g.file.Desc.Path(),
		GoPackage:            string(g.file.GoImportPath),
		Services:             services,
		Config:               g.config,
		JSNamespace:          g.config.GetDefaultJSNamespace(packageName),
		ModuleName:           g.config.GetDefaultModuleName(packageName),
		APIStructure:         g.config.JSStructure,
		Imports:              imports,
		PackageMap:           packageMap,
		TSRelativeImportPath: g.config.GetRelativeTSImportPathForProto(g.file.Desc.Path()),
		TSImports:            tsImports,
	}, nil
}

// collectUniqueImports collects unique package imports from all services and messages
func (g *FileGenerator) collectUniqueImports() ([]ImportInfo, map[string]string) {
	packagePaths := make(map[string]bool)
	packageMap := make(map[string]string)
	var imports []ImportInfo

	// Collect all unique package paths from services across all package files
	for _, file := range g.packageFiles {
		for _, service := range file.Services {
			if !g.config.ShouldGenerateService(string(service.Desc.Name())) {
				continue
			}

			packagePath := string(file.GoImportPath)
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

// buildServiceDataForFile constructs service data for a protobuf service from a specific file
func (g *FileGenerator) buildServiceDataForFile(file *protogen.File, service *protogen.Service, packageMap map[string]string) *ServiceData {
	serviceName := string(service.Desc.Name())

	// Check if this service should be generated
	if !g.config.ShouldGenerateService(serviceName) {
		return nil
	}

	// Get package alias for this service
	packagePath := string(file.GoImportPath)
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
		RequestTSType:  string(method.Input.GoIdent.GoName),
		ResponseTSType: string(method.Output.GoIdent.GoName),
		Comment:        strings.TrimSpace(string(method.Comments.Leading)),
	}
}

// buildTypeScriptImports builds TypeScript import statements by analyzing all request/response types
func (g *FileGenerator) buildTypeScriptImports(services []ServiceData) []TSImportInfo {
	// Map from proto file to types used from that file
	fileToTypes := make(map[string]map[string]bool)
	
	// We need to go back to the original protogen data to get proto file information
	// Iterate through all files in the package to collect method types
	for _, file := range g.packageFiles {
		for _, service := range file.Services {
			// Check if this service should be generated
			if !g.config.ShouldGenerateService(string(service.Desc.Name())) {
				continue
			}
			
			for _, method := range service.Methods {
				// Check if this method should be generated
				if !g.config.ShouldGenerateMethod(string(method.Desc.Name())) {
					continue
				}
				
				// Skip streaming methods
				if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
					continue
				}
				
				// Get the actual proto files for request and response types
				requestProtoFile := string(method.Input.Desc.ParentFile().Path())
				responseProtoFile := string(method.Output.Desc.ParentFile().Path())
				
				// Initialize maps for proto files if needed
				if fileToTypes[requestProtoFile] == nil {
					fileToTypes[requestProtoFile] = make(map[string]bool)
				}
				if fileToTypes[responseProtoFile] == nil {
					fileToTypes[responseProtoFile] = make(map[string]bool)
				}
				
				// Add request and response types to their respective proto files
				requestTSType := string(method.Input.GoIdent.GoName)
				responseTSType := string(method.Output.GoIdent.GoName)
				fileToTypes[requestProtoFile][requestTSType] = true
				fileToTypes[responseProtoFile][responseTSType] = true
			}
		}
	}
	
	// Convert to TSImportInfo slice
	var tsImports []TSImportInfo
	for protoFile, typesMap := range fileToTypes {
		// Convert map to slice
		var types []string
		for typeName := range typesMap {
			types = append(types, typeName)
		}
		
		// Generate import path for this proto file
		importPath := g.buildTSImportPath(protoFile)
		
		tsImports = append(tsImports, TSImportInfo{
			ProtoFile:  protoFile,
			ImportPath: importPath,
			Types:      types,
		})
	}
	
	return tsImports
}

// buildTSImportPath builds the TypeScript import path for a given proto file
func (g *FileGenerator) buildTSImportPath(protoFile string) string {
	// Remove .proto extension and build the path based on TS generator and detected extension
	baseName := strings.TrimSuffix(protoFile, ".proto")
	
	// Calculate relative path from where TS client is generated to TSImportPath
	relativePath := g.config.calculateRelativePath(g.config.WasmExportPath, g.config.TSImportPath)
	
	// Auto-detect file extension by checking what actually exists
	extension := g.detectTSFileExtension(baseName)
	
	var filename string
	if extension == "" {
		filename = baseName + "_pb"
	} else {
		filename = baseName + "_pb." + extension
	}
	
	return relativePath + "/" + filename
}

// detectTSFileExtension detects whether .ts or .js files exist for the given proto
func (g *FileGenerator) detectTSFileExtension(baseName string) string {
	// This is a simplified version - in a real implementation, we might want to
	// check the actual filesystem or use configuration hints
	// For now, we'll use the ts_import_extension if specified, otherwise fall back to heuristics
	
	if g.config.TSImportExtension != "" {
		if g.config.TSImportExtension == "none" {
			return ""
		}
		return g.config.TSImportExtension
	}
	
	// Auto-detect based on ts_generator (backwards compatibility)
	switch g.config.TSGenerator {
	case "protoc-gen-es":
		// protoc-gen-es typically generates .js files, but with target=ts it generates .ts
		// For now, default to no extension for cleaner imports
		return ""
	case "protoc-gen-ts":
		return ""
	default:
		return ""
	}
}

// getQualifiedTypeName returns the fully qualified Go type name for a message
func (g *FileGenerator) getQualifiedTypeName(message *protogen.Message, packageAlias string) string {
	// Return the qualified type name with package alias
	return "*" + packageAlias + "." + string(message.GoIdent.GoName)
}
