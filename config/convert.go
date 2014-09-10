package config

import (
	"fmt"
	"strconv"
)

// Functions of type CheckFunc takes in a string and converts it to a value.
// ok indicates whether the conversion was successful. If ok is true, val hold
// the converted value. If ok is false, err will hold a string indicating what
// prevented the conversion from occuring.
type ConvertFunc func(string) (val interface{}, ok bool, desc string)

func NoConvert(str string) (val interface{}, ok bool, desc string) {
	return str, true, ""
}

func IntConvert(str string) (val interface{}, ok bool, desc string) {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, false, err.Error()
	}
	return num, true, ""
}

func IntRangeConvert(low, high int) ConvertFunc {
	return func(str string) (val interface{}, ok bool, desc string) {
		num, err := strconv.Atoi(str)
		if err != nil {
			return 0, false, err.Error()
		}

		if num < low || num > high {
			desc := fmt.Sprintf("Value %d is outside of range [%d, %d]",
				val, low, high)
			return 0, false, desc
		}

		return num, true, ""
	}
}