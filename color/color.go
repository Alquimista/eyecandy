// Package color
package color

import (
	"fmt"
	"regexp"

	"github.com/Alquimista/eyecandy/utils"
)

var HEXColor = regexp.MustCompile(
	`(?:#)*([0-9A-Fa-f]{2})` + // red component
		`([0-9A-Fa-f]{2})` + // green
		`([0-9A-Fa-f]{2})`) // blue

var SSAColorLong = regexp.MustCompile(
	`&H([0-9A-Fa-f]{2})*` + // alpha
		`([0-9A-Fa-f]{2})` + // blue component
		`([0-9A-Fa-f]{2})` + // green
		`([0-9A-Fa-f]{2})`) // red

func SSALtoHEXAlpha(color string) (string, uint8) {
	// 0: match, 1: alpha, 2: blue, 3: green, 4: red
	clr := SSAColorLong.FindStringSubmatch(color)
	c := "#" + clr[4] + clr[3] + clr[2]
	return c, uint8(utils.Hex2int(clr[1]))
}

func HEXtoSSAL(color string, alpha uint8) string {
	// 0: match, 1: red, 2: green, 3: blue
	clr := HEXColor.FindStringSubmatch(color)
	return fmt.Sprintf("&H%02X%s%s%s", alpha, clr[3], clr[2], clr[1])
}

func HEXtoSSA(color string) string {
	// 0: match, 1: red, 2: green, 3: blue
	clr := HEXColor.FindStringSubmatch(color)
	return fmt.Sprintf("&H%s%s%s&", clr[3], clr[2], clr[1])
}
