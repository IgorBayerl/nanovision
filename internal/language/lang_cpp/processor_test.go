package lang_cpp_test

import (
	"strings"
	"testing"

	"github.com/IgorBayerl/AdlerCov/internal/language/lang_cpp"
	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCppProcessor_Detect(t *testing.T) {
	p := lang_cpp.NewCppProcessor()

	testCases := []struct {
		path     string
		expected bool
	}{
		{"file.c", true},
		{"file.cpp", true},
		{"file.h", true},
		{"file.hpp", true},
		{"file.cxx", true},
		{"file.hxx", true},
		{"file.cc", true},
		{"file.CPP", true},
		{"file.obj", false},
		{"", false},
		{"file.txt", false},
		{"file.go", false},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, p.Detect(tc.path), "Detection failed for path: %s", tc.path)
	}
}

func TestCppProcessor_AnalyzeFile(t *testing.T) {
	testCases := []struct {
		name            string
		sourceCode      string
		expectedMetrics []model.MethodMetrics
	}{
		{
			name: "GoldenPath_ClassMethodAndStandaloneFunction",
			sourceCode: `
#include <iostream>

void standalone_function(int x) {
    std::cout << x;
}

class MyClass {
public:
    int MyMethod(char* param) const {
        // Implementation
    }
};`,
			expectedMetrics: []model.MethodMetrics{
				{Name: "standalone_function", StartLine: 4, EndLine: 6},
				{Name: "MyMethod", StartLine: 10, EndLine: 12},
			},
		},
		{
			name: "ScopedMethodDefinition_OutsideClass",
			sourceCode: `
class MyService {
    MyService();
    void process_data();
};

MyService::MyService() { }

void MyService::process_data()
{
    // logic
}`,
			expectedMetrics: []model.MethodMetrics{
				{Name: "MyService::MyService", StartLine: 7, EndLine: 7},
				{Name: "MyService::process_data", StartLine: 9, EndLine: 12},
			},
		},
		{
			name: "Resilience_BracesInCommentsAndStrings",
			sourceCode: `
// Comment {
int calculate(int val) {
    const char* str = "A string with { and }";
    /*
     * Block comment with }
     */
    return val * 2;
}`,
			expectedMetrics: []model.MethodMetrics{
				{Name: "calculate", StartLine: 3, EndLine: 9},
			},
		},
		{
			name: "PointerAndReferenceInSignature",
			sourceCode: `
#include <vector>
std::vector<int>& get_ids(const Context* ctx) {
    //
}`,
			expectedMetrics: []model.MethodMetrics{
				{Name: "get_ids", StartLine: 3, EndLine: 5},
			},
		},
		{
			name: "Destructor",
			sourceCode: `
class MyClass {
    ~MyClass() {
        // cleanup
    }
};`,
			expectedMetrics: []model.MethodMetrics{
				{Name: "~MyClass", StartLine: 3, EndLine: 5},
			},
		},
		{
			name: "NoFunctions_ShouldReturnEmpty",
			sourceCode: `
#define MY_MACRO 10
struct MyPod { int x; };`,
			expectedMetrics: []model.MethodMetrics{},
		},
		{
			name: "Constructor_Various",
			sourceCode: `
class MyClass {
    MyClass() {
        // default constructor
    }
    
    explicit MyClass(int x) {
        // parameterized constructor
    }
};

MyClass::MyClass(const MyClass& other) {
    // copy constructor
}`,
			expectedMetrics: []model.MethodMetrics{
				{Name: "MyClass", StartLine: 3, EndLine: 5},
				{Name: "MyClass", StartLine: 7, EndLine: 9},
				{Name: "MyClass::MyClass", StartLine: 12, EndLine: 14},
			},
		},
		{
			name: "TemplateFunction",
			sourceCode: `
template<typename T>
T max_value(T a, T b) {
    return a > b ? a : b;
}

template<class T>
void MyClass<T>::process() {
    // template method
}`,
			expectedMetrics: []model.MethodMetrics{
				{Name: "max_value", StartLine: 3, EndLine: 5},
				{Name: "MyClass<T>::process", StartLine: 8, EndLine: 10},
			},
		},
		{
			name: "StaticAndInlineFunctions",
			sourceCode: `
static int helper_function() {
    return 42;
}

inline bool is_valid(const char* str) {
    return str != nullptr;
}

class Utils {
    static void utility_method() {
        // static method
    }
};`,
			expectedMetrics: []model.MethodMetrics{
				{Name: "helper_function", StartLine: 2, EndLine: 4},
				{Name: "is_valid", StartLine: 6, EndLine: 8},
				{Name: "utility_method", StartLine: 11, EndLine: 13},
			},
		},
		{
			name: "ConstAndNoexceptMethods",
			sourceCode: `
class MyClass {
    int getValue() const {
        return value;
    }
    
    void process() noexcept {
        // won't throw
    }
    
    bool isValid() const noexcept override {
        return true;
    }
};`,
			expectedMetrics: []model.MethodMetrics{
				{Name: "getValue", StartLine: 3, EndLine: 5},
				{Name: "process", StartLine: 7, EndLine: 9},
				{Name: "isValid", StartLine: 11, EndLine: 13},
			},
		},
		{
			name: "NamespacedFunctions",
			sourceCode: `
namespace Utils {
    void helper() {
        // helper function
    }
}

void Utils::another_helper() {
    // defined outside namespace
}

namespace Math::Advanced {
    double calculate(double x) {
        return x * 2.0;
    }
}`,
			expectedMetrics: []model.MethodMetrics{
				{Name: "helper", StartLine: 3, EndLine: 5},
				{Name: "Utils::another_helper", StartLine: 8, EndLine: 10},
				{Name: "calculate", StartLine: 13, EndLine: 15},
			},
		},
		{
			name: "VirtualAndOverrideMethods",
			sourceCode: `
class Base {
    virtual void process() {
        // base implementation
    }
    
    virtual ~Base() = default;
};

class Derived : public Base {
    void process() override {
        // derived implementation
    }
};`,
			expectedMetrics: []model.MethodMetrics{
				{Name: "process", StartLine: 3, EndLine: 5},
				{Name: "process", StartLine: 11, EndLine: 13},
			},
		},
		{
			name: "FunctionDeclarations_ShouldBeIgnored",
			sourceCode: `
class MyClass {
    void method();  // declaration only
    int getValue() const;  // declaration only
};

void external_function(int x);  // declaration only

void implemented_function() {
    // this should be detected
}`,
			expectedMetrics: []model.MethodMetrics{
				// Corrected StartLine from 8->9 and EndLine from 10->11
				{Name: "implemented_function", StartLine: 9, EndLine: 11},
			},
		},
		{
			name: "ComplexReturnTypes",
			sourceCode: `
std::vector<std::string> get_names() {
    return {};
}

const MyClass& get_reference() {
    static MyClass instance;
    return instance;
}

auto get_lambda() -> std::function<int(int)> {
    return [](int x) { return x * 2; };
}`,
			expectedMetrics: []model.MethodMetrics{
				{Name: "get_names", StartLine: 2, EndLine: 4},
				{Name: "get_reference", StartLine: 6, EndLine: 9},
				{Name: "get_lambda", StartLine: 11, EndLine: 13},
			},
		},
		{
			name: "OperatorOverloads",
			sourceCode: `
class MyClass {
    bool operator==(const MyClass& other) const {
        return true;
    }
    
    MyClass& operator++() {
        return *this;
    }
};

std::ostream& operator<<(std::ostream& os, const MyClass& obj) {
    return os;
}`,
			expectedMetrics: []model.MethodMetrics{
				{Name: "operator==", StartLine: 3, EndLine: 5},
				{Name: "operator++", StartLine: 7, EndLine: 9},
				{Name: "operator<<", StartLine: 12, EndLine: 14},
			},
		},
		{
			name: "MultilineFunctionSignatures",
			sourceCode: `
void complex_function(
    int param1,
    const std::string& param2,
    bool param3
) {
    // implementation
}

class MyClass {
    void another_complex_method(
        double x,
        double y
    ) const override
    {
        // implementation
    }
};`,
			expectedMetrics: []model.MethodMetrics{
				{Name: "complex_function", StartLine: 2, EndLine: 8},
				{Name: "another_complex_method", StartLine: 11, EndLine: 17},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			p := lang_cpp.NewCppProcessor()
			sourceLines := strings.Split(tc.sourceCode, "\n")

			// Act
			methods, err := p.AnalyzeFile("test.cpp", sourceLines)

			// Assert
			require.NoError(t, err)
			require.NotNil(t, methods)
			assert.ElementsMatch(t, tc.expectedMetrics, methods,
				"The discovered methods did not match the expected metrics.\nActual: %+v\nExpected: %+v",
				methods, tc.expectedMetrics)
		})
	}
}

