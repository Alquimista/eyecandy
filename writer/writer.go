// writer SSA/ASS Subtitle Script Writer
// package writer
package writer

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/Alquimista/eyecandy-go/asstime"
	"github.com/Alquimista/eyecandy-go/utils"
)

const DUMMY_VIDEO_TEMPLATE string = "?dummy:%.6f:%d:%d:%d:%d:%d:%d%s:"
const DUMMY_AUDIO_TEMPLATE string = "dummy-audio:silence?sr=44100&bd=16&" +
	"ch=1&ln=396900000:"
const STYLE_FORMAT string = "Format: Name, Fontname, Fontsize, " +
	"PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, " +
	"Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, " +
	"BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, " +
	"Encoding"
const STYLE_TEMPLATE string = "Style: %s,%s,%d,%s,%s,%s,%s,%s,%s,%s,%s,%.4f," +
	"%.4f,%d,%d,%d,%.4f,%.4f,%d,%04d,%04d,%04d,%d"
const DIALOG_FORMAT string = "Format: Layer, Start, End, Style, Name, " +
	"MarginL, MarginR, MarginV, Effect, Text"
const DIALOG_TEMPLATE string = "%s: %d,%s,%s,%s,%s,%04d,%04d,%04d,%s,%s"
const SCRIPT_TEMPLATE string = `[Script Info]
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

var ALIGN = make(map[string]int)
var ENC = make(map[string]int)

func init() {

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

func (s *Style) toString() string {
	opaquebox := 3
	if !s.OpaqueBox {
		opaquebox = 0
	}
	return fmt.Sprintf(STYLE_TEMPLATE,
		s.Name,
		s.FontName, s.FontSize,
		s.Color[0], s.Color[1], s.Color[2], s.Color[3],
		utils.Bool2str(s.Bold), utils.Bool2str(s.Italic),
		utils.Bool2str(s.Underline), utils.Bool2str(s.StrikeOut),
		s.Scale[0], s.Scale[1],
		s.Spacing,
		s.Angle,
		opaquebox,
		s.Bord, s.Shadow,
		s.Alignment,
		s.Margin[0], s.Margin[1], s.Margin[2],
		s.Encoding)
}

func NewStyle(name string) *Style {
	fontname := "Sans"
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

func (d *Dialog) toString() string {
	text := d.Text
	key := "Dialogue"
	if d.Comment == true {
		key = "Comment"
	}
	if d.Tags != "" {
		text = "{" + d.Tags + "}" + d.Text
	}
	return fmt.Sprintf(DIALOG_TEMPLATE,
		key,
		d.Layer,
		d.Start, d.End,
		d.StyleName, d.Actor,
		d.Margin[0], d.Margin[1], d.Margin[2],
		d.Effect,
		text)
}

func NewDialog(text string) *Dialog {
	return &Dialog{
		StyleName: "Default",
		Start:     "0:00:00.00", End: "0:00:05.00",
		Text: text}
}

// Script SSA/ASS Subtitle Script.
type script struct {
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

func (script *script) GetStyle(name string) *Style {
	style, ok := script.Style[name]
	if !ok {
		style = script.Style["Default"]
		style.Name = name
	}
	return style
}

func (script *script) styleExists(name string) bool {
	_, ok := script.Style[name]
	return ok
}

func (script *script) AddStyle(s *Style) {
	if !script.styleExists(s.Name) {
		script.Style[s.Name] = s
	}
}

func (s *script) AddDialog(d *Dialog) {
	if d.Text != "" {
		s.Dialog = append(s.Dialog, d)
	}
}

func (s *script) ToString() string {

	defaultSTY := NewStyle("Default")
	s.AddStyle(defaultSTY)
	if len(s.Dialog) == 0 {
		defaultDLG := NewDialog("EyecandyFX")
		s.AddDialog(defaultDLG)
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
		framerate := asstime.FPS_NTSC_FILM
		w, h := s.Resolution[0], s.Resolution[1]
		r, g, b := 0, 0, 0
		checkboard := "" // checkbord=True "c", checkboard=False ""
		//TODO: GET THE MAXIMUM TIME DIALOG
		frames := int(framerate) * 60 * 5
		s.VideoPath = fmt.Sprintf(
			DUMMY_VIDEO_TEMPLATE,
			utils.Round(framerate, 3), frames, w, h, r, g, b, checkboard)
	}
	if s.VideoPath != "" {
		// TODO: dummy audio function
		// silence?, noise?
		if strings.HasPrefix(s.VideoPath, "?dummy") {
			s.Audio = DUMMY_AUDIO_TEMPLATE
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
		dialogs.WriteString(d.toString() + "\n")
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
				i = i + 1
			} else {
				i = 1
			}
			if sty.Name == sname && i == 1 {
				styles.WriteString(sty.toString() + "\n")
			}
		}
	}

	return fmt.Sprintf(SCRIPT_TEMPLATE,
		s.Comment,
		s.MetaTitle, s.MetaOriginalScript,
		s.MetaTranslation, s.MetaTiming,
		s.Resolution[0], s.Resolution[1],
		s.VideoPath, s.VideoAR, s.VideoZoom, s.VideoPosition,
		s.Audio,
		STYLE_FORMAT, styles.String(),
		// STYLE_FORMAT, strings.Join(styles, "\n"),
		DIALOG_FORMAT, dialogs.String())

}

// Write an SSA/ASS Subtitle Script.
// FIXME: Aegisub could not narrow down character set to a single one
func (s *script) Save(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(fmt.Errorf("writer: failed saving subtitle file: %s", err))
	}
	defer file.Close()

	s.MetaFilename = filename
	// w := bufio.NewWriter(file)
	// n, err := io.WriteString(file, s.ToString())
	// n, err := w.WriteString(s.ToString())
	n, err := file.WriteString(s.ToString())
	if err != nil {
		fmt.Println(n, err)
	}
	// w.Flush()

	// save changes
	err = file.Sync()
}

func NewScript() *script {
	return &script{
		Comment:   "Script generated by Eyecandy",
		MetaTitle: "Default Eyecandy file",
		VideoZoom: 0.75,
		Style:     map[string]*Style{},
	}
}
