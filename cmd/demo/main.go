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

		C2 := line.Style.Color[1]

		BLUR := 2
		FADEIN := 150
		FADEOUT := 150

		for _, char := range line.Chars() {

			x, y := char.Left, char.Bottom

			// Silabas por cantar
			s := subs.CopyChar(char) // *char
			s.StartTime = line.StartTime
			s.EndTime = char.StartTime
			s.Tags = fmt.Sprintf(`\an1\pos(%g,%g)\blur%d\1c%s`,
				x, y, BLUR, C2.SSA())
			subs.Add(s)

			// Efecto de silaba
			s = subs.CopyChar(char) // *char
			s.EndTime = s.SylEndTime
			s.Layer = 1
			s.Tags = fmt.Sprintf(
				`\an1\pos(%g,%g)\blur%d\fscx%d\fscy%d\t(fscx%d\fscy%d\blur%d\1c%s)`,
				x, y, BLUR, 130, 130, 100, 100, 1, C2.SSA())
			subs.Add(s)

			// Silabas Muertas (cantadas)
			s = subs.CopyChar(char) // *char
			s.StartTime = char.SylEndTime
			s.EndTime = line.EndTime
			s.Tags = fmt.Sprintf(`\an1\pos(%g,%g)\blur%d\1c%s`,
				x, y, BLUR, C2.SSA())
			subs.Add(s)

			// Efecto de entrada
			s = subs.CopyChar(char) // *char
			px := x + utils.RandomFloat(-5, 5)
			py := y + utils.RandomFloat(-5, 5)
			m := fmt.Sprintf("move(%g, %g, %g, %g)", px, py-char.Height/4, x, y)
			s.Tags = fmt.Sprintf(`\an1\blur%s\fade(%d,0)\%s\1c%s`,
				BLUR, FADEIN, m, C2.SSA())
			s.StartTime = line.StartTime - FADEIN - 50
			s.EndTime = line.StartTime
			subs.Add(s)

			// Efecto de salida
			s = subs.CopyChar(char) // *char
			m = fmt.Sprintf("move(%g, %g, %g, %g)", x, y, px, py+char.Height/2)
			s.Tags = fmt.Sprintf(`\an1\blur%d\fade(0,%d)\%s\1c%s`,
				BLUR, FADEOUT, m, C2.SSA())
			s.StartTime = line.EndTime
			s.EndTime = line.EndTime + FADEOUT - 50
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
