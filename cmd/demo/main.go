package main

import (
	"fmt"
	"time"

	"github.com/Alquimista/eyecandy"
	"github.com/Alquimista/eyecandy/utils"
)

const (
	inputScript = "test/test.ass"
	ouputScript = "test/test.fx.ass"
)

func doFX() {
	subs := eyecandy.NewEffect(inputScript)

	for _, line := range subs.Lines() {

		for _, char := range line.Chars() {

			x, y := char.Left, char.Bottom

			// Silabas por cantar
			s := subs.CopyChar(char) // *char
			s.StartTime = line.StartTime
			s.EndTime = char.StartTime
			s.Tags = fmt.Sprintf(`\an1\pos(%g,%g)\blur2\1c&HDFDFDF&`, x, y)
			subs.Add(s)

			// Efecto de silaba
			s = subs.CopyChar(char) // *char
			s.EndTime = s.SylEndTime
			s.Layer = 1
			s.Tags = fmt.Sprintf(`\an1\pos(%g,%g)\blur2\fscx130\fscy130\t(fscx100\fscy100\blur2\1c&HDFDFDF&)`,
				x, y)
			subs.Add(s)

			// Silabas Muertas (cantadas)
			s = subs.CopyChar(char) // *char
			s.StartTime = char.SylEndTime
			s.EndTime = line.EndTime
			s.Tags = fmt.Sprintf(`\an1\pos(%g,%g)\blur2\1c&HDFDFDF&`, x, y)
			subs.Add(s)

			// Efecto de entrada
			s = subs.CopyChar(char) // *char
			px := x + utils.RandomFloat(-5, 5)
			py := y + utils.RandomFloat(-5, 5)
			m := fmt.Sprintf("move(%g, %g, %g, %g)", px, py-char.Height/4, x, y)
			s.Tags = fmt.Sprintf(`\an1\blur2\fade(150,0)\%s\1c&HDFDFDF&`, m)
			s.StartTime = line.StartTime - 100
			s.EndTime = line.StartTime
			subs.Add(s)

			// Efecto de salida
			s = subs.CopyChar(char) // *char
			m = fmt.Sprintf("move(%g, %g, %g, %g)", x, y, px, py+char.Height/2)
			s.Tags = fmt.Sprintf(`\an1\blur2\fade(0,150)\%s\1c&HDFDFDF&`, m)
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
