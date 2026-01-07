package types

import (
	"database/sql/driver"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestUUIDArray_Scan(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    UUIDArray
		expectError bool
		errorMsg    string
	}{
		{
			name:     "empty array from empty string",
			input:    "{}",
			expected: UUIDArray{},
		},
		{
			name:     "empty array from byte",
			input:    []byte("{}"),
			expected: UUIDArray{},
		},
		{
			name:  "single UUID without quotes",
			input: "{12345678-1234-1234-1234-123456789012}",
			expected: UUIDArray{
				uuid.MustParse("12345678-1234-1234-1234-123456789012"),
			},
		},
		{
			name:  "single UUID with quotes",
			input: `{"12345678-1234-1234-1234-123456789012"}`,
			expected: UUIDArray{
				uuid.MustParse("12345678-1234-1234-1234-123456789012"),
			},
		},
		{
			name:  "multiple UUIDs without quotes",
			input: "{12345678-1234-1234-1234-123456789012,87654321-4321-4321-4321-210987654321}",
			expected: UUIDArray{
				uuid.MustParse("12345678-1234-1234-1234-123456789012"),
				uuid.MustParse("87654321-4321-4321-4321-210987654321"),
			},
		},
		{
			name:  "multiple UUIDs with quotes",
			input: `{"12345678-1234-1234-1234-123456789012","87654321-4321-4321-4321-210987654321"}`,
			expected: UUIDArray{
				uuid.MustParse("12345678-1234-1234-1234-123456789012"),
				uuid.MustParse("87654321-4321-4321-4321-210987654321"),
			},
		},
		{
			name:  "multiple UUIDs with spaces",
			input: `{ "12345678-1234-1234-1234-123456789012" , "87654321-4321-4321-4321-210987654321" }`,
			expected: UUIDArray{
				uuid.MustParse("12345678-1234-1234-1234-123456789012"),
				uuid.MustParse("87654321-4321-4321-4321-210987654321"),
			},
		},
		{
			name:        "invalid UUID format",
			input:       "{invalid-uuid-format}",
			expectError: true,
			errorMsg:    "invalid UUID in Array",
		},
		{
			name:        "unsupported data type",
			input:       12345,
			expectError: true,
			errorMsg:    "unsupported data type",
		},
		{
			name:     "empty string with spaces",
			input:    "{  }",
			expected: UUIDArray{},
		},
		{
			name:  "string with empty elements",
			input: `{"12345678-1234-1234-1234-123456789012", ,"87654321-4321-4321-4321-210987654321"}`,
			expected: UUIDArray{
				uuid.MustParse("12345678-1234-1234-1234-123456789012"),
				uuid.MustParse("87654321-4321-4321-4321-210987654321"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var result UUIDArray

			err := result.Scan(tc.input)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tc.errorMsg != "" && !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("Expected error containing '%s', got '%v'", tc.errorMsg, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(result) != len(tc.expected) {
				t.Errorf("Expected length %d, got %d", len(tc.expected), len(result))
				return
			}

			for i, expectedUUID := range tc.expected {
				if result[i] != expectedUUID {
					t.Errorf("At index %d: expected %v, got %v", i, expectedUUID, result[i])
				}
			}
		})
	}
}

func TestUUIDArray_Value(t *testing.T) {
	tests := []struct {
		name     string
		input    UUIDArray
		expected string
	}{
		{
			name:     "empty array",
			input:    UUIDArray{},
			expected: "{}",
		},
		{
			name: "single UUID",
			input: UUIDArray{
				uuid.MustParse("12345678-1234-1234-1234-123456789012"),
			},
			expected: `{"12345678-1234-1234-1234-123456789012"}`,
		},
		{
			name: "multiple UUIDs",
			input: UUIDArray{
				uuid.MustParse("12345678-1234-1234-1234-123456789012"),
				uuid.MustParse("87654321-4321-4321-4321-210987654321"),
				uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			},
			expected: `{"12345678-1234-1234-1234-123456789012","87654321-4321-4321-4321-210987654321","11111111-1111-1111-1111-111111111111"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.input.Value()
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			strResult, ok := result.(string)
			if !ok {
				t.Errorf("Expected string result, got %T", result)
				return
			}

			if strResult != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, strResult)
			}
		})
	}
}

func TestUUIDArray_Value_RoundTrip(t *testing.T) {
	// Test round-trip: Scan -> Value -> Scan should preserve data
	originalUUIDs := []uuid.UUID{
		uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		uuid.MustParse("33333333-3333-3333-3333-333333333333"),
	}

	originalArray := UUIDArray(originalUUIDs)

	// Convert to database value
	dbValue, err := originalArray.Value()
	if err != nil {
		t.Fatalf("Failed to get Value: %v", err)
	}

	// Scan back from database value
	var scannedArray UUIDArray
	err = scannedArray.Scan(dbValue)
	if err != nil {
		t.Fatalf("Failed to Scan: %v", err)
	}

	// Compare
	if len(scannedArray) != len(originalArray) {
		t.Fatalf("Length mismatch: expected %d, got %d", len(originalArray), len(scannedArray))
	}

	for i, original := range originalArray {
		if scannedArray[i] != original {
			t.Errorf("Mismatch at index %d: expected %v, got %v", i, original, scannedArray[i])
		}
	}
}

func TestUUIDArray_GormDataType(t *testing.T) {
	var array UUIDArray
	result := array.GormDataType()

	expected := "uuid[]"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestUUIDArray_DriverValuerImplementation(t *testing.T) {
	// Test that UUIDArray implements driver.Valuer
	var _ driver.Valuer = UUIDArray{}

	// Test with driver.Value usage
	array := UUIDArray{
		uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
		uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
	}

	value, err := array.Value()
	if err != nil {
		t.Errorf("Value() failed: %v", err)
	}

	strValue, ok := value.(string)
	if !ok {
		t.Errorf("Value() did not return string")
	}

	expected := `{"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"}`
	if strValue != expected {
		t.Errorf("Value() returned wrong string: expected %s, got %s", expected, strValue)
	}
}

func TestUUIDArray_Scan_Nil(t *testing.T) {
	var array UUIDArray

	// Test scanning nil
	err := array.Scan(nil)
	if err == nil {
		t.Error("Expected error when scanning nil, got none")
	}
}

func TestUUIDArray_Scan_ByteArray(t *testing.T) {
	testCases := []struct {
		name   string
		input  []byte
		length int
	}{
		{
			name:   "empty byte array",
			input:  []byte("{}"),
			length: 0,
		},
		{
			name:   "single UUID byte array",
			input:  []byte(`{"12345678-1234-1234-1234-123456789012"}`),
			length: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var array UUIDArray
			err := array.Scan(tc.input)
			if err != nil {
				t.Errorf("Scan failed: %v", err)
			}

			if len(array) != tc.length {
				t.Errorf("Expected length %d, got %d", tc.length, len(array))
			}
		})
	}
}

func TestUUIDArray_GormDataType(t *testing.T) {
	var arr UUIDArray
	result := arr.GormDataType()

	expected := "uuid[]"
	if result != expected {
		t.Errorf("GormDataType() = %s, want %s", result, expected)
	}
}
