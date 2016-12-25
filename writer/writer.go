// Package writer create SSA/ASS Subtitle Script
package writer

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/Alquimista/eyecandy/asstime"
	"github.com/Alquimista/eyecandy/color"
	"github.com/Alquimista/eyecandy/utils"
)

// silence?, noise?
const dummyVideoTemplate string = "?dummy:%.6f:%d:%d:%d:%d:%d:%d%s:"
const dummyAudioTemplate string = "dummy-audio:silence?sr=44100&bd=16&" +
	"ch=1&ln=396900000:" // silence?, noise? TODO: dummy audio function
const styleFormat string = "Format: Name, Fontname, Fontsize, " +
	"PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, " +
	"Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, " +
	"BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, " +
	"Encoding"
const styleTemplate string = "Style: %s,%s,%d,%s,%s,%s,%s,%s,%s,%s,%s,%.4f," +
	"%.4f,%.1f,%d,%d,%.4f,%.4f,%d,%04d,%04d,%04d,1"
const dialogFormat string = "Format: Layer, Start, End, Style, Name, " +
	"MarginL, MarginR, MarginV, Effect, Text"
const dialogTemplate string = "%s: %d,%s,%s,%s,%s,0000,0000,0000,%s,%s"
const scriptTemplate string = `[Script Info]
; %s
Title: %s
Original Script: %s
Translation: %s
Timing: %s
ScriptType: v4.00+
PlayResX: %d
PlayResY: %d
WrapStyle: 2
ScaledBorderAndShadow: yes
YCbCr Matrix: TV.601

[Aegisub Project Garbage]
Video File: %s
Video AR Mode: 4
Video AR Value: %.6f
Video Zoom Percent: %.6f
Video Position: %d
Audio File: %s
Active Line: 1

[V4+ Styles]
%s
%s
[Events]
%s
%s`

const (
	// AlignBottomLeft Bottom Left SSA numbered Alignment
	AlignBottomLeft int = 1 + iota
	// AlignBottomCenter Bottom Center SSA numbered Alignment
	AlignBottomCenter
	// AlignBottomRight Bottom Right SSA numbered Alignment
	AlignBottomRight

	// AlignMiddleLeft Middle Left SSA numbered Alignment
	AlignMiddleLeft
	// AlignMiddleCenter Middle Center SSA numbered Alignment
	AlignMiddleCenter
	// AlignMiddleRight Middle Right SSA numbered Alignment
	AlignMiddleRight

	// AlignTopLeft Top Left SSA numbered Alignment
	AlignTopLeft
	// AlignTopCenter Top Left SSA numbered Alignment
	AlignTopCenter
	// AlignTopRight Top Right SSA numbered Alignment
	AlignTopRight
)

// DummyVideo blank video file.
func DummyVideo(framerate float64, w, h int, hexc string, cb bool, timeS int) string {
	c := color.NewFromHex(hexc)
	checkboard := ""
	if cb {
		checkboard = "c"
	}
	frames := asstime.MStoFrames(timeS*asstime.Second, framerate)
	return fmt.Sprintf(
		dummyVideoTemplate,
		utils.Round(framerate, 3), frames, w, h, c.R, c.G, c.B, checkboard)
}

// Style represent subtitle"s styles.
type Style struct {
	Name      string
	FontName  string
	FontSize  int
	Color     [4]*color.Color //Primary, Secondary, Bord, Shadow
	Bold      bool
	Italic    bool
	Underline bool
	StrikeOut bool
	Scale     [2]float64 // WIDTH, HEIGHT map[string]string
	Spacing   float64
	Angle     int
	OpaqueBox bool
	Bord      float64
	Shadow    float64
	Alignment int
	Margin    [3]int // L, R, V map[string]string
	Encoding  int
}

// String get the generated style as a String
func (sty *Style) String() string {

	return fmt.Sprintf(styleTemplate,
		sty.Name,
		sty.FontName, sty.FontSize,
		sty.Color[0].SSAL(),
		sty.Color[1].SSAL(),
		sty.Color[2].SSAL(),
		sty.Color[3].SSAL(),
		utils.Bool2str(sty.Bold), utils.Bool2str(sty.Italic),
		utils.Bool2str(sty.Underline), utils.Bool2str(sty.StrikeOut),
		sty.Scale[0], sty.Scale[1],
		sty.Spacing,
		sty.Angle,
		utils.Bool2Obox(sty.OpaqueBox),
		sty.Bord, sty.Shadow,
		sty.Alignment,
		sty.Margin[0], sty.Margin[1], sty.Margin[2],
	)
}

// NewStyle create a new Style Struct with defaults
func NewStyle(name string) *Style {
	return &Style{
		Name:     name,
		FontName: "Arial",
		FontSize: 35,
		Color: [4]*color.Color{
			color.NewFromHex("#FFFFFF"), //Primary
			color.NewFromHex("#0000FF"), //Secondary
			color.NewFromHex("#000000"), //Bord
			color.NewFromHex("#000000"), //Shadow
		},
		Scale:     [2]float64{100, 100},
		Bord:      2,
		Alignment: AlignBottomCenter,
		Margin:    [3]int{10, 20, 10},
	}
}

