// Package draw
package draw

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/Alquimista/eyecandy/utils"
)

type Filter func(m string) string

type Shape struct {
	draw string
}

// M Draw Move
func (d Shape) M(x, y int) *Shape {
	d.draw += fmt.Sprintf(`m %d %d `, x, y)
	return &d
}

// N Draw Move (no closing)
func (d Shape) N(x, y int) *Shape {
	d.draw += fmt.Sprintf(`n %d %d `, x, y)
	return &d
}

// L Line
func (d Shape) L(x, y int) *Shape {
	d.draw += fmt.Sprintf(`l %d %d `, x, y)
	return &d
}

// B Bézier
func (d Shape) B(args ...int) *Shape {
	lenARGS := len(args)
	if 6 > lenARGS {
		panic("Not enough parameters.")
	} else if lenARGS == 6 {
		x1, y1, x2, y2, x3, y3 := args[0], args[1], args[2],
			args[3], args[4], args[5]
		d.draw += fmt.Sprintf(`b %d %d %d %d %d %d `, x1, y1, x2, y2, x3, y3)
	} else if lenARGS%2 == 0 {
		bspline := "s "
		for _, arg := range args {
			bspline += fmt.Sprintf(`%d `, arg)
		}
		d.draw += bspline + "c"
	} else {
		panic("Wrong parameter count.")
	}
	return &d
}

// Clip Vector Drawing
func (d Shape) Clip(mode int) string {
	if mode < 0 && mode == 3 && mode > 4 {
		panic("Draw mode parameter accept int number in range [1,2,4].")
	}
	return fmt.Sprintf(`\clip(%d,%s)`, mode, d)
}

// IClip Inverse Vector Drawing
func (d Shape) IClip(mode int) string {
	if mode < 0 && mode == 3 && mode > 4 {
		panic("Draw mode parameter accept int number in range [1,2,4].")
	}
	return fmt.Sprintf(`\iclip(%d,%s)`, mode, d)
}

// Draw Drawing command
func (d Shape) Draw(mode int) string {
	if mode < 0 && mode == 3 && mode > 4 {
		panic("Draw mode parameter accept int number in range [1,2,4].")
	}
	return fmt.Sprintf(`{\p%d}%s{\p0}`, mode, d)
}

// String string
func (d Shape) String() string {
	return d.draw
}

// NewScript create a new Script Struct with defaults
func NewShape() *Shape {
	return &Shape{}
}

func Poligon(r int, s int) *Shape {
	iangle := 360.0 / float64(s)
	angle := 90.0 + (iangle / 2.0)
	d := NewShape()
	d = d.M(utils.Polar2Rect(float64(r), angle))
	angle += iangle
	for i := 1; i < s+1; i++ {
		// convert polar to rectangular
		d = d.L(utils.Polar2Rect(float64(r), angle))
		angle += iangle
	}
	d = d.Translate(r, r)
	return d
}

func Pentagon(r int) *Shape {
	return Poligon(r, 5)
}

func Hexagon(r int) *Shape {
	return Poligon(r, 6)
}

func Star(r1 int, r2 int, spikes int) *Shape {
	// the smallest radio is always the inner circle
	if r1 > r2 {
		r1, r2 = r2, r1
	}
	iangle := 360.0 / float64(spikes)
	angle1 := -90.0 + (iangle / 2.0)
	angle2 := angle1 + (iangle / 2.0)

	d := NewShape()
	for i := 0; i < spikes+1; i++ {
		// ass draw commands
		// convert polar to rectangular
		if i == 0 {
			d = d.M(utils.Polar2Rect(float64(r1), angle1))
		} else {
			d = d.L(utils.Polar2Rect(float64(r1), angle1))
		}
		d = d.L(utils.Polar2Rect(float64(r2), angle2))
		angle1 += iangle
		angle2 += iangle
	}
	d = d.Translate(r2, r2)
	return d
}

func Pixel() *Shape {
	return Square(1, 1)
}

func Dot() *Shape {
	return Circle(1, false)
}

func Square(w, h int) *Shape {
	d := NewShape()
	d = d.M(0, 0)
	d = d.L(w, 0)
	d = d.L(w, h)
	d = d.L(0, h)
	return d
}

func Circle(r int, substract bool) *Shape {

	resize := func(m string) string {
		return fmt.Sprintf(`%g`, (utils.Str2float(m)/100.0)*float64(r)*2.0)
	}

	swapCoords := func(m string) string {
		pos := strings.Split(m, " ")
		return fmt.Sprintf(`%s %s`, pos[1], pos[0])
	}

	d := NewShape()
	d.draw = "m 50 0 b 22 0 0 22 0 50 b 0 78 22 100 50 100 b 78 100 100 78 " +
		"100 50 b 100 22 78 0 50 0 "

	if substract {
		d.draw = ShapeFilter(d.draw, swapCoords, "")
	}

	d.draw = ShapeFilter(d.draw, resize, `\d+`)

	return d
}

func Triangle(size int) *Shape {

	h := math.Sqrt(3) * (float64(size) / 2.0)
	base := -h

	d := NewShape()
	d.draw = fmt.Sprintf(`m %g %g l %g %g 0 %g %g %g`,
		float64(size)/2.0, float64(base),
		float64(size), base+h,
		base+h, float64(size)/2.0, base)
	d.Translate(0, int(h+0.5))
	return d
}

func Ring(radio, outlineWidth int) *Shape {
	radio2 := radio - outlineWidth

	circle2 := Circle(radio2, true)
	circle2 = circle2.Translate(-radio2, -radio2)
	circle2 = circle2.Translate(radio, radio)

	d := NewShape()
	d.draw = Circle(radio, false).draw + circle2.draw
	return d
}

func Heart(size int) *Shape {

	resize := func(m string) string {
		return fmt.Sprintf(`%g`, (utils.Str2float(m)/30.0)*float64(size))
	}

	d := NewShape()
	d.draw = "m 15 30 b 27 22 30 18 30 14 30 8 22 " +
		"0 15 10 8 0 0 8 0 14 0 18 3 22 15 30"
	d.draw = ShapeFilter(d.draw, resize, `\d+`)

	return d
}

func ShapeFilter(shape string, f Filter, rx string) string {
	r := regexp.MustCompile(`(-?\d+\.\d+|-?\d+)\s(-?\d+\.\d+|-?\d+)`)
	if rx != "" {
		r = regexp.MustCompile(rx)
	}
	return r.ReplaceAllStringFunc(shape, f)
}

func (d Shape) Scale(x, y float64) *Shape {
	scale := func(m string) string {
		pos := strings.Split(m, " ")
		px := utils.Str2float(pos[0]) * x
		py := utils.Str2float(pos[1]) * y
		return fmt.Sprintf(`%g %g`, px, py)
	}
	d.draw = ShapeFilter(d.draw, scale, "")
	return &d
}

func (d Shape) Translate(x, y int) *Shape {

	move := func(m string) string {
		pos := strings.Split(m, " ")
		px := utils.Str2float(pos[0]) + float64(x)
		py := utils.Str2float(pos[1]) + float64(y)
		return fmt.Sprintf(`%g %g`, px, py)
	}
	d.draw = ShapeFilter(d.draw, move, "")
	return &d
}

func (d Shape) Flip() *Shape {

	flip := func(m string) string {
		pos := strings.Split(m, " ")
		px, py := 0-utils.Str2float(pos[0]), utils.Str2float(pos[1])
		return fmt.Sprintf(`%g %g`, px, py)
	}
	d.draw = ShapeFilter(d.draw, flip, "")
	return &d
}
