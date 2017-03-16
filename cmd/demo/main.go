package main

import (
	"fmt"
	"math"
	"time"

	"github.com/Alquimista/eyecandy"
	. "github.com/Alquimista/eyecandy/asstags"
	"github.com/Alquimista/eyecandy/color"
	"github.com/Alquimista/eyecandy/interpolate"
	"github.com/Alquimista/eyecandy/utils"
)

const (
	inputScript = "test/test.ass"
	ouputScript = "test/test.fx.ass"
)

func Cosine6(t, start, end float64) float64 {
	repeat := 6.0
	t = 0.5 - (math.Cos(repeat*math.Pi*t) / 2)
	return interpolate.Linear(t, start, end)
}

func doFX() {

	subs := eyecandy.NewEffect(inputScript)

	for _, line := range subs.Lines() {

		c1 := line.Style.Color[0]
		c2 := line.Style.Color[1]

		cblue := color.NewFromHTML("#93BEC2")
		colorsb := cblue.Gradient(line.CharN, c1, Cosine6)
		iCustom := interpolate.IRange(line.CharN, 0.0, 1.0, Cosine6)

		// MCOLOR := color.Gradient(
		// 	line.CharN,
		// 	[]*color.Color{
		// 		color.Whitesmoke,
		// 		color.Lightgreen,
		// 		color.Lightsteelblue},
		// 	interpolate.Linear)

		MCOLOR := color.HTMLRange(line.CharN,
			color.Whitesmoke.HTML(),
			color.Lightgreen.HTML(),
			color.Lightsteelblue.HTML())

		for ci, char := range line.Chars() {

			x, y := char.Left, char.Bottom

			// Silabas por cantar
			s := subs.CopyChar(char)
			s.StartTime = line.StartTime
			s.EndTime = char.StartTime
			s.Tags = Pos(x, y) + C(c2.HTML()) + Blur(1) + An(1)
			subs.Add(s)

			// Efecto de silaba
			s = subs.CopyChar(char) // *char
			s.EndTime = char.EndTime + 50
			if s.EndTime >= line.EndTime {
				s.EndTime = char.EndTime
			}
			s.Layer = 1
			s.Tags = Pos(x, y) + Blur(1) + C(MCOLOR[ci]) + Fsc(130.0) +
				T(Fsc(100.0)+Blur(2)+C(colorsb[ci].HTML())) +
				An(1)
			subs.Add(s)

			// Silabas Muertas (cantadas)
			s = subs.CopyChar(char) // *char
			s.StartTime = char.EndTime + 50
			if s.StartTime >= line.EndTime {
				s.StartTime = char.EndTime
			}
			s.EndTime = line.EndTime
			s.Tags = Pos(x, y) + Blur(1) + An(1) + C(colorsb[ci].HTML())
			subs.Add(s)

			// Efecto de entrada
			s = subs.CopyChar(char) // *char
			px := x + float64(utils.RandomInt(-5, 5))
			py := y + float64(utils.RandomInt(-5, 5))
			m := Move(px, py*iCustom[ci]-char.Height/4, x, y)
			s.Tags = Blur(1) + Fade(150, 0) + C(c2.HTML()) + m + An(1)
			s.StartTime = line.StartTime - 100
			s.EndTime = line.StartTime
			subs.Add(s)

			// Efecto de salida
			s = subs.CopyChar(char) // *char
			m = Move(x, y, px, py+char.Height/2)
			s.Tags = Blur(1) + Fade(0, 150) + C(colorsb[ci].HTML()) + m + An(1)
			s.StartTime = line.EndTime
			s.EndTime = line.EndTime + 100
			subs.Add(s)

		}
	}
	subs.Save(ouputScript)
}

func main() {

	t0 := time.Now()
	doFX()
	elapsed := time.Since(t0)

	fmt.Printf("\nTook %s\n", elapsed)

}
