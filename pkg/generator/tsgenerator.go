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

// ====================================================================================
// TypeScript Generation Functions
// ====================================================================================

// ExternalImport represents an external type import
type ExternalImport struct {
	TypeName     string // e.g., "Date" or "Timestamp"
	ImportSource string // e.g., "<native>" or "@bufbuild/protobuf/wkt"
	IsNative     bool   // true if this is a native TypeScript type
}

// InterfaceTemplateData holds data for interface template generation
type InterfaceTemplateData struct {
	Messages        []MessageInfo
	Enums           []EnumInfo
	BaseName        string
	ExternalImports []ExternalImport
}

// ModelTemplateData holds data for model template generation  
type ModelTemplateData struct {
	Messages        []MessageInfo
	Enums           []EnumInfo
	BaseName        string
	ExternalImports []ExternalImport
	DeserializerName string // e.g., "LibraryV2Deserializer"
}

// FactoryDependency represents a dependency on another package's factory
type FactoryDependency struct {
	PackageName  string // e.g., "library.common"
	FactoryName  string // e.g., "LibraryCommonFactory"
	ImportPath   string // e.g., "../common/factory"
	InstanceName string // e.g., "commonFactory"
}

// FactoryTemplateData holds data for factory template generation
type FactoryTemplateData struct {
	Messages        []MessageInfo
	Enums           []EnumInfo
	FactoryName     string
	Dependencies    []FactoryDependency // Cross-package factory dependencies
	ExternalImports []ExternalImport
}

// TSGenerator handles TypeScript-specific generation
type TSGenerator struct {
	fileGen *FileGenerator
}

// collectFactoryDependencies analyzes messages to find cross-package dependencies
func (ts *TSGenerator) collectFactoryDependencies(messages []MessageInfo, currentPackage string) []FactoryDependency {
	dependencyMap := make(map[string]FactoryDependency)
	
	for _, message := range messages {
		for _, field := range message.Fields {
			// Check if field references a message from another package
			if field.MessageType != "" && strings.Contains(field.MessageType, ".") {
				// Skip external types that are mapped to TypeScript types
				if _, isExternal := ts.fileGen.config.GetExternalTypeMapping(field.MessageType); isExternal {
					continue
				}
				
				// Extract package name from fully qualified message type
				parts := strings.Split(field.MessageType, ".")
				if len(parts) >= 2 {
					fieldPackage := strings.Join(parts[:len(parts)-1], ".")
					
					// Skip if it's the same package
					if fieldPackage == currentPackage {
						continue
					}
					
					// Skip wasmjs annotation packages - these are framework types
					if fieldPackage == "wasmjs.v1" {
						continue
					}
					
					// Create dependency if not already exists
					if _, exists := dependencyMap[fieldPackage]; !exists {
						factoryName := ts.getFactoryNameForPackage(fieldPackage)
						importPath := ts.getFactoryImportPath(fieldPackage, currentPackage)
						instanceName := ts.getFactoryInstanceName(fieldPackage)
						
						dependencyMap[fieldPackage] = FactoryDependency{
							PackageName:  fieldPackage,
							FactoryName:  factoryName,
							ImportPath:   importPath,
							InstanceName: instanceName,
						}
					}
				}
			}
		}
	}
	
	// Convert map to slice
	var dependencies []FactoryDependency
	for _, dep := range dependencyMap {
		dependencies = append(dependencies, dep)
	}
	
	return dependencies
}

