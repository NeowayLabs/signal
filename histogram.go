package signal

import (
	"math"
)

type Histogram map[uint64]uint64

func Hist(s Discrete) Histogram {
	H := make(Histogram)
	for _, v := range s {
		vi := math.Float64bits(v)
		if v2, ok := H[vi]; ok {
			H[vi] = v2 + 1
		} else {
			H[vi] = 1
		}
	}
	return H
}

func (h Histogram) Howmany(value float64, ε float64) uint64 {
	// fast path
	if v2, ok := h[math.Float64bits(value)]; ok {
		return v2
	}

	if ε == 0 {
		return 0
	}

	for k, v := range h {
		kf := math.Float64frombits(k)
		if Almost(kf, value, ε) {
			return v
		}
	}
	return 0
}

// HMean calculates the mean from histogram.
// Theoretically, the sames as Mean.
func HMean(hist Histogram) float64 {
	if len(hist) == 0 {
		return 0
	}

	var (
		sum float64
		n   uint64
	)
	for value, count := range hist {
		sum += float64(count) * math.Float64frombits(value)
		n += count
	}
	return sum / float64(n)
}
