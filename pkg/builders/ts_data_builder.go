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
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/core"
	"github.com/panyam/protoc-gen-go-wasmjs/pkg/filters"
)

// TSTemplateData represents all data needed for TypeScript template generation.
// This includes client interfaces, message types, enums, and factory data.
// TSImportGroup represents a group of types imported from a single file
type TSImportGroup struct {
	ImportPath string   // Relative import path (e.g., "./interfaces" or "../models/interfaces")
	Types      []string // Types to import from this path
}

type TSTemplateData struct {
	// Package metadata
	PackageName string // Proto package name (e.g., "library.v1")
	PackagePath string // Directory path (e.g., "library/v1")
	ModuleName  string // Module name for client class (e.g., "library_v1_services")
	SourcePath  string // Primary source proto file (for comment generation)

	// Content data - using TypeScript-specific structures
	Services []ServiceData    // Services for client generation
	Messages []TSMessageInfo // Messages for interface generation (enriched)
	Enums    []TSEnumInfo    // Enums for type generation (enriched)

	// TypeScript type imports for client generation
	TypeImports      []string        // DEPRECATED: Use ImportGroups instead
	ImportGroups     []TSImportGroup // Grouped imports by file path

	// TypeScript-specific configuration
	ImportBasePath string // Base path for imports (e.g., "./library/v1")

	// Cross-package dependencies
	ExternalImports []ExternalImport // Imports from other packages

	// Additional template-specific fields
	BaseName           string // Base name for files (e.g., "library_v1")
	DeserializerName   string // Deserializer class name
	FactoryName        string // Factory class name
	SchemaRegistryName string // Schema registry name
	SchemaImports      []SchemaImport // Schema registry imports for package-level consolidation

	// Computed flags for template conditionals
	HasBrowserServices bool // Whether any browser services exist
	HasBrowserClients  bool // Whether any browser clients exist (alias for HasBrowserServices)
	HasMessages        bool // Whether there are any messages
	HasEnums           bool // Whether there are any enums

	// Client generation specific
	APIStructure string              // API structure (namespaced|flat|service_based)
	JSNamespace  string              // JavaScript namespace
	Dependencies []FactoryDependency // Factory dependencies for cross-package refs
}

// FactoryDependency represents a dependency on another package's factory
type FactoryDependency struct {
	PackageName  string // e.g., "library.common"
	FactoryName  string // e.g., "LibraryCommonFactory"
	ImportPath   string // e.g., "../common/factory"
	InstanceName string // e.g., "commonFactory"
}

// SchemaImport represents an import of a directory-level schema registry
// Used for package-level consolidated schemas that merge multiple directory-level registries
type SchemaImport struct {
	RegistryName string // Name of the schema registry in the source file (e.g., "test_one_package_v1SchemaRegistry")
	Alias        string // Alias to use in the import to avoid naming conflicts (e.g., "modelsSchemas", "models2Schemas")
	ImportPath   string // Relative import path (e.g., "./models/schemas", "./models2/schemas")
}

// TSMessageInfo extends basic message info with TypeScript-specific fields
type TSMessageInfo struct {
	Name               string        // Original proto message name
	TSName             string        // TypeScript interface name
	PackageName        string        // Proto package name
	FullyQualifiedName string        // Full proto name
	ProtoFile          string        // Source proto file
	Comment            string        // Leading comment
	MethodName         string        // Factory method name (e.g., "newBook")
	Fields             []TSFieldInfo // TypeScript field information
	IsNested           bool          // Whether this is nested
	IsMapEntry         bool          // Whether this is a map entry
	OneofGroups        []string      // Oneof group names
}

// TSFieldInfo represents a field in a TypeScript interface
type TSFieldInfo struct {
	Name           string // Proto field name
	TSName         string // TypeScript field name (camelCase)
	TSType         string // TypeScript type
	Number         int32  // Field number
	ProtoFieldID   int32  // Proto field ID (alias for Number, for template compatibility)
	DefaultValue   string // Default value for the field
	IsOptional     bool   // Whether the field is optional
	IsRepeated     bool   // Whether this is repeated
	IsOneof        bool   // Whether this is part of a oneof
	OneofGroup     string // Oneof group name if applicable
	MessageType    string // If this is a message type field (fully qualified name, e.g., "utils.v1.ParentMessage.NestedType")
	MessagePackage string // Package where the message type is defined (e.g., "utils.v1"), extracted from descriptor
	IsNestedType   bool   // Whether the message type is a nested message
	Comment        string // Field comment
}

// TSEnumInfo extends basic enum info with TypeScript-specific fields
type TSEnumInfo struct {
	Name               string          // Original proto enum name
	TSName             string          // TypeScript enum name
	PackageName        string          // Proto package name
	FullyQualifiedName string          // Full proto name
	ProtoFile          string          // Source proto file
	Comment            string          // Leading comment
	Values             []TSEnumValue   // TypeScript enum values
}

// TSEnumValue represents an enum value in TypeScript
type TSEnumValue struct {
	Name    string // Original proto name
	TSName  string // TypeScript name
	Number  int32  // Numeric value
	Comment string // Value comment
}

// ExternalImport represents an import from another package or external library.
type ExternalImport struct {
	ImportPath string   // Import path (e.g., "../common/v1/interfaces")
	Types      []string // Types imported from this path
}

// TSDataBuilder builds template data structures specifically for TypeScript generation.
// This focuses on TypeScript-specific concerns like interfaces, type imports, and client generation.
type TSDataBuilder struct {
	analyzer         *core.ProtoAnalyzer
	pathCalc         *core.PathCalculator
	nameConv         *core.NameConverter
	serviceFilter    *filters.ServiceFilter
	methodFilter     *filters.MethodFilter
	messageCollector *filters.MessageCollector
	enumCollector    *filters.EnumCollector
	wellKnownMapper  *core.WellKnownTypesMapper
}