// collectExternalImports analyzes messages to find external type imports needed
func (ts *TSGenerator) collectExternalImports(messages []MessageInfo) []ExternalImport {
	importMap := make(map[string]ExternalImport)
	
	// Get current package name for comparison
	currentPackage := ""
	if len(messages) > 0 {
		currentPackage = messages[0].PackageName
	}
	
	for _, message := range messages {
		for _, field := range message.Fields {
			// Check if field references an external type
			if field.MessageType != "" && strings.Contains(field.MessageType, ".") {
				// Check if there's an external type mapping for this type (e.g., google.protobuf.Timestamp -> Date)
				if mapping, exists := ts.fileGen.config.GetExternalTypeMapping(field.MessageType); exists {
					// Use the TypeScript type name as the key to avoid duplicates
					importMap[mapping.TypeScript] = ExternalImport{
						TypeName:     mapping.TypeScript,
						ImportSource: mapping.ImportSource,
						IsNative:     mapping.IsNative,
					}
				} else {
					// Check if this is a cross-package message reference within the project
					fieldPackage := ts.extractPackageName(field.MessageType)
					if fieldPackage != currentPackage && fieldPackage != "" {
						// Skip wasmjs annotation packages - these are framework types
						if fieldPackage == "wasmjs.v1" {
							continue
						}
						
						// This is a cross-package reference - generate import for it
						messageName := ts.extractMessageName(field.MessageType)
						importPath := ts.buildCrossPackageImportPath(currentPackage, fieldPackage)
						
						importMap[messageName] = ExternalImport{
							TypeName:     messageName,
							ImportSource: importPath,
							IsNative:     false,
						}
					}
				}
			}
		}
	}
	
	// Convert map to slice
	var imports []ExternalImport
	for _, imp := range importMap {
		imports = append(imports, imp)
	}
	
	return imports
}

// extractPackageName extracts the package name from a fully qualified message type
// e.g., "library.common.BaseMessage" -> "library.common"
func (ts *TSGenerator) extractPackageName(fullMessageType string) string {
	parts := strings.Split(fullMessageType, ".")
	if len(parts) <= 1 {
		return ""
	}
	// Return all parts except the last one (which is the message name)
	return strings.Join(parts[:len(parts)-1], ".")
}

// extractMessageName extracts the message name from a fully qualified message type
// e.g., "library.common.BaseMessage" -> "BaseMessage"
func (ts *TSGenerator) extractMessageName(fullMessageType string) string {
	parts := strings.Split(fullMessageType, ".")
	if len(parts) == 0 {
		return fullMessageType
	}
	// Return the last part (which is the message name)
	return parts[len(parts)-1]
}

// buildCrossPackageImportPath builds the relative import path from current package to target package
// e.g., from "library.v2" to "library.common" -> "../common/interfaces"
func (ts *TSGenerator) buildCrossPackageImportPath(currentPackage, targetPackage string) string {
	currentParts := strings.Split(currentPackage, ".")
	targetParts := strings.Split(targetPackage, ".")
	
	// Find common prefix
	commonPrefixLen := 0
	for i := 0; i < len(currentParts) && i < len(targetParts) && currentParts[i] == targetParts[i]; i++ {
		commonPrefixLen++
	}
	
	// Build relative path
	var pathParts []string
	
	// Go up for each unique part in current package
	for i := commonPrefixLen; i < len(currentParts); i++ {
		pathParts = append(pathParts, "..")
	}
	
	// Go down for each unique part in target package
	for i := commonPrefixLen; i < len(targetParts); i++ {
		pathParts = append(pathParts, targetParts[i])
	}
	
	// Add "interfaces" at the end since we import from interfaces.ts
	pathParts = append(pathParts, "interfaces")
	
	return strings.Join(pathParts, "/")
}

// getFactoryNameForPackage generates factory class name from package name
func (ts *TSGenerator) getFactoryNameForPackage(packageName string) string {
	// Convert "library.common" to "LibraryCommonFactory"
	parts := strings.Split(packageName, ".")
	var camelParts []string
	for _, part := range parts {
		camelParts = append(camelParts, strings.Title(part))
	}
	return strings.Join(camelParts, "") + "Factory"
}

// getFactoryImportPath generates relative import path for factory
func (ts *TSGenerator) getFactoryImportPath(dependencyPackage, currentPackage string) string {
	// Convert package names to relative paths
	// e.g., from "library.v2" to "library.common" -> "../common/factory"
	currentParts := strings.Split(currentPackage, ".")
	depParts := strings.Split(dependencyPackage, ".")
	
	// Find common prefix
	commonLen := 0
	for i := 0; i < len(currentParts) && i < len(depParts); i++ {
		if currentParts[i] == depParts[i] {
			commonLen++
		} else {
			break
		}
	}
	
	// Calculate relative path
	var pathParts []string
	
	// Go up for each unique part in current package
	for i := commonLen; i < len(currentParts); i++ {
		pathParts = append(pathParts, "..")
	}
	
	// Go down for each unique part in dependency package
	for i := commonLen; i < len(depParts); i++ {
		pathParts = append(pathParts, depParts[i])
	}
	
	pathParts = append(pathParts, "factory")
	return strings.Join(pathParts, "/")
}

