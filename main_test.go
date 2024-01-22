package main

import (
	"os"
	"reflect"
	"testing"
)

func TestCalculateMinimumCopies(t *testing.T) {
	tests := []struct {
		name               string
		inputData          []Application
		inputApplicationId string
		expectedResult     int
	}{
		{
			name: "Scenario 1: One laptop and one desktop",
			inputData: []Application{
				{ComputerId: "1", UserId: "1", ApplicationId: "374", ComputerType: "LAPTOP", Comment: "Exported from System A"},
				{ComputerId: "2", UserId: "1", ApplicationId: "374", ComputerType: "DESKTOP", Comment: "Exported from System A"},
			},
			inputApplicationId: "374",
			expectedResult:     1,
		},
		{
			name: "Scenario 2: One laptop and three desktops",
			inputData: []Application{
				{ComputerId: "1", UserId: "1", ApplicationId: "374", ComputerType: "LAPTOP", Comment: "Exported from System A"},
				{ComputerId: "2", UserId: "1", ApplicationId: "374", ComputerType: "DESKTOP", Comment: "Exported from System A"},
				{ComputerId: "3", UserId: "2", ApplicationId: "374", ComputerType: "DESKTOP", Comment: "Exported from System A"},
				{ComputerId: "4", UserId: "2", ApplicationId: "374", ComputerType: "DESKTOP", Comment: "Exported from System A"},
			},
			inputApplicationId: "374",
			expectedResult:     3,
		},
		{
			name: "Scenario 3: Duplicate record",
			inputData: []Application{
				{ComputerId: "1", UserId: "1", ApplicationId: "374", ComputerType: "LAPTOP", Comment: "Exported from System A"},
				{ComputerId: "2", UserId: "2", ApplicationId: "374", ComputerType: "LAPTOP", Comment: "Exported from System A"},
				{ComputerId: "2", UserId: "2", ApplicationId: "374", ComputerType: "desktop", Comment: "Exported from System B"},
			},
			inputApplicationId: "374",
			expectedResult:     2,
		},
		{
			name: "Scenario 4",
			inputData: []Application{
				{ComputerId: "1", UserId: "1", ApplicationId: "374", ComputerType: "DESKTOP", Comment: "Exported from System A"},
				{ComputerId: "2", UserId: "1", ApplicationId: "374", ComputerType: "DESKTOP", Comment: "Exported from System A"},
				{ComputerId: "3", UserId: "2", ApplicationId: "374", ComputerType: "LAPTOP", Comment: "Exported from System A"},
			},
			inputApplicationId: "374",
			expectedResult:     3,
		},
		{
			name: "Scenario 5",
			inputData: []Application{
				{ComputerId: "1", UserId: "1", ApplicationId: "374", ComputerType: "LAPTOP", Comment: "Exported from System A"},
				{ComputerId: "2", UserId: "1", ApplicationId: "374", ComputerType: "LAPTOP", Comment: "Exported from System A"},
			},
			inputApplicationId: "374",
			expectedResult:     2,
		},
		{
			name: "Scenario 6",
			inputData: []Application{
				{ComputerId: "1", UserId: "1", ApplicationId: "374", ComputerType: "DESKTOP", Comment: "Exported from System A"},
				{ComputerId: "2", UserId: "1", ApplicationId: "374", ComputerType: "DESKTOP", Comment: "Exported from System A"},
			},
			inputApplicationId: "374",
			expectedResult:     2,
		},
		{
			name: "Scenario 6",
			inputData: []Application{
				{ComputerId: "1", UserId: "1", ApplicationId: "374", ComputerType: "LAPTOP", Comment: "Exported from System A"},
				{ComputerId: "2", UserId: "1", ApplicationId: "374", ComputerType: "LAPTOP", Comment: "Exported from System A"},
				{ComputerId: "3", UserId: "1", ApplicationId: "374", ComputerType: "LAPTOP", Comment: "Exported from System A"},
				{ComputerId: "4", UserId: "1", ApplicationId: "374", ComputerType: "DESKTOP", Comment: "Exported from System A"},
			},
			inputApplicationId: "374",
			expectedResult:     3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := calculateCopiesRequired(test.inputData, test.inputApplicationId)
			if result != test.expectedResult {
				t.Errorf("Expected %d, but got %d", test.expectedResult, result)
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	// Prepare a sample CSV file content
	csvContent := "ComputerId,UserId,ApplicationId,ComputerType,Comment\n1,1,12,LAPTOP,comment1\n2,1,12,DESKTOP,comment2\n"

	// Create a temporary file for testing
	tempFile := createTempFile(t, csvContent)
	defer tempFile.Close()

	// Call the function being tested
	data, err := readFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	// Define the expected result
	expectedData := []Application{
		{ComputerId: "1", UserId: "1", ApplicationId: "12", ComputerType: "LAPTOP", Comment: "comment1"},
		{ComputerId: "2", UserId: "1", ApplicationId: "12", ComputerType: "DESKTOP", Comment: "comment2"},
	}

	// Compare the actual result with the expected result
	if !reflect.DeepEqual(data, expectedData) {
		t.Fatalf("Expected %v, got %v", expectedData, data)
	}
}

func TestBuildAppCountMap(t *testing.T) {
	// Prepare test data
	applications := []Application{
		{ComputerId: "1", UserId: "1", ApplicationId: "11", ComputerType: "LAPTOP"},
		{ComputerId: "2", UserId: "1", ApplicationId: "11", ComputerType: "DESKTOP"},
		{ComputerId: "3", UserId: "2", ApplicationId: "11", ComputerType: "DESKTOP"},
	}

	// Call the function being tested
	result := buildAppCountMap(applications, "11")

	// Define the expected result
	expectedResult := map[string]Data{
		"1-11": {LaptopCount: 1, DesktopCount: 1},
		"2-11": {LaptopCount: 0, DesktopCount: 1},
	}

	// Compare the actual result with the expected result
	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("Expected %v, got %v", expectedResult, result)
	}
}

func createTempFile(t *testing.T, content string) *os.File {
	t.Helper()

	tempFile, err := os.CreateTemp("", "test*.csv")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}

	_, err = tempFile.WriteString(content)
	if err != nil {
		t.Fatalf("Error writing to temporary file: %v", err)
	}

	return tempFile
}
