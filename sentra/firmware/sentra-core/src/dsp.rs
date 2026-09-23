use crate::error::CoreError;

/// Vibration summary published instead of the raw window.
///
/// Every value is computed on the window after its mean is subtracted.
/// A stationary accelerometer reads about 1 g from gravity; that offset is
/// not motor vibration and must not dominate RMS.
#[derive(Debug, Clone, PartialEq)]
pub struct VibrationFeatures {
    pub accel_rms_g: f64,
    pub accel_peak_g: f64,
    pub crest_factor: f64,
    pub kurtosis: f64,
    pub dominant_freq_hz: f64,
}

/// Population AC-RMS, peak, crest factor, raw kurtosis, and dominant frequency.
///
/// Kurtosis is the non-excess population moment (a pure sine is 1.5, a
/// Gaussian is about 3). Crest factor is peak divided by RMS. Dominant
/// frequency is the positive FFT bin with the largest magnitude, ignoring DC.
pub fn vibration_features(
    samples_g: &[f64],
    sample_rate_hz: u32,
) -> Result<VibrationFeatures, CoreError> {
    if samples_g.is_empty() {
        return Err(CoreError::EmptyWindow);
    }
    if sample_rate_hz == 0 {
        return Err(CoreError::SampleRateZero);
    }
    if !samples_g.len().is_power_of_two() {
        return Err(CoreError::WindowNotPowerOfTwo {
            length: samples_g.len(),
        });
    }
    for sample in samples_g {
        if !sample.is_finite() {
            return Err(CoreError::NonFinite {
                field: "accel_sample_g",
            });
        }
    }

    let mean = samples_g.iter().sum::<f64>() / samples_g.len() as f64;
    let mut demeaned = Vec::with_capacity(samples_g.len());
    let mut second_moment = 0.0;
    let mut fourth_moment = 0.0;
    let mut peak = 0.0;
    for sample in samples_g {
        let centered = sample - mean;
        demeaned.push(centered);
        let abs = centered.abs();
        if abs > peak {
            peak = abs;
        }
        let square = centered * centered;
        second_moment += square;
        fourth_moment += square * square;
    }
    let count = samples_g.len() as f64;
    second_moment /= count;
    fourth_moment /= count;

    let rms = second_moment.sqrt();
    let crest_factor = if rms > 1e-12 { peak / rms } else { 0.0 };
    let kurtosis = if second_moment > 1e-18 {
        fourth_moment / (second_moment * second_moment)
    } else {
        0.0
    };

    Ok(VibrationFeatures {
        accel_rms_g: rms,
        accel_peak_g: peak,
        crest_factor,
        kurtosis,
        dominant_freq_hz: dominant_frequency(&demeaned, sample_rate_hz),
    })
}

fn dominant_frequency(demeaned: &[f64], sample_rate_hz: u32) -> f64 {
    let mut real: Vec<f64> = demeaned.to_vec();
    let mut imag = vec![0.0; demeaned.len()];
    fft_inplace(&mut real, &mut imag);

    let nyquist = demeaned.len() / 2;
    let mut best_bin = 0usize;
    let mut best_magnitude = 0.0;
    for bin in 1..=nyquist {
        let magnitude = real[bin].hypot(imag[bin]);
        if magnitude > best_magnitude {
            best_magnitude = magnitude;
            best_bin = bin;
        }
    }
    if best_bin == 0 {
        return 0.0;
    }
    best_bin as f64 * f64::from(sample_rate_hz) / demeaned.len() as f64
}

