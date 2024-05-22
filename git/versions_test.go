package git_test

import (
	"errors"
	"ollie/git"
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		version       string
		expectedMajor int
		expectedMinor int
		expectedPatch int
		expectedDep   int
		expectedTag   string
		expectedType  string
		expectedErr   error
	}{
		{
			version:       "v1.2.3",
			expectedMajor: 1,
			expectedMinor: 2,
			expectedPatch: 3,
			expectedDep:   0,
			expectedTag:   "",
			expectedType:  "SEMVER",
			expectedErr:   nil,
		},
		{
			version:       "v1.2.3-tag",
			expectedMajor: 1,
			expectedMinor: 2,
			expectedPatch: 3,
			expectedDep:   0,
			expectedTag:   "tag",
			expectedType:  "SEMVER",
			expectedErr:   nil,
		},
		{
			version:       "v2022.1.5.9-rc1",
			expectedMajor: 2022,
			expectedMinor: 1,
			expectedPatch: 5,
			expectedDep:   9,
			expectedTag:   "rc1",
			expectedType:  "DATEVER",
			expectedErr:   nil,
		},
		{
			version:       "v2024.05.10.1",
			expectedMajor: 2024,
			expectedMinor: 5,
			expectedPatch: 10,
			expectedDep:   1,
			expectedTag:   "",
			expectedType:  "DATEVER",
			expectedErr:   nil,
		},
		{
			version:       "v1.2",
			expectedMajor: 0,
			expectedMinor: 0,
			expectedPatch: 0,
			expectedDep:   0,
			expectedTag:   "",
			expectedType:  "",
			expectedErr:   errors.New("invalid version format"),
		},
	}

	for _, test := range tests {
		major, minor, patch, dep, tag, vType, err := git.ParseVersion(test.version)
		if major != test.expectedMajor || minor != test.expectedMinor || patch != test.expectedPatch ||
			dep != test.expectedDep || tag != test.expectedTag || vType != test.expectedType ||
			(err == nil && test.expectedErr != nil) || (err != nil && test.expectedErr == nil) {
			t.Errorf("ParseVersion(%s) = (%d, %d, %d, %d, %s, %s, %v), expected (%d, %d, %d, %d, %s, %s, %v)",
				test.version, major, minor, patch, dep, tag, vType, err,
				test.expectedMajor, test.expectedMinor, test.expectedPatch, test.expectedDep, test.expectedTag, test.expectedType, test.expectedErr)
		}
	}
}

func TestIncrementString_EmptyString(t *testing.T) {
	expected := "1"
	result := git.IncrementString("")
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestIncrementString_NoDigits(t *testing.T) {
	expected := "abc1"
	result := git.IncrementString("abc")
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestIncrementString_SingleDigit(t *testing.T) {
	expected := "2"
	result := git.IncrementString("1")
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestIncrementString_MultipleDigits(t *testing.T) {
	expected := "10"
	result := git.IncrementString("9")
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestIncrementString_LeadingZeros(t *testing.T) {
	expected := "002"
	result := git.IncrementString("001")
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestIncrementString_WordWithDigits(t *testing.T) {
	expected := "test2"
	result := git.IncrementString("test1")
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
