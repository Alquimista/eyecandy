// Package color provides color convention and useful functions
package color

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/Alquimista/eyecandy/interpolate"
	"github.com/Alquimista/eyecandy/utils"
)

// easyrgb.com

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

func clamp(n, min, max float64) float64 {
	if n > max {
		return max
	} else if n < min {
		return min
	}
	return n
}

func ClampRGB1(r, g, b float64) (float64, float64, float64) {
	return clamp(r, 0.0, 1.0), clamp(g, 0.0, 1.0), clamp(b, 0.0, 1.0)
}

func zip(lists ...[]float64) func() []float64 {
	zip := make([]float64, len(lists))
	i := 0
	return func() []float64 {
		for j := range lists {
			if i >= len(lists[j]) {
				return nil
			}
			zip[j] = lists[j][i]
		}
		i++
		return zip
	}
}

func gradient(n int, c *Color, c2 *Color, f interpolate.Interp) (colors []*Color) {

	r1, g1, b1 := c.RGB()
	r2, g2, b2 := c2.RGB()

	red := interpolate.IRange(n, float64(r1), float64(r2), f)
	green := interpolate.IRange(n, float64(g1), float64(g2), f)
	blue := interpolate.IRange(n, float64(b1), float64(b2), f)

	for i := 0; i < n; i++ {
		colors = append(
			colors, NewFromRGB(uint8(red[i]), uint8(green[i]), uint8(blue[i])))
	}
	return
}

func Gradient(n int, clrs []*Color, f interpolate.Interp) []*Color {
	clrn := len(clrs)
	nOut := int(n / (clrn - 1))
	colors := []*Color{}
	if clrn > n {
		panic("The number of colors can not be greater than the steps")
	} else if clrn == n {
		return clrs
	} else if clrn > 2 {
		for i := 0; i < clrn-1; i++ {
			if i == clrn-2 && n%2 == 1 {
				nOut += 1
			}
			colors = append(colors, gradient(nOut, clrs[i], clrs[i+1], f)...)
		}
		return colors
	}
	return gradient(nOut, clrs[0], clrs[1], f)
}