// NewTSDataBuilder creates a new TypeScript data builder with all necessary dependencies.
func NewTSDataBuilder(
	analyzer *core.ProtoAnalyzer,
	pathCalc *core.PathCalculator,
	nameConv *core.NameConverter,
	serviceFilter *filters.ServiceFilter,
	methodFilter *filters.MethodFilter,
	messageCollector *filters.MessageCollector,
	enumCollector *filters.EnumCollector,
) *TSDataBuilder {
	return &TSDataBuilder{
		analyzer:         analyzer,
		pathCalc:         pathCalc,
		nameConv:         nameConv,
		serviceFilter:    serviceFilter,
		methodFilter:     methodFilter,
		messageCollector: messageCollector,
		enumCollector:    enumCollector,
		wellKnownMapper:  core.NewWellKnownTypesMapper(),
	}
}

// BuildServiceClientData creates TypeScript client template data for a single service.
// This builds the data needed for generating a TypeScript client class for one specific service.
func (tb *TSDataBuilder) BuildServiceClientData(
	packageInfo *PackageInfo,
	service *protogen.Service,
	criteria *filters.FilterCriteria,
	config *GenerationConfig,
) (*TSTemplateData, error) {
	
	// Build context
	context := NewBuildContext(nil, config, packageInfo)
	
	// Find the file containing this service
	var serviceFile *protogen.File
	for _, file := range packageInfo.Files {
		for _, fileService := range file.Services {
			if fileService == service {
				serviceFile = file
				break
			}
		}
		if serviceFile != nil {
			break
		}
	}
	
	if serviceFile == nil {
		return nil, fmt.Errorf("could not find file containing service %s", service.GoName)
	}
	
	// Build service result for this service
	serviceResult := tb.serviceFilter.ShouldIncludeService(service, criteria)
	
	// Build service data for this single service
	serviceData, err := tb.buildServiceDataForTS(service, serviceFile, serviceResult, criteria, context)
	if err != nil {
		return nil, fmt.Errorf("failed to build service data: %w", err)
	}
	
	services := []ServiceData{*serviceData}

	log.Printf("TS BuildServiceClientData: Service %s has %d methods", service.GoName, len(serviceData.Methods))
	for _, m := range serviceData.Methods {
		log.Printf("  - Method %s (JSName=%s, ShouldGenerate=%v)", m.Name, m.JSName, m.ShouldGenerate)
	}
	log.Printf("TS BuildServiceClientData: APIStructure=%s, JSNamespace=%s", config.JSStructure, tb.getJSNamespace(packageInfo.Name, config))

	// Collect TypeScript imports grouped by file path
	importGroups := tb.collectServiceTypeImportGroups(service, serviceFile, criteria)

	return &TSTemplateData{
		PackageName:  packageInfo.Name,
		PackagePath:  packageInfo.Path,
		SourcePath:   string(serviceFile.Desc.Path()),
		ModuleName:   tb.getModuleName(packageInfo.Name, config),
		Services:     services,
		ImportGroups: importGroups,
		APIStructure: config.JSStructure,
		JSNamespace:  tb.getJSNamespace(packageInfo.Name, config),
	}, nil
}

// collectServiceTypeImportGroups groups imports by their source file paths
// This handles the case where types are defined in different proto files/directories
func (tb *TSDataBuilder) collectServiceTypeImportGroups(
	service *protogen.Service,
	serviceFile *protogen.File,
	criteria *filters.FilterCriteria,
) []TSImportGroup {
	// Map from import path to set of types
	importMap := make(map[string]map[string]bool)

	// Get the service file's output directory
	serviceDir := filepath.Dir(string(serviceFile.Desc.Path()))

	for _, method := range service.Methods {
		// Check if method should be generated
		methodResult := tb.methodFilter.ShouldIncludeMethod(method, criteria)
		if !methodResult.Include {
			continue
		}

		// Process input message
		if method.Input != nil {
			tb.addMessageToImportMap(method.Input, serviceDir, importMap)
		}

		// Process output message
		if method.Output != nil {
			tb.addMessageToImportMap(method.Output, serviceDir, importMap)
		}
	}

	// Convert map to slice of TSImportGroup
	var groups []TSImportGroup
	for importPath, types := range importMap {
		var typeList []string
		for typeName := range types {
			typeList = append(typeList, typeName)
		}
		// Sort for deterministic output
		sort.Strings(typeList)
		groups = append(groups, TSImportGroup{
			ImportPath: importPath,
			Types:      typeList,
		})
	}

	// Sort groups by import path for deterministic output
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].ImportPath < groups[j].ImportPath
	})

	return groups
}

// addMessageToImportMap adds a message and its import path to the import map
func (tb *TSDataBuilder) addMessageToImportMap(
	message *protogen.Message,
	serviceDir string,
	importMap map[string]map[string]bool,
) {
	// Get the file that defines this message
	messageFile := message.Desc.ParentFile()
	if messageFile == nil {
		log.Printf("Warning: Could not find parent file for message %s", message.GoIdent.GoName)
		return
	}

	// Get the message file's output directory
	messageFilePath := string(messageFile.Path())
	messageDir := filepath.Dir(messageFilePath)

	// Calculate relative path from service directory to message directory
	relativePath, err := filepath.Rel(serviceDir, messageDir)
	if err != nil {
		log.Printf("Warning: Could not calculate relative path from %s to %s: %v", serviceDir, messageDir, err)
		relativePath = messageDir
	}

	// Convert to forward slashes and ensure it starts with ./
	relativePath = filepath.ToSlash(relativePath)
	if !strings.HasPrefix(relativePath, "../") && !strings.HasPrefix(relativePath, "./") {
		if relativePath == "." {
			relativePath = "./interfaces"
		} else {
			relativePath = "./" + relativePath + "/interfaces"
		}
	} else {
		relativePath = relativePath + "/interfaces"
	}

	// Get the TypeScript type name (interface name)
	tsTypeName := message.GoIdent.GoName

	// Add to import map
	if importMap[relativePath] == nil {
		importMap[relativePath] = make(map[string]bool)
	}
	importMap[relativePath][tsTypeName] = true
}

