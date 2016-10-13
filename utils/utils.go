// helper functions
package utils

import (
	"fmt"
	"strconv"
)

// Str2int convert a string to an integer.
func Str2int(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Errorf("reader: failed parsing number: %s", err))
	}
	return i
}

// Str2bool convert a string to a boolean.
func Str2bool(s string) bool {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Errorf("reader: failed parsing boolean: %s", err))
	}
	return i != 0
}

// Str2bool convert a string to an float.
func Str2float(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(fmt.Errorf("reader: failed parsing number: %s", err))
	}
	return f
}
