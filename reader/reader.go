// reader SSA/ASS Subtitle Script Parser
package reader

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Alquimista/eyecandy-go/utils"
)

// Video Subtitle Script Video Settings.
type Video struct {
	Path     string
	Zoom     float64
	Position int
	AR       float64
}

// Meta Subtitle Script Metadata, Informative Text.
type Meta struct {
	Filename, Title, OriginalScript, Translation, Timing string
}

// Resolution Subtitle Script Resolution.
type Resolution struct {
	X, Y int
}

// Represent the subtitle's lines.
type Dialog struct {
	Layer     int
	Start     string
	End       string
	StyleName string
	Style     *Style
	Actor     string
	MarginL   int
	MarginR   int
	MarginV   int
	Effect    string
	Text      string
	Comment   bool
}

// Font font used in a subtitle style.
type Font struct {
	Name string
	Size int
}

// ColorP palette of colors used in a subtitle style.
type ColorP struct {
	Primary, Secondary, Bord, Shadow string
}

// Scale Subtitle scale in percent.
// X: Width,  Y: Height
type Scale struct {
	X, Y float64
}

// Margin Subtitle Margin.
// L: Left,  R: Right, V: Vertical.
type Margin struct {
	L, R, V int
}

// Style represent Subtitle Style.
type Style struct {
	STR       string
	Name      string
	Font      *Font
	Color     *ColorP
	Bold      bool
	Italic    bool
	Underline bool
	StrikeOut bool
	Scale     *Scale
	Spacing   int
	Angle     int
	OpaqueBox bool
	Bord      float64
	Shadow    float64
	Alignment int
	Margin    *Margin
	Encoding  int
}

// Script SSA/ASS Subtitle Script.
type Script struct {
	Dialog            []*Dialog
	DialogWithComment []*Dialog
	Style             map[string]*Style
	StyleUsed         map[string]*Style
	Resolution        *Resolution
	Video             *Video
	Audio             string
	Meta              *Meta
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
	var tempDialog []*Dialog
	var playresx, playresy int
	var videopath string
	var videozoom, videoar float64
	var videopos int
	var title, originalscript, translation, timing string

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
				MarginL:   utils.Str2int(d[5]),
				MarginR:   utils.Str2int(d[6]),
				MarginV:   utils.Str2int(d[7]),
				Effect:    d[8],
				Text:      strings.TrimSpace(d[9]),
				Comment:   key == "comment",
			}
			tempDialog = append(tempDialog, dialog)
		case "style":
			s := strings.SplitN(value, ",", 23)
			style := &Style{
				Name: s[0],
				Font: &Font{
					Name: s[1],
					Size: utils.Str2int(s[2])},
				Color: &ColorP{
					Primary:   s[3],
					Secondary: s[4],
					Bord:      s[5],
					Shadow:    s[6],
				},
				Bold:      utils.Str2bool(s[7]),
				Italic:    utils.Str2bool(s[8]),
				Underline: utils.Str2bool(s[9]),
				StrikeOut: utils.Str2bool(s[10]),
				Scale: &Scale{
					X: utils.Str2float(s[11]),
					Y: utils.Str2float(s[12]),
				},
				Spacing:   utils.Str2int(s[13]),
				Angle:     utils.Str2int(s[14]),
				OpaqueBox: utils.Str2bool(s[15]),
				Bord:      utils.Str2float(s[16]),
				Shadow:    utils.Str2float(s[17]),
				Alignment: utils.Str2int(s[18]),
				Margin: &Margin{
					L: utils.Str2int(s[19]),
					R: utils.Str2int(s[20]),
					V: utils.Str2int(s[21]),
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
			videopath = value
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
		case "video_position":
			videopos = utils.Str2int(value)
		case "title":
			title = value
		case "original_script":
			originalscript = value
		case "translation":
			translation = value
		case "timing":
			timing = value
		default:
			continue
		}
	}

	script.Resolution = &Resolution{playresx, playresy}
	script.Video = &Video{
		Path:     videopath,
		Zoom:     videozoom,
		Position: videopos,
		AR:       videoar,
	}
	script.Meta = &Meta{
		Filename:       filename,
		Title:          title,
		OriginalScript: originalscript,
		Translation:    translation,
		Timing:         timing,
	}

	for _, d := range tempDialog {
		d.Style = script.Style[d.StyleName]
		if !d.Comment {
			// Replace the Dialog Style(Style Name)
			// for the style object of that Style Name
			script.StyleUsed[d.StyleName] = d.Style
			// Capture Subtitles Lines not Commented.
			script.Dialog = append(script.Dialog, d)
		}
		// Capture all the Subtitles Lines Commented or Not.
		script.DialogWithComment = append(script.Dialog, d)
	}

	return
}
