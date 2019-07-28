// Package utils provide some helpers functions for other packages
package utils

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Alquimista/eyecandy/interpolate"
	"github.com/golang/freetype/truetype"

	//"github.com/Alquimista/fonts"
	// "github.com/stephenwithav/fontcache"
	"github.com/Alquimista/eyecandy/fontcache"

	"golang.org/x/image/font"
)

const (
	// RadToDeg Radians to Degrees
	RadToDeg = 180 / math.Pi
	// DegToRad Degrees to Radians
	DegToRad = math.Pi / 180
)

// Rad Degrees to Radians
func Rad(d float64) float64 {
	return d * DegToRad
}

// Deg Radians to Degrees
func Deg(r float64) float64 {
	return r * RadToDeg
}

// DivMod do a division an a modulo operation.
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

// Hex2int convert hexadecimal to decimal.
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
	ssampling := 1.0
	w = float64(d.MeasureString(s)>>6) / ssampling
	h = (float64(ff.Metrics().Height>>6) * 96 / 72) / ssampling
	return
}

// LoadFont load and parse a truetype font.
func LoadFont(fontName string, fontSize int) (face font.Face, err error) {

	fc := fontcache.New()
	fc.Init(fontcache.FontPaths)

	// Retrieve font by name for use in a program.
	f, ok := fc[strings.ToLower(fontName)]
	if !ok {
		return nil, nil
	}
	face = truetype.NewFace(f, &truetype.Options{
		Size: float64(fontSize) * 72.0 / 96.0,
		DPI:  72,
	})
	return face, nil
}

// LenString length of a string.
func LenString(text string) int {
	return len([]rune(text))
}

// TrimSpaceCount trims spaces of text and count the number of Left and Right spaces
func TrimSpaceCount(text string) (string, int, int) {
	leftSpace := LenString(text) -
		LenString(strings.TrimLeft(text, " "))
	rightSpace := LenString(text) -
		LenString(strings.TrimRight(text, " "))
	return strings.TrimSpace(text), leftSpace, rightSpace
}

// ChangeSignN convert possitive to negative and back (integer)
func ChangeSignN(n int) int {
	return n * -1
}

// ChangeSignF convert possitive to negative and back (float)
func ChangeSignF(n float64) float64 {
	return n * -1.0
}

// RandomFloat random decimal number between min and max
func RandomFloat(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()*(max-min) + min
}

// RandomInt random number between min and max
func RandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max+1-min) + min
}

// RandomFloatRange random n numbers between min and max (float)
func RandomFloatRange(n int, min, max float64) (nums []float64) {
	for i := 0; i < n; i++ {
		nums = append(nums, RandomFloat(min, max))
	}
	return
}

// RandomIntRange random n numbers between min and max (integer)
func RandomIntRange(n, min, max int) (nums []int) {
	for i := 0; i < n; i++ {
		nums = append(nums, RandomInt(min, max))
	}
	return
}

//TODO: Generic RandomChoice

// RandomChoiceString select a random choice in a string slice
func RandomChoiceString(list []string) string {
	rand.Seed(time.Now().UnixNano())
	return list[rand.Intn(len(list))]
}

// RandomChoiceInt select a random choice in a int slice
func RandomChoiceInt(list []int) int {
	rand.Seed(time.Now().UnixNano())
	return list[rand.Intn(len(list))]
}

// RandomChoiceFloat select a random choice in a float64 slice
func RandomChoiceFloat(list []float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return list[rand.Intn(len(list))]
}

// Polar2Rect convert polar coords to rectangular.
func Polar2Rect(radius, angle float64) (x, y int) {
	theta := Rad(angle)
	return int(math.Cos(theta) * radius), int(math.Sin(theta) * radius)
}

// Rect2Polar convert rectangolar coords to rectangular.
func Rect2Polar(px, py int) (angle, r float64) {
	if px == 0 && py == 0 {
		return 0, 0 // The angle is actually undefined.
	}
	x, y := float64(px), float64(py)
	return Deg(math.Mod(math.Atan2(x, y)+2*math.Pi, 2*math.Pi)),
		math.Hypot(x, y)
}

// CircleRange generate x, y coords of a circle.
func CircleRange(
	n int, x, y, radius float64, f interpolate.Interp) (nums []float64) {
	for _, angle := range interpolate.ICircleRange(n, f) {
		px, py := Polar2Rect(radius, angle)
		nums = append(nums, x+float64(px), y+float64(py))
	}
	return
}
