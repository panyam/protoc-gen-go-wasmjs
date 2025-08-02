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

// FactoryTemplateData holds data for factory template generation
type FactoryTemplateData struct {
	Messages           []MessageInfo
	FactoryName        string
	ImportsByFile      map[string][]string
	ModelImportsByFile map[string][]string
}

// TSGenerator handles TypeScript-specific generation
type TSGenerator struct {
	fileGen *FileGenerator
}

// NewTSGenerator creates a new TypeScript generator
func NewTSGenerator(fileGen *FileGenerator) *TSGenerator {
	return &TSGenerator{fileGen: fileGen}
}

// GenerateAll generates all TypeScript artifacts (interfaces, models, factory)
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

// getBaseFileName converts proto file path to base filename
func (ts *TSGenerator) getBaseFileName(protoFile string) string {
	// Remove .proto extension and convert path separators
	baseName := strings.TrimSuffix(protoFile, ".proto")
	baseName = strings.ReplaceAll(baseName, "/", "_")
	baseName = strings.ReplaceAll(baseName, "\\", "_")
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
	
	data := FactoryTemplateData{
		Messages:           messagesWithMethods,
		FactoryName:        factoryName,
		ImportsByFile:      importsByFile,
		ModelImportsByFile: modelImportsByFile,
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

// buildFactoryName converts package name to factory class name
func (ts *TSGenerator) buildFactoryName(packageName string) string {
	parts := strings.Split(packageName, ".")
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}
	return strings.Join(parts, "") + "Factory"
}