/// In-place radix-2 Cooley–Tukey FFT. `real` and `imag` must share a power-of-two length.
fn fft_inplace(real: &mut [f64], imag: &mut [f64]) {
    let width = real.len();
    debug_assert!(width.is_power_of_two() && width == imag.len());

    let mut reversed = 0usize;
    for index in 1..width {
        let mut bit = width >> 1;
        while reversed & bit != 0 {
            reversed ^= bit;
            bit >>= 1;
        }
        reversed ^= bit;
        if index < reversed {
            real.swap(index, reversed);
            imag.swap(index, reversed);
        }
    }

    let mut length = 2usize;
    while length <= width {
        let theta = -2.0 * std::f64::consts::PI / length as f64;
        let width_re = theta.cos();
        let width_im = theta.sin();
        let mut start = 0usize;
        while start < width {
            let mut twiddle_re = 1.0;
            let mut twiddle_im = 0.0;
            for offset in 0..(length / 2) {
                let even = start + offset;
                let odd = even + length / 2;
                let odd_re = real[odd] * twiddle_re - imag[odd] * twiddle_im;
                let odd_im = real[odd] * twiddle_im + imag[odd] * twiddle_re;
                real[odd] = real[even] - odd_re;
                imag[odd] = imag[even] - odd_im;
                real[even] += odd_re;
                imag[even] += odd_im;
                let next_re = twiddle_re * width_re - twiddle_im * width_im;
                twiddle_im = twiddle_re * width_im + twiddle_im * width_re;
                twiddle_re = next_re;
            }
            start += length;
        }
        length <<= 1;
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use serde::Deserialize;

    #[derive(Deserialize)]
    struct DspCase {
        sine: SineCase,
        fixed: FixedCase,
    }

    #[derive(Deserialize)]
    struct SineCase {
        sample_rate_hz: u32,
        frequency_hz: f64,
        amplitude: f64,
        offset: f64,
        count: usize,
        expected_rms: f64,
        expected_peak: f64,
        expected_crest_factor: f64,
        expected_kurtosis: f64,
        expected_dominant_freq_hz: f64,
    }

    #[derive(Deserialize)]
    struct FixedCase {
        sample_rate_hz: u32,
        samples: Vec<f64>,
        expected_rms: f64,
        expected_peak: f64,
        expected_crest_factor: f64,
        expected_kurtosis: f64,
        expected_dominant_freq_hz: f64,
    }

    fn load_case() -> DspCase {
        serde_json::from_str(include_str!("../../../testdata/dsp_case.json")).expect("dsp case")
    }

    fn assert_close(actual: f64, expected: f64) {
        let tolerance = 1e-9;
        assert!(
            (actual - expected).abs() <= tolerance,
            "actual {actual} expected {expected}"
        );
    }

    #[test]
    fn sine_with_gravity_offset_matches_shared_oracle() {
        let case = load_case().sine;
        let samples: Vec<f64> = (0..case.count)
            .map(|index| {
                let phase = 2.0 * std::f64::consts::PI * case.frequency_hz * index as f64
                    / f64::from(case.sample_rate_hz);
                case.offset + case.amplitude * phase.sin()
            })
            .collect();
        let features = vibration_features(&samples, case.sample_rate_hz).unwrap();
        assert_close(features.accel_rms_g, case.expected_rms);
        assert_close(features.accel_peak_g, case.expected_peak);
        assert_close(features.crest_factor, case.expected_crest_factor);
        assert_close(features.kurtosis, case.expected_kurtosis);
        assert_close(features.dominant_freq_hz, case.expected_dominant_freq_hz);
    }

    #[test]
    fn fixed_block_wave_matches_shared_oracle() {
        let case = load_case().fixed;
        let features = vibration_features(&case.samples, case.sample_rate_hz).unwrap();
        assert_close(features.accel_rms_g, case.expected_rms);
        assert_close(features.accel_peak_g, case.expected_peak);
        assert_close(features.crest_factor, case.expected_crest_factor);
        assert_close(features.kurtosis, case.expected_kurtosis);
        assert_close(features.dominant_freq_hz, case.expected_dominant_freq_hz);
    }

    #[test]
    fn impulse_fft_spreads_energy_across_bins() {
        let mut real = vec![1.0, 0.0, 0.0, 0.0];
        let mut imag = vec![0.0; 4];
        fft_inplace(&mut real, &mut imag);
        for (re, im) in real.iter().zip(imag.iter()) {
            assert_close(*re, 1.0);
            assert_close(*im, 0.0);
        }
    }

    #[test]
    fn rejects_odd_window_lengths() {
        let err = vibration_features(&[0.0, 0.1, 0.2], 1000).unwrap_err();
        assert_eq!(err, CoreError::WindowNotPowerOfTwo { length: 3 });
    }

    #[test]
    fn silent_window_is_zero_not_nan() {
        let features = vibration_features(&[1.0; 16], 1000).unwrap();
        assert_eq!(features.accel_rms_g, 0.0);
        assert_eq!(features.accel_peak_g, 0.0);
        assert_eq!(features.crest_factor, 0.0);
        assert_eq!(features.kurtosis, 0.0);
        assert_eq!(features.dominant_freq_hz, 0.0);
    }
}
