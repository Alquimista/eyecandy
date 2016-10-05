package main

import (
	"fmt"

	"github.com/Alquimista/eyecandy-go/reader"
	// "os"
	"time"
)

func main() {
	// filename := os.Args[1]
	t0 := time.Now()

	script := reader.Read("test/test.ass")

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
			fmt.Println("Comment:", dialog.Style.Name)
		} else {
			fmt.Println("Dialog:", dialog.Style.Name)
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
