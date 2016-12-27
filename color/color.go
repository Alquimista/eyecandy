// Package color provides color convention and useful functions
package color

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/Alquimista/eyecandy/utils"
)

var reColorHEX = regexp.MustCompile(
	`(?:#)*([0-9A-Fa-f]{2})` + // red component
		`([0-9A-Fa-f]{2})` + // green
		`([0-9A-Fa-f]{2})`) // blue

var reColorHEXSTRING = regexp.MustCompile(
	`(?:#)*([a-fA-F0-9]{2}\z|[a-fA-F0-9]{3}\z|[a-fA-F0-9]{6}\z)`)

var reSSAColor = regexp.MustCompile(
	`&H([0-9A-Fa-f]{2})*` + // alpha
		`([0-9A-Fa-f]{2})` + // blue component
		`([0-9A-Fa-f]{2})` + // green
		`([0-9A-Fa-f]{2})` + // red
		`(?:&)*`)

func hexToComponents(color string) []string {
	hexstring := reColorHEXSTRING.FindStringSubmatch(color)[1]
	if len(hexstring) == 3 {
		clr := []string{}
		for _, s := range hexstring {
			clr = append(clr, strings.Repeat(string(s), 2))
		}
		return clr
	} else if len(hexstring) == 2 {
		return []string{hexstring, hexstring, hexstring}
	}
	return reColorHEX.FindStringSubmatch(hexstring)[1:]
}

type Color struct {
	R, G, B, A uint8
}

func (c *Color) Min() (min uint8) {
	min = c.R
	if c.G < min {
		min = c.G
	}
	if c.B < min {
		min = c.B
	}
	return min
}

func (c *Color) Max() (max uint8) {
	max = c.R
	if c.G > max {
		max = c.G
	}
	if c.B > max {
		max = c.B
	}
	return max
}

func (c *Color) MinRGB1() (min float64) {
	r, g, b := c.RGB1()
	min = r
	if g < min {
		min = g
	}
	if b < min {
		min = b
	}
	return min
}

func (c *Color) MaxRGB1() (max float64) {
	r, g, b := c.RGB1()
	max = r
	if g > max {
		max = g
	}
	if b > max {
		max = b
	}
	return max
}

func (c Color) RGB() (uint8, uint8, uint8) {
	return c.R, c.G, c.B
}

func (c Color) RGB1() (float64, float64, float64) {
	return float64(c.R) / 255.0, float64(c.G) / 255.0, float64(c.B) / 255.0
}

func (c Color) RGBA() (uint8, uint8, uint8, uint8) {
	return c.R, c.G, c.B, c.A
}

func (c Color) SSA() string {
	return fmt.Sprintf("&H%02X%02X%02X&", c.B, c.G, c.R)
}

func (c Color) SSAL() string {
	return fmt.Sprintf("&H%02X%02X%02X%02X", c.A, c.B, c.G, c.R)
}

func (c Color) HTML() string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

func (c Color) String() string {
	return c.SSA()
}

func (c Color) HEX() uint32 {
	return uint32(c.R)<<16 + uint32(c.G)<<8 + uint32(c.B)
}

func (c Color) HSL() (h int, s, l float64) {
	r, g, b := c.RGB1()

	rgbMax := c.MaxRGB1()
	rgbMin := c.MinRGB1()
	delta := rgbMax - rgbMin

	l = (rgbMax + rgbMin) / 2
	if delta != 0 {
		s = delta / (1 - math.Abs(2*l-1))
	}
	if rgbMax == r {
		h = int(60*(g-b)/delta+360) % 360
	} else if rgbMax == g {
		h = int(60*(b-r)/delta + 120)
	} else {
		h = int(60*(r-g)/delta + 240)
	}
	return
}

func (c Color) HSB() (h int, s, b float64) {
	h, s, l := c.HSL()
	b = (2*l + s*(1-math.Abs(2*l-1))) / 2
	s = (s * (1 - math.Abs(2*l-1))) / b
	return
}

func (c Color) HSV() (h int, s, b float64) {
	return c.HSB()
}

// NewRGB
func NewFromRGB(r, g, b uint8) *Color {
	return &Color{R: r, G: g, B: b}
}

// NewRGBA
func NewFromRGBA(r, g, b, a uint8) *Color {
	return &Color{R: r, G: g, B: b, A: a}
}

// NewRGB1
func NewFromRGB1(r, g, b float64) *Color {
	return &Color{
		R: uint8(r*255 + 0.5),
		G: uint8(g*255 + 0.5),
		B: uint8(b*255 + 0.5)}
}

// NewHSL
func NewFromHSL(h int, s, l float64) *Color {

	C := (1 - math.Abs(2*l-1)) * s
	X := C * (1 - math.Abs(float64((h/60)%2-1)))
	m := l - C/2.0

	if 0 <= h && h < 60 {
		return NewFromRGB1(C+m, X+m, m)
	} else if 60 <= h && h < 120 {
		return NewFromRGB1(X+m, C+m, m)
	} else if 120 <= h && h < 180 {
		return NewFromRGB1(m, C+m, X+m)
	} else if 180 <= h && h < 240 {
		return NewFromRGB1(m, X+m, C+m)
	} else if 240 <= h && h < 300 {
		return NewFromRGB1(X+m, m, C+m)
	} else if 300 <= h && h < 360 {
		return NewFromRGB1(C+m, m, X+m)
	}
	return NewFromRGB(0, 0, 0)
}

// NewHSB
func NewFromHSB(h int, s, b float64) *Color {
	l := b * (2 - s) / 2
	s = b * s / (1 - math.Abs(2*l-1))
	return NewFromHSL(h, s, l)
}

func NewFromHSV(h int, s, b float64) *Color {
	return NewFromHSB(h, s, b)
}

// NewFromHEX
func NewFromHTML(hexc string) *Color {
	clr := hexToComponents(hexc)
	if clr[0] == "" {
		return &Color{}
	}
	return &Color{
		R: uint8(utils.Hex2int(clr[0])),
		G: uint8(utils.Hex2int(clr[1])),
		B: uint8(utils.Hex2int(clr[2])),
	}
}

func NewFromHEX(x uint32) *Color {
	return NewFromRGB(
		uint8((x>>16)&0xFF),
		uint8((x>>8)&0xFF),
		uint8(x&0xFF))
}

// NewFromHEXAlpha
func NewFromHTMLAlpha(hexc string, a uint8) *Color {
	clr := hexToComponents(hexc)
	if clr[0] == "" {
		return &Color{}
	}
	return &Color{
		R: uint8(utils.Hex2int(clr[0])),
		G: uint8(utils.Hex2int(clr[1])),
		B: uint8(utils.Hex2int(clr[2])),
		A: a,
	}
}

// NewFromHEXAlpha
func NewFromSSA(ssac string) *Color {
	// 0: match, 1: alpha, 2: blue, 3: green, 4: red
	clr := reSSAColor.FindStringSubmatch(ssac)
	if clr[0] == "" {
		return &Color{}
	}
	return &Color{
		R: uint8(utils.Hex2int(clr[4])),
		G: uint8(utils.Hex2int(clr[3])),
		B: uint8(utils.Hex2int(clr[2])),
		A: uint8(utils.Hex2int(clr[1])),
	}
}
