// Package color provides color convention and useful functions
package color

import (
	"fmt"
	"regexp"

	"github.com/Alquimista/eyecandy/utils"
)

var hexColor = regexp.MustCompile(
	`(?:#)*([0-9A-Fa-f]{2})` + // red component
		`([0-9A-Fa-f]{2})` + // green
		`([0-9A-Fa-f]{2})`) // blue

var ssaColorLong = regexp.MustCompile(
	`&H([0-9A-Fa-f]{2})*` + // alpha
		`([0-9A-Fa-f]{2})` + // blue component
		`([0-9A-Fa-f]{2})` + // green
		`([0-9A-Fa-f]{2})`) // red

// SSALtoHEXAlpha convert SSA color format (&HAABBGGRR&) to hexcolor + alpha (decimal)
func SSALtoHEXAlpha(color string) (string, uint8) {
	// 0: match, 1: alpha, 2: blue, 3: green, 4: red
	clr := ssaColorLong.FindStringSubmatch(color)
	c := "#" + clr[4] + clr[3] + clr[2]
	return c, uint8(utils.Hex2int(clr[1]))
}

// HEXtoSSAL hexcolor to convert SSA color format (&HAABBGGRR&) with alpha = 0 (opaque)
func HEXtoSSAL(color string, alpha uint8) string {
	// 0: match, 1: red, 2: green, 3: blue
	clr := hexColor.FindStringSubmatch(color)
	return fmt.Sprintf("&H%02X%s%s%s", alpha, clr[3], clr[2], clr[1])
}

// HEXtoSSA hexcolor to convert SSA color format (for color tag: &HBBGGRR) ignoring alpha
func HEXtoSSA(color string) string {
	// 0: match, 1: red, 2: green, 3: blue
	clr := hexColor.FindStringSubmatch(color)
	// &HBBGGRR
	return fmt.Sprintf("&H%s%s%s&", clr[3], clr[2], clr[1])
}
