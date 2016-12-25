// Package color provides color convention and useful functions
package color

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Alquimista/eyecandy/utils"
)

var reColorHEX = regexp.MustCompile(
	`(?:#)*([0-9A-Fa-f]{2})` + // red component
		`([0-9A-Fa-f]{2})` + // green
		`([0-9A-Fa-f]{2})`) // blue

var reColorHEXSTRING = regexp.MustCompile(
	`(?:#)*([a-fA-F0-9]{2}\z|[a-fA-F0-9]{3}\z|[a-fA-F0-9]{6}\z)`)

var reSSAColor = regexp.MustCompile(
	`&H([0-9A-Fa-f]{2})*` + // alpha
		`([0-9A-Fa-f]{2})` + // blue component
		`([0-9A-Fa-f]{2})` + // green
		`([0-9A-Fa-f]{2})` + // red
		`(?:&)*`)

func hexToComponents(color string) []string {
	hexstring := reColorHEXSTRING.FindStringSubmatch(color)[1]
	if len(hexstring) == 3 {
		clr := []string{}
		for _, s := range hexstring {
			clr = append(clr, strings.Repeat(string(s), 2))
		}
		return clr
	} else if len(hexstring) == 2 {
		return []string{hexstring, hexstring, hexstring}
	}
	return reColorHEX.FindStringSubmatch(hexstring)[1:]
}

type Color struct {
	R, G, B, A uint8
}

func (c Color) RGB() (uint8, uint8, uint8) {
	return c.R, c.G, c.B
}

func (c Color) RGBA() (uint8, uint8, uint8, uint8) {
	return c.R, c.G, c.B, c.A
}

func (c Color) SSA() string {
	return fmt.Sprintf("&H%02X%02X%02X&", c.B, c.G, c.R)
}

func (c Color) SSAL() string {
	return fmt.Sprintf("&H%02X%02X%02X%02X", c.A, c.B, c.G, c.R)
}

func (c Color) HEX() string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

// NewRGB
func NewFromRGB(r, g, b uint8) *Color {
	return &Color{R: r, G: g, B: b}
}

// NewRGB
func NewFromRGBA(r, g, b, a uint8) *Color {
	return &Color{R: r, G: g, B: b, A: a}
}

// NewFromHex
func NewFromHex(hexc string) *Color {
	clr := hexToComponents(hexc)
	return &Color{
		R: uint8(utils.Hex2int(clr[0])),
		G: uint8(utils.Hex2int(clr[1])),
		B: uint8(utils.Hex2int(clr[2])),
	}
}

// NewFromHexAlpha
func NewFromHexAlpha(hexc string, a uint8) *Color {
	clr := hexToComponents(hexc)
	return &Color{
		R: uint8(utils.Hex2int(clr[0])),
		G: uint8(utils.Hex2int(clr[1])),
		B: uint8(utils.Hex2int(clr[2])),
		A: a,
	}
}

// NewFromHexAlpha
func NewFromSSA(ssac string) *Color {
	// 0: match, 1: alpha, 2: blue, 3: green, 4: red
	clr := reSSAColor.FindStringSubmatch(ssac)
	return &Color{
		R: uint8(utils.Hex2int(clr[4])),
		G: uint8(utils.Hex2int(clr[3])),
		B: uint8(utils.Hex2int(clr[2])),
		A: uint8(utils.Hex2int(clr[1])),
	}
}

