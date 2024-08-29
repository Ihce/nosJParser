package deserializer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// func TestDeserializeRootMap(t *testing.T) {
// 	test := []struct {
// 		input       string
// 		shouldPanic bool
// 	}{
// 		{"(<>)", false},
// 		{"(<>", true},
// 		{"<>", true},
// 		{"(<", true},
// 		{"", true},
// 		{"()", true},
// 		{"<>)", true},
// 		{"(<)", true},
// 	}

// 	for _, tt := range test {
// 		t.Run(tt.input, func(t *testing.T) {
// 			defer func() {
// 				if r := recover(); r != nil {
// 					if !tt.shouldPanic {
// 						t.Errorf("Deserialize() panicked for input %q, but it shouldn't have", tt.input)
// 					}
// 				} else {
// 					if tt.shouldPanic {
// 						t.Errorf("Deserialize() did not panic for input %q, but it should have", tt.input)
// 					}

// 				}
// 			}()
// 			Deserialize(tt.input)
// 		})
// 	}
// }

func TestDeserializeValidData(t *testing.T) {
	dir := "../../testdata/valid/"
	// Discover all test cases in the docs/spec-testcases/valid directory
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("Failed to read test cases directory: %v", err)
	}

	var testCases []struct {
		inputFile  string
		outputFile string
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".input") {
			baseName := strings.TrimSuffix(file.Name(), ".input")
			testCases = append(testCases, struct {
				inputFile  string
				outputFile string
			}{
				inputFile:  filepath.Join(dir, baseName+".input"),
				outputFile: filepath.Join(dir, baseName+".output"),
			})
		}
	}

	for _, tc := range testCases {
		t.Run(tc.inputFile, func(t *testing.T) {
			// Read input from the input file
			inputData, err := os.ReadFile(tc.inputFile)
			if err != nil {
				t.Fatalf("Failed to read input file %s: %v", tc.inputFile, err)
			}

			// Read expected output from the output file
			expectedOutput, err := os.ReadFile(tc.outputFile)
			if err != nil {
				t.Fatalf("Failed to read output file %s: %v", tc.outputFile, err)
			}

			// Deserialize the input
			result, err := Deserialize(string(inputData))
			if err != nil {
				t.Fatalf("Deserialize failed for input file %s: %v", tc.inputFile, err)
			}

			// Compare the result with the expected output
			if result != string(expectedOutput) {
				t.Errorf("Expected %s, but got %s for input file %s", string(expectedOutput), result, tc.inputFile)
			}
		})
	}
}
