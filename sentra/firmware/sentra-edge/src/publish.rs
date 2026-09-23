use std::io::Write;

use crate::config::parse_broker;
use crate::mqtt::{LastWill, MqttSession};

pub trait Publisher {
    fn publish(&mut self, topic: &str, payload: &[u8]) -> Result<(), PublishError>;
}

#[derive(Debug)]
pub struct PublishError {
    pub message: String,
}

impl PublishError {
    pub fn new(message: impl Into<String>) -> Self {
        Self {
            message: message.into(),
        }
    }
}

impl std::fmt::Display for PublishError {
    fn fmt(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(formatter, "{}", self.message)
    }
}

impl std::error::Error for PublishError {}

/// Serial-port stand-in used by tests. The process itself writes JSON with [`write_serial`].
#[cfg(test)]
pub struct StdoutPublisher<W: Write> {
    output: W,
}

#[cfg(test)]
impl<W: Write> StdoutPublisher<W> {
    pub fn new(output: W) -> Self {
        Self { output }
    }
}

#[cfg(test)]
impl<W: Write> Publisher for StdoutPublisher<W> {
    fn publish(&mut self, _topic: &str, payload: &[u8]) -> Result<(), PublishError> {
        self.output
            .write_all(payload)
            .and_then(|_| self.output.write_all(b"\n"))
            .and_then(|_| self.output.flush())
            .map_err(|err| PublishError::new(err.to_string()))
    }
}

pub struct MqttPublisher {
    session: MqttSession,
    broker: String,
    client_id: String,
    status_topic: String,
    firmware_version: String,
}

impl MqttPublisher {
    pub fn connect(
        broker: &str,
        client_id: &str,
        status_topic: &str,
        firmware_version: &str,
    ) -> Result<Self, PublishError> {
        let (host, port) = parse_broker(broker).map_err(PublishError::new)?;
        let session = MqttSession::connect(
            &host,
            port,
            client_id,
            &LastWill {
                topic: status_topic.to_owned(),
                message: status_payload(false, firmware_version),
            },
        )?;
        Ok(Self {
            session,
            broker: broker.to_owned(),
            client_id: client_id.to_owned(),
            status_topic: status_topic.to_owned(),
            firmware_version: firmware_version.to_owned(),
        })
    }

    pub fn reconnect(&mut self) -> Result<(), PublishError> {
        let replacement = Self::connect(
            &self.broker,
            &self.client_id,
            &self.status_topic,
            &self.firmware_version,
        )?;
        *self = replacement;
        Ok(())
    }
}

impl Publisher for MqttPublisher {
    fn publish(&mut self, topic: &str, payload: &[u8]) -> Result<(), PublishError> {
        self.session.publish(topic, payload)
    }
}

pub fn status_payload(online: bool, firmware_version: &str) -> Vec<u8> {
    format!(
        "{{\"online\":{online},\"firmware_version\":\"{firmware_version}\",\"detail\":\"edge node\"}}"
    )
    .into_bytes()
}

pub fn write_serial(payload: &[u8]) -> std::io::Result<()> {
    let mut stdout = std::io::stdout().lock();
    stdout.write_all(payload)?;
    stdout.write_all(b"\n")?;
    stdout.flush()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn stdout_publisher_writes_one_json_line() {
        let mut buffer = Vec::new();
        {
            let mut publisher = StdoutPublisher::new(&mut buffer);
            publisher
                .publish("sentra/v1/motor-lab-01/telemetry", b"{\"ok\":true}")
                .unwrap();
        }
        assert_eq!(buffer, b"{\"ok\":true}\n");
    }

    #[test]
    fn status_payload_marks_the_node_online() {
        let payload = status_payload(true, "0.1.0");
        let text = String::from_utf8(payload).unwrap();
        assert!(text.contains("\"online\":true"));
        assert!(text.contains("0.1.0"));
    }
}
