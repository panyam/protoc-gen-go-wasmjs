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
	"strings"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/core"
	"github.com/panyam/protoc-gen-go-wasmjs/pkg/filters"
)

// GoTemplateData represents all data needed for Go WASM template generation.
// This is the complete data structure that Go templates will consume.
type GoTemplateData struct {
	// Package metadata
	PackageName string // Proto package name (e.g., "library.v1")
	SourcePath  string // Primary proto file path
	GoPackage   string // Go import path
	ModuleName  string // WASM module name (e.g., "library_v1_services")

	// Message data (always generated)
	Messages []MessageInfo // All message types in this package
	Enums    []EnumInfo    // All enum types in this package

	// Service implementations (for non-browser-provided services)
	Services []ServiceData // Services that Go implements

	// Client interfaces (for browser-provided services)
	BrowserClients []ServiceData // Clients for browser-provided services

	// JavaScript API configuration
	JSNamespace  string // Global namespace (e.g., "library_v1")
	APIStructure string // namespaced|flat|service_based

	// Import management
	Imports    []ImportInfo      // Go package imports
	PackageMap map[string]string // Import path to alias mapping

	// Flags
	HasMessages        bool // Whether any messages exist
	HasEnums           bool // Whether any enums exist
	HasServices        bool // Whether any services to implement exist
	HasBrowserClients  bool // Whether any browser clients exist
}

// GoDataBuilder builds template data structures specifically for Go WASM generation.
// This handles the complex logic of transforming filtered proto data into Go template data.
type GoDataBuilder struct {
	analyzer      *core.ProtoAnalyzer
	pathCalc      *core.PathCalculator
	nameConv      *core.NameConverter
	serviceFilter *filters.ServiceFilter
	methodFilter  *filters.MethodFilter
	msgCollector  *filters.MessageCollector
	enumCollector *filters.EnumCollector
}

// NewGoDataBuilder creates a new Go data builder with all necessary dependencies.
func NewGoDataBuilder(
	analyzer *core.ProtoAnalyzer,
	pathCalc *core.PathCalculator,
	nameConv *core.NameConverter,
	serviceFilter *filters.ServiceFilter,
	methodFilter *filters.MethodFilter,
	msgCollector *filters.MessageCollector,
	enumCollector *filters.EnumCollector,
) *GoDataBuilder {
	return &GoDataBuilder{
		analyzer:      analyzer,
		pathCalc:      pathCalc,
		nameConv:      nameConv,
		serviceFilter: serviceFilter,
		methodFilter:  methodFilter,
		msgCollector:  msgCollector,
		enumCollector: enumCollector,
	}
}

// BuildTemplateData creates Go template data from filtered proto files and configuration.
// This is the main entry point for building Go WASM template data structures.
func (gb *GoDataBuilder) BuildTemplateData(
	packageInfo *PackageInfo,
	allBrowserServices []*protogen.Service,
	criteria *filters.FilterCriteria,
	config *GenerationConfig,
) (*GoTemplateData, error) {

	// Build context for this generation
	context := NewBuildContext(nil, config, packageInfo)

	// STEP 1: Always collect messages and enums from this package
	messages := gb.collectMessages(packageInfo.Files, criteria, context)
	enums := gb.collectEnums(packageInfo.Files, criteria, context)

	// STEP 2: Build service implementations (for non-browser-provided services)
	serviceImplementations, err := gb.buildServiceImplementations(packageInfo.Files, criteria, context)
	if err != nil {
		return nil, err
	}

	// STEP 3: Build browser clients (for browser-provided services in THIS package)
	browserClients, err := gb.buildBrowserClients(packageInfo.Files, criteria, context)
	if err != nil {
		return nil, err
	}

	// Check if we have anything to generate
	hasContent := len(messages) > 0 || len(enums) > 0 ||
		len(serviceImplementations) > 0 || len(browserClients) > 0

	if !hasContent {
		return nil, nil // Nothing to generate for this package
	}

	// Always add wasm package import if we have services or clients
	if len(serviceImplementations) > 0 || len(browserClients) > 0 {
		context.AddImport("github.com/panyam/protoc-gen-go-wasmjs/pkg/wasm", "wasm")
	}

	// Determine names and structure
	moduleName := gb.getModuleName(packageInfo.Name, config)
	jsNamespace := gb.getJSNamespace(packageInfo.Name, config)

	return &GoTemplateData{
		PackageName:        packageInfo.Name,
		SourcePath:         gb.getPrimarySourcePath(packageInfo.Files),
		GoPackage:          packageInfo.GoPackage,
		ModuleName:         moduleName,
		Messages:           messages,
		Enums:              enums,
		Services:           serviceImplementations,
		BrowserClients:     browserClients,
		JSNamespace:        jsNamespace,
		APIStructure:       config.JSStructure,
		Imports:            context.GetImports(),
		PackageMap:         context.ImportMap,
		HasMessages:        len(messages) > 0,
		HasEnums:           len(enums) > 0,
		HasServices:        len(serviceImplementations) > 0,
		HasBrowserClients:  len(browserClients) > 0,
	}, nil
}

