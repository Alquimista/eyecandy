// Package asstime provide time conversions an framerate functions
package asstime

// FpsNtscFilm Frame per second rate NTSC film standard (23.976)
const FpsNtscFilm float64 = float64(24000) / float64(1001)

// FpsNtsc Frame per second rate NTSC standard (30)
const FpsNtsc float64 = float64(30000) / float64(1001)

// FpsNtscDouble Frame per second rate NTSC Double standard (60)
const FpsNtscDouble float64 = float64(60000) / float64(1001)

// FpsNtscQuad Frame per second rate NTSC Quad standard (120)
const FpsNtscQuad float64 = float64(120000) / float64(1001)

// FpsFilm Frame per second rate Film standard
const FpsFilm int = 24

// FpsPal Frame per second rate PAL standard
const FpsPal int = 25

// FpsPalDouble Frame per second rate PAL Double standard (50)
const FpsPalDouble int = 50

// FpsPalQuad Frame per second rate PAL Quad standard (100)
const FpsPalQuad int = 100
