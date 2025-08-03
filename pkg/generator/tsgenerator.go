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

// InterfaceTemplateData holds data for interface template generation
type InterfaceTemplateData struct {
	Messages []MessageInfo
	BaseName string
}

// ModelTemplateData holds data for model template generation  
type ModelTemplateData struct {
	Messages []MessageInfo
	BaseName string
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
	Messages           []MessageInfo
	FactoryName        string
	ImportsByFile      map[string][]string
	ModelImportsByFile map[string][]string
	Dependencies       []FactoryDependency // Cross-package factory dependencies
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
				// Extract package name from fully qualified message type
				parts := strings.Split(field.MessageType, ".")
				if len(parts) >= 2 {
					fieldPackage := strings.Join(parts[:len(parts)-1], ".")
					
					// Skip if it's the same package
					if fieldPackage == currentPackage {
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
func (ts *TSGenerator) GenerateAll(messages []MessageInfo, packagePath string) error {
	if len(messages) == 0 {
		return nil
	}
	
	// Generate interfaces
	if err := ts.generateInterfaces(messages, packagePath); err != nil {
		return err
	}
	
	// Generate model classes
	if err := ts.generateModels(messages, packagePath); err != nil {
		return err
	}
	
	// Generate factory
	if err := ts.generateFactory(messages, packagePath); err != nil {
		return err
	}
	
	// Generate schemas
	if err := ts.generateSchemas(messages, packagePath); err != nil {
		return err
	}
	
	// Generate deserializer
	if err := ts.generateDeserializers(messages, packagePath); err != nil {
		return err
	}
	
	return nil
}

// generateInterfaces generates TypeScript interface files for all messages
func (ts *TSGenerator) generateInterfaces(messages []MessageInfo, packagePath string) error {
	// Group messages by proto file
	fileMessages := ts.groupMessagesByProtoFile(messages)
	
	for protoFile, msgs := range fileMessages {
		if err := ts.generateInterfaceFile(protoFile, msgs, packagePath); err != nil {
			return err
		}
	}
	
	return nil
}

// generateModels generates TypeScript model class files for all messages
func (ts *TSGenerator) generateModels(messages []MessageInfo, packagePath string) error {
	// Group messages by proto file
	fileMessages := ts.groupMessagesByProtoFile(messages)
	
	for protoFile, msgs := range fileMessages {
		if err := ts.generateModelFile(protoFile, msgs, packagePath); err != nil {
			return err
		}
	}
	
	return nil
}

// generateFactory generates a factory file for creating message instances
func (ts *TSGenerator) generateFactory(messages []MessageInfo, packagePath string) error {
	return ts.generateFactoryFile(messages, packagePath)
}

// generateSchemas generates schema files for runtime type information
func (ts *TSGenerator) generateSchemas(messages []MessageInfo, packagePath string) error {
	// Group messages by proto file
	fileMessages := ts.groupMessagesByProtoFile(messages)
	
	for protoFile, msgs := range fileMessages {
		if err := ts.generateSchemaFile(protoFile, msgs, packagePath); err != nil {
			return err
		}
	}
	
	return nil
}

// generateDeserializers generates deserializer files for schema-aware deserialization
func (ts *TSGenerator) generateDeserializers(messages []MessageInfo, packagePath string) error {
	// Group messages by proto file
	fileMessages := ts.groupMessagesByProtoFile(messages)
	
	for protoFile, msgs := range fileMessages {
		if err := ts.generateDeserializerFile(protoFile, msgs, packagePath); err != nil {
			return err
		}
	}
	
	return nil
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
func (ts *TSGenerator) generateFactoryFile(messages []MessageInfo, packagePath string) error {
	filename := filepath.Join(packagePath, "factory.ts")
	
	// Create generated file
	generatedFile := ts.fileGen.plugin.NewGeneratedFile(filename, "")
	
	// Generate factory content
	content, err := ts.generateFactoryContent(messages)
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
		Messages: messages,
		BaseName: baseName,
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
		Messages: messages,
		BaseName: baseName,
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
func (ts *TSGenerator) generateFactoryContent(messages []MessageInfo) (string, error) {
	// Collect all interface imports organized by file
	importsByFile := make(map[string][]string)
	modelImportsByFile := make(map[string][]string)
	
	for _, msg := range messages {
		baseName := ts.getBaseFileName(msg.ProtoFile)
		importsByFile[baseName] = append(importsByFile[baseName], msg.TSName+" as "+msg.TSName+"Interface")
		modelImportsByFile[baseName] = append(modelImportsByFile[baseName], msg.TSName+" as Concrete"+msg.TSName)
	}
	
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
		Messages:           messagesWithMethods,
		FactoryName:        factoryName,
		ImportsByFile:      importsByFile,
		ModelImportsByFile: modelImportsByFile,
		Dependencies:       dependencies,
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
	Messages         []MessageInfo
	BaseName         string
	PackageName      string
	DeserializerName string
}

// generateSchemaContent generates the TypeScript schema file content using templates
func (ts *TSGenerator) generateSchemaContent(messages []MessageInfo, baseName string) (string, error) {
	// Get package name from first message (all messages in a file share the same package)
	packageName := ""
	if len(messages) > 0 {
		packageName = messages[0].PackageName
	}
	
	// Generate registry name from package (e.g., "library.v1" -> "LibraryV1SchemaRegistry")
	registryName := ts.buildSchemaRegistryName(packageName)
	
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