package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/phil-mansfield/rogue/error"
)

// Should contain only public fields and only those of type string.
//
// TODO: consider whether or not it would be exceptionally paintful to do
// otherwise. Since we're using reflection anyway, there's no reason to do this
// other than laziness. There's also no reason to lock it in to the Info
// type.
type Info struct {
	FramesPerSecond int

	FavoriteQuote string
	FavoriteNumber int
}

type varInfo struct {
	defaultValue interface{}
	convert      ConvertFunc
}

var (
	varInfos = map[string]varInfo{
		"FramesPerSecond": {20, IntRangeConvert(1, 1000)},
		"FavoriteQuote": {"What I cannot create, I do not understand.", NoConvert},
		"FavoriteNumber": {1729, IntConvert},
	}
)

// TODO: add Value error handling to this. (e.g. if the use gives a blank
// string).

// Parse returns an Info instance corresponding to the configuration variables
// set in fileName.
//
// Parse can return Configuration, Library, MissingFile, and Sanity errors.
// The Info instance is valid (although perhaps not what the user wanted) as
// long as a Sanity error is not returned.
func Parse(filePath string) (*Info, *error.Error) {
	info, err := Default()
	if err != nil {
		return nil, err
	}

	lines, err := readFile(filePath)
	if err != nil {
		return info, err
	}

	reflectedInfo := reflect.Indirect(reflect.ValueOf(info))

	for i, line := range lines {

		// Parse current line.

		field, value, desc := parseLine(line)
		if desc != "" {
			def, err := Default()
			if err != nil {
				return nil, err
			}

			return def, propogateParseError(i, filePath, desc)
		} else if field == "" {
			continue
		}

		// Place converted value in info.

		vInfo, ok := varInfos[field]
		if !ok {
			def, err := Default()
			if err != nil { 
				return nil, err
			}

			desc := fmt.Sprintf("Unknown variable '%s'.", field)
			return def, propogateParseError(i, filePath, desc)
		}

		if  val, ok, desc := vInfo.convert(value); !ok {
			def, err := Default()
			if err != nil {
				return nil, err
			}

			return def, propogateParseError(i, filePath, desc)
		} else if !setField(reflectedInfo, field, val) {
			desc := fmt.Sprintf("setField returned false for field %s.",
				field)
			return nil, error.New(error.Sanity, desc)
		}
	}

	return info, nil
}

// Default returns an Info instance containing the default values for all
// parameters.
//
// Default can return a Sanity error if there is an inconsistency in
// configuration variable bookkeeping. The returned Info instance will not be
// valid if this error is returned.
func Default() (*Info, *error.Error) {
	info := new(Info)
	v := reflect.Indirect(reflect.ValueOf(info))

	for key, vInfo := range varInfos {
		field := v.FieldByName(key)
		if !field.IsValid() {
			desc := fmt.Sprintf("Key '%s' in varInfos unmatched by field "+
				"in Info struct.", key)
			return nil, error.New(error.Sanity, desc)
		}
		field.Set(reflect.ValueOf(vInfo.defaultValue))
	}

	return info, nil
}

// readFile takes care of the boiler plate required to decompose a file into a
// slice ofthe lines inside of that file, and potentially return an error.
func readFile(filePath string) ([]string, *error.Error) {
	if bytes, err := ioutil.ReadFile(filePath); err != nil {
		if _, err := os.Stat(filePath); err == nil {
			return []string{}, error.New(error.Library, err.Error())
		}
		desc := fmt.Sprintf("Config file '%s' does not exist.", filePath)
		return []string{}, error.New(error.MissingFile, desc)
	} else {
		return strings.Split(string(bytes), "\n"), nil
	}
}

// This function is awful. Rewrite as more tests are done.
//
// An empty string in the desc return value indicates no errors.
func parseLine(line string) (field, value string, desc string) {
	line = strings.Trim(line, " \t")
	if line == "" {
		return "", "", ""
	}

	subs := strings.SplitN(line, "=", 2)
	if len(subs) < 2 {
		return "", "", "No variable assignment in non-empty line."
	}

	return strings.Trim(subs[0], " \t"), strings.Trim(subs[1], " \t"), ""
}

// propogateParseError takes care of the boiler-plate code required to create
// an Error instance from a description of a parsing error.
//
// This will almost always return Configuration errors, but if something goes
// horribly, horribly wrong, it could return a Library error.
func propogateParseError(lineNum int, filePath, desc string) *error.Error {
	if absPath, pathErr := filepath.Abs(filePath); pathErr != nil {
		return error.New(error.Library, pathErr.Error())
	} else {
		fullDesc := fmt.Sprintf("%s On line %d of '%s'.",
			desc, lineNum+1, absPath)
		return error.New(error.Configuration, fullDesc)
	}
}

// setField attempts to set the field with name fieldName to value. The function
// returns true on success and false if v had no such field.
//
// TODO: This function can panic if things go wrong due to design choices in
// the standard library. Maybe fix that.
func setField(v reflect.Value, fieldName string, value interface{}) bool {
	if field := v.FieldByName(fieldName); !field.IsValid() {
		return false
	} else if !field.CanSet() {
		return false
	} else {
		field.Set(reflect.ValueOf(value))
		return true
	}
}
