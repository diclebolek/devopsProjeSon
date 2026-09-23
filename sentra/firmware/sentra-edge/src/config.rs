use std::env;
use std::time::Duration;

#[derive(Debug, Clone)]
pub struct EdgeConfig {
    pub device_id: String,
    pub firmware_version: String,
    pub sample_rate_hz: u32,
    pub window_samples: usize,
    pub publish_interval: Duration,
    pub max_frames: Option<u64>,
    pub mqtt_broker: Option<String>,
    pub wifi_rssi_dbm: Option<i32>,
}

impl EdgeConfig {
    pub fn from_env() -> Result<Self, String> {
        let device_id = env_or("SENTRA_DEVICE_ID", "motor-lab-01");
        let firmware_version = env_or("SENTRA_FIRMWARE_VERSION", "0.1.0");
        let sample_rate_hz = parse_u32("SENTRA_SAMPLE_RATE_HZ", 1_000)?;
        let window_samples = parse_u32("SENTRA_WINDOW_SAMPLES", 256)? as usize;
        let interval_ms = parse_u64("SENTRA_PUBLISH_INTERVAL_MS", 1_000)?;
        let max_frames = match env::var("SENTRA_MAX_FRAMES") {
            Ok(value) if value != "0" && !value.is_empty() => {
                Some(parse_u64("SENTRA_MAX_FRAMES", 0)?)
            }
            _ => None,
        };
        let mqtt_broker = match env::var("SENTRA_MQTT_BROKER") {
            Ok(value) if !value.trim().is_empty() => Some(value),
            _ => None,
        };
        let wifi_rssi_dbm = match env::var("SENTRA_WIFI_RSSI_DBM") {
            Ok(value) if !value.is_empty() => Some(
                value
                    .parse::<i32>()
                    .map_err(|_| "SENTRA_WIFI_RSSI_DBM is not an integer".to_owned())?,
            ),
            _ => Some(-65),
        };
        Ok(Self {
            device_id,
            firmware_version,
            sample_rate_hz,
            window_samples,
            publish_interval: Duration::from_millis(interval_ms),
            max_frames,
            mqtt_broker,
            wifi_rssi_dbm,
        })
    }
}

fn env_or(key: &str, default: &str) -> String {
    env::var(key)
        .ok()
        .filter(|value| !value.is_empty())
        .unwrap_or_else(|| default.to_owned())
}

fn parse_u32(key: &str, default: u32) -> Result<u32, String> {
    match env::var(key) {
        Ok(value) if !value.is_empty() => value
            .parse::<u32>()
            .map_err(|_| format!("{key} must be an unsigned integer")),
        _ => Ok(default),
    }
}

fn parse_u64(key: &str, default: u64) -> Result<u64, String> {
    match env::var(key) {
        Ok(value) if !value.is_empty() => value
            .parse::<u64>()
            .map_err(|_| format!("{key} must be an unsigned integer")),
        _ => Ok(default),
    }
}

/// Accept `mqtt://host:1883`, `tcp://host:1883`, and `host:1883`.
pub fn parse_broker(raw: &str) -> Result<(String, u16), String> {
    let trimmed = raw.trim();
    let without_scheme = trimmed
        .strip_prefix("mqtt://")
        .or_else(|| trimmed.strip_prefix("tcp://"))
        .unwrap_or(trimmed);
    let (host, port_text) = without_scheme
        .rsplit_once(':')
        .ok_or_else(|| format!("broker {raw} needs a host and port"))?;
    if host.is_empty() {
        return Err(format!("broker {raw} is missing a host"));
    }
    let port = port_text
        .parse::<u16>()
        .map_err(|_| format!("broker {raw} has no numeric port"))?;
    Ok((host.to_owned(), port))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn parses_paho_and_mqtt_urls() {
        assert_eq!(
            parse_broker("tcp://127.0.0.1:1883").unwrap(),
            ("127.0.0.1".to_owned(), 1883)
        );
        assert_eq!(
            parse_broker("mqtt://broker.local:1884").unwrap(),
            ("broker.local".to_owned(), 1884)
        );
        assert!(parse_broker("127.0.0.1").is_err());
    }
}
