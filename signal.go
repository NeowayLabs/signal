package signal

import "math"

type (
	// Discrete signal
	Discrete []float64

	// Continuous does not exists inside the computer
	// but let the bullshit continue.
	Continuous func(x float64) float64
)

// Almost asserts that x is close to y with some precision ε.
// Mathematically speaking: |(x - y)| <= ε.
func Almost(x, y, ε float64) bool {
	return math.Abs(x-y) <= ε
}

// Mean (μ or average) value of a discrete signal.
// μ = (1/N) * Σx(i)
func Mean(sig Discrete) float64 {
	var sum float64
	if len(sig) == 0 {
		return 0.0
	}

	for _, v := range sig {
		sum += v
	}
	return sum / float64(len(sig))
}

// StdDeviation returns the standard deviation of the signal
// (Represented as σ in DSP formulas).
func StdDeviation(sig Discrete) float64 {
	return StdDeviation2(sig, Mean(sig))
}

// StdDeviation2 returns the standatd deviation from mean μ.
func StdDeviation2(sig Discrete, μ float64) float64 {
	return math.Sqrt(Variance2(sig, μ))
}

// Variance is the power of the standard deviation σ (represented by σ² in DSP).
// TODO: Call to this function result in excessive round-off errors because
// of the substracting of very close values (x(i)-μ).
// DSP Book, Chapter 2, pag. 16, shows an alternate method called
// 'running statistics' that should be used instead.
func Variance(sig Discrete) float64 {
	return Variance2(sig, Mean(sig))
}

// Variance2 is the variance from the mean μ
func Variance2(sig Discrete, μ float64) float64 {
	if len(sig) == 0 {
		return 0.0
	}

	var sum float64
	for _, v := range sig {
		sum += math.Pow(v-μ, 2)
	}
	return sum / float64(len(sig)-1)
}