// collectServiceTypeImports collects unique TypeScript types needed for method signatures
func (tb *TSDataBuilder) collectServiceTypeImports(services []ServiceData) []string {
	importSet := make(map[string]bool)

	for _, service := range services {
		for _, method := range service.Methods {
			if !method.ShouldGenerate {
				continue
			}

			// Add request type
			if method.RequestTSType != "" {
				importSet[method.RequestTSType] = true
			}

			// Add response type
			if method.ResponseTSType != "" {
				importSet[method.ResponseTSType] = true
			}
		}
	}
	
	// Convert set to sorted slice
	var imports []string
	for importType := range importSet {
		imports = append(imports, importType)
	}
	
	// Sort for consistent output
	// Note: using simple append for now, could add sorting if needed
	return imports
}

// BuildClientData creates TypeScript client template data for services.
// This builds the data needed for generating TypeScript client classes that call WASM.
func (tb *TSDataBuilder) BuildClientData(
    packageInfo *PackageInfo,
    criteria *filters.FilterCriteria,
    config *GenerationConfig,
) (*TSTemplateData, error) {

	// Build context
	context := NewBuildContext(nil, config, packageInfo)

	// Filter and build services for client generation
	services, err := tb.buildClientServices(packageInfo.Files, criteria, context)
	if err != nil {
		return nil, err
	}

	// Log services and methods for debugging
	for _, svc := range services {
		log.Printf("TS BuildClientData: Service %s has %d methods", svc.Name, len(svc.Methods))
		for _, m := range svc.Methods {
			log.Printf("  - Method %s (JSName=%s, ShouldGenerate=%v)", m.Name, m.JSName, m.ShouldGenerate)
		}
	}
	log.Printf("TS BuildClientData: APIStructure=%s, JSNamespace=%s (config.JSNamespace=%s)", config.JSStructure, tb.getJSNamespace(packageInfo.Name, config), config.JSNamespace)

	// Skip if no services
	if len(services) == 0 {
		return nil, nil
	}

	// Collect TypeScript imports needed for typed method signatures
	requiredImports := tb.collectServiceTypeImports(services)
	
	// Generate names for TypeScript artifacts
	baseName := strings.ReplaceAll(packageInfo.Name, ".", "_")

	// Determine if there are any browser services
	hasBrowserServices := false
	for _, svc := range services {
		if svc.IsBrowserProvided {
			hasBrowserServices = true
			break
		}
	}

	return &TSTemplateData{
		PackageName:        packageInfo.Name,
		PackagePath:        packageInfo.Path,
		ModuleName:         tb.getModuleName(packageInfo.Name, config),
		SourcePath:         tb.getPrimarySourcePath(packageInfo.Files),
		Services:           services,
		TypeImports:        requiredImports,
		ImportBasePath:     tb.getImportBasePath(packageInfo.Path, config),
		ExternalImports:    []ExternalImport{}, // TODO: Implement external imports
		BaseName:           baseName,
		APIStructure:       config.JSStructure,
		JSNamespace:        tb.getJSNamespace(packageInfo.Name, config),
		HasBrowserServices: hasBrowserServices,
		HasBrowserClients:  hasBrowserServices, // Same as HasBrowserServices
	}, nil
}

// BuildTypeData creates TypeScript type template data for messages and enums.
// This builds the data needed for generating TypeScript interfaces, classes, and enums.
func (tb *TSDataBuilder) BuildTypeData(
	packageInfo *PackageInfo,
	criteria *filters.FilterCriteria,
	config *GenerationConfig,
) (*TSTemplateData, error) {

	// Collect messages
	messageResult := tb.messageCollector.CollectMessages(packageInfo.Files, criteria)

	// Collect enums
	enumResult := tb.enumCollector.CollectEnums(packageInfo.Files, criteria)

	// Skip if no content
	if len(messageResult.Items) == 0 && len(enumResult.Items) == 0 {
		return nil, nil
	}

	// Transform to TypeScript-specific structures
	tsMessages := tb.transformMessages(messageResult.Items, packageInfo.Files)
	tsEnums := tb.transformEnums(enumResult.Items)

	// Build external imports for cross-package references
	externalImports := tb.buildExternalImportsFromTSMessages(tsMessages, packageInfo, config)

	// Generate names for TypeScript artifacts
	baseName := strings.ReplaceAll(packageInfo.Name, ".", "_")

	return &TSTemplateData{
		PackageName:        packageInfo.Name,
		PackagePath:        packageInfo.Path,
		SourcePath:         tb.getPrimarySourcePath(packageInfo.Files),
		Messages:           tsMessages,
		Enums:              tsEnums,
		ImportBasePath:     tb.getImportBasePath(packageInfo.Path, config),
		ExternalImports:    externalImports,
		BaseName:           baseName,
		DeserializerName:   tb.nameConv.ToPascalCase(baseName) + "Deserializer",
		FactoryName:        tb.nameConv.ToPascalCase(baseName) + "Factory",
		SchemaRegistryName: tb.nameConv.ToCamelCase(baseName) + "SchemaRegistry",
		HasMessages:        len(tsMessages) > 0,
		HasEnums:           len(tsEnums) > 0,
		HasBrowserServices: false, // This is for type generation, no browser services here
		HasBrowserClients:  false, // This is for type generation, no browser clients here
	}, nil
}

