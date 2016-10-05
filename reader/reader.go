// reader SSA/ASS Subtitle Script Parser
package reader

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	styleName string
	Style     Style
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

// Color palette of colors used in a subtitle style.
type Color struct {
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
	Font      Font
	Color     Color
	Bold      bool
	Italic    bool
	Underline bool
	StrikeOut bool
	Scale     Scale
	Spacing   int
	Angle     int
	OpaqueBox bool
	Bord      float64
	Shadow    float64
	Alignment int
	Margin    Margin
	Encoding  int
}

// Script SSA/ASS Subtitle Script.
type Script struct {
	Dialog            []Dialog
	DialogWithComment []Dialog
	Style             map[string]Style
	StyleUsed         map[string]Style
	Resolution        Resolution
	Video             Video
	Audio             string
	Meta              Meta
}

// str2int convert a string to an integer.
func str2int(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Errorf("reader: failed parsing number: %s", err))
	}
	return i
}

// str2bool convert a string to a boolean.
func str2bool(s string) bool {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Errorf("reader: failed parsing boolean: %s", err))
	}
	return i != 0
}

// str2bool convert a string to an float.
func str2float(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(fmt.Errorf("reader: failed parsing number: %s", err))
	}
	return f
}

// Read parse and read an SSA/ASS Subtitle Script.
func Read(filename string) (script Script) {

	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Errorf("reader: failed opening subtitle file: %s", err))
	}
	defer file.Close()

	script.Style = make(map[string]Style)
	var tempDialog []Dialog
	script.StyleUsed = make(map[string]Style)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, ";") ||
			strings.HasPrefix(line, "!:") ||
			strings.HasPrefix(line, "Format:") {
			continue
		}

		splitline := strings.SplitN(line, ":", 2)
		if len(splitline) != 2 {
			continue
		}
		key, value := splitline[0], splitline[1]
		key = strings.ToLower(key)
		key = strings.Replace(key, " ", "_", -1)
		value = strings.TrimSpace(value)

		switch key {
		case "dialogue", "comment":
			d := strings.SplitN(value, ",", 10)
			dialog := Dialog{
				Layer:     str2int(d[0]),
				Start:     d[1],
				End:       d[2],
				styleName: d[3],
				Actor:     d[4],
				MarginL:   str2int(d[5]),
				MarginR:   str2int(d[6]),
				MarginV:   str2int(d[7]),
				Effect:    d[8],
				Text:      strings.TrimSpace(d[9]),
				Comment:   key == "comment",
			}
			tempDialog = append(tempDialog, dialog)
		case "style":
			s := strings.SplitN(value, ",", 23)
			style := Style{
				Name: s[0],
				Font: Font{
					Name: s[1],
					Size: str2int(s[2])},
				Color: Color{
					Primary:   s[3],
					Secondary: s[4],
					Bord:      s[5],
					Shadow:    s[6],
				},
				Bold:      str2bool(s[7]),
				Italic:    str2bool(s[8]),
				Underline: str2bool(s[9]),
				StrikeOut: str2bool(s[10]),
				Scale: Scale{
					X: str2float(s[11]),
					Y: str2float(s[12]),
				},
				Spacing:   str2int(s[13]),
				Angle:     str2int(s[14]),
				OpaqueBox: str2bool(s[15]),
				Bord:      str2float(s[16]),
				Shadow:    str2float(s[17]),
				Alignment: str2int(s[18]),
				Margin: Margin{
					L: str2int(s[19]),
					R: str2int(s[20]),
					V: str2int(s[21]),
				},
				Encoding: str2int(s[22]),
			}
			script.Style[style.Name] = style
		case "playresx":
			script.Resolution.X = str2int(value)
		case "playresy":
			script.Resolution.Y = str2int(value)
		case "video_file":
			script.Video.Path = value
		case "audio_uri", "audio_file":
			script.Audio = value
		case "video_zoom_percent":
			script.Video.Zoom = str2float(value)
		case "video_zoom":
			// Use "video_zoom_percent" key if present
			// else use "video_zoom" key
			if script.Video.Zoom == 0 {
				zoom := strings.Replace(value, "%", "", -1)
				script.Video.Zoom = str2float(zoom) / 100.0
			}
		case "video_aspect_ratio", "video_ar_value",
			"aegisub_video_aspect_ratio":
			ar := strings.Replace(value, "c", "", -1)
			arsplit := strings.SplitN(ar, ":", 2)
			if len(arsplit) == 2 {
				num, den := str2float(arsplit[0]), str2float(arsplit[1])
				script.Video.AR = num / den
			} else {
				script.Video.AR = str2float(ar)
			}
		case "title":
			script.Meta.Title = value
		case "original_script":
			script.Meta.OriginalScript = value
		case "translation":
			script.Meta.Translation = value
		case "timing":
			script.Meta.Timing = value
		case "video_position":
			script.Video.Position = str2int(value)
		default:
			continue
		}
	}

	for _, d := range tempDialog {
		d.Style = script.Style[d.styleName]
		if !d.Comment {
			// Replace the Dialog Style(Style Name)
			// for the style object of that Style Name
			script.StyleUsed[d.styleName] = d.Style
			// Capture Subtitles Lines not Commented.
			script.Dialog = append(script.Dialog, d)
		}
		// Capture all the Subtitles Lines Commented or Not.
		script.DialogWithComment = append(script.Dialog, d)
	}

	return
}
