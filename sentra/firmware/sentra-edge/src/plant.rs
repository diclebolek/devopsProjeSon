use sentra_core::{Frame, FrameIdentity, MotorSnapshot};

/// Deterministic stand-in for the 12 V motor and the MPU6050.
///
/// The sine sits on a 1 g offset on purpose: the DSP contract has to remove
/// gravity before it publishes RMS. Fault regimes live in the Go simulator;
/// this plant only proves the firmware loop.
#[derive(Debug, Clone)]
pub struct MockPlant {
    pub rpm: f64,
    pub current_a: f64,
    pub voltage_v: f64,
    pub temperature_c: f64,
    sample_index: u64,
}

impl Default for MockPlant {
    fn default() -> Self {
        Self {
            rpm: 1500.0,
            current_a: 0.82,
            voltage_v: 12.05,
            temperature_c: 36.4,
            sample_index: 0,
        }
    }
}

impl MockPlant {
    pub fn next_window(&mut self, sample_rate_hz: u32, window_samples: usize) -> Vec<f64> {
        let frequency_hz = self.rpm / 60.0;
        let mut samples = Vec::with_capacity(window_samples);
        for _ in 0..window_samples {
            let time_s = self.sample_index as f64 / f64::from(sample_rate_hz);
            let vibration = (2.0 * std::f64::consts::PI * frequency_hz * time_s).sin() * 0.18;
            samples.push(1.0 + vibration);
            self.sample_index = self.sample_index.wrapping_add(1);
        }
        samples
    }

    pub fn frame(
        &mut self,
        device_id: &str,
        firmware_version: &str,
        captured_at_ms: u64,
        sample_rate_hz: u32,
        window_samples: usize,
        wifi_rssi_dbm: Option<i32>,
    ) -> Result<Frame, sentra_core::CoreError> {
        let samples = self.next_window(sample_rate_hz, window_samples);
        Frame::from_window(
            FrameIdentity {
                device_id,
                firmware_version,
                captured_at_ms,
                sample_rate_hz,
            },
            MotorSnapshot {
                rpm: self.rpm,
                current_a: self.current_a,
                voltage_v: self.voltage_v,
                temperature_c: self.temperature_c,
            },
            &samples,
            true,
            wifi_rssi_dbm,
        )
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn mock_motor_publishes_a_valid_compact_frame() {
        let mut plant = MockPlant::default();
        let frame = plant
            .frame(
                "motor-lab-01",
                "0.1.0",
                1_710_000_000_000,
                1_000,
                256,
                Some(-65),
            )
            .unwrap();
        frame.validate().unwrap();
        // 0.18 g sine → RMS near 0.127 g, far below the 1 g gravity offset.
        assert!(frame.accel_rms_g > 0.12 && frame.accel_rms_g < 0.14);
        assert!(frame.dominant_freq_hz > 20.0 && frame.dominant_freq_hz < 30.0);
    }
}
