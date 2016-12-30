// Package color provides color convention and useful functions
package color

import (
	"math"
	"math/rand"
	"reflect"
	"time"

	"github.com/Alquimista/eyecandy/utils"
)

// Grayscale desaturate the color
func (c Color) Grayscale() *Color {
	// http://bit.ly/ce5Kps
	Y := uint8((77*int(c.R) + 151*int(c.G) + 28*int(c.B)) >> 8)
	return NewFromRGB(Y, Y, Y)
}

func (c Color) Invert() *Color {
	return NewFromRGB(^c.R, ^c.G, ^c.B)
}

func (c Color) Complementary() *Color {
	h, s, v := c.HSV()
	h = (h + 180) % 360
	return NewFromHSV(h, s, v)
}

func (c Color) Analog(n int, separation int) (colors []*Color) {
	h, s, v := c.HSV()
	sep := separation
	pl := -1

	for i := 1; i <= n; i++ {
		h = (h + sep*pl) % 360
		if pl < 0 {
			colors = append([]*Color{NewFromHSV(h, s, v)}, colors...)
		} else {
			colors = append(colors, NewFromHSV(h, s, v))
		}
		if i%2 != 0 {
			sep += separation
		}
		pl *= -1
	}
	return colors
}

// Lighter return a lighter version of this color
func (c Color) Lighter(amt int) *Color {
	h, s, l := c.HSL()
	l = int(math.Min(float64(l)+float64(amt), 100))
	return NewFromHSL(h, s, l)
}

// Darker return a darker version of this color
func (c Color) Darker(amt int) *Color {
	h, s, l := c.HSL()
	l = int(math.Max(float64(l)-float64(amt), 0))
	return NewFromHSL(h, s, l)
}

// BlendRGB return a new color, interpolated between this color and
// other by an amount specified by t, ranges from 0 (entirely this color)
// to 1.0 (entirely other.)
func (c Color) BlendRGB(c2 *Color, t float64) *Color {
	r, g, b := c.RGB1()
	r2, g2, b2 := c2.RGB1()
	return NewFromRGB1(r+t*(r2-r), g+t*(g2-g), b+t*(b2-b))
}

func (c Color) MixRGB(c2 *Color) *Color {
	return c.BlendRGB(c2, 0.5)
}

func (c1 Color) DistanceRgb(c2 *Color) float64 {
	deltaR := c1.R - c2.R
	deltaG := c1.G - c2.G
	deltaB := c1.B - c2.B
	return math.Sqrt(float64(deltaR*deltaR + deltaG*deltaG + deltaB*deltaB))
}

func Equal(c, c2 *Color) bool {
	return reflect.DeepEqual(c, c2)
}

func Rainbow(n, s, v int) (colors []*Color) {
	for i := 0; i < n; i++ {
		h := float64(i) * 360 / float64(n)
		colors = append(colors, NewFromHSV(int(h+0.5), s, v))
	}
	return
}

// Algorithm from here:
// http://gamedev.stackexchange.com/questions/46463/is-there-an-optimum-set-of-colors-for-10-players
// Golden Ratio

func DistinguishableColor(n, s, v int) (colors []*Color) {
	for i := 0; i < n; i++ {
		h := math.Mod(360*0.618033988749895*float64(i), 360.0)
		colors = append(colors, NewFromHSV(int(h+0.5), s, v))
	}
	return
}

type rnd func() int

// RandomColorHSV
func RandomColorHSV(s, v int, f rnd) *Color {
	if f == nil {
		f = RGoldenHue
	}
	return NewFromHSV(f(), s, v)
}

func RGoldenHue() int {
	rand.Seed(time.Now().UnixNano())
	return int(math.Mod(360*0.618033988749895*rand.Float64(), 360.0) + 0.5)
}
func RHue() int {
	return utils.RandomInt(0, 360)
}