// BuildPackageSchemaData builds template data for package-level consolidated schemas.
// This collects all directory-level schema files from the package and generates imports
// to merge them into a single package-level schema registry.
func (tb *TSDataBuilder) BuildPackageSchemaData(
	packageInfo *PackageInfo,
	config *GenerationConfig,
) (*TSTemplateData, error) {
	// Find all unique directories that contain proto files in this package
	dirMap := make(map[string]bool)
	for _, file := range packageInfo.Files {
		protoPath := string(file.Desc.Path())
		dir := filepath.Dir(protoPath)
		dirMap[dir] = true
	}

	// Build schema imports for each directory
	var schemaImports []SchemaImport
	packageRootDir := strings.ReplaceAll(packageInfo.Name, ".", "/")
	baseName := strings.ReplaceAll(packageInfo.Name, ".", "_")
	schemaRegistryName := tb.nameConv.ToCamelCase(baseName) + "SchemaRegistry"

	for dir := range dirMap {
		// Calculate relative path from package root to directory
		relativePath, err := filepath.Rel(packageRootDir, dir)
		if err != nil {
			log.Printf("Warning: Could not calculate relative path from %s to %s: %v", packageRootDir, dir, err)
			continue
		}

		// Convert to forward slashes and ensure it starts with ./
		relativePath = filepath.ToSlash(relativePath)
		if relativePath == "." {
			// Skip if this is the package root itself (no subdirectory)
			continue
		}
		if !strings.HasPrefix(relativePath, "../") && !strings.HasPrefix(relativePath, "./") {
			relativePath = "./" + relativePath
		}

		// Create import path to schemas file
		importPath := relativePath + "/schemas"

		// Create an alias based on the directory name to avoid conflicts
		// e.g., "models" -> "modelsSchemas", "models2" -> "models2Schemas"
		dirName := filepath.Base(dir)
		alias := dirName + "Schemas"

		schemaImports = append(schemaImports, SchemaImport{
			RegistryName: schemaRegistryName,
			Alias:        alias,
			ImportPath:   importPath,
		})
	}

	// Sort schema imports by import path for deterministic output
	sort.Slice(schemaImports, func(i, j int) bool {
		return schemaImports[i].ImportPath < schemaImports[j].ImportPath
	})

	return &TSTemplateData{
		PackageName:        packageInfo.Name,
		PackagePath:        packageInfo.Path,
		SchemaRegistryName: schemaRegistryName,
		SchemaImports:      schemaImports,
	}, nil
}

// buildClientServices filters and builds services for TypeScript client generation.
func (tb *TSDataBuilder) buildClientServices(
	files []*protogen.File,
	criteria *filters.FilterCriteria,
	context *BuildContext,
) ([]ServiceData, error) {

	var services []ServiceData

	for _, file := range files {
		for _, service := range file.Services {
			// Filter the service
			serviceResult := tb.serviceFilter.ShouldIncludeService(service, criteria)
			if !serviceResult.Include {
				continue
			}

			// Build service data
			serviceData, err := tb.buildServiceDataForTS(service, file, serviceResult, criteria, context)
			if err != nil {
				return nil, err
			}

			if serviceData != nil {
				services = append(services, *serviceData)
			}
		}
	}

	return services, nil
}

// buildServiceDataForTS creates ServiceData specifically for TypeScript templates.
func (tb *TSDataBuilder) buildServiceDataForTS(
	service *protogen.Service,
	file *protogen.File,
	serviceResult filters.ServiceFilterResult,
	criteria *filters.FilterCriteria,
	context *BuildContext,
) (*ServiceData, error) {

	serviceName := string(service.Desc.Name())

	// Filter and build methods
	var methods []MethodData
	for _, method := range service.Methods {
		methodResult := tb.methodFilter.ShouldIncludeMethod(method, criteria)
		if !methodResult.Include {
			continue
		}

		methodData := tb.buildMethodDataForTS(method, serviceName, methodResult)
		methods = append(methods, methodData)
	}

	// Continue with empty methods array - still generate client with 0 methods
	// This handles services like: service TestService { option (wasmjs.v1.browser_provided) = true; }
	// where the service exists but has no RPC methods defined

	// JavaScript name
	jsName := tb.nameConv.ToCamelCase(serviceName)
	if serviceResult.CustomName != "" {
		jsName = serviceResult.CustomName
	}

	return &ServiceData{
		Name:              serviceName,
		JSName:            jsName,
		IsBrowserProvided: serviceResult.IsBrowserProvided,
		CustomName:        serviceResult.CustomName,
		Comment:           strings.TrimSpace(string(service.Comments.Leading)),
		Methods:           methods,
	}, nil
}

// buildMethodDataForTS creates MethodData specifically for TypeScript templates.
func (tb *TSDataBuilder) buildMethodDataForTS(
	method *protogen.Method,
	serviceName string,
	methodResult filters.MethodFilterResult,
) MethodData {

	methodName := string(method.Desc.Name())

	// JavaScript name
	jsName := methodResult.CustomJSName
	if jsName == "" {
		jsName = tb.nameConv.ToCamelCase(methodName)
	}

	return MethodData{
		Name:              methodName,
		JSName:            jsName,
		ShouldGenerate:    true, // Method passed filtering, should be generated
		Comment:           strings.TrimSpace(string(method.Comments.Leading)),
		RequestTSType:     string(method.Input.GoIdent.GoName),
		ResponseTSType:    string(method.Output.GoIdent.GoName),
		IsAsync:           methodResult.IsAsync,
		IsServerStreaming: methodResult.IsServerStreaming,
	}
}

// buildExternalImportsFromTSMessages analyzes TypeScript messages to determine cross-package type dependencies.
func (tb *TSDataBuilder) buildExternalImportsFromTSMessages(
	messages []TSMessageInfo,
	packageInfo *PackageInfo,
	config *GenerationConfig,
) []ExternalImport {
	// Map to track unique imports by source
	importMap := make(map[string]map[string]bool)

	// Analyze all messages to find external type references
	for _, msg := range messages {
		tb.collectTSMessageExternalTypes(msg, importMap, packageInfo)
	}

	// Convert map to sorted ExternalImport list
	var imports []ExternalImport
	for importSource, types := range importMap {
		// Convert type set to sorted list
		var typeList []string
		for typeName := range types {
			typeList = append(typeList, typeName)
		}
		// Sort for consistent output
		sort.Strings(typeList)

		imports = append(imports, ExternalImport{
			ImportPath: importSource,
			Types:      typeList,
		})
	}

	// Sort imports by path for consistent output
	sort.Slice(imports, func(i, j int) bool {
		return imports[i].ImportPath < imports[j].ImportPath
	})

	return imports
}

