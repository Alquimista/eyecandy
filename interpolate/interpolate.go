// Package interpolate Simple interpolation
package interpolate

import "math"

type Interp func(float64, float64, float64) float64

// Interpolate funcs
// https://play.golang.org/p/OKSM_h0zn-

func Linear(t, start, end float64) float64 {
	return t*(end-start) + start
}

// For gradient Color correction
// http://youtu.be/LKnqECcg6Gw

func LinearSqr(t, start, end float64) float64 {
	return math.Sqrt(Linear(t, math.Pow(start, 2), math.Pow(end, 2)))
}

func Cosine(t, start, end float64) float64 {
	t = 0.5 - (math.Cos(math.Pi*t) / 2)
	return Linear(t, start, end)
}

func Sine(t, start, end float64) float64 {
	t = math.Sin((math.Pi * t) / 2)
	return Linear(t, start, end)
}

func SmoothStep(t, start, end float64) float64 {
	t = math.Pow(t, 2) * (3 - (2 * t))
	return Linear(t, start, end)
}

func SmoothStepDouble(t, start, end float64) float64 {
	t = SmoothStep(t, start, end)
	return SmoothStep(t, start, end)
}

func Acceleration(t, start, end float64) float64 {
	t = math.Pow(t, 2)
	return Linear(t, start, end)
}

func CubicAcceleration(t, start, end float64) float64 {
	t = math.Pow(t, 3)
	return Linear(t, start, end)
}

func Deccelaration(t, start, end float64) float64 {
	t = 1 - math.Pow(1-t, 2)
	return Linear(t, start, end)
}

func CubicDeccelaration(t, start, end float64) float64 {
	t = 1 - math.Pow(1-t, 3)
	return Linear(t, start, end)
}

func Sigmoid(t, start, end float64) float64 {
	t = 1 / (1 + math.Exp(-t))
	return Linear(t, start, end)
}

//http://cubic-bezier.com
// NyuFx

func fact(x int) int {
	if x == 0 {
		return 1
	}
	return x * fact(x-1)
}

// Binomial coefficient
func binomial(i, n int) float64 {
	return float64(fact(n) / (fact(i) * fact(n-i)))
}

// Bernstein polynom
func bernstein(t float64, i, n int) float64 {
	return binomial(i, n) * math.Pow(t, float64(i)) * math.Pow((1-t), float64(n-i))
}

// Bezier Curve
func BezierCurve(t float64, p []float64) float64 {
	// Calculate coordinate
	n := len(p) - 1
	num := 0.0
	for i, position := range p {
		num += position * bernstein(t, i, n)
	}
	return num
}

func CustomCurve(t float64, curve []float64, start, end float64) float64 {
	t = BezierCurve(t, curve)
	return Linear(t, start, end)
}

// http://matthewlein.com/ceaser/

func Ease(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.25, 0.1, 0.25, 1}, start, end)
}

func EaseIn(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.42, 0, 1, 1}, start, end)
}

func EaseOut(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0, 0, 0.58, 1}, start, end)
}

func EaseInOut(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.420, 0.000, 0.580, 1.000}, start, end)

}

// Penner Equation (aproximated)

func EaseInQuad(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.550, 0.085, 0.680, 0.530}, start, end)
}

func EaseInCubic(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.550, 0.055, 0.675, 0.190}, start, end)
}

func EaseInQuart(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.895, 0.030, 0.685, 0.220}, start, end)
}

func EaseInQuint(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.755, 0.050, 0.855, 0.060}, start, end)

}
func EaseInSine(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.470, 0.000, 0.745, 0.715}, start, end)

}
func EaseInExpo(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.950, 0.050, 0.795, 0.035}, start, end)

}
func EaseInCirc(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.600, 0.040, 0.980, 0.335}, start, end)

}
func EaseOutQuad(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.250, 0.460, 0.450, 0.940}, start, end)

}
func EaseOutCubic(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.215, 0.610, 0.355, 1.000}, start, end)
}

func EaseOutQuart(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.165, 0.840, 0.440, 1.000}, start, end)

}
func EaseOutQuint(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.230, 1.000, 0.320, 1.000}, start, end)
}

func EaseOutSine(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.390, 0.575, 0.565, 1.000}, start, end)
}

func EaseOutExpo(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.190, 1.000, 0.220, 1.000}, start, end)
}

func EaseOutCirc(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.075, 0.820, 0.165, 1.000}, start, end)

}
func EaseInOutQuad(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.455, 0.030, 0.515, 0.955}, start, end)
}

func EaseInOutCubic(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.645, 0.045, 0.355, 1.000}, start, end)
}

func EaseInOutQuart(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.770, 0.000, 0.175, 1.000}, start, end)
}

func EaseInOutQuint(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.860, 0.000, 0.070, 1.000}, start, end)
}

func EaseInOutSine(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.445, 0.050, 0.550, 0.950}, start, end)
}

func EaseInOutExpo(t, start, end float64) float64 {
	return CustomCurve(t, []float64{1.000, 0.000, 0.000, 1.000}, start, end)
}

func EaseInOutCirc(t, start, end float64) float64 {
	return CustomCurve(t, []float64{0.785, 0.135, 0.150, 0.860}, start, end)
}

// KAFX Equations

func Backstart(t, start, end float64) float64 {
	return CustomCurve(
		t, []float64{0, 0, 0.2, -0.3, 0.6, 0.26, 1, 1}, start, end)
}

func Boing(t, start, end float64) float64 {
	return CustomCurve(
		t, []float64{0, 0, 0.42, 0.0, 0.58, 1.5, 1, 1}, start, end)
}

// IRange
func IRange(n int, start, end float64, f Interp) (rng []float64) {
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		rng = append(rng, f(t, start, end))
	}
	return rng
}

func ICircleRange(n int, f Interp) []float64 {
	return IRange(n+1, 0.0, 360.0, f)[:n]
}

func BezierCurveRange(n int, points []float64) (rng []float64) {
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		rng = append(rng, BezierCurve(t, points))
	}
	return rng
}
