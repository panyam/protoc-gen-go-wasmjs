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
	Imports    []ImportInfo      // Unique package imports with aliases
	PackageMap map[string]string // Maps full package path to alias

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

	// Check if package has services
	hasServices := false
	for _, file := range g.packageFiles {
		if len(file.Services) > 0 {
			hasServices = true
			break
		}
	}

	// Check if package has messages (for TypeScript generation)
	messages := g.collectAllMessages()
	hasMessages := len(messages) > 0

	// Skip packages with no services AND no messages
	if !hasServices && !hasMessages {
		return nil
	}

	// Build template data
	templateData, err := g.buildTemplateData()
	if err != nil {
		return err
	}

	// Generate WASM wrapper if enabled (only when there are services)
	if g.config.GenerateWasm && hasServices {
		if err := g.generateWasmWrapper(templateData); err != nil {
			return err
		}
	}

	// Generate TypeScript artifacts if enabled
	if g.config.GenerateTypeScript {
		// Generate old-style TypeScript client (for now, during transition)
		if templateData != nil {
			if err := g.generateTypeScriptClient(templateData); err != nil {
				return err
			}
		}

		// Generate new TypeScript interfaces, models, factory, schemas, deserializers
		// Generate these for ANY package that has messages (not just services)
		if hasMessages {
			// Get package name for TypeScript generation
			packageName := string(g.file.Desc.Package())
			packagePath := g.buildPackagePath(packageName)

			// Generate TypeScript artifacts for all messages
			tsGen := NewTSGenerator(g)
			if err := tsGen.GenerateAll(messages, packagePath); err != nil {
				return err
			}
		}
	}

	// Generate build script and main example only when generating WASM and have services
	if g.config.GenerateWasm && hasServices {
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

	return &TemplateData{
		PackageName:  packageName,
		SourcePath:   g.file.Desc.Path(),
		GoPackage:    string(g.file.GoImportPath),
		Services:     services,
		Config:       g.config,
		JSNamespace:  g.config.GetDefaultJSNamespace(packageName),
		ModuleName:   g.config.GetDefaultModuleName(packageName),
		APIStructure: g.config.JSStructure,
		Imports:      imports,
		PackageMap:   packageMap,
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

// getQualifiedTypeName returns the fully qualified Go type name for a message
func (g *FileGenerator) getQualifiedTypeName(message *protogen.Message, packageAlias string) string {
	// Return the qualified type name with package alias
	return packageAlias + "." + string(message.GoIdent.GoName)
}

// ====================================================================================
// New TypeScript Generation Functions
// ====================================================================================

// MessageInfo represents a proto message for TypeScript generation
type MessageInfo struct {
	Name        string      // Message name (e.g., "Book")
	GoName      string      // Go struct name (e.g., "Book")
	TSName      string      // TypeScript interface name (e.g., "Book")
	Fields      []FieldInfo // All fields in the message
	PackageName string      // Proto package name
	ProtoFile   string      // Source proto file path
	IsNested    bool        // Whether this is a nested message
	Comment     string      // Leading comment from proto
	MethodName  string      // Factory method name (e.g., "newBook") - for template use
	OneofGroups []string    // List of oneof group names in this message
}

// FieldInfo represents a proto field for TypeScript generation
type FieldInfo struct {
	Name         string // Original proto field name (e.g., "user_id")
	JSONName     string // JSON field name (e.g., "userId")
	TSName       string // TypeScript property name (e.g., "userId")
	ProtoType    string // Original proto type (e.g., "string", "int32")
	TSType       string // TypeScript type (e.g., "string", "number")
	GoType       string // Go type for reference
	IsRepeated   bool   // Whether this is a repeated field
	IsOptional   bool   // Whether this is an optional field
	IsOneof      bool   // Whether this field is part of a oneof
	OneofName    string // Name of the oneof group (if applicable)
	OneofGroup   string // Alias for OneofName (for template compatibility)
	MessageType  string // For message fields, the message type name
	DefaultValue string // Default value for the field
	Comment      string // Field comment from proto
	ProtoFieldID int32  // Proto field number (e.g., text_query = 1)
}

// collectAllMessages collects all message definitions from package files
func (g *FileGenerator) collectAllMessages() []MessageInfo {
	var messages []MessageInfo

	for _, file := range g.packageFiles {
		// Collect top-level messages
		for _, message := range file.Messages {
			// Skip map entry messages (synthetic messages for map fields)
			if message.Desc.IsMapEntry() {
				continue
			}

			messageInfo := g.buildMessageInfo(message, file, false)
			messages = append(messages, messageInfo)

			// Collect nested messages recursively
			nestedMessages := g.collectNestedMessages(message, file)
			messages = append(messages, nestedMessages...)
		}
	}

	return messages
}

// buildMessageInfo constructs MessageInfo from a protogen.Message
func (g *FileGenerator) buildMessageInfo(message *protogen.Message, file *protogen.File, isNested bool) MessageInfo {
	messageName := string(message.Desc.Name())

	// Build field information
	var fields []FieldInfo
	for _, field := range message.Fields {
		fieldInfo := g.buildFieldInfo(field)
		fields = append(fields, fieldInfo)
	}

	// Collect oneof groups
	var oneofGroups []string
	oneofMap := make(map[string]bool)
	for _, oneof := range message.Oneofs {
		oneofName := string(oneof.Desc.Name())
		if !oneofMap[oneofName] {
			oneofGroups = append(oneofGroups, oneofName)
			oneofMap[oneofName] = true
		}
	}

	return MessageInfo{
		Name:        messageName,
		GoName:      string(message.GoIdent.GoName),
		TSName:      messageName, // Same as proto name for interfaces
		Fields:      fields,
		PackageName: string(file.Desc.Package()),
		ProtoFile:   file.Desc.Path(),
		IsNested:    isNested,
		Comment:     strings.TrimSpace(string(message.Comments.Leading)),
		OneofGroups: oneofGroups,
	}
}

// collectNestedMessages recursively collects nested message definitions
func (g *FileGenerator) collectNestedMessages(message *protogen.Message, file *protogen.File) []MessageInfo {
	var nestedMessages []MessageInfo

	for _, nested := range message.Messages {
		// Skip map entry messages (synthetic messages for map fields)
		if nested.Desc.IsMapEntry() {
			continue
		}

		nestedInfo := g.buildMessageInfo(nested, file, true)
		nestedMessages = append(nestedMessages, nestedInfo)

		// Recursively collect deeply nested messages
		deeplyNested := g.collectNestedMessages(nested, file)
		nestedMessages = append(nestedMessages, deeplyNested...)
	}

	return nestedMessages
}

// buildFieldInfo constructs FieldInfo from a protogen.Field
func (g *FileGenerator) buildFieldInfo(field *protogen.Field) FieldInfo {
	fieldName := string(field.Desc.Name())
	jsonName := field.Desc.JSONName() // Proto provides JSON name conversion

	// Convert proto type to TypeScript type
	protoType := g.getProtoFieldType(field)
	tsType := g.convertProtoTypeToTS(protoType, field)
	goType := g.getGoFieldType(field)

	// Check if field is part of a oneof
	isOneof := field.Oneof != nil
	oneofName := ""
	if isOneof {
		oneofName = string(field.Oneof.Desc.Name())
	}

	// For message types, get the fully qualified message type name
	messageType := ""
	if field.Message != nil {
		// Get the full package name + message name
		packageName := string(field.Message.Desc.ParentFile().Package())
		messageName := string(field.Message.Desc.Name())
		if packageName != "" {
			messageType = packageName + "." + messageName
		} else {
			messageType = messageName
		}
	}

	// Get proto field number
	protoFieldID := int32(field.Desc.Number())

	return FieldInfo{
		Name:         fieldName,
		JSONName:     jsonName,
		TSName:       g.convertToTSPropertyName(jsonName),
		ProtoType:    protoType,
		TSType:       tsType,
		GoType:       goType,
		IsRepeated:   field.Desc.IsList(),
		IsOptional:   field.Desc.HasOptionalKeyword(),
		IsOneof:      isOneof,
		OneofName:    oneofName,
		OneofGroup:   oneofName, // Alias for template compatibility
		MessageType:  messageType,
		DefaultValue: g.getDefaultValue(field),
		Comment:      strings.TrimSpace(string(field.Comments.Leading)),
		ProtoFieldID: protoFieldID,
	}
}

// getProtoFieldType returns the proto field type as a string
func (g *FileGenerator) getProtoFieldType(field *protogen.Field) string {
	kind := field.Desc.Kind()
	switch kind.String() {
	case "double", "float", "int32", "int64", "uint32", "uint64", "sint32", "sint64", "fixed32", "fixed64", "sfixed32", "sfixed64":
		return kind.String()
	case "bool":
		return "bool"
	case "string":
		return "string"
	case "bytes":
		return "bytes"
	case "message":
		return string(field.Message.Desc.Name())
	case "enum":
		return string(field.Enum.Desc.Name())
	default:
		return kind.String()
	}
}

// convertProtoTypeToTS converts proto types to TypeScript types
func (g *FileGenerator) convertProtoTypeToTS(protoType string, field *protogen.Field) string {
	// Check if this is a map field
	if g.isMapField(field) {
		mapKeyType, mapValueType := g.getMapKeyValueTypes(field)
		return "Map<" + mapKeyType + ", " + mapValueType + ">"
	}

	// Handle repeated fields
	baseType := g.getBaseTSType(protoType, field)

	if field.Desc.IsList() {
		return baseType + "[]"
	}

	return baseType
}

// getBaseTSType returns the base TypeScript type (without array notation)
func (g *FileGenerator) getBaseTSType(protoType string, field *protogen.Field) string {
	switch protoType {
	case "double", "float", "int32", "int64", "uint32", "uint64", "sint32", "sint64", "fixed32", "fixed64", "sfixed32", "sfixed64":
		return "number"
	case "bool":
		return "boolean"
	case "string":
		return "string"
	case "bytes":
		return "Uint8Array"
	default:
		// For message types, check for external type mappings first
		if field.Message != nil {
			// Get the fully qualified message type name
			packageName := string(field.Message.Desc.ParentFile().Package())
			messageName := string(field.Message.Desc.Name())
			fullTypeName := packageName + "." + messageName
			
			// Check if there's an external type mapping for this type
			if mapping, exists := g.config.GetExternalTypeMapping(fullTypeName); exists {
				return mapping.TypeScript
			}
			
			// Use the message name as interface reference
			return string(field.Message.Desc.Name())
		}
		// For enum types
		if field.Enum != nil {
			return string(field.Enum.Desc.Name())
		}
		// Fallback
		return "any"
	}
}

// getGoFieldType returns the Go field type for reference
func (g *FileGenerator) getGoFieldType(field *protogen.Field) string {
	if field.Message != nil {
		return string(field.Message.GoIdent.GoName)
	}
	if field.Enum != nil {
		return string(field.Enum.GoIdent.GoName)
	}

	// For primitive types, return the Go equivalent
	kind := field.Desc.Kind()
	switch kind.String() {
	case "double":
		return "float64"
	case "float":
		return "float32"
	case "int32", "sint32", "sfixed32":
		return "int32"
	case "int64", "sint64", "sfixed64":
		return "int64"
	case "uint32", "fixed32":
		return "uint32"
	case "uint64", "fixed64":
		return "uint64"
	case "bool":
		return "bool"
	case "string":
		return "string"
	case "bytes":
		return "[]byte"
	default:
		return "interface{}"
	}
}

// convertToTSPropertyName converts JSON names to TypeScript property names
func (g *FileGenerator) convertToTSPropertyName(jsonName string) string {
	// JSON names from proto are already in camelCase, so we can use them directly
	return jsonName
}

// getDefaultValue returns the default value for a field
func (g *FileGenerator) getDefaultValue(field *protogen.Field) string {
	// Handle repeated fields first
	if field.Desc.IsList() {
		return "[]"
	}

	kind := field.Desc.Kind()
	switch kind.String() {
	case "double", "float", "int32", "int64", "uint32", "uint64", "sint32", "sint64", "fixed32", "fixed64", "sfixed32", "sfixed64":
		return "0"
	case "bool":
		return "false"
	case "string":
		return "\"\""
	case "bytes":
		return "new Uint8Array()"
	default:
		if field.Message != nil {
			// For message types, use undefined as default - instances created via factory
			return "undefined"
		}
		if field.Enum != nil {
			// For enums, use the first value (typically 0)
			return "0"
		}
		return "undefined"
	}
}

// buildPackagePath creates the nested directory structure from proto package name
func (g *FileGenerator) buildPackagePath(packageName string) string {
	// Convert "library.v1" to "library/v1" structure
	packagePath := strings.ReplaceAll(packageName, ".", "/")
	return packagePath
}

// isMapField checks if a protobuf field represents a map
func (g *FileGenerator) isMapField(field *protogen.Field) bool {
	// Check if the field descriptor has IsMap method
	if field.Message != nil {
		// For map fields, protogen might not set IsList but the message will be a map entry
		return field.Message.Desc.IsMapEntry()
	}
	
	return false
}

// getMapKeyValueTypes extracts the key and value types from a map field
func (g *FileGenerator) getMapKeyValueTypes(field *protogen.Field) (string, string) {
	if !g.isMapField(field) {
		return "any", "any"
	}

	mapEntry := field.Message

	// Map entry messages have exactly 2 fields: key (field 1) and value (field 2)
	var keyField, valueField *protogen.Field
	for _, f := range mapEntry.Fields {
		if f.Desc.Number() == 1 {
			keyField = f
		} else if f.Desc.Number() == 2 {
			valueField = f
		}
	}

	if keyField == nil || valueField == nil {
		return "any", "any"
	}

	// Convert key and value types to TypeScript
	keyType := g.getBaseTSType(g.getProtoFieldType(keyField), keyField)
	valueType := g.getBaseTSType(g.getProtoFieldType(valueField), valueField)

	return keyType, valueType
}