// collectTSMessageExternalTypes collects external type dependencies from a TypeScript message.
// This uses the MessagePackage field which was extracted from the protobuf descriptor API.
func (tb *TSDataBuilder) collectTSMessageExternalTypes(msg TSMessageInfo, importMap map[string]map[string]bool, currentPackage *PackageInfo) {
	// Check each field for external type references
	for _, field := range msg.Fields {
		if field.MessageType == "" {
			continue
		}

		// Check if this is a well-known type
		if mapping, exists := tb.wellKnownMapper.GetMapping(field.MessageType); exists {
			if !mapping.IsNative && mapping.ImportSource != "" {
				// Add well-known type to import map
				if importMap[mapping.ImportSource] == nil {
					importMap[mapping.ImportSource] = make(map[string]bool)
				}
				importMap[mapping.ImportSource][mapping.TSType] = true
			}
			continue
		}

		// Handle cross-package message types
		// Use MessagePackage which was extracted directly from field.Message.Desc.ParentFile().Package()
		// This is accurate and handles nested types correctly
		if field.MessagePackage != "" && field.MessagePackage != currentPackage.Name {
			// This is a cross-package reference - calculate import path
			importPath := tb.calculateCrossPackageImportPath(currentPackage.Path, field.MessagePackage)
			typeName := tb.extractTypeNameFromFullyQualified(field.MessageType)

			// Add to import map
			if importMap[importPath] == nil {
				importMap[importPath] = make(map[string]bool)
			}
			importMap[importPath][typeName] = true
		}
	}
}

// extractPackageFromTypeName extracts the package name from a fully qualified type name.
// This handles both simple and nested types correctly.
// Examples:
//   - "utils.v1.HelperUtilType" -> "utils.v1"
//   - "utils.v1.ParentMessage.NestedType" -> "utils.v1"
//   - "google.protobuf.Timestamp" -> "google.protobuf"
//
// The heuristic: Proto package names are lowercase, message names are PascalCase.
// We find the package by taking parts up to (but not including) the first PascalCase part.
func (tb *TSDataBuilder) extractPackageFromTypeName(fullyQualifiedName string) string {
	parts := strings.Split(fullyQualifiedName, ".")
	if len(parts) <= 1 {
		return ""
	}

	// Find the first part that starts with an uppercase letter (message name)
	// Everything before that is the package
	packageParts := []string{}
	for _, part := range parts {
		// Check if this part starts with uppercase (indicates a message name)
		if len(part) > 0 && part[0] >= 'A' && part[0] <= 'Z' {
			// This is a message name, stop here
			break
		}
		// This is part of the package name
		packageParts = append(packageParts, part)
	}

	if len(packageParts) == 0 {
		// Fallback: if no lowercase parts found, return all but last part
		// This handles edge cases where package names might not follow conventions
		return strings.Join(parts[:len(parts)-1], ".")
	}

	return strings.Join(packageParts, ".")
}

// extractTypeNameFromFullyQualified extracts the full type name (including parent messages) from a fully qualified name.
// For nested types, this returns the complete message hierarchy for proper TypeScript import.
// Examples:
//   - "utils.v1.HelperUtilType" -> "HelperUtilType"
//   - "utils.v1.ParentMessage.NestedType" -> "ParentMessage_NestedType" (for flattened TypeScript exports)
//   - "google.protobuf.Timestamp" -> "Timestamp"
//
// Note: Some TypeScript generators flatten nested types with underscores, others keep them nested.
// For now, we return the last part only, but this can be adjusted based on the generation strategy.
func (tb *TSDataBuilder) extractTypeNameFromFullyQualified(fullyQualifiedName string) string {
	parts := strings.Split(fullyQualifiedName, ".")
	if len(parts) == 0 {
		return fullyQualifiedName
	}

	// Extract package parts (lowercase) vs message parts (PascalCase)
	messageParts := []string{}
	for _, part := range parts {
		// If this part starts with uppercase, it's a message name
		if len(part) > 0 && part[0] >= 'A' && part[0] <= 'Z' {
			messageParts = append(messageParts, part)
		}
	}

	if len(messageParts) == 0 {
		// Fallback: return last part
		return parts[len(parts)-1]
	}

	// For nested types, join with underscore: ParentMessage_NestedType
	// This is a common TypeScript pattern for flattened nested types
	return strings.Join(messageParts, "_")
}

// calculateCrossPackageImportPath calculates the relative import path from current package to target package.
// For example, from "presenter/v1" to "utils.v1" -> "../../utils/v1/interfaces"
func (tb *TSDataBuilder) calculateCrossPackageImportPath(currentPackagePath string, targetPackageName string) string {
	// Convert target package name to path (e.g., "utils.v1" -> "utils/v1")
	targetPackagePath := tb.pathCalc.BuildPackagePath(targetPackageName)

	// Calculate relative path from current package directory to target package directory
	relativePath := tb.pathCalc.CalculateRelativePath(currentPackagePath, targetPackagePath)

	// Append "/interfaces" since types are defined in interfaces.ts
	return relativePath + "/interfaces"
}

// getModuleName determines the TypeScript module name.
func (tb *TSDataBuilder) getModuleName(packageName string, config *GenerationConfig) string {
	if config.ModuleName != "" {
		return config.ModuleName
	}
	return tb.nameConv.ToModuleName(packageName)
}

// getImportBasePath determines the base path for TypeScript imports.
func (tb *TSDataBuilder) getImportBasePath(packagePath string, config *GenerationConfig) string {
	// This would calculate the relative import path based on the output structure
	// For now, use the package path directly
	return "./" + packagePath
}