// collectMessages collects all messages from the package files.
// Messages are always generated regardless of service presence.
func (gb *GoDataBuilder) collectMessages(
	files []*protogen.File,
	criteria *filters.FilterCriteria,
	context *BuildContext,
) []MessageInfo {
	var messages []MessageInfo

	// Iterate through files to collect actual protogen.Message objects
	for _, file := range files {
		packagePath := string(file.GoImportPath)
		packageAlias := gb.pathCalc.GetGoPackageAlias(packagePath)

		// Register import for this package
		if packagePath != "" {
			context.AddImport(packagePath, packageAlias)
		}

		// Collect all messages from this file
		for _, message := range file.Messages {
			// Build message info
			msgInfo := MessageInfo{
				Name:        string(message.Desc.Name()),
				GoType:      packageAlias + "." + string(message.GoIdent.GoName),
				PackagePath: packagePath,
			}
			messages = append(messages, msgInfo)
		}
	}

	return messages
}

// collectEnums collects all enums from the package files.
// Enums are always generated regardless of service presence.
func (gb *GoDataBuilder) collectEnums(
	files []*protogen.File,
	criteria *filters.FilterCriteria,
	context *BuildContext,
) []EnumInfo {
	var enums []EnumInfo

	// Iterate through files to collect actual protogen.Enum objects
	for _, file := range files {
		packagePath := string(file.GoImportPath)
		packageAlias := gb.pathCalc.GetGoPackageAlias(packagePath)

		// Register import for this package
		if packagePath != "" {
			context.AddImport(packagePath, packageAlias)
		}

		// Collect all enums from this file
		for _, enum := range file.Enums {
			// Build enum info
			enumInfo := EnumInfo{
				Name:        string(enum.Desc.Name()),
				GoType:      packageAlias + "." + string(enum.GoIdent.GoName),
				PackagePath: packagePath,
			}
			enums = append(enums, enumInfo)
		}
	}

	return enums
}

