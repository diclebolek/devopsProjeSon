mod config;
mod mqtt;
mod plant;
mod publish;

use std::process::ExitCode;
use std::thread;
use std::time::{SystemTime, UNIX_EPOCH};

use config::EdgeConfig;
use plant::MockPlant;
use publish::{status_payload, write_serial, MqttPublisher, Publisher};
use sentra_core::Backoff;

fn main() -> ExitCode {
    let config = match EdgeConfig::from_env() {
        Ok(config) => config,
        Err(message) => {
            eprintln!("sentra-edge config: {message}");
            return ExitCode::from(2);
        }
    };

    let telemetry_topic = format!("sentra/v1/{}/telemetry", config.device_id);
    let status_topic = format!("sentra/v1/{}/status", config.device_id);
    let mut mqtt = match config.mqtt_broker.as_deref() {
        Some(broker) => {
            let client_id = format!("sentra-edge-{}", config.device_id);
            match MqttPublisher::connect(
                broker,
                &client_id,
                &status_topic,
                &config.firmware_version,
            ) {
                Ok(publisher) => {
                    eprintln!("sentra-edge mqtt broker {broker}");
                    Some(publisher)
                }
                Err(err) => {
                    eprintln!("sentra-edge mqtt setup failed: {err}");
                    eprintln!("sentra-edge continuing with serial output only");
                    None
                }
            }
        }
        None => {
            eprintln!("sentra-edge serial only; set SENTRA_MQTT_BROKER to publish");
            None
        }
    };

    if let Some(publisher) = mqtt.as_mut() {
        let payload = status_payload(true, &config.firmware_version);
        if let Err(err) = publisher.publish(&status_topic, &payload) {
            eprintln!("sentra-edge status publish failed: {err}");
        }
    }

    let mut plant = MockPlant::default();
    let mut published = 0u64;
    let mut backoff = Backoff::default();
    loop {
        let captured_at_ms = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .map(|duration| duration.as_millis() as u64)
            .unwrap_or(0);
        match plant.frame(
            &config.device_id,
            &config.firmware_version,
            captured_at_ms.max(1),
            config.sample_rate_hz,
            config.window_samples,
            config.wifi_rssi_dbm,
        ) {
            Ok(frame) => match frame.to_json_vec() {
                Ok(payload) => {
                    if let Err(err) = write_serial(&payload) {
                        eprintln!("sentra-edge serial write failed: {err}");
                        return ExitCode::from(1);
                    }
                    if let Some(publisher) = mqtt.as_mut() {
                        if let Err(err) = publisher.publish(&telemetry_topic, &payload) {
                            let delay = std::time::Duration::from_millis(backoff.next_delay_ms());
                            eprintln!("sentra-edge telemetry publish failed: {err}; retrying in {delay:?}");
                            thread::sleep(delay);
                            if let Err(reconnect_err) = publisher.reconnect() {
                                eprintln!("sentra-edge mqtt reconnect failed: {reconnect_err}");
                            }
                        } else {
                            backoff.reset();
                        }
                    }
                }
                Err(err) => eprintln!("sentra-edge encode failed: {err}"),
            },
            Err(err) => eprintln!("sentra-edge frame rejected: {err}"),
        }

        published += 1;
        if config.max_frames.is_some_and(|limit| published >= limit) {
            return ExitCode::SUCCESS;
        }
        thread::sleep(config.publish_interval);
    }
}