// getJSNamespace determines the JavaScript namespace for the TypeScript client.
func (tb *TSDataBuilder) getJSNamespace(packageName string, config *GenerationConfig) string {
	if config.JSNamespace != "" {
		return config.JSNamespace
	}
	// Default to package name with underscores
	return strings.ReplaceAll(packageName, ".", "_")
}

// getPrimarySourcePath gets the primary source file path from the package files.
func (tb *TSDataBuilder) getPrimarySourcePath(files []*protogen.File) string {
	if len(files) > 0 {
		return files[0].Desc.Path()
	}
	return ""
}

// transformMessages converts basic MessageInfo to TypeScript-enriched structures.
func (tb *TSDataBuilder) transformMessages(messages []filters.MessageInfo, protoFiles []*protogen.File) []TSMessageInfo {
	result := make([]TSMessageInfo, 0, len(messages))

	// Create a map for quick lookup of protogen.Message by fully qualified name
	protoMessageMap := tb.buildProtoMessageMap(protoFiles)

	for _, msg := range messages {
		// Find the corresponding protogen.Message for field extraction
		protoMessage := protoMessageMap[msg.FullyQualifiedName]

		// Calculate TypeScript name - for nested types, flatten with parent name
		tsName := msg.Name
		if msg.IsNested {
			// For nested types, use flattened name: ParentMessage_NestedType
			// Extract all message names from fully qualified name
			tsName = tb.extractTypeNameFromFullyQualified(msg.FullyQualifiedName)
		}

		tsMsg := TSMessageInfo{
			Name:               msg.Name,
			TSName:             tsName, // Flattened for nested types
			PackageName:        msg.PackageName,
			FullyQualifiedName: msg.FullyQualifiedName,
			ProtoFile:          msg.ProtoFile,
			Comment:            msg.Comment,
			MethodName:         "new" + tsName, // Factory method name uses flattened name
			Fields:             tb.extractFieldInfo(protoMessage),
			IsNested:           msg.IsNested,
			IsMapEntry:         msg.IsMapEntry,
			OneofGroups:        tb.extractOneofGroups(protoMessage),
		}
		result = append(result, tsMsg)
	}

	return result
}

// buildProtoMessageMap creates a lookup map of protogen.Message by fully qualified name
func (tb *TSDataBuilder) buildProtoMessageMap(protoFiles []*protogen.File) map[string]*protogen.Message {
	messageMap := make(map[string]*protogen.Message)
	
	for _, file := range protoFiles {
		packageName := string(file.Desc.Package())
		
		// Add top-level messages
		for _, msg := range file.Messages {
			messageName := string(msg.Desc.Name())
			fullyQualifiedName := packageName + "." + messageName
			messageMap[fullyQualifiedName] = msg
			
			// Add nested messages recursively
			tb.addNestedMessages(messageMap, msg, packageName)
		}
	}
	
	return messageMap
}

// addNestedMessages recursively adds nested messages to the map
func (tb *TSDataBuilder) addNestedMessages(messageMap map[string]*protogen.Message, parentMsg *protogen.Message, packageName string) {
	for _, nestedMsg := range parentMsg.Messages {
		parentName := string(parentMsg.Desc.Name())
		nestedName := string(nestedMsg.Desc.Name())
		fullyQualifiedName := packageName + "." + parentName + "." + nestedName
		messageMap[fullyQualifiedName] = nestedMsg
		
		// Recursively add nested messages
		tb.addNestedMessages(messageMap, nestedMsg, packageName)
	}
}

