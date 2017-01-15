// Package draw
package draw

import (
	"fmt"
)

type Cursor struct {
	draw string
	mode int
}

// M Draw Move
func (d Cursor) M(x, y int) *Cursor {
	d.draw += fmt.Sprintf(`m %d %d `, x, y)
	return &d
}

// N Draw Move (no closing)
func (d Cursor) N(x, y int) *Cursor {
	d.draw += fmt.Sprintf(`n %d %d `, x, y)
	return &d
}

// L Line
func (d Cursor) L(x, y int) *Cursor {
	d.draw += fmt.Sprintf(`l %d %d `, x, y)
	return &d
}

// B Bézier
func (d Cursor) B(args ...int) *Cursor {
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
		d.draw += bspline + " c"
	} else {
		panic("Wrong parameter count.")
	}
	return &d
}

// Clip Vector Drawing
func (d Cursor) Clip() string {
	if d.mode < 0 && d.mode == 3 && d.mode > 4 {
		panic("Draw mode parameter accept int number in range [1,2,4].")
	}
	return fmt.Sprintf(`\clip(%d,%s)`, d.mode, d)
}

// IClip Inverse Vector Drawing
func (d Cursor) IClip() string {
	if d.mode < 0 && d.mode == 3 && d.mode > 4 {
		panic("Draw mode parameter accept int number in range [1,2,4].")
	}
	return fmt.Sprintf(`\iclip(%d,%s)`, d.mode, d)
}

// Draw Drawing command
func (d Cursor) Draw() string {
	if d.mode < 0 && d.mode == 3 && d.mode > 4 {
		panic("Draw mode parameter accept int number in range [1,2,4].")
	}
	return fmt.Sprintf(`{\p%d}%s{\p0}`, d.mode, d)
}

// String string
func (d Cursor) String() string {
	return d.draw
}

// NewScript create a new Script Struct with defaults
func NewCursor(mode int) *Cursor {
	return &Cursor{mode: mode}
}