// getFactoryInstanceName generates instance variable name from package name
func (ts *TSGenerator) getFactoryInstanceName(packageName string) string {
	// Convert "library.common" to "commonFactory"
	parts := strings.Split(packageName, ".")
	lastPart := parts[len(parts)-1]
	return lastPart + "Factory"
}

// NewTSGenerator creates a new TypeScript generator
func NewTSGenerator(fileGen *FileGenerator) *TSGenerator {
	return &TSGenerator{fileGen: fileGen}
}

// GenerateAll generates all TypeScript artifacts (interfaces, models, factory, schemas, deserializer)
func (ts *TSGenerator) GenerateAll(messages []MessageInfo, enums []EnumInfo, packagePath string) error {
	if len(messages) == 0 && len(enums) == 0 {
		return nil
	}
	
	// Generate interfaces (includes both messages and enums)
	if err := ts.generateInterfaces(messages, enums, packagePath); err != nil {
		return err
	}
	
	// Generate model classes (only for messages)
	if len(messages) > 0 {
		if err := ts.generateModels(messages, enums, packagePath); err != nil {
			return err
		}
	}
	
	// Generate factory (only for messages)
	if len(messages) > 0 {
		if err := ts.generateFactory(messages, enums, packagePath); err != nil {
			return err
		}
	}
	
	// Generate deserializer schemas (framework types) - only if we have messages
	if len(messages) > 0 {
		if err := ts.generateDeserializerSchemas(packagePath); err != nil {
			return err
		}
	}
	
	// Generate schemas (only for messages)
	if len(messages) > 0 {
		if err := ts.generateSchemas(messages, packagePath); err != nil {
			return err
		}
	}
	
	// Generate deserializer (only for messages)
	if len(messages) > 0 {
		if err := ts.generateDeserializers(messages, packagePath); err != nil {
			return err
		}
	}
	
	return nil
}

// generateInterfaces generates a single TypeScript interface file for the entire package
func (ts *TSGenerator) generateInterfaces(messages []MessageInfo, enums []EnumInfo, packagePath string) error {
	// Generate a single interfaces file per package to avoid import issues
	if err := ts.generatePackageInterfaceFile(messages, enums, packagePath); err != nil {
		return err
	}
	
	return nil
}

// generateModels generates a single TypeScript model file for the entire package
func (ts *TSGenerator) generateModels(messages []MessageInfo, enums []EnumInfo, packagePath string) error {
	// Generate a single models file per package to avoid import issues
	if err := ts.generatePackageModelFile(messages, enums, packagePath); err != nil {
		return err
	}
	
	return nil
}

// generateFactory generates a factory file for creating message instances
func (ts *TSGenerator) generateFactory(messages []MessageInfo, enums []EnumInfo, packagePath string) error {
	return ts.generateFactoryFile(messages, enums, packagePath)
}

// generateSchemas generates a single schema file for the entire package
func (ts *TSGenerator) generateSchemas(messages []MessageInfo, packagePath string) error {
	// Generate a single schemas file per package to avoid import issues
	if err := ts.generatePackageSchemaFile(messages, packagePath); err != nil {
		return err
	}
	
	return nil
}

// generateDeserializers generates a single deserializer file for the entire package
func (ts *TSGenerator) generateDeserializers(messages []MessageInfo, packagePath string) error {
	// Generate a single deserializer file per package to avoid import issues
	if err := ts.generatePackageDeserializerFile(messages, packagePath); err != nil {
		return err
	}
	
	return nil
}