// extractFieldInfo extracts field information from a protogen.Message
func (tb *TSDataBuilder) extractFieldInfo(protoMessage *protogen.Message) []TSFieldInfo {
	if protoMessage == nil {
		return []TSFieldInfo{}
	}
	
	fields := make([]TSFieldInfo, 0, len(protoMessage.Fields))
	
	for _, field := range protoMessage.Fields {
		fieldNumber := int32(field.Desc.Number())
		fieldInfo := TSFieldInfo{
			Name:         string(field.Desc.Name()),
			TSName:       field.Desc.JSONName(), // camelCase name
			Number:       fieldNumber,
			ProtoFieldID: fieldNumber, // Alias for template compatibility
			IsRepeated:   field.Desc.IsList(),
			IsOptional:   field.Desc.HasOptionalKeyword(),
			Comment:      strings.TrimSpace(string(field.Comments.Leading)),
		}
		
		// Handle oneof fields
		if field.Oneof != nil {
			fieldInfo.IsOneof = true
			fieldInfo.OneofGroup = string(field.Oneof.Desc.Name())
		}
		
		// Determine field type and TypeScript type
		kind := field.Desc.Kind()
		switch kind.String() {
		case "string":
			fieldInfo.TSType = "string"
			fieldInfo.DefaultValue = `""`
		case "int32", "int64", "uint32", "uint64", "sint32", "sint64", "fixed32", "fixed64", "sfixed32", "sfixed64", "double", "float":
			fieldInfo.TSType = "number"
			fieldInfo.DefaultValue = "0"
		case "bool":
			fieldInfo.TSType = "boolean"
			fieldInfo.DefaultValue = "false"
		case "bytes":
			fieldInfo.TSType = "Uint8Array"
			fieldInfo.DefaultValue = "new Uint8Array()"
		case "message":
			if field.Message != nil {
				// Check if this is a map field
				if field.Message.Desc.IsMapEntry() {
					// Handle map types: map<K,V> → Record<K,V>
					mapFields := field.Message.Fields
					if len(mapFields) >= 2 {
						keyField := mapFields[0]   // Key field
						valueField := mapFields[1] // Value field
						
						keyType := tb.protoKindToTSType(keyField.Desc.Kind())
						valueType := tb.protoKindToTSType(valueField.Desc.Kind())
						
						// Handle message value types
						if valueField.Desc.Kind().String() == "message" && valueField.Message != nil {
							valueType = string(valueField.Message.Desc.Name())
						}
						
						fieldInfo.TSType = fmt.Sprintf("Record<%s, %s>", keyType, valueType)
						fieldInfo.DefaultValue = "{}"
					} else {
						// Fallback for malformed map
						fieldInfo.TSType = "Record<string, any>"
						fieldInfo.DefaultValue = "{}"
					}
				} else {
					// Regular message type
					// Use FullName() to get the complete qualified name including parent messages
					// e.g., "utils.v1.ParentMessage.NestedType" for nested types
					fullTypeName := string(field.Message.Desc.FullName())
					fieldInfo.MessageType = fullTypeName

					// Extract package using the descriptor API
					fieldInfo.MessagePackage = string(field.Message.Desc.ParentFile().Package())

					// Check if this is a nested type using Parent() descriptor method
					parent := field.Message.Desc.Parent()
					_, fieldInfo.IsNestedType = parent.(protoreflect.MessageDescriptor)

					// Check if this is a well-known type
					if mapping, exists := tb.wellKnownMapper.GetMapping(fullTypeName); exists {
						fieldInfo.TSType = mapping.TSType
						// For well-known types, we typically don't set a default value
						fieldInfo.DefaultValue = "undefined"
					} else {
						// For regular message types, determine the TypeScript type name
						// For nested types, use flattened name (e.g., "ParentMessage_NestedType")
						// For top-level types, use simple name (e.g., "HelperUtilType")
						if fieldInfo.IsNestedType {
							// Use flattened name for nested types
							fieldInfo.TSType = tb.extractTypeNameFromFullyQualified(fullTypeName)
						} else {
							// Use simple name for top-level types
							fieldInfo.TSType = string(field.Message.Desc.Name())
						}
						fieldInfo.DefaultValue = "undefined"
					}
				}
			}
		case "enum":
			if field.Enum != nil {
				enumName := string(field.Enum.Desc.Name())
				fieldInfo.TSType = enumName
				// Find first enum value for default
				if len(field.Enum.Values) > 0 {
					fieldInfo.DefaultValue = enumName + "." + string(field.Enum.Values[0].Desc.Name())
				}
			}
		default:
			fieldInfo.TSType = "any"
			fieldInfo.DefaultValue = "undefined"
		}
		
		// Handle repeated fields
		if fieldInfo.IsRepeated {
			fieldInfo.TSType = fieldInfo.TSType + "[]"
			fieldInfo.DefaultValue = "[]"
		}
		
		// Handle optional fields
		if fieldInfo.IsOptional {
			fieldInfo.TSType = fieldInfo.TSType + " | undefined"
		}
		
		fields = append(fields, fieldInfo)
	}
	
	return fields
}

// extractOneofGroups extracts oneof group names from a protogen.Message
func (tb *TSDataBuilder) extractOneofGroups(protoMessage *protogen.Message) []string {
	if protoMessage == nil {
		return []string{}
	}
	
	oneofMap := make(map[string]bool)
	
	for _, field := range protoMessage.Fields {
		if field.Oneof != nil {
			oneofName := string(field.Oneof.Desc.Name())
			oneofMap[oneofName] = true
		}
	}
	
	groups := make([]string, 0, len(oneofMap))
	for group := range oneofMap {
		groups = append(groups, group)
	}
	
	return groups
}

// protoKindToTSType converts protobuf field kind to TypeScript type
func (tb *TSDataBuilder) protoKindToTSType(kind protoreflect.Kind) string {
	switch kind.String() {
	case "string":
		return "string"
	case "int32", "int64", "uint32", "uint64", "sint32", "sint64", "fixed32", "fixed64", "sfixed32", "sfixed64", "double", "float":
		return "number"
	case "bool":
		return "boolean"
	case "bytes":
		return "Uint8Array"
	default:
		return "any"
	}
}

// transformEnums converts basic EnumInfo to TypeScript-enriched structures.
func (tb *TSDataBuilder) transformEnums(enums []filters.EnumInfo) []TSEnumInfo {
	result := make([]TSEnumInfo, 0, len(enums))

	for _, enum := range enums {
		tsEnum := TSEnumInfo{
			Name:               enum.Name,
			TSName:             enum.Name, // TypeScript uses same name as proto
			PackageName:        enum.PackageName,
			FullyQualifiedName: enum.FullyQualifiedName,
			ProtoFile:          enum.ProtoFile,
			Comment:            enum.Comment,
			Values:             tb.transformEnumValues(enum.Values),
		}
		result = append(result, tsEnum)
	}

	return result
}

// transformEnumValues converts enum values to TypeScript format.
func (tb *TSDataBuilder) transformEnumValues(values []filters.EnumValueInfo) []TSEnumValue {
	result := make([]TSEnumValue, 0, len(values))

	for _, val := range values {
		tsVal := TSEnumValue{
			Name:    val.Name,
			TSName:  val.Name, // TypeScript uses same name as proto
			Number:  val.Number,
			Comment: val.Comment,
		}
		result = append(result, tsVal)
	}

	return result
}

