use serde::{Deserialize, Serialize};

use crate::dsp::vibration_features;
use crate::error::CoreError;

pub const SCHEMA_VERSION: u32 = 1;

/// One compact telemetry frame. This is the MQTT payload and the database row.
///
/// Raw samples stay on the device. The fields below are the contract shared
/// with the Go ingest service; both sides decode `testdata/frame_normal.json`.
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub struct Frame {
    pub schema_version: u32,
    pub device_id: String,
    pub firmware_version: String,
    pub captured_at_ms: u64,
    pub sample_rate_hz: u32,
    pub window_samples: u32,
    pub rpm: f64,
    pub current_a: f64,
    pub voltage_v: f64,
    pub temperature_c: f64,
    pub accel_rms_g: f64,
    pub accel_peak_g: f64,
    pub crest_factor: f64,
    pub kurtosis: f64,
    pub dominant_freq_hz: f64,
    pub sensor_ok: bool,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub wifi_rssi_dbm: Option<i32>,
}

/// Slow channels sampled once per vibration window.
#[derive(Debug, Clone, Copy)]
pub struct MotorSnapshot {
    pub rpm: f64,
    pub current_a: f64,
    pub voltage_v: f64,
    pub temperature_c: f64,
}

/// Identity and timing that accompany a window.
#[derive(Debug, Clone)]
pub struct FrameIdentity<'a> {
    pub device_id: &'a str,
    pub firmware_version: &'a str,
    pub captured_at_ms: u64,
    pub sample_rate_hz: u32,
}

impl Frame {
    pub fn from_window(
        identity: FrameIdentity<'_>,
        motor: MotorSnapshot,
        samples_g: &[f64],
        sensor_ok: bool,
        wifi_rssi_dbm: Option<i32>,
    ) -> Result<Self, CoreError> {
        let features = vibration_features(samples_g, identity.sample_rate_hz)?;
        let frame = Self {
            schema_version: SCHEMA_VERSION,
            device_id: identity.device_id.to_owned(),
            firmware_version: identity.firmware_version.to_owned(),
            captured_at_ms: identity.captured_at_ms,
            sample_rate_hz: identity.sample_rate_hz,
            window_samples: samples_g.len() as u32,
            rpm: motor.rpm,
            current_a: motor.current_a,
            voltage_v: motor.voltage_v,
            temperature_c: motor.temperature_c,
            accel_rms_g: features.accel_rms_g,
            accel_peak_g: features.accel_peak_g,
            crest_factor: features.crest_factor,
            kurtosis: features.kurtosis,
            dominant_freq_hz: features.dominant_freq_hz,
            sensor_ok,
            wifi_rssi_dbm,
        };
        frame.validate()?;
        Ok(frame)
    }

    pub fn validate(&self) -> Result<(), CoreError> {
        if self.schema_version != SCHEMA_VERSION {
            return Err(CoreError::SchemaVersion);
        }
        if !valid_device_id(&self.device_id) {
            return Err(CoreError::InvalidDeviceId);
        }
        if !valid_firmware_version(&self.firmware_version) {
            return Err(CoreError::InvalidFirmwareVersion);
        }
        if self.captured_at_ms == 0 {
            return Err(CoreError::CapturedAtMissing);
        }
        if !(100..=8_000).contains(&self.sample_rate_hz) {
            return Err(CoreError::OutOfRange {
                field: "sample_rate_hz",
            });
        }
        if !(16..=2048).contains(&self.window_samples) || !self.window_samples.is_power_of_two() {
            return Err(CoreError::OutOfRange {
                field: "window_samples",
            });
        }
        check_range(self.rpm, 0.0, 20_000.0, "rpm")?;
        check_range(self.current_a, 0.0, 20.0, "current_a")?;
        check_range(self.voltage_v, 0.0, 30.0, "voltage_v")?;
        check_range(self.temperature_c, -40.0, 150.0, "temperature_c")?;
        check_range(self.accel_rms_g, 0.0, 50.0, "accel_rms_g")?;
        check_range(self.accel_peak_g, 0.0, 50.0, "accel_peak_g")?;
        check_range(self.crest_factor, 0.0, 100.0, "crest_factor")?;
        check_range(self.kurtosis, 0.0, 100.0, "kurtosis")?;
        let nyquist = f64::from(self.sample_rate_hz) / 2.0;
        check_range(self.dominant_freq_hz, 0.0, nyquist, "dominant_freq_hz")?;
        if let Some(rssi) = self.wifi_rssi_dbm {
            if !(-120..=0).contains(&rssi) {
                return Err(CoreError::OutOfRange {
                    field: "wifi_rssi_dbm",
                });
            }
        }
        Ok(())
    }

    pub fn to_json_vec(&self) -> Result<Vec<u8>, serde_json::Error> {
        serde_json::to_vec(self)
    }
}

