// helper functions
package utils

import (
	"fmt"
	"math"
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

func Bool2str(b bool) string {
	if b {
		return "-1"
	}
	return "0"
}

// Str2bool convert a string to an float.
func Str2float(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(fmt.Errorf("reader: failed parsing number: %s", err))
	}
	return f
}

func _round(f float64) float64 {
	return math.Floor(f + .5)
}

func Round(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return _round(f*shift) / shift
}

// FromScale returns a converted input from min, max to 0.0...1.0.
func FromScale(input, min, max float64) float64 {
	return (input - math.Min(min, max)) / (math.Max(min, max) - math.Min(min, max))
}

// ToScale returns a converted input from 0...1 to min...max scale.
// If input is less than min then ToScale returns min.
// If input is greater than max then ToScale returns max
func ToScale(input, min, max float64) float64 {
	i := input*(math.Max(min, max)-math.Min(min, max)) + math.Min(min, max)
	if i < math.Min(min, max) {
		return math.Min(min, max)
	} else if i > math.Max(min, max) {
		return math.Max(min, max)
	} else {
		return i
	}
}