func (c *Color) Gradient(n int, c2 *Color, f interpolate.Interp) []*Color {
	return gradient(n, c, c2, f)
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

func (c *Color) MinMaxRGB1() (min, max float64) {
	return c.MinRGB1(), c.MaxRGB1()
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

func (c Color) XYZ() (x, y, z float64) {

	r, g, b := c.RGB1()
	rangeRGB := []float64{r, g, b}

	rgb := []float64{}
	for _, v := range rangeRGB {
		value := v / 12.92
		if v > 0.04045 {
			value = math.Pow(v+0.055, 2.4)
		}
		rgb = append(rgb, value)
	}

	r, g, b = rgb[0]*100, rgb[1]*100, rgb[2]*100

	// Observer. = 2°, Illuminant = D65
	x = (r * 0.4124) + (g * 0.3576) + (b * 0.1805)
	y = (r * 0.2126) + (g * 0.7152) + (b * 0.0722)
	z = (r * 0.0193) + (g * 0.1192) + (b * 0.9505)
	return
}

func (c Color) LAB() (l, a, b float64) {

	x, y, z := c.XYZ()
	//Observer= 2°, Illuminant= D65
	XYZ := []float64{x / 95.047, y / 100.000, z / 108.883}

	components := []float64{}
	for _, v := range XYZ {
		value := (7.787 * v) + (16.0 / 116.0)
		if v > 0.008856 {
			value = math.Pow(v, 1.0/3.0)
		}
		components = append(components, value)
	}

	x, y, z = components[0], components[1], components[2]
	l = (116 * y) - 16
	a = 500 * (x - y)
	b = 200 * (y - z)

	return
}

func (c1 Color) HCL() (h int, c, l float64) {

	l, a, b := c1.LAB()

	H := math.Atan2(float64(b), float64(a)) //Quadrant by signs

	if H > 0 {
		h = int(((H / math.Pi) * 180) + 0.5)
	} else {
		h = int((360 - (math.Abs(H)/math.Pi)*180) + 0.5)
	}

	c = math.Sqrt(a*a + b*b)

	return
}

func (c Color) HSL() (h, s, l int) {
	r, g, b := c.RGB1()

	rgbMin, rgbMax := c.MinMaxRGB1()
	delta := rgbMax - rgbMin

	L := (rgbMax + rgbMin) / 2.0
	if delta != 0 {
		s = int((delta/(1-math.Abs(2*L-1)))*100 + 0.5)
	}
	if rgbMax == r {
		h = int(60*(g-b)/delta+360) % 360
	} else if rgbMax == g {
		h = int(60*(b-r)/delta + 120)
	} else {
		h = int(60*(r-g)/delta + 240)
	}
	l = int(L*100 + 0.5)
	return
}

func (c Color) HSB() (h, s, b int) {
	h, s, l := c.HSL()
	S := float64(s) / 100.0
	L := float64(l) / 100.0
	B := (2*L + S*(1-math.Abs(2*L-1))) / 2.0
	S = (S * (1 - math.Abs(2*L-1))) / B

	s = int(S*100 + 0.5)
	b = int(B*100 + 0.5)
	return
}

func (c Color) HSV() (h, s, b int) {
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
	r, g, b = ClampRGB1(r, g, b)
	return &Color{
		R: uint8(r*255 + 0.5),
		G: uint8(g*255 + 0.5),
		B: uint8(b*255 + 0.5)}
}

// NewHSL
func NewFromHSL(h, s, l int) *Color {

	S := float64(s) / 100
	L := float64(l) / 100

	C := (1 - math.Abs(2*L-1)) * S
	X := C * (1 - math.Abs(float64((h/60.0)%2-1)))
	m := L - C/2.0

	switch {
	case 0 <= h && h < 60:
		return NewFromRGB1(C+m, X+m, m)
	case 60 <= h && h < 120:
		return NewFromRGB1(X+m, C+m, m)
	case 120 <= h && h < 180:
		return NewFromRGB1(m, C+m, X+m)
	case 180 <= h && h < 240:
		return NewFromRGB1(m, X+m, C+m)
	case 240 <= h && h < 300:
		return NewFromRGB1(X+m, m, C+m)
	case 300 <= h && h < 360:
		return NewFromRGB1(C+m, m, X+m)
	}

	return NewFromRGB(0, 0, 0)
}

// NewHSB
func NewFromHSB(h, s, b int) *Color {
	S := float64(s) / 100.0
	B := float64(b) / 100.0
	L := B * (2 - S) / 2
	S = B * S / (1 - math.Abs(2*L-1))
	return NewFromHSL(h, int(S*100+0.5), int(L*100+0.5))
}

func NewFromHSV(h, s, b int) *Color {
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

func NewFromXYZ(x, y, z float64) *Color {

	// Observer = 2°, Illuminant = D65
	x /= 100 //X from 0 to  95.047
	y /= 100 //Y from 0 to 100.000
	z /= 100 //Z from 0 to 108.883

	// fmt.Println(x, y, z)

	r := x*3.2406 + y*-1.5372 + z*-0.4986
	g := x*-0.9689 + y*1.8758 + z*0.0415
	b := x*0.0557 + y*-0.2040 + z*1.0570

	// fmt.Println(r)

	RGB := []float64{r, g, b}
	components := []float64{}
	for _, v := range RGB {
		value := 12.92 * v
		if v > 0.0031308 {
			value = (1.055 * math.Pow(v, (1.0/2.4))) - 0.055
		}
		components = append(components, value)
	}

	// fmt.Println(components[0])

	return NewFromRGB1(
		components[0],
		components[1],
		components[2])
}

func NewFromLAB(l, a, b float64) *Color {

	y := (l + 16) / 116.0
	x := a/500.0 + y
	z := y - b/200.0

	XYZ := []float64{x, y, z}
	components := []float64{}
	for _, v := range XYZ {
		value := (v - 16.0/116.0) / 7.787
		v = math.Pow(v, 3)
		if v > 0.008856 {
			value = v
		}
		components = append(components, value)
	}

	return NewFromXYZ(
		95.047*components[0],
		100.000*components[1],
		108.883*components[2])
}

func NewFromHCL(h int, c, l float64) *Color {
	a := math.Cos(utils.Deg(float64(h))) * c
	b := math.Sin(utils.Deg(float64(h))) * c
	return NewFromLAB(l, a, b)
}