// buildServiceImplementations builds service data for non-browser-provided services.
// These services are implemented in Go and exposed to JavaScript via WASM.
func (gb *GoDataBuilder) buildServiceImplementations(
	files []*protogen.File,
	criteria *filters.FilterCriteria,
	context *BuildContext,
) ([]ServiceData, error) {

	var services []ServiceData

	for _, file := range files {
		for _, service := range file.Services {
			// Filter the service
			serviceResult := gb.serviceFilter.ShouldIncludeService(service, criteria)
			if !serviceResult.Include || serviceResult.IsBrowserProvided {
				continue // Skip excluded or browser services
			}

			// Build service data
			serviceData, err := gb.buildServiceData(service, file, serviceResult, criteria, context)
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

// buildBrowserClients builds client interfaces for browser-provided services in this package.
// These services are implemented in JavaScript and called from WASM via clients.
func (gb *GoDataBuilder) buildBrowserClients(
	files []*protogen.File,
	criteria *filters.FilterCriteria,
	context *BuildContext,
) ([]ServiceData, error) {

	var clients []ServiceData

	for _, file := range files {
		for _, service := range file.Services {
			// Filter the service
			serviceResult := gb.serviceFilter.ShouldIncludeService(service, criteria)
			if !serviceResult.Include || !serviceResult.IsBrowserProvided {
				continue // Skip excluded or non-browser services
			}

			// Build client data for browser service
			serviceData, err := gb.buildServiceData(service, file, serviceResult, criteria, context)
			if err != nil {
				return nil, err
			}

			if serviceData != nil {
				clients = append(clients, *serviceData)
			}
		}
	}

	return clients, nil
}



// buildServiceData creates ServiceData from a protogen.Service and filter result.
func (gb *GoDataBuilder) buildServiceData(
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
		methodResult := gb.methodFilter.ShouldIncludeMethod(method, criteria)
		if !methodResult.Include {
			continue
		}

		methodData := gb.buildMethodData(method, serviceName, methodResult, context)
		methods = append(methods, methodData)
	}

	// Skip services with no methods
	if len(methods) == 0 {
		return nil, nil
	}

	// Register package import
	packagePath := string(file.GoImportPath)
	packageAlias := gb.pathCalc.GetGoPackageAlias(packagePath)
	context.AddImport(packagePath, packageAlias)

	// Determine Go type (Server vs Client for browser services)
	goType := packageAlias + "." + string(service.GoName) + "Server"
	if serviceResult.IsBrowserProvided {
		goType = packageAlias + "." + string(service.GoName) + "Client"
	}

	// Determine JavaScript name
	jsName := gb.nameConv.ToCamelCase(serviceName)
	if serviceResult.CustomName != "" {
		jsName = serviceResult.CustomName
	}

	return &ServiceData{
		Name:              serviceName,
		GoType:            goType,
		JSName:            jsName,
		PackagePath:       packagePath,
		PackageAlias:      packageAlias,
		IsBrowserProvided: serviceResult.IsBrowserProvided,
		CustomName:        serviceResult.CustomName,
		Comment:           strings.TrimSpace(string(service.Comments.Leading)),
		Methods:           methods,
	}, nil
}

// buildMethodData creates MethodData from a protogen.Method and filter result.
func (gb *GoDataBuilder) buildMethodData(
	method *protogen.Method,
	serviceName string,
	methodResult filters.MethodFilterResult,
	context *BuildContext,
) MethodData {

	methodName := string(method.Desc.Name())

	// Determine JavaScript name
	jsName := methodResult.CustomJSName
	if jsName == "" {
		jsName = gb.nameConv.ToCamelCase(methodName)
	}

	// Build Go function name for WASM wrapper
	goFuncName := gb.nameConv.ToGoFuncName(serviceName, methodName)

	// Ensure imports are registered for request/response types
	// Request and response types might be from the same package as the service or different
	requestPackagePath := string(method.Input.GoIdent.GoImportPath)
	responsePackagePath := string(method.Output.GoIdent.GoImportPath)

	// Get or create package aliases
	requestPackageAlias := ""
	responsePackageAlias := ""

	if requestPackagePath != "" {
		requestPackageAlias = gb.pathCalc.GetGoPackageAlias(requestPackagePath)
		context.AddImport(requestPackagePath, requestPackageAlias)
	}

	if responsePackagePath != "" {
		responsePackageAlias = gb.pathCalc.GetGoPackageAlias(responsePackagePath)
		context.AddImport(responsePackagePath, responsePackageAlias)
	}

	// Build fully qualified type names
	requestType := string(method.Input.GoIdent.GoName)
	responseType := string(method.Output.GoIdent.GoName)

	if requestPackageAlias != "" {
		requestType = requestPackageAlias + "." + requestType
	}

	if responsePackageAlias != "" {
		responseType = responsePackageAlias + "." + responseType
	}

	return MethodData{
		Name:              methodName,
		JSName:            jsName,
		GoFuncName:        goFuncName,
		ShouldGenerate:    true, // This method passed filtering, so it should be generated
		Comment:           strings.TrimSpace(string(method.Comments.Leading)),
		RequestType:       requestType,
		ResponseType:      responseType,
		RequestTSType:     string(method.Input.GoIdent.GoName),
		ResponseTSType:    string(method.Output.GoIdent.GoName),
		IsAsync:           methodResult.IsAsync,
		IsServerStreaming: methodResult.IsServerStreaming,
	}
}

// getModuleName determines the WASM module name from package and configuration.
func (gb *GoDataBuilder) getModuleName(packageName string, config *GenerationConfig) string {
	if config.ModuleName != "" {
		return config.ModuleName
	}
	return gb.nameConv.ToModuleName(packageName)
}

// getJSNamespace determines the JavaScript namespace from package and configuration.
func (gb *GoDataBuilder) getJSNamespace(packageName string, config *GenerationConfig) string {
	if config.JSNamespace != "" {
		return config.JSNamespace
	}
	return gb.nameConv.ToJSNamespace(packageName)
}

// getPrimarySourcePath returns the primary proto file path for the package.
// This is used in generated file headers to indicate the source.
func (gb *GoDataBuilder) getPrimarySourcePath(files []*protogen.File) string {
	if len(files) == 0 {
		return ""
	}
	return files[0].Desc.Path()
}

// findFileForService finds the protogen.File that contains the given service.
// This is needed for browser services that may come from different packages.
func (gb *GoDataBuilder) findFileForService(service *protogen.Service) *protogen.File {
	// The service descriptor should have a reference to its parent file
	// This is a simple implementation - in practice, we might need more sophisticated lookup

	// For now, we'll implement a basic search
	// In a real implementation, this might be passed as context or cached
	return nil // TODO: Implement proper file lookup
}