// BuildFactoryData creates TypeScript factory template data for a factory artifact.
// This builds the data needed for generating a combined factory + deserializer file
// with correct relative import paths to message interfaces and models.
func (tb *TSDataBuilder) BuildFactoryData(
	factoryFile *protogen.File,
	packageInfo *PackageInfo,
	importedMessages []filters.MessageInfo,
	config *GenerationConfig,
) (*TSTemplateData, error) {
	// Transform imported messages to TypeScript structures
	tsMessages := tb.transformMessages(importedMessages, packageInfo.Files)

	// Get the factory file's output directory
	factoryFilePath := string(factoryFile.Desc.Path())
	factoryDir := filepath.Dir(factoryFilePath)

	// Collect import groups for both interfaces and models
	interfaceImportGroups := tb.collectFactoryTypeImports(tsMessages, factoryDir, "interfaces")
	modelImportGroups := tb.collectFactoryTypeImports(tsMessages, factoryDir, "models")

	// Combine both into a single import groups list
	// We'll use a map to track which paths we've seen to avoid duplicates
	allImportGroups := make([]TSImportGroup, 0)
	importPathMap := make(map[string]*TSImportGroup)

	// Add interface imports
	for _, group := range interfaceImportGroups {
		importPathMap[group.ImportPath] = &group
		allImportGroups = append(allImportGroups, group)
	}

	// Add model imports (with different paths)
	for _, group := range modelImportGroups {
		if _, exists := importPathMap[group.ImportPath]; !exists {
			importPathMap[group.ImportPath] = &group
			allImportGroups = append(allImportGroups, group)
		}
	}

	// Generate names for TypeScript artifacts
	baseName := strings.ReplaceAll(packageInfo.Name, ".", "_")

	// Build schema imports from each directory containing messages
	// Group messages by their directory to create per-directory schema imports
	schemaImports := tb.buildFactorySchemaImports(tsMessages, factoryFilePath, baseName)

	return &TSTemplateData{
		PackageName:        packageInfo.Name,
		PackagePath:        packageInfo.Path,
		SourcePath:         factoryFilePath,
		Messages:           tsMessages,
		ImportGroups:       allImportGroups,
		BaseName:           baseName,
		FactoryName:        tb.nameConv.ToPascalCase(baseName) + "Factory",
		DeserializerName:   tb.nameConv.ToPascalCase(baseName) + "Deserializer",
		SchemaRegistryName: tb.nameConv.ToCamelCase(baseName) + "SchemaRegistry",
		HasMessages:        len(tsMessages) > 0,
		SchemaImports:      schemaImports,
	}, nil
}

// buildFactorySchemaImports builds schema import statements for factory's aggregated schemas file.
// This creates imports from each directory's schemas file with unique aliases.
func (tb *TSDataBuilder) buildFactorySchemaImports(
	messages []TSMessageInfo,
	factoryFilePath string,
	baseName string,
) []SchemaImport {
	factoryDir := filepath.Dir(factoryFilePath)

	// Group messages by directory
	dirMap := make(map[string]bool)
	for _, msg := range messages {
		messageDir := filepath.Dir(msg.ProtoFile)
		dirMap[messageDir] = true
	}

	// Build schema imports for each directory
	var schemaImports []SchemaImport
	for messageDir := range dirMap {
		// Calculate relative path from factory to message directory
		relativePath, err := filepath.Rel(factoryDir, messageDir)
		if err != nil {
			log.Printf("Warning: Could not calculate relative path from %s to %s: %v", factoryDir, messageDir, err)
			continue
		}

		// Convert to forward slashes
		relativePath = filepath.ToSlash(relativePath)
		if !strings.HasPrefix(relativePath, "../") && !strings.HasPrefix(relativePath, "./") {
			if relativePath == "." {
				relativePath = "./schemas"
			} else {
				relativePath = "./" + relativePath + "/schemas"
			}
		} else {
			relativePath = relativePath + "/schemas"
		}

		// Create unique alias based on directory name
		// e.g., "models" -> "modelsSchemas", "models2" -> "models2Schemas"
		dirName := filepath.Base(messageDir)
		alias := tb.nameConv.ToCamelCase(dirName) + "Schemas"

		schemaImports = append(schemaImports, SchemaImport{
			RegistryName: tb.nameConv.ToCamelCase(baseName) + "SchemaRegistry",
			Alias:        alias,
			ImportPath:   relativePath,
		})
	}

	// Sort for deterministic output
	sort.Slice(schemaImports, func(i, j int) bool {
		return schemaImports[i].ImportPath < schemaImports[j].ImportPath
	})

	return schemaImports
}

// collectFactoryTypeImports collects import groups for factory generation.
// This is similar to collectServiceTypeImportGroups but handles messages defined across
// multiple directories within the same package.
func (tb *TSDataBuilder) collectFactoryTypeImports(
	messages []TSMessageInfo,
	factoryDir string,
	importType string, // "interfaces" or "models"
) []TSImportGroup {
	// Map from import path to set of types
	importMap := make(map[string]map[string]bool)

	for _, msg := range messages {
		// Get the file that defines this message
		// We need to find the proto file from ProtoFile field
		messageFilePath := msg.ProtoFile
		messageDir := filepath.Dir(messageFilePath)

		// Calculate relative path from factory directory to message directory
		relativePath, err := filepath.Rel(factoryDir, messageDir)
		if err != nil {
			log.Printf("Warning: Could not calculate relative path from %s to %s: %v", factoryDir, messageDir, err)
			relativePath = messageDir
		}

		// Convert to forward slashes and ensure it starts with ./
		relativePath = filepath.ToSlash(relativePath)
		if !strings.HasPrefix(relativePath, "../") && !strings.HasPrefix(relativePath, "./") {
			if relativePath == "." {
				relativePath = "./" + importType
			} else {
				relativePath = "./" + relativePath + "/" + importType
			}
		} else {
			relativePath = relativePath + "/" + importType
		}

		// Add message type to import map
		if importMap[relativePath] == nil {
			importMap[relativePath] = make(map[string]bool)
		}

		// For models, we import the concrete class
		// For interfaces, we import the interface type
		if importType == "models" {
			// Import concrete model class (e.g., "SampleRequest")
			importMap[relativePath][msg.TSName] = true
		} else {
			// Import interface type with "as" alias (e.g., "SampleRequest as SampleRequestInterface")
			importMap[relativePath][msg.TSName] = true
		}
	}

	// Convert map to slice of TSImportGroup
	var groups []TSImportGroup
	for importPath, types := range importMap {
		var typeList []string
		for typeName := range types {
			typeList = append(typeList, typeName)
		}
		// Sort for deterministic output
		sort.Strings(typeList)
		groups = append(groups, TSImportGroup{
			ImportPath: importPath,
			Types:      typeList,
		})
	}

	// Sort groups by import path for deterministic output
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].ImportPath < groups[j].ImportPath
	})

	return groups
}
