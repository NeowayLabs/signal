package signal_test

import (
	"testing"

	format "fmt"

	"github.com/NeowayLabs/signal"
)

const precision = 1e-09

type testcase struct {
	sig signal.Discrete
	μ   float64 // mean
	σ   float64 // std deviation
	σ2  float64 // variance
}

var fmt = format.Sprintf
var testcases = []testcase{
	{sig: signal.Discrete{1, 1, 1, 1}, μ: 1, σ: 0},
	{sig: signal.Discrete{1, 1, 1, 0}, μ: 0.75, σ: 0.5, σ2: 0.25},
	{sig: signal.Discrete{1, 1, 0, 0}, μ: 0.5, σ: 0.577350269, σ2: 0.333333333},
	{sig: signal.Discrete{1, 0, 0, 0}, μ: 0.25, σ: 0.5, σ2: 0.25},
	{sig: signal.Discrete{0, 0, 0, 0}, μ: 0, σ: 0, σ2: 0},
	{}, // division by zero
	{
		sig: signal.Discrete{
			0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
			10, 20, 30, 40, 50, 60, 70, 80, 90, 100,
			200, 300, 400, 500, 600, 700, 800, 900, 1000, 1100,
			200, 300, 400, 500, 600, 700, 800, 900, 1000, 1100,
			10, 20, 30, 40, 50, 60, 70, 80, 90, 100,
			0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
		},
		μ:  236.5,
		σ:  340.030781258,
		σ2: 115620.932203389,
	},
}

func assert(t *testing.T, b bool, msg string) {
	t.Helper()
	if !b {
		t.Fatal(msg)
	}
}

func assertAlmost(t *testing.T, x, y, ε float64, msg string) {
	t.Helper()
	assert(t, signal.Almost(x, y, ε), fmt("Fail: %s. Differs: %.12f != %.12f",
		msg, x, y))
}

func testmean(t *testing.T, s signal.Discrete, expected float64) {
	assertAlmost(t, signal.Mean(s), expected, precision, "mean")
}

func testdeviation(t *testing.T, s signal.Discrete, expected float64) {
	assertAlmost(t, signal.StdDeviation(s), expected, precision, "std deviation")
}

func testvariance(t *testing.T, s signal.Discrete, expected float64) {
	assertAlmost(t, signal.Variance(s), expected, precision, "variance")
}

func testhmean(t *testing.T, s signal.Discrete, expected float64) {
	assertAlmost(t, signal.HMean(signal.Hist(s)), expected, precision, "histogram mean")
}

func TestMean(t *testing.T) {
	for _, tc := range testcases {
		tc := tc
		testmean(t, tc.sig, tc.μ)
	}
}
func TestStdDeviation(t *testing.T) {
	for _, tc := range testcases {
		tc := tc
		testdeviation(t, tc.sig, tc.σ)
	}
}

func TestVariance(t *testing.T) {
	for _, tc := range testcases {
		tc := tc
		testvariance(t, tc.sig, tc.σ2)
	}
}

func TestHMean(t *testing.T) {
	for _, tc := range testcases {
		tc := tc
		testhmean(t, tc.sig, tc.μ)
	}
}
