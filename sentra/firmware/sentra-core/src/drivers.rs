use crate::error::CoreError;

/// ±2 g full-scale setting of the MPU6050: 16 384 counts per g.
pub const MPU6050_COUNTS_PER_G_2G: f64 = 16384.0;

/// Convert a raw accelerometer count to g.
pub fn accel_g_from_raw(raw_count: i16, counts_per_g: f64) -> Result<f64, CoreError> {
    if counts_per_g == 0.0 {
        return Err(CoreError::ZeroSensitivity);
    }
    Ok(f64::from(raw_count) / counts_per_g)
}

/// Quantize g back to a register count. Used by the host plant and by tests.
pub fn raw_from_accel_g(accel_g: f64, counts_per_g: f64) -> Result<i16, CoreError> {
    if counts_per_g == 0.0 || !accel_g.is_finite() {
        return Err(CoreError::ZeroSensitivity);
    }
    let scaled = (accel_g * counts_per_g).round();
    if scaled < f64::from(i16::MIN) || scaled > f64::from(i16::MAX) {
        return Err(CoreError::OutOfRange { field: "accel_raw" });
    }
    Ok(scaled as i16)
}

/// DS18B20 12-bit scratchpad temperature. 1 LSB = 1/16 °C.
pub fn ds18b20_celsius_from_raw(raw: i16) -> f64 {
    f64::from(raw) / 16.0
}

/// Hall-effect tachometer. `pulses` is the count over `window_seconds`.
pub fn rpm_from_pulses(
    pulses: u32,
    pulses_per_revolution: u32,
    window_seconds: f64,
) -> Result<f64, CoreError> {
    if pulses_per_revolution == 0 {
        return Err(CoreError::ZeroPulsesPerRevolution);
    }
    if !(window_seconds.is_finite() && window_seconds > 0.0) {
        return Err(CoreError::WindowDuration);
    }
    Ok(f64::from(pulses) / f64::from(pulses_per_revolution) / window_seconds * 60.0)
}

/// INA219 configuration for the 12 V lab motor.
///
/// The shunt and the expected full-scale current fix the current LSB and the
/// calibration register, following the equations in the TI INA219 datasheet.
#[derive(Debug, Clone, Copy, PartialEq)]
pub struct Ina219Config {
    pub shunt_ohms: f64,
    pub max_expected_current_a: f64,
}

impl Ina219Config {
    /// 0.1 Ω shunt, 3.2 A full scale. Comfortable for a small 12 V DC motor,
    /// and it keeps the calibration register inside the 16-bit range.
    pub const LAB_MOTOR: Self = Self {
        shunt_ohms: 0.1,
        max_expected_current_a: 3.2,
    };
}

pub fn ina219_current_lsb(config: Ina219Config) -> Result<f64, CoreError> {
    if !(config.shunt_ohms > 0.0 && config.max_expected_current_a > 0.0) {
        return Err(CoreError::OutOfRange {
            field: "ina219_config",
        });
    }
    Ok(config.max_expected_current_a / 32768.0)
}

pub fn ina219_calibration_register(config: Ina219Config) -> Result<u16, CoreError> {
    let current_lsb = ina219_current_lsb(config)?;
    let calibration = (0.04096 / (current_lsb * config.shunt_ohms)).trunc();
    if !(calibration.is_finite() && (1.0..65535.0).contains(&calibration)) {
        return Err(CoreError::OutOfRange {
            field: "ina219_calibration",
        });
    }
    Ok(calibration as u16)
}

/// Bus-voltage register to volts. Bits 3..15 are the voltage, 4 mV per LSB.
pub fn ina219_bus_voltage_volts(raw_register: u16) -> f64 {
    f64::from(raw_register >> 3) * 0.004
}

pub fn ina219_current_amps(raw_current: i16, current_lsb: f64) -> f64 {
    f64::from(raw_current) * current_lsb
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn mpu6050_roundtrip_stays_within_one_count() {
        let original = 0.37;
        let raw = raw_from_accel_g(original, MPU6050_COUNTS_PER_G_2G).unwrap();
        let decoded = accel_g_from_raw(raw, MPU6050_COUNTS_PER_G_2G).unwrap();
        assert!((decoded - original).abs() < 1.0 / MPU6050_COUNTS_PER_G_2G);
    }

    #[test]
    fn ds18b20_known_scratchpad() {
        // 0x0191 is the datasheet example for 25.0625 °C.
        assert!((ds18b20_celsius_from_raw(0x0191) - 25.0625).abs() < 1e-9);
    }

    #[test]
    fn hall_rpm_at_one_pulse_per_revolution() {
        let rpm = rpm_from_pulses(25, 1, 1.0).unwrap();
        assert!((rpm - 1500.0).abs() < 1e-9);
    }

    #[test]
    fn ina219_lab_calibration_matches_datasheet_equation() {
        let calibration = ina219_calibration_register(Ina219Config::LAB_MOTOR).unwrap();
        assert_eq!(calibration, 4194);
        let current_lsb = ina219_current_lsb(Ina219Config::LAB_MOTOR).unwrap();
        let amps = ina219_current_amps(10_000, current_lsb);
        assert!((amps - 0.9765625).abs() < 1e-9);
    }

    #[test]
    fn ina219_bus_voltage_twelve_volts() {
        // 12.000 V / 4 mV = 3000, stored in bits 3..15.
        let register = 3000u16 << 3;
        assert!((ina219_bus_voltage_volts(register) - 12.0).abs() < 1e-9);
    }
}
