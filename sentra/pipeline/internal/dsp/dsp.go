// Package dsp mirrors sentra-core. Keep the formulas aligned with docs/dsp.md
// and testdata/dsp_case.json. The ingest path does not recompute features; the
// simulator uses this package so a laptop can emit the same contract as the device.
package dsp

import "math"

type Features struct {
	AccelRMS       float64
	AccelPeak      float64
	CrestFactor    float64
	Kurtosis       float64
	DominantFreqHz float64
}

func VibrationFeatures(samples []float64, sampleRateHz int) (Features, error) {
	if len(samples) == 0 {
		return Features{}, errString("vibration window is empty")
	}
	if sampleRateHz <= 0 {
		return Features{}, errString("sample rate must be greater than zero")
	}
	if !isPowerOfTwo(len(samples)) {
		return Features{}, errString("window length is not a power of two")
	}
	mean := 0.0
	for _, sample := range samples {
		if math.IsNaN(sample) || math.IsInf(sample, 0) {
			return Features{}, errString("accel sample is not finite")
		}
		mean += sample
	}
	mean /= float64(len(samples))

	demeaned := make([]float64, len(samples))
	second := 0.0
	fourth := 0.0
	peak := 0.0
	for i, sample := range samples {
		centered := sample - mean
		demeaned[i] = centered
		abs := math.Abs(centered)
		if abs > peak {
			peak = abs
		}
		square := centered * centered
		second += square
		fourth += square * square
	}
	count := float64(len(samples))
	second /= count
	fourth /= count
	rms := math.Sqrt(second)
	crest := 0.0
	if rms > 1e-12 {
		crest = peak / rms
	}
	kurtosis := 0.0
	if second > 1e-18 {
		kurtosis = fourth / (second * second)
	}
	return Features{
		AccelRMS:       rms,
		AccelPeak:      peak,
		CrestFactor:    crest,
		Kurtosis:       kurtosis,
		DominantFreqHz: dominantFrequency(demeaned, sampleRateHz),
	}, nil
}

func dominantFrequency(demeaned []float64, sampleRateHz int) float64 {
	real := append([]float64(nil), demeaned...)
	imag := make([]float64, len(demeaned))
	fftInPlace(real, imag)
	nyquist := len(demeaned) / 2
	bestBin := 0
	bestMagnitude := 0.0
	for bin := 1; bin <= nyquist; bin++ {
		magnitude := math.Hypot(real[bin], imag[bin])
		if magnitude > bestMagnitude {
			bestMagnitude = magnitude
			bestBin = bin
		}
	}
	if bestBin == 0 {
		return 0
	}
	return float64(bestBin) * float64(sampleRateHz) / float64(len(demeaned))
}

func fftInPlace(real, imag []float64) {
	width := len(real)
	reversed := 0
	for index := 1; index < width; index++ {
		bit := width >> 1
		for reversed&bit != 0 {
			reversed ^= bit
			bit >>= 1
		}
		reversed ^= bit
		if index < reversed {
			real[index], real[reversed] = real[reversed], real[index]
			imag[index], imag[reversed] = imag[reversed], imag[index]
		}
	}
	for length := 2; length <= width; length <<= 1 {
		theta := -2 * math.Pi / float64(length)
		widthRe := math.Cos(theta)
		widthIm := math.Sin(theta)
		for start := 0; start < width; start += length {
			twiddleRe := 1.0
			twiddleIm := 0.0
			for offset := 0; offset < length/2; offset++ {
				even := start + offset
				odd := even + length/2
				oddRe := real[odd]*twiddleRe - imag[odd]*twiddleIm
				oddIm := real[odd]*twiddleIm + imag[odd]*twiddleRe
				real[odd] = real[even] - oddRe
				imag[odd] = imag[even] - oddIm
				real[even] += oddRe
				imag[even] += oddIm
				nextRe := twiddleRe*widthRe - twiddleIm*widthIm
				twiddleIm = twiddleRe*widthIm + twiddleIm*widthRe
				twiddleRe = nextRe
			}
		}
	}
}

func isPowerOfTwo(value int) bool {
	return value > 0 && value&(value-1) == 0
}

type errString string

func (e errString) Error() string { return string(e) }
