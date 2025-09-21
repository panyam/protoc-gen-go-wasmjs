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

package renderers

import (
	"fmt"
	"log"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/builders"
)

// getTemplateFuncMap returns the template function map used by all generators.
// This provides common template helper functions for both Go and TypeScript templates.
func getTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		// String manipulation functions
		"title": func(s string) string {
			if len(s) == 0 {
				return s
			}
			return strings.ToUpper(s[:1]) + s[1:]
		},
		"lower": func(s string) string {
			return strings.ToLower(s)
		},
		"upper": func(s string) string {
			return strings.ToUpper(s)
		},

		// String replacement functions
		"replaceAll": func(old, new, s string) string {
			return strings.ReplaceAll(s, old, new)
		},
		"replace": func(old, new string, n int, s string) string {
			return strings.Replace(s, old, new, n)
		},

		// String testing functions
		"hasPrefix": func(prefix, s string) bool {
			return strings.HasPrefix(s, prefix)
		},
		"hasSuffix": func(suffix, s string) bool {
			return strings.HasSuffix(s, suffix)
		},
		"contains": func(substr, s string) bool {
			return strings.Contains(s, substr)
		},

		// String trimming functions
		"trim": func(s string) string {
			return strings.TrimSpace(s)
		},
		"trimPrefix": func(prefix, s string) string {
			return strings.TrimPrefix(s, prefix)
		},
		"trimSuffix": func(suffix, s string) string {
			return strings.TrimSuffix(s, suffix)
		},

		// Array/slice functions
		"join": func(sep string, elems []string) string {
			return strings.Join(elems, sep)
		},
		"split": func(sep, s string) []string {
			return strings.Split(s, sep)
		},

		// Conditional functions
		"eq": func(a, b any) bool {
			return a == b
		},
		"ne": func(a, b any) bool {
			return a != b
		},
		"and": func(a, b bool) bool {
			return a && b
		},
		"or": func(a, b bool) bool {
			return a || b
		},
		"not": func(a bool) bool {
			return !a
		},

		// Utility functions
		"default": func(defaultValue, value any) any {
			if value == nil || value == "" {
				return defaultValue
			}
			return value
		},
	}
}

// ExecuteTemplate executes a template with the given data and function map.
// This is a helper function that provides consistent template execution with error handling.
func ExecuteTemplate(templateName, templateContent string, data any, output interface{ Write([]byte) (int, error) }) error {
	tmpl, err := template.New(templateName).Funcs(getTemplateFuncMap()).Parse(templateContent)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	if err := tmpl.Execute(output, data); err != nil {
		log.Printf("ExecuteTEMPLATE: EXECUTION ERROR: %v", err)
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return nil
}

// ExecuteTemplateToFile executes a template directly to a GeneratedFile like the old generator does.
// This matches the old generator's pattern exactly.
func ExecuteTemplateToFile(templateName, templateContent string, data any, file *protogen.GeneratedFile) error {
	log.Printf("TEMPLATE: Parsing template %s (%d bytes)", templateName, len(templateContent))
	tmpl, err := template.New(templateName).Funcs(getTemplateFuncMap()).Parse(templateContent)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	log.Printf("TEMPLATE: Executing template %s to GeneratedFile %p", templateName, file)
	log.Printf("TEMPLATE: Data type: %T", data)
	if goData, ok := data.(*builders.GoTemplateData); ok {
		log.Printf("TEMPLATE: GoTemplateData - Package=%s, ModuleName=%s, Services=%d, JSNamespace=%s, APIStructure=%s",
			goData.PackageName, goData.ModuleName, len(goData.Services), goData.JSNamespace, goData.APIStructure)
	}
	if tsData, ok := data.(*builders.TSTemplateData); ok {
		log.Printf("TEMPLATE: TSTemplateData - Package=%s, ModuleName=%s, Services=%d, JSNamespace=%s, APIStructure=%s",
			tsData.PackageName, tsData.ModuleName, len(tsData.Services), tsData.JSNamespace, tsData.APIStructure)
	}

	if err := tmpl.Execute(file, data); err != nil {
		log.Printf("ExecuteTEMPLATEToFile: EXECUTION ERROR: %v", err)
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	log.Printf("TEMPLATE: Template %s executed successfully to file %p", templateName, file)
	return nil
}

// ValidateTemplate validates that a template can be parsed without executing it.
// This is useful for template validation during development and testing.
func ValidateTemplate(templateName, templateContent string) error {
	_, err := template.New(templateName).Funcs(getTemplateFuncMap()).Parse(templateContent)
	if err != nil {
		return fmt.Errorf("template %s validation failed: %w", templateName, err)
	}
	return nil
}
