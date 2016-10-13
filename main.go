package main

import (
	"fmt"
	"time"
	//"os"

	"github.com/Alquimista/eyecandy-go/reader"
)

// compile passing -ldflags "-X main.Build <build sha1>
// git rev-parse --short HEAD
// go list -f '{{join .Deps "\n"}}' |
//   xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}'
var Build string

func main() {
	// filename := os.Args[1]
	t0 := time.Now()

	fmt.Printf("Using build: %s\n", Build)

	script := reader.Read("test/test.ass")
	fmt.Println(script)

	fmt.Println("\nDIALOG")
	for _, dialog := range script.Dialog {
		if dialog.Comment {
			fmt.Println("Comment:", dialog.Style.Name)
		} else {
			fmt.Println("Dialog:", dialog.Style.Name)
		}
	}

	fmt.Println("\nALL DIALOG")
	for _, dialog := range script.DialogWithComment {
		if dialog.Comment {
			fmt.Println("Comment:", dialog.Style.Color)
		} else {
			fmt.Println("Dialog:", dialog.Style.Color.Primary)
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

	elapsed := time.Since(t0)
	fmt.Printf("\nTook %s\n", elapsed)
}