func TestCppProcessor_Name(t *testing.T) {
	p := lang_cpp.NewCppProcessor()
	assert.Equal(t, "C/C++", p.Name())
}

// Helper test to verify individual regex patterns
func TestIndividualRegexPatterns(t *testing.T) {
	p := lang_cpp.NewCppProcessor()

	testCases := []struct {
		name        string
		sourceCode  string
		expectCount int
		expectNames []string
	}{
		{
			name: "OnlyStandaloneFunctions",
			sourceCode: `
void func1() { }
int func2(int x) { return x; }
bool func3(const char* str) { return true; }`,
			expectCount: 3,
			expectNames: []string{"func1", "func2", "func3"},
		},
		{
			name: "OnlyScopedMethods",
			sourceCode: `
void MyClass::method1() { }
int MyClass::method2(int x) { return x; }
MyClass::MyClass() { }`,
			expectCount: 3,
			expectNames: []string{"MyClass::method1", "MyClass::method2", "MyClass::MyClass"},
		},
		{
			name: "OnlyConstructors",
			sourceCode: `
class Test {
    Test() { }
    explicit Test(int x) { }
};`,
			expectCount: 2,
			expectNames: []string{"Test", "Test"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sourceLines := strings.Split(tc.sourceCode, "\n")
			methods, err := p.AnalyzeFile("test.cpp", sourceLines)

			require.NoError(t, err)
			assert.Len(t, methods, tc.expectCount)

			actualNames := make([]string, len(methods))
			for i, method := range methods {
				actualNames[i] = method.Name
			}
			assert.ElementsMatch(t, tc.expectNames, actualNames)
		})
	}
}
