package main

import (
	"fmt"
	"time"

	"github.com/Alquimista/eyecandy"
)

const (
	InputScript = "test/test.ass"
	OuputScript = "test/test.fx.ass"
)

func doFX() {
	subs := eyecandy.NewEffect(InputScript)
	for _, line := range subs.Lines() {

		for _, syl := range line.Syls() {

			s := subs.CopySyl(syl) // *syl
			s.StartTime = line.StartTime - 50
			s.EndTime = line.EndTime + 50
			s.Tags = fmt.Sprintf(`\pos(%g,%g)\fad(100,100)\blur0.3\1c&H705F05&\3c&HFFFAB5&\bord1.2\shad1`, s.X, s.Y)
			subs.Add(s)

			s = subs.CopySyl(syl) // *syl
			s.Layer = 1
			s.Tags = fmt.Sprintf(`\pos(%g,%g)\blur0.3\bord0`, s.X, s.Y)
			subs.Add(s)
		}
	}
	subs.Save(OuputScript)
}

func main() {

	t0 := time.Now()
	doFX()
	elapsed := time.Since(t0)

	fmt.Printf("\nTook %s\n", elapsed)

}
