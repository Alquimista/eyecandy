package main

import (
	"fmt"
	"time"
	//"os"

	"github.com/Alquimista/eyecandy-go/reader"
	"github.com/Alquimista/eyecandy-go/writer"
)

// compile passing -ldflags "-X main.Build <build sha1>
// git rev-parse --short HEAD
var Build string

func main() {
	fmt.Printf("Using build: %s\n", Build)

	// filename := os.Args[1]
	t0 := time.Now()

	scriptw := writer.NewScript()
	// scriptw.Resolution = [2]int{848, 480}

	d := writer.NewDialog("EyecandyFX")
	scriptw.AddDialog(d)

	d = writer.NewDialog("EyecandyFXComment1")
	d.Style = "Default2"
	d.Tags = "\\blur5\\3c&HFFFAB5&\\bord0.5"
	scriptw.AddDialog(d)

	d = writer.NewDialog("EyecandyFXComment2")
	d.Comment = false
	d.Style = "Default3"
	scriptw.AddDialog(d)

	d = writer.NewDialog("EyecandyFXComment3")
	d.Comment = true
	d.Style = "Default4"
	scriptw.AddDialog(d)

	s := writer.NewStyle("Default2")
	s.Color[0] = "&H00FFFAB5&"
	scriptw.AddStyle(s)

	s = writer.NewStyle("Default5")
	scriptw.AddStyle(s)

	// fmt.Println(scriptw.ToString())
	scriptw.Save("test/test.ass")

	script := reader.Read("test/test.ass")
	// fmt.Println(script)

	fmt.Println("\nDIALOG")
	for _, dialog := range script.Dialog.GetNotCommented() {
		fmt.Println("Dialog:", dialog.Style.Name)
	}

	fmt.Println("\nALL DIALOG")
	for _, dialog := range script.Dialog.GetAll() {
		if dialog.Comment {
			fmt.Println("Comment:", dialog.Style.Color)
		} else {
			fmt.Println("Dialog:", dialog.Style.Color)
		}
	}

	fmt.Println("\nALL STYLES")
	for styleName, style := range script.Style {
		fmt.Println("Style:", styleName, style)
	}

	fmt.Println("\nUSED STYLES")
	for styleName, style := range script.StyleUsed {
		fmt.Println("StyleUsed:", styleName, style)
	}

	// fmt.Println(utils.FromScale(127, 0, 255))
	// fmt.Println(utils.ToScale(0.5, 0, 255))
	// fmt.Println(writer.ALIGN["top left"])

	elapsed := time.Since(t0)
	fmt.Printf("\nTook %s\n", elapsed)
}
