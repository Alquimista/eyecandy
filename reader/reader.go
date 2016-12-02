// reader SSA/ASS Subtitle Script Parser
package reader

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Alquimista/eyecandy-go/utils"
)

// Dialog Represent the subtitle's lines.
type Dialog struct {
	Layer     int
	Start     string
	End       string
	StyleName string
	Style     *Style
	Actor     string
	Effect    string
	Text      string
	Tags      string
	Margin    [3]int // L, R, V map[string]string
	Comment   bool
}

type DialogCollection []*Dialog

func (dlgs DialogCollection) GetAll() DialogCollection {
	return dlgs
}

func (dlgs DialogCollection) GetCommented() DialogCollection {
	var dialogs DialogCollection
	for _, d := range dlgs {
		if d.Comment {
			dialogs = append(dialogs, d)
		}
	}
	return dialogs
}

func (dlgs DialogCollection) GetNotCommented() DialogCollection {
	var dialogs DialogCollection
	for _, d := range dlgs {
		if !d.Comment {
			dialogs = append(dialogs, d)
		}
	}
	return dialogs
}

// Style represent Subtitle Style.
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

// Script SSA/ASS Subtitle Script.
type Script struct {
	Dialog             DialogCollection
	Style              map[string]*Style
	StyleUsed          map[string]*Style
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

// Read parse and read an SSA/ASS Subtitle Script.
func Read(filename string) (script Script) {

	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Errorf("reader: failed opening subtitle file: %s", err))
	}
	defer file.Close()

	script.Style = make(map[string]*Style)
	script.StyleUsed = make(map[string]*Style)
	var playresx, playresy int
	var videozoom, videoar float64

	scanner := bufio.NewScanner(file)
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
		key = strings.ToLower(key)
		key = strings.Replace(key, " ", "_", -1)
		value = strings.TrimSpace(value)

		switch key {
		case "dialogue", "comment":
			d := strings.SplitN(value, ",", 10)
			dialog := &Dialog{
				Layer:     utils.Str2int(d[0]),
				Start:     d[1],
				End:       d[2],
				StyleName: d[3],
				Actor:     d[4],
				Margin: [3]int{
					utils.Str2int(d[5]),
					utils.Str2int(d[6]),
					utils.Str2int(d[7]),
				},
				Effect:  d[8],
				Text:    strings.TrimSpace(d[9]),
				Comment: key == "comment",
			}
			script.Dialog = append(script.Dialog, dialog)
		case "style":
			s := strings.SplitN(value, ",", 23)
			style := &Style{
				Name:      s[0],
				FontName:  s[1],
				FontSize:  utils.Str2int(s[2]),
				Color:     [4]string{s[3], s[4], s[5], s[6]},
				Bold:      utils.Str2bool(s[7]),
				Italic:    utils.Str2bool(s[8]),
				Underline: utils.Str2bool(s[9]),
				StrikeOut: utils.Str2bool(s[10]),
				Scale: [2]float64{
					utils.Str2float(s[11]),
					utils.Str2float(s[12]),
				},
				Spacing:   utils.Str2int(s[13]),
				Angle:     utils.Str2int(s[14]),
				OpaqueBox: utils.Obox2bool(s[15]),
				Bord:      utils.Str2float(s[16]),
				Shadow:    utils.Str2float(s[17]),
				Alignment: utils.Str2int(s[18]),
				Margin: [3]int{
					utils.Str2int(s[19]),
					utils.Str2int(s[20]),
					utils.Str2int(s[21]),
				},
				Encoding: utils.Str2int(s[22]),
			}
			script.Style[style.Name] = style
		case "playresx":
			playresx = utils.Str2int(value)
		case "playresy":
			playresy = utils.Str2int(value)
		case "audio_uri", "audio_file":
			script.Audio = value
		case "video_file":
			script.VideoPath = value
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
			script.VideoAR = videoar
		case "video_position":
			script.VideoPosition = utils.Str2int(value)
		case "title":
			script.MetaTitle = value
		case "original_script":
			script.MetaOriginalScript = value
		case "translation":
			script.MetaTranslation = value
		case "timing":
			script.MetaTiming = value
		default:
			continue
		}
	}

	script.Resolution = [2]int{playresx, playresy}
	script.VideoZoom = videozoom
	script.MetaFilename = filename

	// Get only the styles used in dialogs
	for _, d := range script.Dialog {
		style, ok := script.Style[d.StyleName]
		if !ok {
			style = script.Style["Default"]
		}
		d.Style = style
		if !d.Comment {
			script.StyleUsed[d.StyleName] = style
		}
	}

	return
}
