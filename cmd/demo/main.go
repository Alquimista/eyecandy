package main

import (
	"fmt"
	"time"

	"github.com/Alquimista/eyecandy"
	// "github.com/Alquimista/eyecandy/utils"
)

func main() {
	// fmt.Printf("Using build: %s\n", Build)

	// filename := os.Args[1]
	t0 := time.Now()

	subs := eyecandy.NewEffect("test/test.ass")

	// fmt.Println(subs.Resolution)
	// fmt.Println(subs.VideoPath)
	// fmt.Println(subs.VideoZoom)
	// fmt.Println(subs.VideoPosition)
	// fmt.Println(subs.VideoAR)
	// fmt.Println(subs.MetaFilename)
	// fmt.Println(subs.MetaTitle)
	// fmt.Println(subs.MetaOriginalScript)
	// fmt.Println(subs.MetaTranslation)
	// fmt.Println(subs.MetaTiming)
	// fmt.Println(subs.Audio)

	// subs.Resolution = [2]int{640, 480}

	for _, dlg := range subs.Dialogs() {

		d := dlg
		d.Tags = "\\blur0.3\\1c&H705F05&\\2c&H2B0C00&\\3c&HFFFAB5&\\bord1.6\\shad1"
		subs.AddDialog(d)

	}
	subs.Save("test/test.fx.ass")

	elapsed := time.Since(t0)
	fmt.Printf("\nTook %s\n", elapsed)
}
