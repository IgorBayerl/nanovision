package csharp_test

import (
	"testing"

	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/language/csharp"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	// Arrange
	formatter := csharp.NewCSharpProcessor()

	// Act
	name := formatter.Name()

	// Assert
	assert.Equal(t, "C#", name, "The formatter name should be 'C#'")
}

func TestDetect(t *testing.T) {
	testCases := []struct {
		name     string
		filePath string
		expected bool
	}{
		{
			name:     "CSharpFile_LowerCase_ShouldReturnTrue",
			filePath: "/path/to/file.cs",
			expected: true,
		},
		{
			name:     "CSharpFile_UpperCase_ShouldReturnTrue",
			filePath: "/path/to/file.CS",
			expected: true,
		},
		{
			name:     "FSharpFile_LowerCase_ShouldReturnTrue",
			filePath: "/path/to/another.fs",
			expected: true,
		},
		{
			name:     "GoFile_ShouldReturnFalse",
			filePath: "main.go",
			expected: false,
		},
		{
			name:     "FileWithNoExtension_ShouldReturnFalse",
			filePath: "somefile",
			expected: false,
		},
		{
			name:     "EmptyFilePath_ShouldReturnFalse",
			filePath: "",
			expected: false,
		},
		{
			name:     "FileWithOtherExtension_ShouldReturnFalse",
			filePath: "style.css",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			formatter := csharp.NewCSharpProcessor()

			// Act
			result := formatter.Detect(tc.filePath)

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetLogicalClassName(t *testing.T) {
	testCases := []struct {
		name         string
		rawClassName string
		expected     string
	}{
		{
			name:         "StandardClassName_ShouldReturnSame",
			rawClassName: "MyProject.Core.MyService",
			expected:     "MyProject.Core.MyService",
		},
		{
			name:         "NestedClassWithPlus_ShouldReturnParent",
			rawClassName: "MyProject.Core.MyService+NestedHelper",
			expected:     "MyProject.Core.MyService",
		},
		{
			name:         "NestedClassWithSlash_ShouldReturnParent",
			rawClassName: "MyProject.Core.MyService/NestedHelper",
			expected:     "MyProject.Core.MyService",
		},
		{
			name:         "AsyncStateMachine_ShouldReturnParent",
			rawClassName: "MyProject.Core.MyService+<MyAsyncMethod>d__10",
			expected:     "MyProject.Core.MyService",
		},
		{
			name:         "DeeplyNestedClass_ShouldReturnTopLevelParent",
			rawClassName: "MyProject.Core.MyService+Nested+DeeplyNested",
			expected:     "MyProject.Core.MyService",
		},
		{
			name:         "EmptyString_ShouldReturnEmpty",
			rawClassName: "",
			expected:     "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			formatter := csharp.NewCSharpProcessor()

			// Act
			result := formatter.GetLogicalClassName(tc.rawClassName)

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsCompilerGeneratedClass(t *testing.T) {
	testCases := []struct {
		name       string
		classInput model.Class
		expected   bool
	}{
		{
			name:       "LambdaCacheClass_WithPlus_ShouldReturnTrue",
			classInput: model.Class{Name: "MyProject.Services.MyService+<>c"},
			expected:   true,
		},
		{
			name:       "LambdaCacheClass_WithSlash_ShouldReturnTrue",
			classInput: model.Class{Name: "MyProject.Services.MyService/<>c"},
			expected:   true,
		},
		{
			name:       "AsyncStateMachineClass_ShouldReturnTrue",
			classInput: model.Class{Name: "MyProject.Services.MyService+<MyAsyncMethod>d__10"},
			expected:   true,
		},
		{
			name:       "StandardClass_ShouldReturnFalse",
			classInput: model.Class{Name: "MyProject.Services.MyService"},
			expected:   false,
		},
		{
			name:       "StandardNestedClass_ShouldReturnFalse",
			classInput: model.Class{Name: "MyProject.Services.MyService+MyNestedClass"},
			expected:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			formatter := csharp.NewCSharpProcessor()

			// Act
			result := formatter.IsCompilerGeneratedClass(&tc.classInput)

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFormatClassName(t *testing.T) {
	testCases := []struct {
		name       string
		classInput model.Class
		expected   string
	}{
		{
			name:       "NonGenericClass_ShouldReturnSame",
			classInput: model.Class{Name: "MyProject.MyService"},
			expected:   "MyProject.MyService",
		},
		{
			name:       "SingleGenericClass_ShouldFormatToBrackets",
			classInput: model.Class{Name: "System.Collections.Generic.List`1"},
			expected:   "System.Collections.Generic.List<T>",
		},
		{
			name:       "MultiGenericClass_ShouldFormatToBracketsWithNumbers",
			classInput: model.Class{Name: "System.Collections.Generic.Dictionary`2"},
			expected:   "System.Collections.Generic.Dictionary<T1, T2>",
		},
		{
			name:       "NestedClassWithSlash_ShouldConvertToDot",
			classInput: model.Class{Name: "MyProject.MyService/Nested"},
			expected:   "MyProject.MyService.Nested",
		},
		{
			name:       "NestedGenericClass_ShouldFormatCorrectly",
			classInput: model.Class{Name: "MyProject.MyService/MyGenericHelper`1"},
			expected:   "MyProject.MyService.MyGenericHelper<T>",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			formatter := csharp.NewCSharpProcessor()

			// Act
			result := formatter.FormatClassName(&tc.classInput)

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFormatMethodName(t *testing.T) {
	testCases := []struct {
		name         string
		methodInput  model.Method
		classInput   model.Class
		expectedName string
	}{
		{
			name:         "StandardMethod_ShouldReturnSame",
			methodInput:  model.Method{Name: "MyMethod", Signature: "()"},
			classInput:   model.Class{Name: "MyClass"},
			expectedName: "MyMethod()",
		},
		{
			name:         "AsyncMethod_ShouldReturnOriginalMethodName",
			methodInput:  model.Method{Name: "MoveNext", Signature: "()"},
			classInput:   model.Class{Name: "MyNamespace.MyService+<ProcessDataAsync>d__5"},
			expectedName: "ProcessDataAsync()",
		},
		{
			name:         "LocalFunction_ShouldReturnLocalFunctionName",
			methodInput:  model.Method{Name: "<Execute>g__ProcessItem|0_0", Signature: "()"},
			classInput:   model.Class{Name: "MyNamespace.MyService"},
			expectedName: "ProcessItem()",
		},
		{
			name:         "LocalFunction_AlternativeName_ShouldReturnLocalFunctionName",
			methodInput:  model.Method{Name: "MyLocalFunc|1_12", Signature: "(int)"},
			classInput:   model.Class{Name: "MyNamespace.MyService"},
			expectedName: "MyLocalFunc()", // Signature from model is not used, method formats to "()"
		},
		{
			name:         "NegativeTest_NonAsyncMoveNext_ShouldReturnSame",
			methodInput:  model.Method{Name: "MoveNext", Signature: "()"},
			classInput:   model.Class{Name: "MyNamespace.CustomIterator"},
			expectedName: "MoveNext()",
		},
		{
			name:         "PropertyGetter_ShouldReturnSame",
			methodInput:  model.Method{Name: "get_MyProperty", Signature: "()"},
			classInput:   model.Class{Name: "MyClass"},
			expectedName: "get_MyProperty()",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			formatter := csharp.NewCSharpProcessor()

			// Act
			result := formatter.FormatMethodName(&tc.methodInput, &tc.classInput)

			// Assert
			assert.Equal(t, tc.expectedName, result)
		})
	}
}

func TestCategorizeCodeElement(t *testing.T) {
	testCases := []struct {
		name         string
		methodInput  model.Method
		expectedType model.CodeElementType
	}{
		{
			name:         "PropertyGetter_ShouldReturnPropertyType",
			methodInput:  model.Method{DisplayName: "get_Name"},
			expectedType: model.PropertyElementType,
		},
		{
			name:         "PropertySetter_ShouldReturnPropertyType",
			methodInput:  model.Method{DisplayName: "set_Name"},
			expectedType: model.PropertyElementType,
		},
		{
			name:         "StandardMethod_ShouldReturnMethodType",
			methodInput:  model.Method{DisplayName: "DoWork()"},
			expectedType: model.MethodElementType,
		},
		{
			name:         "FormattedAsyncMethod_ShouldReturnMethodType",
			methodInput:  model.Method{DisplayName: "MyAsyncMethod()"}, // Input is the already formatted name
			expectedType: model.MethodElementType,
		},
		{
			name:         "MethodNameContainingGet_ButNotPrefix_ShouldReturnMethodType",
			methodInput:  model.Method{DisplayName: " toget_Name()"},
			expectedType: model.MethodElementType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			formatter := csharp.NewCSharpProcessor()

			// Act
			result := formatter.CategorizeCodeElement(&tc.methodInput)

			// Assert
			assert.Equal(t, tc.expectedType, result)
		})
	}
}
