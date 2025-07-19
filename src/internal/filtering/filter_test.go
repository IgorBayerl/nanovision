package filtering

import (
	"testing"
)

// TestIsElementIncludedInReport tests the core filtering logic of the DefaultFilter.
func TestIsElementIncludedInReport(t *testing.T) {
	testCases := []struct {
		name               string
		filters            []string
		pathSeparator      bool
		elementName        string
		expectedIsIncluded bool
	}{
		// --- Basic Include/Exclude ---
		{
			name:               "SimpleInclude_MatchingElement_ReturnsTrue",
			filters:            []string{"+MyProject.Core"},
			elementName:        "MyProject.Core",
			expectedIsIncluded: true,
		},
		{
			name:               "SimpleInclude_NonMatchingElement_ReturnsFalse",
			filters:            []string{"+MyProject.Core"},
			elementName:        "MyProject.Data",
			expectedIsIncluded: false,
		},
		{
			name:               "SimpleExclude_MatchingElement_ReturnsFalse",
			filters:            []string{"-MyProject.Tests"},
			elementName:        "MyProject.Tests",
			expectedIsIncluded: false,
		},
		{
			name:               "SimpleExclude_NonMatchingElement_ReturnsTrue",
			filters:            []string{"-MyProject.Tests"},
			elementName:        "MyProject.Core",
			expectedIsIncluded: true,
		},

		// --- Wildcard Logic ---
		{
			name:               "WildcardInclude_MatchingElement_ReturnsTrue",
			filters:            []string{"+MyProject.*"},
			elementName:        "MyProject.Data",
			expectedIsIncluded: true,
		},
		{
			name:               "WildcardInclude_NonMatchingParent_ReturnsFalse",
			filters:            []string{"+MyProject.*"},
			elementName:        "AnotherProject.Core",
			expectedIsIncluded: false,
		},
		{
			name:               "WildcardExclude_MatchingElement_ReturnsFalse",
			filters:            []string{"-*Tests"},
			elementName:        "MyProject.Core.Tests",
			expectedIsIncluded: false,
		},

		// --- Precedence Rules ---
		{
			name:               "ExcludeWins_WhenElementMatchesBothIncludeAndExclude_ReturnsFalse",
			filters:            []string{"+MyProject.*", "-*Tests"},
			elementName:        "MyProject.Tests",
			expectedIsIncluded: false,
		},

		// --- Default Behavior ---
		{
			name:               "NoFilters_AnyElement_ReturnsTrue",
			filters:            []string{},
			elementName:        "Any.Element.Can.Be.Here",
			expectedIsIncluded: true,
		},
		{
			name:               "OnlyExcludeFilters_NonMatchingElement_ReturnsTrue",
			filters:            []string{"-Excluded.*"},
			elementName:        "MyProject.Core",
			expectedIsIncluded: true,
		},
		{
			name:               "OnlyExcludeFilters_MatchingElement_ReturnsFalse",
			filters:            []string{"-Excluded.*"},
			elementName:        "Excluded.Data",
			expectedIsIncluded: false,
		},

		// --- Case Insensitivity ---
		{
			name:               "CaseInsensitiveMatch_OnInclude_ReturnsTrue",
			filters:            []string{"+myproject.core"},
			elementName:        "MyProject.Core",
			expectedIsIncluded: true,
		},
		{
			name:               "CaseInsensitiveMatch_OnExclude_ReturnsFalse",
			filters:            []string{"-*tests"},
			elementName:        "MyProject.Core.Tests",
			expectedIsIncluded: false,
		},

		// --- Path Separator Logic ---
		{
			name:               "PathSeparator_WithUnixPath_ReturnsTrue",
			filters:            []string{"+*/Tests/*"},
			pathSeparator:      true,
			elementName:        "C:/Projects/MyProject/Tests/test.cs",
			expectedIsIncluded: true,
		},
		{
			name:               "PathSeparator_WithWindowsPath_ReturnsTrue",
			filters:            []string{"+*/Tests/*"},
			pathSeparator:      true,
			elementName:        "C:\\Projects\\MyProject\\Tests\\test.cs",
			expectedIsIncluded: true,
		},
		{
			name:               "PathSeparator_DisabledWithWindowsPath_ReturnsFalse",
			filters:            []string{"+*/Tests/*"},
			pathSeparator:      false, // Disabled
			elementName:        "C:\\Projects\\MyProject\\Tests\\test.cs",
			expectedIsIncluded: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			filter, err := NewDefaultFilter(tc.filters, tc.pathSeparator)
			if err != nil {
				t.Fatalf("Test setup failed. NewDefaultFilter returned an unexpected error: %v", err)
			}

			// Act
			isIncluded := filter.IsElementIncludedInReport(tc.elementName)

			// Assert
			if isIncluded != tc.expectedIsIncluded {
				t.Errorf("Expected IsElementIncludedInReport to be %v, but got %v", tc.expectedIsIncluded, isIncluded)
			}
		})
	}
}

// TestNewDefaultFilter tests the constructor for validation and state initialization.
func TestNewDefaultFilter(t *testing.T) {
	testCases := []struct {
		name              string
		filters           []string
		expectError       bool
		expectedHasCustom bool
	}{
		{
			name:              "ValidIncludeAndExclude_NoError_HasCustomIsTrue",
			filters:           []string{"+Inc", "-Exc"},
			expectError:       false,
			expectedHasCustom: true,
		},
		{
			name:              "NoFilters_NoError_HasCustomIsFalse",
			filters:           []string{},
			expectError:       false,
			expectedHasCustom: false,
		},
		{
			name:              "EmptyFilterStrings_AreIgnored_NoError",
			filters:           []string{"+Inc", "", "-Exc"},
			expectError:       false,
			expectedHasCustom: true,
		},
		{
			name:              "MalformedFilter_NoPrefix_ReturnsError",
			filters:           []string{"NoPrefix"},
			expectError:       true,
			expectedHasCustom: false,
		},
		{
			name:              "InvalidRegexInFilter_ReturnsError",
			filters:           []string{"+[UnclosedBracket"},
			expectError:       true,
			expectedHasCustom: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange & Act
			filter, err := NewDefaultFilter(tc.filters)

			// Assert
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error for filters %v, but got none", tc.filters)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for filters %v, but got: %v", tc.filters, err)
				}
				if filter == nil {
					t.Fatal("Expected filter to be non-nil on success")
				}
				if filter.HasCustomFilters() != tc.expectedHasCustom {
					t.Errorf("Expected HasCustomFilters() to be %v, but got %v", tc.expectedHasCustom, filter.HasCustomFilters())
				}
			}
		})
	}
}
