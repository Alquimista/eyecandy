// Package asstime provide time conversions an framerate functions
package asstime

import (
	"fmt"
	"math"
	"regexp"

	"github.com/Alquimista/eyecandy/utils"
)

const (
	// FpsNtscFilm Frame per second rate NTSC film standard (23.976)
	FpsNtscFilm float64 = float64(24000) / float64(1001)
	// FpsNtsc Frame per second rate NTSC standard (30)
	FpsNtsc float64 = float64(30000) / float64(1001)
	// FpsNtscDouble Frame per second rate NTSC Double standard (60)
	FpsNtscDouble float64 = float64(60000) / float64(1001)
	// FpsNtscQuad Frame per second rate NTSC Quad standard (120)
	FpsNtscQuad float64 = float64(120000) / float64(1001)
	// FpsFilm Frame per second rate Film standard
	FpsFilm float64 = 24.0
	// FpsPal Frame per second rate PAL standard
	FpsPal float64 = 25.0
	// FpsPalDouble Frame per second rate PAL Double standard (50)
	FpsPalDouble float64 = 50.0
	// FpsPalQuad Frame per second rate PAL Quad standard (100)
	FpsPalQuad int = 100.0
)

const (
	Millisecond = 1
	Centisecond = 10 * Millisecond
	Second      = 1000 * Millisecond
	Minute      = 60 * Second
	Hour        = 60 * Minute
)

// TODO: Try scanf
var ReSSAfmt = regexp.MustCompile(`(\d):(\d+):(\d+).(\d+)`)

// MStoFrames Convert Frames to Milliseconds
func MStoFrames(ms int, framerate float64) int {
	return int(math.Ceil(framerate * (float64(ms) / float64(1000))))
}

// FramesToMS Convert Frames to Milliseconds
func FramesToMS(frames int, framerate float64) int {
	return int((float64(frames) / float64(framerate)) * Second)
}

func MStoSSA(milli int) string {
	sec, ms := utils.DivMod(milli, 1000)
	min, s := utils.DivMod(sec, 60)
	h, m := utils.DivMod(min, 60)
	cs, _ := utils.DivMod(ms, 10)
	return fmt.Sprintf("%01d:%02d:%02d.%02d", h, m, s, cs)
}

func SSAtoMS(t string) int {
	h, m, s, cs := SSAtoSplit(t)
	return (h*Hour + m*Minute + s*Second + cs*Centisecond)
}

func SSAtoSplit(t string) (h, m, s, cs int) {
	//H:MM:SS.CC (H=Hour, M=Minute, S=Second, C=centisecond)
	tm := ReSSAfmt.FindStringSubmatch(t)
	return utils.Str2int(tm[1]), utils.Str2int(tm[2]),
		utils.Str2int(tm[3]), utils.Str2int(tm[4])
}
