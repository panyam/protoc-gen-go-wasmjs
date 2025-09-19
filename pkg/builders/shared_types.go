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

package builders

import (
	"google.golang.org/protobuf/compiler/protogen"
)

// GenerationConfig holds configuration options common to both Go and TypeScript generators.
// This represents the subset of configuration that affects template data building.
type GenerationConfig struct {
	// Core generation settings
	WasmExportPath string // Where to write WASM artifacts
	TSExportPath   string // Where to write TypeScript artifacts

	// JavaScript API configuration
	JSStructure string // namespaced|flat|service_based
	JSNamespace string // Global JavaScript namespace
	ModuleName  string // WASM module name

	// Build integration
	WasmPackageSuffix   string // Package suffix for WASM wrapper
	GenerateBuildScript bool   // Whether to generate build scripts
	
	// TypeScript generation control
	GenerateClients   bool // Whether to generate TypeScript clients
	GenerateTypes     bool // Whether to generate TypeScript types
	GenerateFactories bool // Whether to generate TypeScript factory classes
}

// ImportInfo represents a Go package import with alias for template generation.
type ImportInfo struct {
	Path  string // Full import path (e.g., "github.com/example/proto/gen/go/library/v1")
	Alias string // Package alias (e.g., "libraryv1")
}

// MessageInfo represents a protobuf message type for template generation.
type MessageInfo struct {
	Name        string // Message name (e.g., "FindBooksRequest")
	GoType      string // Fully qualified Go type (e.g., "libraryv1.FindBooksRequest")
	PackagePath string // Go import path for this message
	Fields      []FieldInfo // Field information (optional, for detailed generation)
}

// EnumInfo represents a protobuf enum type for template generation.
type EnumInfo struct {
	Name        string // Enum name (e.g., "BookStatus")
	GoType      string // Fully qualified Go type (e.g., "libraryv1.BookStatus")
	PackagePath string // Go import path for this enum
	Values      []string // Enum value names (optional)
}

// FieldInfo represents a message field for detailed generation.
type FieldInfo struct {
	Name       string // Field name
	Type       string // Field type
	IsRepeated bool   // Whether field is repeated
	IsMap      bool   // Whether field is a map
}

// ServiceData represents a gRPC service prepared for template generation.
// This is the processed form that templates consume.
type ServiceData struct {
	// Basic service information
	Name         string // Service name (e.g., "LibraryService")
	GoType       string // Go interface type (e.g., "libraryv1.LibraryServiceServer")
	JSName       string // JavaScript name (e.g., "libraryService")
	PackagePath  string // Go package import path
	PackageAlias string // Go package alias for qualified types

	// Service metadata
	IsBrowserProvided bool   // Whether service is implemented by browser
	CustomName        string // Custom JavaScript name from annotations
	Comment           string // Service comment from proto

	// Methods in this service
	Methods []MethodData // All methods that passed filtering
}

// MethodData represents a gRPC method prepared for template generation.
type MethodData struct {
	// Basic method information
	Name           string // Original protobuf method name (e.g., "FindBooks")
	JSName         string // JavaScript method name (e.g., "findBooks" or "searchBooks")
	GoFuncName     string // Go function name for WASM wrapper (e.g., "libraryServiceFindBooks")
	ShouldGenerate bool   // Whether to generate this method based on filters
	Comment        string // Method comment from protobuf

	// Method types (Go)
	RequestType  string // Fully qualified Go request type (e.g., "libraryv1.FindBooksRequest")
	ResponseType string // Fully qualified Go response type (e.g., "libraryv1.FindBooksResponse")

	// Method types (TypeScript)
	RequestTSType  string // TypeScript request type name (e.g., "FindBooksRequest")
	ResponseTSType string // TypeScript response type name (e.g., "FindBooksResponse")

	// Method behavior
	IsAsync           bool // Whether method requires async/callback handling
	IsServerStreaming bool // Whether method uses server-side streaming
}

// PackageInfo represents metadata about a protobuf package for generation.
type PackageInfo struct {
	Name      string           // Package name (e.g., "library.v1")
	Path      string           // Package directory path (e.g., "library/v1")
	GoPackage string           // Go import path
	Files     []*protogen.File // Proto files in this package

	// Content analysis
	HasServices bool // Whether package contains any services
	HasMessages bool // Whether package contains any messages
	HasEnums    bool // Whether package contains any enums
}

// BuildContext provides context information for template data building.
// This carries state and configuration needed during the building process.
type BuildContext struct {
	// Generation configuration
	Config *GenerationConfig

	// Package being processed
	CurrentPackage *PackageInfo
	AllPackages    []*PackageInfo

	// Proto generation context
	Plugin *protogen.Plugin

	// Import management
	ImportMap map[string]string // Maps import path to alias
}

// NewBuildContext creates a new build context for template data building.
func NewBuildContext(plugin *protogen.Plugin, config *GenerationConfig, currentPackage *PackageInfo) *BuildContext {
	return &BuildContext{
		Config:         config,
		CurrentPackage: currentPackage,
		Plugin:         plugin,
		ImportMap:      make(map[string]string),
	}
}

// AddImport registers an import path with an alias for template generation.
// This ensures imports are managed consistently across template data building.
func (bc *BuildContext) AddImport(path, alias string) {
	bc.ImportMap[path] = alias
}

// GetImports returns all registered imports as ImportInfo structures.
// This is used by templates to generate import statements.
func (bc *BuildContext) GetImports() []ImportInfo {
	var imports []ImportInfo
	for path, alias := range bc.ImportMap {
		imports = append(imports, ImportInfo{
			Path:  path,
			Alias: alias,
		})
	}
	return imports
}

// SetAllPackages sets the complete list of packages being processed.
// This is used for cross-package dependency analysis.
func (bc *BuildContext) SetAllPackages(packages []*PackageInfo) {
	bc.AllPackages = packages
}
