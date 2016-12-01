// reader SSA/ASS Subtitle Script Parser
package reader

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Alquimista/eyecandy-go/utils"
)

type DialogCollection []*dialog

// video Subtitle Script Video Settings.
type video struct {
	Path     string
	Zoom     float64
	Position int
	AR       float64
}

// meta Subtitle Script Metadata, Informative Text.
type meta struct {
	Filename, Title, OriginalScript, Translation, Timing string
}

// resolution Subtitle Script Resolution.
type resolution struct {
	X, Y int
}

// dialog Represent the subtitle's lines.
type dialog struct {
	Layer     int
	Start     string
	End       string
	styleName string
	Style     *style
	Actor     string
	Margin    *margin
	Effect    string
	Text      string
	Comment   bool
}

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

// font font used in a subtitle style.
type font struct {
	Name string
	Size int
}

// colorp palette of colors used in a subtitle style.
type colorp struct {
	Primary, Secondary, Bord, Shadow string
}

// scale Subtitle scale in percent.
// X: Width,  Y: Height
type scale struct {
	X, Y float64
}

// margin Subtitle Margin.
// L: Left,  R: Right, V: Vertical.
type margin struct {
	L, R, V int
}

// style represent Subtitle Style.
type style struct {
	Name      string
	Font      *font
	Color     *colorp
	Bold      bool
	Italic    bool
	Underline bool
	StrikeOut bool
	Scale     *scale
	Spacing   int
	Angle     int
	OpaqueBox bool
	Bord      float64
	Shadow    float64
	Alignment int
	Margin    *margin
	Encoding  int
}

// script SSA/ASS Subtitle Script.
type script struct {
	Dialog     DialogCollection
	Style      map[string]*style
	StyleUsed  map[string]*style
	Resolution *resolution
	Video      *video
	Audio      string
	Meta       *meta
}

// Read parse and read an SSA/ASS Subtitle Script.
func Read(filename string) (script script) {

	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Errorf("reader: failed opening subtitle file: %s", err))
	}
	defer file.Close()

	script.Style = make(map[string]*style)
	script.StyleUsed = make(map[string]*style)
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
			dialog := &dialog{
				Layer: utils.Str2int(d[0]),
				Start: d[1],
				End:   d[2],
				// Style:     &style{},
				styleName: d[3],
				Actor:     d[4],
				Margin: &margin{
					L: utils.Str2int(d[5]),
					R: utils.Str2int(d[6]),
					V: utils.Str2int(d[7]),
				},
				Effect:  d[8],
				Text:    strings.TrimSpace(d[9]),
				Comment: key == "comment",
			}
			script.Dialog = append(script.Dialog, dialog)
		case "style":
			s := strings.SplitN(value, ",", 23)
			style := &style{
				Name: s[0],
				Font: &font{
					Name: s[1],
					Size: utils.Str2int(s[2])},
				Color: &colorp{
					Primary:   s[3],
					Secondary: s[4],
					Bord:      s[5],
					Shadow:    s[6],
				},
				Bold:      utils.Str2bool(s[7]),
				Italic:    utils.Str2bool(s[8]),
				Underline: utils.Str2bool(s[9]),
				StrikeOut: utils.Str2bool(s[10]),
				Scale: &scale{
					X: utils.Str2float(s[11]),
					Y: utils.Str2float(s[12]),
				},
				Spacing:   utils.Str2int(s[13]),
				Angle:     utils.Str2int(s[14]),
				OpaqueBox: utils.Str2bool(s[15]),
				Bord:      utils.Str2float(s[16]),
				Shadow:    utils.Str2float(s[17]),
				Alignment: utils.Str2int(s[18]),
				Margin: &margin{
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

	script.Resolution = &resolution{playresx, playresy}
	script.Video = &video{
		Path:     videopath,
		Zoom:     videozoom,
		Position: videopos,
		AR:       videoar,
	}
	script.Meta = &meta{
		Filename:       filename,
		Title:          title,
		OriginalScript: originalscript,
		Translation:    translation,
		Timing:         timing,
	}

	// Get only the styles used in dialogs
	for _, d := range script.Dialog {
		style, ok := script.Style[d.styleName]
		if !ok {
			style = script.Style["Default"]
		}
		d.Style = style
		if !d.Comment {
			script.StyleUsed[d.styleName] = style
		}
	}

	return
}
