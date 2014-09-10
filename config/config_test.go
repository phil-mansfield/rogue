package config

import (
	"reflect"
	"testing"

	"github.com/phil-mansfield/rogue/error"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		in           string
		field, value string
		isValid      bool
	}{
		{"A=B", "A", "B", true},
		{" A = B ", "A", "B", true},
		{"A =	B = C", "A", "B = C", true},
		{"", "", "", true},
		{" 		 ", "", "", true},
		{"What I cannot create, I do not understand.", "", "", false},
	}

	for i, test := range tests {
		field, value, desc := parseLine(test.in)
		if desc != "" && test.isValid {
			t.Errorf("Test %d: Expected valid output for parseLine('%s'), "+
				"but got desc = '%s'", i, test.in, desc)
		} else if desc == "" && !test.isValid {
			t.Errorf("Test %d: Expected invalid output for parseLine('%s'), "+
				"but got field = '%s' and value = '%s'.",
				i, test.in, field, value)
		} else if field != test.field || value != test.value {
			t.Errorf("Test %d: Expected field = '%s' and value = '%s' for "+
				"parseLine('%s'), but got field = '%s' and value = '%s'.",
				i, test.field, test.value, test.in, field, value)
		}
	}
}

// This is not really the best, since it relies on the test being run in the
// directory where it's source code is.

func TestReadFile(t *testing.T) {
	tests := []struct {
		in      string
		lines   []string
		isValid bool
		code    error.ErrorCode
	}{
		{"test_config_files/no_assignment.txt",
			[]string{"What I cannot create, I do not understand.", ""}, true,0},
		{"test_config_files/does_not_exist.txt",
			[]string{}, false, error.MissingFile},
		{"test_config_files/double_assignment.txt",
			[]string{"A = B", "", "C = D", ""}, true, 0},

	}

	for i, test := range tests {
		lines, err := readFile(test.in)

		if !test.isValid && err == nil {
			t.Errorf("Test %d: Expected %s for config file '%s', but got " +
				"lines %v.", i, test.code.String(), test.in, lines)
		} else if !test.isValid && (test.code != err.Code) {
			t.Errorf("Test %d: Expected %s for config file '%s', but got " +
				"%s. Verbose error is:\n%s", i, test.code.String(),
				test.in, err.Code, err.VerboseError())
		} else if test.isValid && err != nil {
			t.Errorf("Test %d: Expected lines %v for config file '%s', but " +
				"got %s. Verbose error is:\n%s",
				i, test.lines, test.in, err.Code.String(), err.VerboseError()) 
		} else if test.isValid && !stringSliceEqual(test.lines, lines) {
			t.Errorf("Test %d: Expected lines to be %v for config file " + 
				"'%s', but got %v.", i, test.lines, test.in, lines)
		}
	}
}

type setFieldTestStruct struct {
	Ned, Ed string
}

func TestSetField(t *testing.T) {
	tests := [] struct {
		inField, inValue string
		isValid bool
	} {
		{"Ned", "Meow", true},
		{"Rabina", "Meow", false},
	}

	testStruct :=  &setFieldTestStruct{"meow", "meow"}
	val := reflect.Indirect(reflect.ValueOf(testStruct))

	for i, test := range tests {
		if test.isValid != setField(val, test.inField, test.inValue) {
			t.Errorf("Test %d: Expected testStruct.%s = '%s' to be a %s " +
				" operation.", i, test.inField, test.inValue,
				boolToValid(test.isValid))
		} else if test.isValid {
			fieldStr := val.FieldByName(test.inField).String()
			if fieldStr != test.inValue {
				t.Errorf("Test %d: Expected struct.%s to be set to '%s', " + 
					"but found '%s'", i , test.inField, test.inValue, fieldStr)
			}
		}
	}
}

func TestParse(t *testing.T) {
	var parseTests = [] struct {
		filePath string
		favoriteQuote string
		favoriteNumber int
		isValid bool
		code error.ErrorCode
	} {
		{"test_config_files/valid_config.txt",
		"Do not destroy what you cannot create.", 3, true, 0},
		{"test_config_files/invalid_var_config.txt",
		"What I cannot create, I do not understand.", 1729,
		false, error.Configuration},
		{"test_config_files/does_not_exist.txt",
		"What I cannot create, I do not understand.", 1729,
		false, error.MissingFile},
		{"test_config_files/partial_config.txt",
		"Do not destroy what you cannot create.", 1729, true, 0},
	}

	for i, test := range parseTests {
		info, err := Parse(test.filePath)

		if test.isValid && err != nil {
			t.Errorf("Test %d: Expected '%s' to be valid, but gave %s. " +
				"Verbose error:\n%s", i, test.filePath, err.Code.String(),
				err.VerboseError())
		} else if !test.isValid && err == nil {
			t.Errorf("Test %d: Expected '%s' to give %s, but gave the struct" +
				"%v instead.", i, test.filePath, test.code.String(), info)
		} else if info.FavoriteQuote != test.favoriteQuote {
			t.Errorf("Test %d: Expected '%s' to give FavoriteQuote = '%s'," +
				"but got '%s' instead.", i, test.filePath, test.favoriteQuote,
				info.FavoriteQuote)
		} else if info.FavoriteNumber != test.favoriteNumber {
			t.Errorf("Test %d: Expected '%s' to give FavoriteNumber = '%s'," +
				"but got '%s' instead.", i, test.filePath, test.favoriteNumber,
				info.FavoriteNumber)
		}
	}
}

func boolToValid(b bool) string {
	if b {
		return "valid"
	}
	return "invalid"
}

func stringSliceEqual(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}