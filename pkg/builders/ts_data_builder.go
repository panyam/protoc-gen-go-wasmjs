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
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/core"
	"github.com/panyam/protoc-gen-go-wasmjs/pkg/filters"
)

// TSTemplateData represents all data needed for TypeScript template generation.
// This includes client interfaces, message types, enums, and factory data.
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
	TypeImports []string // TypeScript types to import for method signatures

	// TypeScript-specific configuration
	ImportBasePath string // Base path for imports (e.g., "./library/v1")

	// Cross-package dependencies
	ExternalImports []ExternalImport // Imports from other packages

	// Additional template-specific fields
	BaseName           string // Base name for files (e.g., "library_v1")
	DeserializerName   string // Deserializer class name
	FactoryName        string // Factory class name
	SchemaRegistryName string // Schema registry name

	// Computed flags for template conditionals
	HasBrowserServices bool // Whether any browser services exist
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
	Name         string // Proto field name
	TSName       string // TypeScript field name (camelCase)
	TSType       string // TypeScript type
	Number       int32  // Field number
	ProtoFieldID int32  // Proto field ID (alias for Number, for template compatibility)
	DefaultValue string // Default value for the field
	IsOptional   bool   // Whether the field is optional
	IsRepeated   bool   // Whether this is repeated
	IsOneof      bool   // Whether this is part of a oneof
	OneofGroup   string // Oneof group name if applicable
	MessageType  string // If this is a message type field
	Comment      string // Field comment
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
	
	// Collect TypeScript imports needed for typed method signatures
	requiredImports := tb.collectServiceTypeImports(services)
	
	// Generate names for TypeScript artifacts
	baseName := strings.ReplaceAll(packageInfo.Name, ".", "_")
	
	return &TSTemplateData{
		PackageName:     packageInfo.Name,
		PackagePath:     packageInfo.Path,
		SourcePath:      serviceFile.Desc.Path(),
		ModuleName:      baseName,
		Services:        services,
		TypeImports:     requiredImports,
		APIStructure:    config.JSStructure,
		JSNamespace:     tb.getJSNamespace(packageInfo.Name, config),
	}, nil
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
		HasBrowserServices: false, // TODO: Determine from services
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
	externalImports := tb.buildExternalImports(messageResult.Items, packageInfo, config)

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

	// Skip services with no methods
	if len(methods) == 0 {
		return nil, nil
	}

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

// buildExternalImports analyzes messages to determine cross-package type dependencies.
func (tb *TSDataBuilder) buildExternalImports(
	messages []filters.MessageInfo,
	packageInfo *PackageInfo,
	config *GenerationConfig,
) []ExternalImport {

	// For now, return empty - this would be implemented based on message field analysis
	// TODO: Analyze message fields to find cross-package type references
	return []ExternalImport{}
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
		
		tsMsg := TSMessageInfo{
			Name:               msg.Name,
			TSName:             msg.Name, // TypeScript uses same name as proto
			PackageName:        msg.PackageName,
			FullyQualifiedName: msg.FullyQualifiedName,
			ProtoFile:          msg.ProtoFile,
			Comment:            msg.Comment,
			MethodName:         "new" + msg.Name, // Factory method name
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
					msgPackage := string(field.Message.Desc.ParentFile().Package())
					msgName := string(field.Message.Desc.Name())
					fieldInfo.MessageType = msgPackage + "." + msgName
					fieldInfo.TSType = msgName // Simple name for TypeScript
					fieldInfo.DefaultValue = "undefined"
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
