// Package asstags
package asstags

// Mov Set the position and Movement of the line (incremental)
func Mov(args ...interface{}) string {
	lenARGS := len(args)
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
		return Pos(x, y)
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
		return Move(
			x1, y1,
			x1+x2,
			y1+y2)
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
		return Move(
			x1, y1,
			x1+x2,
			y1+y2,
			args[4], args[5])
	} else {
		return ""
	}
}

// Ejemplo:
// CycleTags(200, 1000, 100, be(1), be(2))
// >>> '\\t(200,300,\\be1)\t(300,400,\\be2)..\\t(900,1000,\\be2)'

func CycleTags(start, dur, interval int, tags ...string) (ttags string) {
	i := 0
	n := len(tags)
	startTime := start
	endTime := startTime + interval
	for {
		ttags += T(startTime, endTime, tags[i%n])
		startTime = endTime
		endTime += interval
		if endTime >= dur {
			ttags += T(startTime, dur, tags[i%n])
			break
		}
		i++
	}
	return
}

func FscAR(x, y float64, ar float64, res [2]int) string {
	return Fscx(x) + Fscy(y*float64(res[1])*(ar)/float64(res[0]))
}

func FscScale(x, y float64, scale [2]float64) string {
	return Fscx(x*scale[0]/100) + Fscy(y*scale[1]/100)
}
