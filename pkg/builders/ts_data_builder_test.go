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
	"testing"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/core"
	"github.com/panyam/protoc-gen-go-wasmjs/pkg/filters"
)

func TestExtractPackageFromTypeName(t *testing.T) {
	// Create a TSDataBuilder instance
	analyzer := core.NewProtoAnalyzer()
	pathCalc := core.NewPathCalculator()
	nameConv := core.NewNameConverter()
	serviceFilter := filters.NewServiceFilter(analyzer)
	methodFilter := filters.NewMethodFilter(analyzer)
	msgCollector := filters.NewMessageCollector(analyzer)
	enumCollector := filters.NewEnumCollector(analyzer)

	builder := NewTSDataBuilder(
		analyzer,
		pathCalc,
		nameConv,
		serviceFilter,
		methodFilter,
		msgCollector,
		enumCollector,
	)

	tests := []struct {
		name               string
		fullyQualifiedName string
		expectedPackage    string
	}{
		{
			name:               "simple type",
			fullyQualifiedName: "utils.v1.HelperUtilType",
			expectedPackage:    "utils.v1",
		},
		{
			name:               "nested type",
			fullyQualifiedName: "utils.v1.ParentMessage.NestedType",
			expectedPackage:    "utils.v1",
		},
		{
			name:               "deeply nested type",
			fullyQualifiedName: "utils.v1.GrandParent.Parent.Child",
			expectedPackage:    "utils.v1",
		},
		{
			name:               "well-known type",
			fullyQualifiedName: "google.protobuf.Timestamp",
			expectedPackage:    "google.protobuf",
		},
		{
			name:               "single level package",
			fullyQualifiedName: "mypackage.MyType",
			expectedPackage:    "mypackage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := builder.extractPackageFromTypeName(tt.fullyQualifiedName)
			if result != tt.expectedPackage {
				t.Errorf("extractPackageFromTypeName(%s) = %s, want %s",
					tt.fullyQualifiedName, result, tt.expectedPackage)
			}
		})
	}
}

func TestExtractTypeNameFromFullyQualified(t *testing.T) {
	// Create a TSDataBuilder instance
	analyzer := core.NewProtoAnalyzer()
	pathCalc := core.NewPathCalculator()
	nameConv := core.NewNameConverter()
	serviceFilter := filters.NewServiceFilter(analyzer)
	methodFilter := filters.NewMethodFilter(analyzer)
	msgCollector := filters.NewMessageCollector(analyzer)
	enumCollector := filters.NewEnumCollector(analyzer)

	builder := NewTSDataBuilder(
		analyzer,
		pathCalc,
		nameConv,
		serviceFilter,
		methodFilter,
		msgCollector,
		enumCollector,
	)

	tests := []struct {
		name               string
		fullyQualifiedName string
		expectedTypeName   string
	}{
		{
			name:               "simple type",
			fullyQualifiedName: "utils.v1.HelperUtilType",
			expectedTypeName:   "HelperUtilType",
		},
		{
			name:               "nested type",
			fullyQualifiedName: "utils.v1.ParentMessage.NestedType",
			expectedTypeName:   "ParentMessage_NestedType", // Flattened for TypeScript imports
		},
		{
			name:               "deeply nested type",
			fullyQualifiedName: "utils.v1.GrandParent.Parent.Child",
			expectedTypeName:   "GrandParent_Parent_Child", // Flattened for TypeScript imports
		},
		{
			name:               "well-known type",
			fullyQualifiedName: "google.protobuf.Timestamp",
			expectedTypeName:   "Timestamp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := builder.extractTypeNameFromFullyQualified(tt.fullyQualifiedName)
			if result != tt.expectedTypeName {
				t.Errorf("extractTypeNameFromFullyQualified(%s) = %s, want %s",
					tt.fullyQualifiedName, result, tt.expectedTypeName)
			}
		})
	}
}

func TestCalculateCrossPackageImportPath(t *testing.T) {
	// Create a TSDataBuilder instance
	analyzer := core.NewProtoAnalyzer()
	pathCalc := core.NewPathCalculator()
	nameConv := core.NewNameConverter()
	serviceFilter := filters.NewServiceFilter(analyzer)
	methodFilter := filters.NewMethodFilter(analyzer)
	msgCollector := filters.NewMessageCollector(analyzer)
	enumCollector := filters.NewEnumCollector(analyzer)

	builder := NewTSDataBuilder(
		analyzer,
		pathCalc,
		nameConv,
		serviceFilter,
		methodFilter,
		msgCollector,
		enumCollector,
	)

	tests := []struct {
		name                string
		currentPackagePath  string
		targetPackageName   string
		expectedImportPath  string
	}{
		{
			name:               "same level packages",
			currentPackagePath: "presenter/v1",
			targetPackageName:  "utils.v1",
			expectedImportPath: "../../utils/v1/interfaces",
		},
		{
			name:               "parent to child",
			currentPackagePath: "common",
			targetPackageName:  "common.types.v1",
			expectedImportPath: "../common/types/v1/interfaces",
		},
		{
			name:               "different top-level packages",
			currentPackagePath: "services/auth/v1",
			targetPackageName:  "models.user.v1",
			expectedImportPath: "../../../models/user/v1/interfaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := builder.calculateCrossPackageImportPath(tt.currentPackagePath, tt.targetPackageName)
			if result != tt.expectedImportPath {
				t.Errorf("calculateCrossPackageImportPath(%s, %s) = %s, want %s",
					tt.currentPackagePath, tt.targetPackageName, result, tt.expectedImportPath)
			}
		})
	}
}