// generatePackageInterfaceFile generates a single interfaces file for the entire package
func (ts *TSGenerator) generatePackageInterfaceFile(messages []MessageInfo, enums []EnumInfo, packagePath string) error {
	// Get package name from first message or enum
	packageName := ""
	if len(messages) > 0 {
		packageName = messages[0].PackageName
	} else if len(enums) > 0 {
		packageName = enums[0].PackageName
	}
	
	// Create filename based on package name (e.g., "interfaces.ts")
	filename := filepath.Join(packagePath, "interfaces.ts")
	
	// Create generated file
	generatedFile := ts.fileGen.plugin.NewGeneratedFile(filename, "")
	
	// Generate interface content
	content, err := ts.generatePackageInterfaceContent(messages, enums, packageName)
	if err != nil {
		return err
	}
	
	_, err = generatedFile.Write([]byte(content))
	return err
}

// generatePackageModelFile generates a single models file for the entire package
func (ts *TSGenerator) generatePackageModelFile(messages []MessageInfo, enums []EnumInfo, packagePath string) error {
	// Get package name from first message
	packageName := ""
	if len(messages) > 0 {
		packageName = messages[0].PackageName
	}
	
	// Create filename based on package name (e.g., "models.ts")
	filename := filepath.Join(packagePath, "models.ts")
	
	// Create generated file
	generatedFile := ts.fileGen.plugin.NewGeneratedFile(filename, "")
	
	// Generate model content
	content, err := ts.generatePackageModelContent(messages, enums, packageName)
	if err != nil {
		return err
	}
	
	_, err = generatedFile.Write([]byte(content))
	return err
}

// generatePackageSchemaFile generates a single schemas file for the entire package
func (ts *TSGenerator) generatePackageSchemaFile(messages []MessageInfo, packagePath string) error {
	// Get package name from first message
	packageName := ""
	if len(messages) > 0 {
		packageName = messages[0].PackageName
	}
	
	// Create filename based on package name (e.g., "schemas.ts")
	filename := filepath.Join(packagePath, "schemas.ts")
	
	// Create generated file
	generatedFile := ts.fileGen.plugin.NewGeneratedFile(filename, "")
	
	// Generate schema content
	content, err := ts.generateSchemaContent(messages, packageName)
	if err != nil {
		return err
	}
	
	_, err = generatedFile.Write([]byte(content))
	return err
}

// generatePackageDeserializerFile generates a single deserializer file for the entire package  
func (ts *TSGenerator) generatePackageDeserializerFile(messages []MessageInfo, packagePath string) error {
	// Get package name from first message
	packageName := ""
	if len(messages) > 0 {
		packageName = messages[0].PackageName
	}
	
	// Create filename based on package name (e.g., "deserializer.ts")
	filename := filepath.Join(packagePath, "deserializer.ts")
	
	// Create generated file
	generatedFile := ts.fileGen.plugin.NewGeneratedFile(filename, "")
	
	// Generate deserializer content
	content, err := ts.generatePackageDeserializerContent(messages, packageName)
	if err != nil {
		return err
	}
	
	_, err = generatedFile.Write([]byte(content))
	return err
}

// groupMessagesByProtoFile groups messages by their source proto file
func (ts *TSGenerator) groupMessagesByProtoFile(messages []MessageInfo) map[string][]MessageInfo {
	fileMessages := make(map[string][]MessageInfo)
	
	for _, msg := range messages {
		protoFile := msg.ProtoFile
		fileMessages[protoFile] = append(fileMessages[protoFile], msg)
	}
	
	return fileMessages
}

// generateInterfaceFile generates a TypeScript interface file for messages from one proto file
func (ts *TSGenerator) generateInterfaceFile(protoFile string, messages []MessageInfo, packagePath string) error {
	// Create filename: proto/path/file.proto -> proto_path_file_interfaces.ts
	baseName := ts.getBaseFileName(protoFile)
	filename := filepath.Join(packagePath, baseName+"_interfaces.ts")
	
	// Create generated file
	generatedFile := ts.fileGen.plugin.NewGeneratedFile(filename, "")
	
	// Generate TypeScript interfaces
	content, err := ts.generateInterfaceContent(messages, baseName)
	if err != nil {
		return err
	}
	
	_, err = generatedFile.Write([]byte(content))
	return err
}

