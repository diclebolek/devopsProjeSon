//! Edge contract for SENTRA.
//!
//! This crate is deliberately free of hardware crates so the DSP, the sensor
//! decoders, and the telemetry schema can be tested on a laptop. The ESP32
//! port calls the same functions; it does not reimplement the math.

mod drivers;
mod dsp;
mod error;
mod frame;

pub use drivers::{
    accel_g_from_raw, ds18b20_celsius_from_raw, ina219_bus_voltage_volts,
    ina219_calibration_register, ina219_current_amps, ina219_current_lsb, raw_from_accel_g,
    rpm_from_pulses, Ina219Config, MPU6050_COUNTS_PER_G_2G,
};
pub use dsp::{vibration_features, VibrationFeatures};
pub use error::CoreError;
pub use frame::{
    valid_device_id, valid_firmware_version, Backoff, Frame, FrameIdentity, MotorSnapshot,
    SCHEMA_VERSION,
};
