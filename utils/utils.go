// Package utils provide some helpers functions for other packages
package utils

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Alquimista/fonts"
	"github.com/golang/freetype/truetype"

	"golang.org/x/image/font"
)

func DivMod(a, b int) (q, r int) {
	return int(a / b), a % b
}

// Str2int convert a string to an integer.
func Str2int(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Errorf("reader: failed parsing number: %s", err))
	}
	return i
}

func Hex2int(hexStr string) int {
	// base 16 for hexadecimal
	i, err := strconv.ParseInt(hexStr, 16, 64)
	if err != nil {
		panic(fmt.Errorf("reader: failed parsing hex: %s", err))
	}
	return int(i)
}

// Str2bool convert a string (ssa) to a boolean.
func Str2bool(s string) bool {
	return s == "-1"
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

// AppendStrUnique append to a string slice if the value doesn't already exist
func AppendStrUnique(slice []string, s string) []string {
	for _, ele := range slice {
		if ele == s {
			return slice
		}
	}
	return append(slice, s)
}

// MeasureString returns the rendered width and height of the specified text
// given the current font face.
func MeasureString(ff font.Face, s string) (w float64, h float64) {
	d := &font.Drawer{Face: ff}
	return float64(d.MeasureString(s) >> 6),
		float64(ff.Metrics().Height>>6) * 96.0 / 72.0
}

func FontLoad(fontName string, fontSize int) (font.Face, error) {

	// TODO: select the correct font path
	fontBytes, err := fonts.LoadFont(fontName)
	if err != nil {
		return nil, err
	}
	fontFace, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(fontFace, &truetype.Options{
		Size: float64(fontSize) * 72.0 / 96.0,
		DPI:  72,
		// Hinting: font.HintingNone,
		Hinting: font.HintingFull,
	})
	return face, nil
}

func LenString(text string) int {
	return len([]rune(text))
}

func TrimSpaceCount(text string) (string, int, int) {
	preSpace := LenString(text) -
		LenString(strings.TrimLeft(text, " "))
	postSpace := LenString(text) -
		LenString(strings.TrimRight(text, " "))
	return strings.TrimSpace(text), preSpace, postSpace
}

func RandomFloat(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()*(max-min) + min
}

func RandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func RandomChoiceString(list []string) string {
	rand.Seed(time.Now().UnixNano())
	return list[rand.Intn(len(list))]
}

func RandomChoiceInt(list []int) int {
	rand.Seed(time.Now().UnixNano())
	return list[rand.Intn(len(list))]
}