// generateModelFile generates a TypeScript model class file for messages from one proto file
func (ts *TSGenerator) generateModelFile(protoFile string, messages []MessageInfo, packagePath string) error {
	// Create filename: proto/path/file.proto -> proto_path_file_models.ts
	baseName := ts.getBaseFileName(protoFile)
	filename := filepath.Join(packagePath, baseName+"_models.ts")
	
	// Create generated file
	generatedFile := ts.fileGen.plugin.NewGeneratedFile(filename, "")
	
	// Generate TypeScript model classes
	content, err := ts.generateModelContent(messages, baseName)
	if err != nil {
		return err
	}
	
	_, err = generatedFile.Write([]byte(content))
	return err
}

// generateFactoryFile generates a TypeScript factory file
func (ts *TSGenerator) generateFactoryFile(messages []MessageInfo, enums []EnumInfo, packagePath string) error {
	filename := filepath.Join(packagePath, "factory.ts")
	
	// Create generated file
	generatedFile := ts.fileGen.plugin.NewGeneratedFile(filename, "")
	
	// Generate factory content
	content, err := ts.generateFactoryContent(messages, enums)
	if err != nil {
		return err
	}
	
	_, err = generatedFile.Write([]byte(content))
	return err
}

// generateSchemaFile generates a TypeScript schema file for messages from one proto file
func (ts *TSGenerator) generateSchemaFile(protoFile string, messages []MessageInfo, packagePath string) error {
	// Create filename: proto/path/file.proto -> proto_path_file_schemas.ts
	baseName := ts.getBaseFileName(protoFile)
	filename := filepath.Join(packagePath, baseName+"_schemas.ts")
	
	// Create generated file
	generatedFile := ts.fileGen.plugin.NewGeneratedFile(filename, "")
	
	// Generate TypeScript schemas
	content, err := ts.generateSchemaContent(messages, baseName)
	if err != nil {
		return err
	}
	
	_, err = generatedFile.Write([]byte(content))
	return err
}

// generateDeserializerFile generates a TypeScript deserializer file for messages from one proto file
func (ts *TSGenerator) generateDeserializerFile(protoFile string, messages []MessageInfo, packagePath string) error {
	// Create filename: proto/path/file.proto -> proto_path_file_deserializer.ts
	baseName := ts.getBaseFileName(protoFile)
	filename := filepath.Join(packagePath, baseName+"_deserializer.ts")
	
	// Create generated file
	generatedFile := ts.fileGen.plugin.NewGeneratedFile(filename, "")
	
	// Generate TypeScript deserializer
	content, err := ts.generateDeserializerContent(messages, baseName)
	if err != nil {
		return err
	}
	
	_, err = generatedFile.Write([]byte(content))
	return err
}

// getBaseFileName extracts just the filename from proto file path
func (ts *TSGenerator) getBaseFileName(protoFile string) string {
	// Get just the filename without path and extension
	// e.g., "library/v1/library.proto" -> "library"
	baseName := filepath.Base(protoFile)
	baseName = strings.TrimSuffix(baseName, ".proto")
	return baseName
}

// generateInterfaceContent generates the TypeScript interface file content using templates
func (ts *TSGenerator) generateInterfaceContent(messages []MessageInfo, baseName string) (string, error) {
	data := InterfaceTemplateData{
		Messages:        messages,
		BaseName:        baseName,
		ExternalImports: ts.collectExternalImports(messages),
	}
	
	tmpl, err := template.New("interfaces").Funcs(templateFuncMap).Parse(interfacesTemplate)
	if err != nil {
		return "", err
	}
	
	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", err
	}
	
	return result.String(), nil
}

// generateModelContent generates the TypeScript model class file content using templates
func (ts *TSGenerator) generateModelContent(messages []MessageInfo, baseName string) (string, error) {
	data := ModelTemplateData{
		Messages:        messages,
		BaseName:        baseName,
		ExternalImports: ts.collectExternalImports(messages),
	}
	
	tmpl, err := template.New("models").Funcs(templateFuncMap).Parse(modelsTemplate)
	if err != nil {
		return "", err
	}
	
	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", err
	}
	
	return result.String(), nil
}

