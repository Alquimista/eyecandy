// writer SSA/ASS Subtitle Script Writer
// package writer
package main

import (
	"fmt"
)

const STYLE_FORMAT string = "Format: Name, Fontname, Fontsize, " +
	"PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, " +
	"Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, " +
	"BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, " +
	"Encoding"
const DIALOG_FORMAT string = "Format: Layer, Start, End, Style, Name, " +
	"MarginL, MarginR, MarginV, Effect, Text"

func main() {
	fmt.Println(STYLE_FORMAT)
}

// `[Script Info]
// ; Script generated by Eyecandy
// ScriptType: v4.00+
// {meta}
// {resolution}
// WrapStyle: 2
// ScaledBorderAndShadow: yes
// YCbCr Matrix: TV.601

// [Aegisub Project Garbage]
// {video}
// Active Line: 1

// [V4+ Styles]
// {style_format}
// {styles}

// [Events]
// {dialog_format}
// {dialogs}

// `