fn check_range(value: f64, min: f64, max: f64, field: &'static str) -> Result<(), CoreError> {
    if !value.is_finite() {
        return Err(CoreError::NonFinite { field });
    }
    if !(min..=max).contains(&value) {
        return Err(CoreError::OutOfRange { field });
    }
    Ok(())
}

pub fn valid_device_id(device_id: &str) -> bool {
    let bytes = device_id.as_bytes();
    if bytes.is_empty() || bytes.len() > 64 {
        return false;
    }
    if !is_slug_char(bytes[0]) || bytes[0] == b'-' {
        return false;
    }
    if bytes[bytes.len() - 1] == b'-' {
        return false;
    }
    bytes.iter().all(|byte| is_slug_char(*byte))
}

fn is_slug_char(byte: u8) -> bool {
    byte.is_ascii_lowercase() || byte.is_ascii_digit() || byte == b'-'
}

pub fn valid_firmware_version(version: &str) -> bool {
    let bytes = version.as_bytes();
    (1..=32).contains(&bytes.len())
        && bytes
            .iter()
            .all(|byte| byte.is_ascii_alphanumeric() || matches!(byte, b'.' | b'_' | b'-'))
}

/// Reconnect delay for Wi-Fi and MQTT. Caps at 30 s so a dead broker cannot
/// spin the radio, and resets after a successful publish.
#[derive(Debug, Clone)]
pub struct Backoff {
    attempt: u32,
    base_ms: u64,
    cap_ms: u64,
}

impl Default for Backoff {
    fn default() -> Self {
        Self {
            attempt: 0,
            base_ms: 1_000,
            cap_ms: 30_000,
        }
    }
}

impl Backoff {
    pub fn next_delay_ms(&mut self) -> u64 {
        let shift = self.attempt.min(5);
        self.attempt = self.attempt.saturating_add(1);
        (self.base_ms.saturating_mul(1_u64 << shift)).min(self.cap_ms)
    }

    pub fn reset(&mut self) {
        self.attempt = 0;
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn golden_frame_decodes_and_validates() {
        let frame: Frame =
            serde_json::from_str(include_str!("../../../testdata/frame_normal.json")).unwrap();
        frame.validate().unwrap();
        assert_eq!(frame.device_id, "motor-lab-01");
        assert_eq!(frame.wifi_rssi_dbm, Some(-61));
        let encoded = frame.to_json_vec().unwrap();
        let again: Frame = serde_json::from_slice(&encoded).unwrap();
        assert_eq!(again, frame);
    }

    #[test]
    fn missing_rssi_is_accepted() {
        let mut frame: Frame =
            serde_json::from_str(include_str!("../../../testdata/frame_normal.json")).unwrap();
        frame.wifi_rssi_dbm = None;
        let encoded = serde_json::to_string(&frame).unwrap();
        assert!(!encoded.contains("wifi_rssi_dbm"));
        frame.validate().unwrap();
    }

    #[test]
    fn rejects_bad_identity() {
        let mut frame: Frame =
            serde_json::from_str(include_str!("../../../testdata/frame_normal.json")).unwrap();
        frame.device_id = "Motor_01".to_owned();
        assert_eq!(frame.validate(), Err(CoreError::InvalidDeviceId));
        frame.device_id = "motor-lab-01".to_owned();
        frame.schema_version = 2;
        assert_eq!(frame.validate(), Err(CoreError::SchemaVersion));
    }

    #[test]
    fn from_window_removes_gravity_before_publishing() {
        let samples: Vec<f64> = (0..256)
            .map(|index| {
                let phase = 2.0 * std::f64::consts::PI * 8.0 * index as f64 / 256.0;
                1.0 + phase.sin()
            })
            .collect();
        let frame = Frame::from_window(
            FrameIdentity {
                device_id: "motor-lab-01",
                firmware_version: "0.1.0",
                captured_at_ms: 1_710_000_000_000,
                sample_rate_hz: 256,
            },
            MotorSnapshot {
                rpm: 1500.0,
                current_a: 0.8,
                voltage_v: 12.0,
                temperature_c: 35.0,
            },
            &samples,
            true,
            Some(-70),
        )
        .unwrap();
        assert!((frame.accel_rms_g - 0.7071067811865476).abs() < 1e-9);
        assert!((frame.dominant_freq_hz - 8.0).abs() < 1e-9);
        assert_eq!(frame.window_samples, 256);
    }

    #[test]
    fn backoff_doubles_until_the_cap_then_resets() {
        let mut backoff = Backoff::default();
        assert_eq!(backoff.next_delay_ms(), 1_000);
        assert_eq!(backoff.next_delay_ms(), 2_000);
        assert_eq!(backoff.next_delay_ms(), 4_000);
        assert_eq!(backoff.next_delay_ms(), 8_000);
        assert_eq!(backoff.next_delay_ms(), 16_000);
        assert_eq!(backoff.next_delay_ms(), 30_000);
        assert_eq!(backoff.next_delay_ms(), 30_000);
        backoff.reset();
        assert_eq!(backoff.next_delay_ms(), 1_000);
    }
}
