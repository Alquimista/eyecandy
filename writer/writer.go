// Package writer create SSA/ASS Subtitle Script
package writer

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/Alquimista/eyecandy/asstime"
	"github.com/Alquimista/eyecandy/utils"
)

const dummyVideoTemplate string = "?dummy:%.6f:%d:%d:%d:%d:%d:%d%s:"
const dummyAudioTemplate string = "dummy-audio:silence?sr=44100&bd=16&" +
	"ch=1&ln=396900000:"
const styleFormat string = "Format: Name, Fontname, Fontsize, " +
	"PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, " +
	"Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, " +
	"BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, " +
	"Encoding"
const styleTemplate string = "Style: %s,%s,%d,%s,%s,%s,%s,%s,%s,%s,%s,%.4f," +
	"%.4f,%d,%d,%d,%.4f,%.4f,%d,%04d,%04d,%04d,%d"
const dialogFormat string = "Format: Layer, Start, End, Style, Name, " +
	"MarginL, MarginR, MarginV, Effect, Text"
const dialogTemplate string = "%s: %d,%s,%s,%s,%s,%04d,%04d,%04d,%s,%s"
const scriptTemplate string = `[Script Info]
; %s
ScriptType: v4.00+
Title: %s
Original Script: %s
Translation: %s
Timing: %s
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

// ALIGN SSA/ASS alignment map
var ALIGN = make(map[string]int)

// ENC SSA/ASS encoding map
var ENC = make(map[string]int)

func init() {

	// ENC := NewCIMap()
	// ALIGN := NewCIMap()

	ALIGN["top left"] = 7
	ALIGN["top center"] = 8
	ALIGN["top right"] = 9
	ALIGN["middle left"] = 4
	ALIGN["middle center"] = 5
	ALIGN["middle right"] = 6
	ALIGN["bottom left"] = 1
	ALIGN["bottom center"] = 2
	ALIGN["bottom right"] = 3

	ENC["ansi"] = 0
	ENC["default"] = 1
	ENC["symbol"] = 2
	ENC["mac"] = 77
	ENC["shift_jis"] = 128
	ENC["hangeul"] = 129
	ENC["johab"] = 130
	ENC["gb2312"] = 134
	ENC["chinese big5"] = 136
	ENC["greek"] = 161
	ENC["turkish"] = 162
	ENC["Vietnamese"] = 163
	ENC["hebrew"] = 177
	ENC["arabic"] = 178
	ENC["baltic"] = 186
	ENC["russian"] = 204
	ENC["thai"] = 222
	ENC["east european"] = 238
	ENC["oem"] = 255

}

// Style represent subtitle"s styles.
type Style struct {
	Name      string
	FontName  string
	FontSize  int
	Color     [4]string //Primary, Secondary, Bord, Shadow map[string]string
	Bold      bool
	Italic    bool
	Underline bool
	StrikeOut bool
	Scale     [2]float64 // WIDTH, HEIGHT map[string]string
	Spacing   int
	Angle     int
	OpaqueBox bool
	Bord      float64
	Shadow    float64
	Alignment int
	Margin    [3]int // L, R, V map[string]string
	Encoding  int
}

// String get the generated Style as a String
func (sty *Style) String() string {
	return fmt.Sprintf(styleTemplate,
		sty.Name,
		sty.FontName, sty.FontSize,
		sty.Color[0], sty.Color[1], sty.Color[2], sty.Color[3],
		utils.Bool2str(sty.Bold), utils.Bool2str(sty.Italic),
		utils.Bool2str(sty.Underline), utils.Bool2str(sty.StrikeOut),
		sty.Scale[0], sty.Scale[1],
		sty.Spacing,
		sty.Angle,
		utils.Bool2Obox(sty.OpaqueBox),
		sty.Bord, sty.Shadow,
		sty.Alignment,
		sty.Margin[0], sty.Margin[1], sty.Margin[2],
		sty.Encoding)
}

// NewStyle create a new Style Struct with defaults
func NewStyle(name string) *Style {
	fontname := "Sans"
	// TODO: GO GENERATE FOR DIFERENT PLATFORM CASES
	if runtime.GOOS == "windows" {
		fontname = "Arial"
	}
	return &Style{
		Name:     name,
		FontName: fontname,
		FontSize: 35,
		Color: [4]string{
			"&H00FFFFFF", //Primary
			"&H000000FF", //Secondary
			"&H00000000", //Bord
			"&H00000000", //Shadow
		},
		Scale:     [2]float64{100, 100},
		Bord:      2,
		Alignment: ALIGN["bottom center"],
		Margin:    [3]int{10, 20, 10},
		Encoding:  ENC["default"],
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
	Margin    [3]int // L, R, V map[string]string
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
		d.Margin[0], d.Margin[1], d.Margin[2],
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
		// TODO: dummy video function
		// Dummy video
		framerate := asstime.FpsNtscFilm
		w, h := s.Resolution[0], s.Resolution[1]
		r, g, b := 0, 0, 0
		checkboard := "" // checkbord=True "c", checkboard=False ""
		//TODO: GET THE MAXIMUM TIME DIALOG
		frames := int(framerate) * 60 * 5
		s.VideoPath = fmt.Sprintf(
			dummyVideoTemplate,
			utils.Round(framerate, 3), frames, w, h, r, g, b, checkboard)
	}
	if s.VideoPath != "" {
		// TODO: dummy audio function
		// silence?, noise?
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

// FIXME: Aegisub could not narrow down character set to a single one

// Save write an SSA/ASS Subtitle Script.
func (s *Script) Save(fn string) {
	f, err := os.Create(fn)
	if err != nil {
		panic(fmt.Errorf("writer: failed saving subtitle file: %s", err))
	}
	defer f.Close()

	s.MetaFilename = fn

	n, err := f.WriteString(s.String())
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
