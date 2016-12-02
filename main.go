package main

import (
	"fmt"
	"time"

	"github.com/Alquimista/eyecandy-go/reader"
	"github.com/Alquimista/eyecandy-go/writer"
)

// compile passing -ldflags "-X main.Build <build sha1>
// git rev-parse --short HEAD
var Build string

func main() {
	// fmt.Printf("Using build: %s\n", Build)

	// filename := os.Args[1]
	t0 := time.Now()

	input := reader.Read("test/test.ass")
	output := writer.NewScript()
	output.Resolution = input.Resolution
	output.VideoPath = input.VideoPath
	output.VideoZoom = input.VideoZoom
	output.VideoPosition = input.VideoPosition
	output.VideoAR = input.VideoAR
	output.MetaFilename = input.MetaFilename
	output.MetaTitle = input.MetaTitle
	output.MetaOriginalScript = input.MetaOriginalScript
	output.MetaTranslation = input.MetaTranslation
	output.MetaTiming = input.MetaTiming
	output.Audio = input.Audio

	for _, style := range input.StyleUsed {
		s := writer.NewStyle(style.Name)
		s.Name = style.Name
		s.FontName = style.FontName
		s.FontSize = style.FontSize
		s.Color = style.Color
		s.Bold = style.Bold
		s.Italic = style.Italic
		s.Underline = style.Underline
		s.StrikeOut = style.StrikeOut
		s.Scale = style.Scale
		s.Spacing = style.Spacing
		s.Angle = style.Angle
		s.OpaqueBox = style.OpaqueBox
		s.Bord = style.Bord
		s.Shadow = style.Shadow
		s.Alignment = style.Alignment
		s.Margin = style.Margin
		s.Encoding = style.Encoding
		output.AddStyle(s)
	}

	for _, dlg := range input.Dialog.GetNotCommented() {
		d := writer.NewDialog(dlg.Text)
		d.Layer = dlg.Layer
		d.Start = dlg.Start
		d.End = dlg.End
		d.StyleName = dlg.StyleName
		d.Actor = dlg.Actor
		d.Effect = dlg.Effect
		d.Tags = "\\blur5\\1c&HD0C9AD&\\3c&HFFFAB5&\\bord0.5"
		d.Margin = dlg.Margin
		d.Comment = dlg.Comment
		output.AddDialog(d)
	}

	output.Save("test/test.fx.ass")

	elapsed := time.Since(t0)
	fmt.Printf("\nTook %s\n", elapsed)
}