// generateFactoryContent generates the TypeScript factory file content using templates
func (ts *TSGenerator) generateFactoryContent(messages []MessageInfo, enums []EnumInfo) (string, error) {
	// No longer need file-based imports since we use package-level imports in template
	
	// Generate factory class name from package
	factoryName := ""
	if len(messages) > 0 {
		factoryName = ts.buildFactoryName(messages[0].PackageName)
	}
	
	// Add method names to messages for template
	messagesWithMethods := make([]MessageInfo, len(messages))
	for i, msg := range messages {
		msgCopy := msg
		msgCopy.MethodName = "new" + msg.TSName
		messagesWithMethods[i] = msgCopy
	}
	
	// Collect factory dependencies
	currentPackage := ""
	if len(messages) > 0 {
		currentPackage = messages[0].PackageName
	}
	dependencies := ts.collectFactoryDependencies(messages, currentPackage)
	
	data := FactoryTemplateData{
		Messages:        messagesWithMethods,
		Enums:           enums,
		FactoryName:     factoryName,
		Dependencies:    dependencies,
		ExternalImports: ts.collectExternalImports(messages),
	}
	
	tmpl, err := template.New("factory").Funcs(templateFuncMap).Parse(factoryTemplate)
	if err != nil {
		return "", err
	}
	
	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", err
	}
	
	return result.String(), nil
}

// SchemaTemplateData holds data for schema template generation
type SchemaTemplateData struct {
	Messages        []MessageInfo
	BaseName        string
	PackageName     string
	RegistryName    string
}

// DeserializerTemplateData holds data for deserializer template generation
type DeserializerTemplateData struct {
	Messages            []MessageInfo
	BaseName            string
	PackageName         string
	DeserializerName    string
	FactoryName         string // e.g., "LibraryV2Factory"
	SchemaRegistryName  string // e.g., "libraryV2SchemaRegistry"
}

// generatePackageInterfaceContent generates TypeScript interface content for the entire package
func (ts *TSGenerator) generatePackageInterfaceContent(messages []MessageInfo, enums []EnumInfo, packageName string) (string, error) {
	data := InterfaceTemplateData{
		Messages:        messages,
		Enums:           enums,
		BaseName:        "interfaces",
		ExternalImports: ts.collectExternalImports(messages),
	}
	
	tmpl, err := template.New("interfaces").Funcs(templateFuncMap).Parse(interfacesTemplate)
	if err != nil {
		return "", err
	}
	
	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", err
	}
	
	return result.String(), nil
}

// generatePackageModelContent generates TypeScript model content for the entire package
func (ts *TSGenerator) generatePackageModelContent(messages []MessageInfo, enums []EnumInfo, packageName string) (string, error) {
	data := ModelTemplateData{
		Messages:         messages,
		Enums:            enums,
		BaseName:         "interfaces", // Changed from "models" to "interfaces" for package-based imports
		ExternalImports:  ts.collectExternalImports(messages),
		DeserializerName: ts.buildDeserializerName(packageName),
	}
	
	tmpl, err := template.New("models").Funcs(templateFuncMap).Parse(modelsTemplate)
	if err != nil {
		return "", err
	}
	
	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", err
	}
	
	return result.String(), nil
}

// generatePackageDeserializerContent generates TypeScript deserializer content for the entire package
func (ts *TSGenerator) generatePackageDeserializerContent(messages []MessageInfo, packageName string) (string, error) {
	// Generate deserializer name from package (e.g., "library.v1" -> "LibraryV1Deserializer")
	deserializerName := ts.buildDeserializerName(packageName)
	factoryName := ts.buildFactoryName(packageName)
	schemaRegistryName := ts.buildSchemaRegistryName(packageName)
	
	data := DeserializerTemplateData{
		Messages:           messages,
		BaseName:           "deserializer",
		PackageName:        packageName,
		DeserializerName:   deserializerName,
		FactoryName:        factoryName,
		SchemaRegistryName: schemaRegistryName,
	}
	
	tmpl, err := template.New("deserializer").Funcs(templateFuncMap).Parse(deserializerTemplate)
	if err != nil {
		return "", err
	}
	
	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", err
	}
	
	return result.String(), nil
}

