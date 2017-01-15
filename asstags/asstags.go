// Package asstags
package asstags

import (
	"fmt"

	"github.com/Alquimista/eyecandy/color"
)

// Bord Border size
func Bord(size float64) string {
	return fmt.Sprintf(`\bord%g`, size)
}

// XBord Border size X
func XBord(size float64) string {
	return fmt.Sprintf(`\xbord%g`, size)
}

// XBord Border size Y
func YBord(depth float64) string {
	return fmt.Sprintf(`\ybord%g`, depth)
}

// Shad Shadow distance
func Shad(depth float64) string {
	return fmt.Sprintf(`\shad%g`, depth)
}

// XShad Shadow distance X
func XShad(depth float64) string {
	return fmt.Sprintf(`\xshad%g`, depth)
}

// YShad Shadow distance Y
func YShad(depth float64) string {
	return fmt.Sprintf(`\yshad%g`, depth)
}

// Be Blur edges
func Be(strength int) string {
	return fmt.Sprintf(`\be%d`, strength)
}

// Blur edges (Gaussian kernel)
func Blur(strength float64) string {
	return fmt.Sprintf(`\be%g`, strength)
}

// Fsc Font scale
func Fsc(args ...interface{}) string {

	lenARGS := len(args)
	if 1 > lenARGS {
		panic("Not enough parameters.")
	}
	if lenARGS == 1 {
		scale := args[0].(float64)
		return Fscx(scale) + Fscy(scale)
	} else if lenARGS == 2 {
		scalex := args[0].(float64)
		scaley := args[1].(float64)
		return Fscx(scalex) + Fscy(scaley)
	} else {
		panic("Wrong parameter count.")
	}
	return ""
}

// Fsc Font scale X
func Fscx(scale float64) string {
	return fmt.Sprintf(`\fscx%g`, scale)
}

// Fsc Font scale Y
func Fscy(scale float64) string {
	return fmt.Sprintf(`\fscy%g`, scale)
}

// Frx Text rotation X
func Frx(amount float64) string {
	return fmt.Sprintf(`\frx%g`, amount)
}

// Fry Text rotation Y
func Fry(amount float64) string {
	return fmt.Sprintf(`\fry%g`, amount)
}

// Fry Text rotation Z
func Frz(amount float64) string {
	return fmt.Sprintf(`\frz%g`, amount)
}

// Fr Text rotation Z
func Fr(amount float64) string {
	return Frz(amount)
}

// Text shearing X
func Fax(factor float64) string {
	return fmt.Sprintf(`\fax%g`, factor)
}

// Text shearing Y
func Fay(factor float64) string {
	return fmt.Sprintf(`\fay%g`, factor)
}

// C Set the color of the following text.
func C(args ...interface{}) string {
	lenARGS := len(args)
	if 1 > lenARGS {
		panic("Not enough parameters.")
	}
	if lenARGS == 1 {
		hextring, ok := args[0].(string)
		if !ok {
			panic("1st parameter not type string.")
		}
		return fmt.Sprintf(`\c%s`, color.NewFromHTML(hextring).SSA())
	} else if lenARGS == 2 {
		i, ok := args[0].(int)
		if !ok {
			panic("1st parameter not type int.")
		}
		if i > 4 && i < 1 {
			panic("1st parameter accept int number in range [1..4].")
		}
		hextring, ok1 := args[1].(string)
		if !ok1 {
			panic("2d parameter not type string.")
		}
		return fmt.Sprintf(`\%dc%s`, i, color.NewFromHTML(hextring).SSA())
	} else {
		panic("Wrong parameter count.")
	}
}

// A Sets the alpha of all components
func A(args ...interface{}) string {
	lenARGS := len(args)
	if 1 > lenARGS {
		panic("Not enough parameters.")
	}
	if lenARGS == 1 {
		alpha, ok := args[0].(int)
		if !ok {
			panic("1st parameter not type int.")
		}
		if alpha > 255 && alpha < 0 {
			panic("alpha parameter accept int number in range [0..255].")
		}
		return fmt.Sprintf(`\alpha&H%02X`, alpha)
	} else if lenARGS == 2 {
		i, ok := args[0].(int)
		if !ok {
			panic("1st parameter not type int.")
		}
		if i > 4 && i < 1 {
			panic("1st parameter accept int number in range [1..4].")
		}
		alpha, ok1 := args[1].(int)
		if !ok1 {
			panic("2d parameter not type int.")
		}
		if alpha > 255 && alpha < 0 {
			panic("alpha parameter accept int number in range [0..255].")
		}
		return fmt.Sprintf(`\%d&H%02X`, i, alpha)
	} else {
		panic("Wrong parameter count.")
	}
}

// An Line alignment
func An(align int) string {
	if align > 9 && align < 1 {
		panic("align parameter accept int number in range [1..9].")
	}
	return fmt.Sprintf(`\an%d`, align)
}

// Pos Set the position and Movement of the line
func Pos(x, y float64) string {
	return fmt.Sprintf(`\pos(%g,%g)`, x, y)
}

