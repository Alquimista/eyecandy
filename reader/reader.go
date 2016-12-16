// Package reader read an SSA/ASS Subtitle Script
package reader

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Alquimista/eyecandy/utils"
)

var SSAColorLong = regexp.MustCompile(
	`&H([0-9A-Fa-f]{2})` + // alpha
		`([0-9A-Fa-f]{2})` + // blue component
		`([0-9A-Fa-f]{2})` + // green
		`([0-9A-Fa-f]{2})`) // red
// var SSAColor = regexp.MustCompile(
// 	`&H([0-9A-Fa-f]{2})` + // blue component
// 		`([0-9A-Fa-f]{2})` + // green
// 		`([0-9A-Fa-f]{2})&`) // red

// Dialog Represent the subtitle's lines.
type Dialog struct {
	Layer     int
	StartTime string
	EndTime   string
	StyleName string
	Style     *Style
	Actor     string
	Effect    string
	Text      string
	Tags      string
	Comment   bool
}

// DialogCollection collection of Dialog's in a SSA/ASS Script
type DialogCollection []*Dialog

// get list dialog's in a SSA/ASS Script
func (dlgs DialogCollection) get(commented bool) (dialogs DialogCollection) {
	for _, d := range dlgs {
		if d.Comment == commented {
			dialogs = append(dialogs, d)
		}
	}
	return dialogs
}

// Commented list only the commented dialog's in a SSA/ASS Script
func (dlgs DialogCollection) Commented() DialogCollection {
	return dlgs.get(true)
}

// NotCommented list only the not commented dialog's in a SSA/ASS Script
func (dlgs DialogCollection) NotCommented() DialogCollection {
	return dlgs.get(false)
}

// Style represent Subtitle Style.
type Style struct {
	Name      string
	FontName  string
	FontSize  int
	Color     [4]string //Primary, Secondary, Bord, Shadow
	Alpha     [4]int    //Primary, Secondary, Bord, Shadow
	Bold      bool
	Italic    bool
	Underline bool
	StrikeOut bool
	Scale     [2]float64 // WIDTH, HEIGHT
	Spacing   float64
	Angle     int
	OpaqueBox bool
	Bord      float64
	Shadow    float64
	Alignment int
	Margin    [3]int // L, R, V
	Encoding  int
}