// generateSchemaContent generates the TypeScript schema file content using templates
func (ts *TSGenerator) generateSchemaContent(messages []MessageInfo, packageNameOverride string) (string, error) {
	// Use the provided package name (for package-level generation)
	packageName := packageNameOverride
	if packageName == "" && len(messages) > 0 {
		packageName = messages[0].PackageName
	}
	
	// Generate registry name from package (e.g., "library.v1" -> "LibraryV1SchemaRegistry")
	registryName := ts.buildSchemaRegistryName(packageName)
	
	// For package-level generation, use "schemas" as base name
	baseName := "schemas"
	
	data := SchemaTemplateData{
		Messages:     messages,
		BaseName:     baseName,
		PackageName:  packageName,
		RegistryName: registryName,
	}
	
	tmpl, err := template.New("schemas").Funcs(templateFuncMap).Parse(schemasTemplate)
	if err != nil {
		return "", err
	}
	
	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", err
	}
	
	return result.String(), nil
}

// generateDeserializerContent generates the TypeScript deserializer file content using templates
func (ts *TSGenerator) generateDeserializerContent(messages []MessageInfo, baseName string) (string, error) {
	// Get package name from first message (all messages in a file share the same package)
	packageName := ""
	if len(messages) > 0 {
		packageName = messages[0].PackageName
	}
	
	// Generate deserializer name from package (e.g., "library.v1" -> "LibraryV1Deserializer")
	deserializerName := ts.buildDeserializerName(packageName)
	
	data := DeserializerTemplateData{
		Messages:         messages,
		BaseName:         baseName,
		PackageName:      packageName,
		DeserializerName: deserializerName,
	}
	
	tmpl, err := template.New("deserializer").Funcs(templateFuncMap).Parse(deserializerTemplate)
	if err != nil {
		return "", err
	}
	
	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", err
	}
	
	return result.String(), nil
}

// generateDeserializerSchemas generates the framework schema types for deserializer
func (ts *TSGenerator) generateDeserializerSchemas(packagePath string) error {
	// Create filename for deserializer schemas
	filename := filepath.Join(packagePath, "deserializer_schemas.ts")
	
	// Create generated file
	generatedFile := ts.fileGen.plugin.NewGeneratedFile(filename, "")
	
	// Generate deserializer schema content (no template data needed - just static types)
	content, err := ts.generateDeserializerSchemaContent()
	if err != nil {
		return err
	}
	
	_, err = generatedFile.Write([]byte(content))
	return err
}

// generateDeserializerSchemaContent generates the deserializer schema framework types
func (ts *TSGenerator) generateDeserializerSchemaContent() (string, error) {
	// Use empty data since this template doesn't need any dynamic content
	var emptyData struct{}
	
	tmpl, err := template.New("deserializer_schemas").Funcs(templateFuncMap).Parse(deserializerSchemasTemplate)
	if err != nil {
		return "", err
	}
	
	var result strings.Builder
	err = tmpl.Execute(&result, emptyData)
	if err != nil {
		return "", err
	}
	
	return result.String(), nil
}

// buildFactoryName converts package name to factory class name
func (ts *TSGenerator) buildFactoryName(packageName string) string {
	parts := strings.Split(packageName, ".")
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}
	return strings.Join(parts, "") + "Factory"
}

// buildSchemaRegistryName converts package name to schema registry name
func (ts *TSGenerator) buildSchemaRegistryName(packageName string) string {
	parts := strings.Split(packageName, ".")
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}
	return strings.Join(parts, "") + "SchemaRegistry"
}

// buildDeserializerName converts package name to deserializer class name
func (ts *TSGenerator) buildDeserializerName(packageName string) string {
	parts := strings.Split(packageName, ".")
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}
	return strings.Join(parts, "") + "Deserializer"
}