// Package utils provide some helpers functions for other packages
package utils

import (
	"fmt"
	"math"
	"strconv"
	// "github.com/flopp/go-findfont"
	// "github.com/golang/freetype/truetype"
	// "golang.org/x/image/font"
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

// Obox2bool convert a string to a boolean (Opaquebox ass format).
func Obox2bool(s string) bool {
	return s == "3"
}

// Bool2Obox convert a boolen to a int (Opaquebox ass format).
func Bool2Obox(b bool) int {
	if !b {
		return 0
	}
	return 3
}

// Bool2str convert a bool to an string (ass format).
func Bool2str(b bool) string {
	if b {
		return "-1"
	}
	return "0"
}

// Str2float convert a string to an float.
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

// Round returns a rounded input with  p places
func Round(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return _round(f*shift) / shift
}

// // FromScale returns a converted input from min, max to 0.0...1.0.
// func FromScale(input, min, max float64) float64 {
// 	return (input - math.Min(min, max)) / (math.Max(min, max) - math.Min(min, max))
// }

// // ToScale returns a converted input from 0...1 to min...max scale.
// // If input is less than min then ToScale returns min.
// // If input is greater than max then ToScale returns max
// func ToScale(input, min, max float64) float64 {
// 	i := input*(math.Max(min, max)-math.Min(min, max)) + math.Min(min, max)
// 	if i < math.Min(min, max) {
// 		return math.Min(min, max)
// 	} else if i > math.Max(min, max) {
// 		return math.Max(min, max)
// 	} else {
// 		return i
// 	}
// }

// AppendStrUnique append to a string slice if the value doesn't already exist
func AppendStrUnique(slice []string, s string) []string {
	for _, ele := range slice {
		if ele == s {
			return slice
		}
	}
	return append(slice, s)
}

// // /usr/share/fonts
// // ~/.fonts
// // MeasureString returns the rendered width and height of the specified text
// // given the current font face.
// func MeasureString(ff font.Face, s string) (w float64) {
// 	d := &font.Drawer{Face: ff}
// 	a := d.MeasureString(s)
// 	// return float64(a >> 6), dc.fontHeight
// 	return float64(a >> 6)
// }

// func FontLoad(fontName string, fontSize float64) (font.Face, error) {

// 	fontPath, err := findfont.Find(fontName)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fontBytes, err := ioutil.ReadFile(fontPath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fontFace, err := truetype.Parse(fontBytes)
// 	if err != nil {
// 		return nil, err
// 	}
// 	face := truetype.NewFace(fontFace, &truetype.Options{
// 		Size:    fontSize,
// 		DPI:     72,
// 		Hinting: font.HintingNone,
// 		// Hinting: font.HintingFull,
// 	})
// 	return face, nil
// }

// type CIMap struct {
// 	m map[string]string
// }

// func NewCIMap() CIMap {
// 	return CIMap{m: make(map[string]int)}
// }

// func (m CIMap) Set(k string, v int) {
// 	m.m[strings.ToLower(k)] = v
// }

// func (m CIMap) Get(s string) (i int, ok bool) {
// 	i, ok = m.m[strings.ToLower(s)]
// 	return
// }