// Dialog Represent the subtitle"s lines.
type Dialog struct {
	Layer     int
	Start     string
	End       string
	StyleName string
	Actor     string
	Effect    string
	Text      string
	Tags      string
	Comment   bool
}

// String get the generated Dialog as a String
func (d *Dialog) String() string {
	text := d.Text
	key := "Dialogue"
	if d.Comment {
		key = "Comment"
	}
	if d.Tags != "" {
		text = "{" + d.Tags + "}" + d.Text
	}
	return fmt.Sprintf(dialogTemplate,
		key,
		d.Layer,
		d.Start, d.End,
		d.StyleName, d.Actor,
		d.Effect,
		text)
}

// NewDialog create a new Dialog Struct with defaults
func NewDialog(text string) *Dialog {
	return &Dialog{
		StyleName: "Default",
		Start:     "0:00:00.00", End: "0:00:05.00",
		Text: text}
}

// Script SSA/ASS Subtitle Script.
type Script struct {
	Dialog             []*Dialog
	Style              map[string]*Style
	Comment            string
	Resolution         [2]int // WIDTH, HEIGHT map[string]string
	VideoPath          string
	VideoZoom          float64
	VideoPosition      int
	VideoAR            float64
	MetaFilename       string
	MetaTitle          string
	MetaOriginalScript string
	MetaTranslation    string
	MetaTiming         string
	Audio              string
}

// GetStyle get the Style matching the argument name if exist
// else return the Default Style
func (s *Script) GetStyle(name string) *Style {
	style, ok := s.Style[name]
	if !ok {
		style = s.Style["Default"]
		style.Name = name
	}
	return style
}

// StyleExists get if a Style exists matching the argument name
func (s *Script) StyleExists(name string) bool {
	_, ok := s.Style[name]
	return ok
}

// AddStyle add a Style to SSA/ASS Script.
func (s *Script) AddStyle(sty *Style) {
	if !s.StyleExists(sty.Name) {
		s.Style[sty.Name] = sty
	}
}

// AddDialog add a Dialog to SSA/ASS Script.
func (s *Script) AddDialog(d *Dialog) {
	if d.Text != "" {
		s.Dialog = append(s.Dialog, d)
	}
}

// String get the generated SSA/ASS Script as a String
func (s *Script) String() string {

	// Add default dialog and style
	s.AddStyle(NewStyle("Default"))
	if len(s.Dialog) == 0 {
		s.AddDialog(NewDialog("EyecandyFX"))
	}

	if s.Resolution[0] == 0 || s.Resolution[1] == 0 {
		s.Resolution = [2]int{1280, 720}
	}

	if s.MetaOriginalScript == "" {
		s.MetaOriginalScript = s.MetaFilename
	}

	if s.VideoAR == 0 {
		s.VideoAR = float64(s.Resolution[0]) / float64(s.Resolution[1])
	}
	if s.VideoPath == "" {
		s.VideoPath = DummyVideo(
			asstime.FpsNtscFilm,
			s.Resolution[0], s.Resolution[1],
			"#000",
			false,
			600)
	}
	if s.VideoPath != "" {
		if strings.HasPrefix(s.VideoPath, "?dummy") {
			s.Audio = dummyAudioTemplate
		} else {
			s.Audio = s.VideoPath
		}
	}

	var dialogStyleNames []string
	var styles bytes.Buffer
	var dialogs bytes.Buffer

	for _, d := range s.Dialog {
		if !d.Comment {
			dialogStyleNames = utils.AppendStrUnique(
				dialogStyleNames, d.StyleName)
		}
		dialogs.WriteString(d.String() + "\n")
	}

	// Write only used styles in dialogs
	// If doesn't exist create it
	for _, sname := range dialogStyleNames {
		_, ok := s.Style[sname]
		i := 0
		for _, sty := range s.Style {
			if !ok {
				sty = s.Style["Default"]
				sty.Name = sname
				i++
			} else {
				i = 1
			}
			if sty.Name == sname && i == 1 {
				styles.WriteString(sty.String() + "\n")
			}
		}
	}

	return fmt.Sprintf(scriptTemplate,
		s.Comment,
		s.MetaTitle, s.MetaOriginalScript,
		s.MetaTranslation, s.MetaTiming,
		s.Resolution[0], s.Resolution[1],
		s.VideoPath, s.VideoAR, s.VideoZoom, s.VideoPosition,
		s.Audio,
		styleFormat, styles.String(),
		dialogFormat, dialogs.String())

}

// Save write an SSA/ASS Subtitle Script.
func (s *Script) Save(fn string) {

	BOM := "\uFEFF"
	f, err := os.Create(fn)
	if err != nil {
		panic(fmt.Errorf("writer: failed saving subtitle file: %s", err))
	}
	defer f.Close()

	s.MetaFilename = fn

	n, err := f.WriteString(BOM + s.String())
	if err != nil {
		fmt.Println(n, err)
	}

	// save changes
	err = f.Sync()
}

// NewScript create a new Script Struct with defaults
func NewScript() *Script {
	return &Script{
		Comment:   "Script generated by Eyecandy",
		MetaTitle: "Default Eyecandy file",
		VideoZoom: 0.75,
		Style:     map[string]*Style{},
	}
}