// Script SSA/ASS Subtitle Script.
type Script struct {
	Dialog             DialogCollection
	Style              map[string]*Style
	StyleUsed          map[string]*Style
	Resolution         [2]int // WIDTH, HEIGHT
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

// parseStyle parse an SSA/ASS Subtitle Dialog.
func parseDialog(key, value string) *Dialog {
	// TODO ?: use sprintf
	d := strings.SplitN(value, ",", 10)
	return &Dialog{
		Layer:     utils.Str2int(d[0]),
		StartTime: d[1],
		EndTime:   d[2],
		StyleName: d[3],
		Actor:     d[4],
		Effect:    d[8],
		Text:      strings.TrimSpace(d[9]),
		Comment:   key == "comment",
	}
}

// parseStyle parse an SSA/ASS Subtitle Style.
func parseStyle(value string) *Style {
	// TODO ?: use sprintf
	sty := strings.SplitN(value, ",", 23)

	// 0: match, 1: alpha, 2: blue, 3: green, 4: red
	color1 := SSAColorLong.FindStringSubmatch(sty[3])
	c1 := "#" + color1[4] + color1[3] + color1[2]
	a1 := utils.Hex2int(color1[1])

	color2 := SSAColorLong.FindStringSubmatch(sty[4])
	c2 := "#" + color2[4] + color2[3] + color2[2]
	a2 := utils.Hex2int(color2[1])

	color3 := SSAColorLong.FindStringSubmatch(sty[5])
	c3 := "#" + color3[4] + color3[3] + color3[2]
	a3 := utils.Hex2int(color3[1])

	color4 := SSAColorLong.FindStringSubmatch(sty[6])
	c4 := "#" + color4[4] + color4[3] + color4[2]
	a4 := utils.Hex2int(color4[1])

	return &Style{
		Name:     sty[0],
		FontName: sty[1],
		FontSize: utils.Str2int(sty[2]),
		Color: [4]string{
			c1,  // Primary
			c2,  // Secondary
			c3,  // Bord
			c4}, // Shadow
		Alpha: [4]int{
			a1,  // Primary
			a2,  // Secondary
			a3,  // Bord
			a4}, // Shadow
		Bold:      utils.Str2bool(sty[7]),
		Italic:    utils.Str2bool(sty[8]),
		Underline: utils.Str2bool(sty[9]),
		StrikeOut: utils.Str2bool(sty[10]),
		Scale: [2]float64{
			utils.Str2float(sty[11]), // X
			utils.Str2float(sty[12]), // Y
		},
		Spacing:   utils.Str2float(sty[13]),
		Angle:     utils.Str2int(sty[14]),
		OpaqueBox: utils.Obox2bool(sty[15]),
		Bord:      utils.Str2float(sty[16]),
		Shadow:    utils.Str2float(sty[17]),
		Alignment: utils.Str2int(sty[18]),
		Margin: [3]int{
			utils.Str2int(sty[19]), // L
			utils.Str2int(sty[20]), // R
			utils.Str2int(sty[21]), // V
		},
	}
}

// Read parse and read an SSA/ASS Subtitle Script.
func Read(fn string) *Script {

	s := &Script{}
	f, err := os.Open(fn)
	if err != nil {
		panic(fmt.Errorf("reader: failed opening subtitle file: %s", err))
	}
	defer f.Close()

	s.Style = make(map[string]*Style)
	s.StyleUsed = make(map[string]*Style)
	var playresx, playresy int
	var videozoom, videoar float64

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, ";") ||
			strings.HasPrefix(line, "!:") ||
			strings.HasPrefix(line, "Format:") {
			continue
		}

		keyvalue := strings.SplitN(line, ":", 2)
		if len(keyvalue) != 2 {
			continue
		}
		key, value := keyvalue[0], keyvalue[1]
		key = strings.TrimSpace(key)
		key = strings.ToLower(key)
		key = strings.Replace(key, " ", "_", -1)
		value = strings.TrimSpace(value)

		switch key {
		case "dialogue", "comment":
			s.Dialog = append(s.Dialog, parseDialog(key, value))
		case "style":
			style := parseStyle(value)
			s.Style[style.Name] = style
		case "playresx":
			playresx = utils.Str2int(value)
		case "playresy":
			playresy = utils.Str2int(value)
		case "audio_uri", "audio_file":
			s.Audio = value
		case "video_file":
			s.VideoPath = value
		case "video_zoom_percent":
			videozoom = utils.Str2float(value)
		case "video_zoom":
			// Use "video_zoom_percent" key if present
			// else use "video_zoom" key
			if videozoom == 0 {
				zoom := strings.Replace(value, "%", "", -1)
				videozoom = utils.Str2float(zoom) / 100.0
			}
		case "video_aspect_ratio", "video_ar_value",
			"aegisub_video_aspect_ratio":
			ar := strings.Replace(value, "c", "", -1)
			numden := strings.SplitN(ar, ":", 2)
			if len(numden) == 2 {
				num, den := utils.Str2float(numden[0]),
					utils.Str2float(numden[1])
				videoar = num / den
			} else {
				videoar = utils.Str2float(ar)
			}
			s.VideoAR = videoar
		case "video_position":
			s.VideoPosition = utils.Str2int(value)
		case "title":
			s.MetaTitle = value
		case "original_script":
			s.MetaOriginalScript = value
		case "translation":
			s.MetaTranslation = value
		case "timing":
			s.MetaTiming = value
		default:
			continue
		}
	}

	s.Resolution = [2]int{playresx, playresy}
	s.VideoZoom = videozoom
	s.MetaFilename = fn

	// Get only the styles used in dialogs
	for _, d := range s.Dialog {
		sty, ok := s.Style[d.StyleName]
		if !ok {
			sty = s.Style["Default"]
		}
		d.Style = sty
		if !d.Comment {
			s.StyleUsed[d.StyleName] = sty
		}
	}

	return s
}
