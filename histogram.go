package signal

import (
	"math"
)

// Histogram of a discrete signal.
// It's just a map of the amplitudes/ordinate values to how many times
// it occurs in the signal.
// Note that, as Discrete signal  stores floats but isn't sane to use
// floats as Go map keys, then the keys are stored as uint64 bits from
// the float64 values.
// Use Howmany() to ask for the counter of each amplitude value.
type Histogram map[uint64]uint64

// Hist creates a new histogram from the discrete signal s.
func Hist(s Discrete) Histogram {
	h := make(Histogram)
	for _, v := range s {
		vi := math.Float64bits(v)
		if v2, ok := h[vi]; ok {
			h[vi] = v2 + 1
		} else {
			h[vi] = 1
		}
	}
	return h
}

// Howmany returns how many samples has value aproximate by ε.
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