// Move Set the position and Movement of the line
func Move(args ...interface{}) string {
	lenARGS := len(args)
	movestr := ""
	if lenARGS > 6 || lenARGS%2 != 0 {
		panic("Wrong parameter count.")
	}
	if lenARGS == 2 {
		x, ok := args[0].(float64)
		if !ok {
			panic("1st parameter not type float64.")
		}
		y, ok1 := args[1].(float64)
		if !ok1 {
			panic("2d parameter not type float64.")
		}
		movestr = Pos(x, y)
	} else if lenARGS == 4 {
		x1, ok1 := args[0].(float64)
		if !ok1 {
			panic("1st parameter not type float64.")
		}
		y1, ok2 := args[1].(float64)
		if !ok2 {
			panic("2d parameter not type float64.")
		}
		x2, ok3 := args[2].(float64)
		if !ok3 {
			panic("3rd parameter not type float64.")
		}
		y2, ok4 := args[3].(float64)
		if !ok4 {
			panic("4th parameter not type float64.")
		}
		movestr = fmt.Sprintf(`\move(%g,%g,%g,%g)`, x1, y1, x2, y2)
	} else if lenARGS == 6 {
		x1, ok1 := args[0].(float64)
		if !ok1 {
			panic("1st parameter not type float64.")
		}
		y1, ok2 := args[1].(float64)
		if !ok2 {
			panic("2d parameter not type float64.")
		}
		x2, ok3 := args[2].(float64)
		if !ok3 {
			panic("3rd parameter not type float64.")
		}
		y2, ok4 := args[3].(float64)
		if !ok4 {
			panic("4th parameter not type float64.")
		}
		t1, ok5 := args[4].(int)
		if !ok5 {
			panic("5ft parameter not type int.")
		}
		t2, ok6 := args[5].(int)
		if !ok6 {
			panic("6st parameter not type int.")
		}
		movestr = fmt.Sprintf(`\move(%g,%g,%g,%g,%d,%d)`,
			x1, y1, x2, y2, t1, t2)
	}
	return movestr
}

// Org Rotation origin
func Org(x, y float64) string {
	return fmt.Sprintf(`\org(%g,%g)`, x, y)
}

// Fad Fade
func Fad(fadein, fadeout int) string {
	return fmt.Sprintf(`\fad(%d,%d)`, fadein, fadeout)
}

// Fade (complex)
func Fade(args ...interface{}) string {
	lenARGS := len(args)
	if 2 > lenARGS {
		panic("Not enough parameters.")
	}
	if lenARGS == 2 {
		fadein, ok := args[0].(int)
		if !ok {
			panic("1st parameter not type int.")
		}
		fadeout, ok1 := args[1].(int)
		if !ok1 {
			panic("2d parameter not type int.")
		}
		return Fad(fadein, fadeout)
	} else if lenARGS == 7 {
		a1, ok := args[0].(int)
		if !ok {
			panic("1st parameter not type int.")
		}
		if a1 > 255 && a1 < 0 {
			panic("a1 parameter accept int number in range [0..255].")
		}
		a2, ok1 := args[1].(int)
		if !ok1 {
			panic("2d parameter not type int.")
		}
		if a2 > 255 && a2 < 0 {
			panic("a2 parameter accept int number in range [0..255].")
		}
		a3, ok2 := args[2].(int)
		if !ok2 {
			panic("3rd parameter not type int.")
		}
		if a3 > 255 && a3 < 0 {
			panic("a3 parameter accept int number in range [0..255].")
		}
		t1, ok3 := args[3].(int)
		if !ok3 {
			panic("4rt parameter not type int.")
		}
		t2, ok4 := args[4].(int)
		if !ok4 {
			panic("5th parameter not type int.")
		}
		t3, ok5 := args[5].(int)
		if !ok5 {
			panic("6d parameter not type int.")
		}
		t4, ok6 := args[5].(int)
		if !ok6 {
			panic("7en parameter not type int.")
		}
		return fmt.Sprintf(`\fade(%d,%d,%d,%d,%d,%d,%d)`,
			a1, a2, a3, t1, t2, t3, t4)
	} else {
		panic("Wrong parameter count.")
	}
}

func T(args ...interface{}) string {
	lenARGS := len(args)
	if 1 > lenARGS {
		panic("Not enough parameters.")
	}
	if lenARGS == 1 {
		// T(modifiers)
		m, ok := args[0].(string)
		if !ok {
			panic("1st parameter not type string.")
		}
		return fmt.Sprintf(`\t(%s)`, m)
	} else if lenARGS == 2 {
		// T(accel, modifiers)
		accel, ok := args[0].(float64)
		if !ok {
			panic("1st parameter not type float64.")
		}
		m, ok1 := args[1].(string)
		if !ok1 {
			panic("2nd parameter not type string.")
		}
		return fmt.Sprintf(`\t(%0.2f,%s)`, m, accel)
	} else if lenARGS == 3 {
		// T(t1, t2, style)
		t1, ok := args[0].(int)
		if !ok {
			panic("1st parameter not type int.")
		}
		t2, ok1 := args[1].(int)
		if !ok1 {
			panic("2nd parameter not type int.")
		}
		m, ok2 := args[2].(string)
		if !ok2 {
			panic("3rd parameter not type string.")
		}
		return fmt.Sprintf(`\t(%d,%d,%s)`, t1, t2, m)
	} else if lenARGS == 4 {
		// T(t1, t2, accel, modifiers)
		t1, ok := args[0].(int)
		if !ok {
			panic("1st parameter not type int.")
		}
		t2, ok1 := args[1].(int)
		if !ok1 {
			panic("2nd parameter not type int.")
		}
		accel, ok2 := args[0].(float64)
		if !ok2 {
			panic("3rd parameter not type float64.")
		}
		m, ok3 := args[3].(string)
		if !ok3 {
			panic("4th parameter not type string.")
		}
		return fmt.Sprintf(`\t(%d,%d,%s)`, t1, t2, accel, m)
	} else {
		panic("Wrong parameter count.")
	}
	return ""
}

func Clip(x1, y1, x2, y2 int) string {
	return fmt.Sprintf(`\clip(%d,%d,%d,%d)`, x1, y1, x2, y2)
}

func IClip(x1, y1, x2, y2 int) string {
	return fmt.Sprintf(`\clip(%d,%d,%d,%d)`, x1, y1, x2, y2)
}