var Aliceblue = NewFromHex("#f0f8ff")
var Antiquewhite = NewFromHex("#faebd7")
var Aqua = NewFromHex("#00ffff")
var Aquamarine = NewFromHex("#7fffd4")
var Azure = NewFromHex("#f0ffff")
var Beige = NewFromHex("#f5f5dc")
var Bisque = NewFromHex("#ffe4c4")
var Black = NewFromHex("#000000")
var Blanchedalmond = NewFromHex("#ffebcd")
var Blue = NewFromHex("#0000ff")
var Blueviolet = NewFromHex("#8a2be2")
var Brown = NewFromHex("#a52a2a")
var Burlywood = NewFromHex("#deb887")
var Cadetblue = NewFromHex("#5f9ea0")
var Chartreuse = NewFromHex("#7fff00")
var Chocolate = NewFromHex("#d2691e")
var Coral = NewFromHex("#ff7f50")
var Cornflowerblue = NewFromHex("#6495ed")
var Cornsilk = NewFromHex("#fff8dc")
var Crimson = NewFromHex("#dc143c")
var Cyan = NewFromHex("#00ffff")
var Darkblue = NewFromHex("#00008b")
var Darkcyan = NewFromHex("#008b8b")
var Darkgoldenrod = NewFromHex("#b8860b")
var Darkgray = NewFromHex("#a9a9a9")
var Darkgrey = NewFromHex("#a9a9a9")
var Darkgreen = NewFromHex("#006400")
var Darkkhaki = NewFromHex("#bdb76b")
var Darkmagenta = NewFromHex("#8b008b")
var Darkolivegreen = NewFromHex("#556b2f")
var Darkorange = NewFromHex("#ff8c00")
var Darkorchid = NewFromHex("#9932cc")
var Darkred = NewFromHex("#8b0000")
var Darksalmon = NewFromHex("#e9967a")
var Darkseagreen = NewFromHex("#8fbc8f")
var Darkslateblue = NewFromHex("#483d8b")
var Darkslategray = NewFromHex("#2f4f4f")
var Darkslategrey = NewFromHex("#2f4f4f")
var Darkturquoise = NewFromHex("#00ced1")
var Darkviolet = NewFromHex("#9400d3")
var Deeppink = NewFromHex("#ff1493")
var Deepskyblue = NewFromHex("#00bfff")
var Dimgray = NewFromHex("#696969")
var Dimgrey = NewFromHex("#696969")
var Dodgerblue = NewFromHex("#1e90ff")
var Firebrick = NewFromHex("#b22222")
var Floralwhite = NewFromHex("#fffaf0")
var Forestgreen = NewFromHex("#228b22")
var Fuchsia = NewFromHex("#ff00ff")
var Gainsboro = NewFromHex("#dcdcdc")
var Ghostwhite = NewFromHex("#f8f8ff")
var Gold = NewFromHex("#ffd700")
var Goldenrod = NewFromHex("#daa520")
var Gray = NewFromHex("#808080")
var Grey = NewFromHex("#808080")
var Green = NewFromHex("#008000")
var Greenyellow = NewFromHex("#adff2f")
var Honeydew = NewFromHex("#f0fff0")
var Hotpink = NewFromHex("#ff69b4")
var Indianred = NewFromHex("#cd5c5c")
var Indigo = NewFromHex("#4b0082")
var Ivory = NewFromHex("#fffff0")
var Khaki = NewFromHex("#f0e68c")
var Lavender = NewFromHex("#e6e6fa")
var Lavenderblush = NewFromHex("#fff0f5")
var Lawngreen = NewFromHex("#7cfc00")
var Lemonchiffon = NewFromHex("#fffacd")
var Lightblue = NewFromHex("#add8e6")
var Lightcoral = NewFromHex("#f08080")
var Lightcyan = NewFromHex("#e0ffff")
var Lightgoldenrodyellow = NewFromHex("#fafad2")
var Lightgray = NewFromHex("#d3d3d3")
var Lightgrey = NewFromHex("#d3d3d3")
var Lightgreen = NewFromHex("#90ee90")
var Lightpink = NewFromHex("#ffb6c1")
var Lightsalmon = NewFromHex("#ffa07a")
var Lightseagreen = NewFromHex("#20b2aa")
var Lightskyblue = NewFromHex("#87cefa")
var Lightslategray = NewFromHex("#778899")
var Lightslategrey = NewFromHex("#778899")
var Lightsteelblue = NewFromHex("#b0c4de")
var Lightyellow = NewFromHex("#ffffe0")
var Lime = NewFromHex("#00ff00")
var LimegreAliceblueen = NewFromHex("#32cd32")
var Linen = NewFromHex("#faf0e6")
var Magenta = NewFromHex("#ff00ff")
var Maroon = NewFromHex("#800000")
var Mediumaquamarine = NewFromHex("#66cdaa")
var Mediumblue = NewFromHex("#0000cd")
var Mediumorchid = NewFromHex("#ba55d3")
var Mediumpurple = NewFromHex("#9370d8")
var Mediumseagreen = NewFromHex("#3cb371")
var Mediumslateblue = NewFromHex("#7b68ee")
var Mediumspringgreen = NewFromHex("#00fa9a")
var Mediumturquoise = NewFromHex("#48d1cc")
var Mediumvioletred = NewFromHex("#c71585")
var Midnightblue = NewFromHex("#191970")
var Mintcream = NewFromHex("#f5fffa")
var Mistyrose = NewFromHex("#ffe4e1")
var Moccasin = NewFromHex("#ffe4b5")
var Navajowhite = NewFromHex("#ffdead")
var Navy = NewFromHex("#000080")
var Oldlace = NewFromHex("#fdf5e6")
var Olive = NewFromHex("#808000")
var Olivedrab = NewFromHex("#6b8e23")
var Orange = NewFromHex("#ffa500")
var Orangered = NewFromHex("#ff4500")
var Orchid = NewFromHex("#da70d6")
var Palegoldenrod = NewFromHex("#eee8aa")
var Palegreen = NewFromHex("#98fb98")
var Paleturquoise = NewFromHex("#afeeee")
var Palevioletred = NewFromHex("#d87093")
var Papayawhip = NewFromHex("#ffefd5")
var Peachpuff = NewFromHex("#ffdab9")
var Peru = NewFromHex("#cd853f")
var Pink = NewFromHex("#ffc0cb")
var Plum = NewFromHex("#dda0dd")
var Powderblue = NewFromHex("#b0e0e6")
var Purple = NewFromHex("#800080")
var Red = NewFromHex("#ff0000")
var Rosybrown = NewFromHex("#bc8f8f")
var Royalblue = NewFromHex("#4169e1")
var Saddlebrown = NewFromHex("#8b4513")
var Salmon = NewFromHex("#fa8072")
var Sandybrown = NewFromHex("#f4a460")
var Seagreen = NewFromHex("#2e8b57")
var Seashell = NewFromHex("#fff5ee")
var Sienna = NewFromHex("#a0522d")
var Silver = NewFromHex("#c0c0c0")
var Skyblue = NewFromHex("#87ceeb")
var Slateblue = NewFromHex("#6a5acd")
var Slategray = NewFromHex("#708090")
var Slategrey = NewFromHex("#708090")
var Snow = NewFromHex("#fffafa")
var Springgreen = NewFromHex("#00ff7f")
var Steelblue = NewFromHex("#4682b4")
var Tan = NewFromHex("#d2b48c")
var Teal = NewFromHex("#008080")
var Thistle = NewFromHex("#d8bfd8")
var Tomato = NewFromHex("#ff6347")
var Turquoise = NewFromHex("#40e0d0")
var Violet = NewFromHex("#ee82ee")
var Wheat = NewFromHex("#f5deb3")
var White = NewFromHex("#ffffff")
var Whitesmoke = NewFromHex("#f5f5f5")
var Yellow = NewFromHex("#ffff00")
var Yellowgreen = NewFromHex("#9acd32")
