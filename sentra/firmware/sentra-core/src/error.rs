use std::fmt;

/// Recoverable failure in the edge contract.
///
/// Library code returns this instead of panicking. The host binary and the
/// future ESP32 port both surface it on the serial log and skip the frame.
#[derive(Debug, Clone, PartialEq, Eq)]
pub enum CoreError {
    EmptyWindow,
    SampleRateZero,
    WindowNotPowerOfTwo { length: usize },
    SchemaVersion,
    InvalidDeviceId,
    InvalidFirmwareVersion,
    CapturedAtMissing,
    NonFinite { field: &'static str },
    OutOfRange { field: &'static str },
    ZeroPulsesPerRevolution,
    WindowDuration,
    ZeroSensitivity,
}

impl fmt::Display for CoreError {
    fn fmt(&self, formatter: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Self::EmptyWindow => write!(formatter, "vibration window is empty"),
            Self::SampleRateZero => write!(formatter, "sample rate must be greater than zero"),
            Self::WindowNotPowerOfTwo { length } => {
                write!(formatter, "window length {length} is not a power of two")
            }
            Self::SchemaVersion => write!(formatter, "unsupported schema_version"),
            Self::InvalidDeviceId => write!(formatter, "device_id is not a valid slug"),
            Self::InvalidFirmwareVersion => {
                write!(formatter, "firmware_version is empty or too long")
            }
            Self::CapturedAtMissing => {
                write!(formatter, "captured_at_ms must be greater than zero")
            }
            Self::NonFinite { field } => write!(formatter, "{field} is not a finite number"),
            Self::OutOfRange { field } => {
                write!(formatter, "{field} is outside the accepted range")
            }
            Self::ZeroPulsesPerRevolution => {
                write!(formatter, "hall pulses_per_revolution is zero")
            }
            Self::WindowDuration => write!(formatter, "rpm window duration must be positive"),
            Self::ZeroSensitivity => write!(formatter, "accelerometer counts_per_g is zero"),
        }
    }
}

impl std::error::Error for CoreError {}
