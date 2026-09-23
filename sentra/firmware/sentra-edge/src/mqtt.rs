use std::io::{Read, Write};
use std::net::{TcpStream, ToSocketAddrs};
use std::time::Duration;

use crate::publish::PublishError;

/// Lab MQTT 3.1.1 client. QoS 1 publish plus a last-will in CONNECT.
///
/// The packet encoder is dependency-free so the edge binary stays buildable
/// on the Rust toolchain used for host tests. TLS is intentionally absent:
/// the lab broker in docker-compose is plaintext, and production TLS belongs
/// in the broker profile, not in an untested stack.
pub struct MqttSession {
    stream: TcpStream,
    next_packet_id: u16,
}

pub struct LastWill {
    pub topic: String,
    pub message: Vec<u8>,
}

impl MqttSession {
    pub fn connect(
        host: &str,
        port: u16,
        client_id: &str,
        will: &LastWill,
    ) -> Result<Self, PublishError> {
        let address = format!("{host}:{port}");
        let socket = address
            .to_socket_addrs()
            .map_err(|err| PublishError::new(err.to_string()))?
            .next()
            .ok_or_else(|| PublishError::new(format!("no address for {address}")))?;
        let mut stream = TcpStream::connect_timeout(&socket, Duration::from_secs(3))
            .map_err(|err| PublishError::new(err.to_string()))?;
        stream
            .set_read_timeout(Some(Duration::from_secs(3)))
            .map_err(|err| PublishError::new(err.to_string()))?;
        stream
            .set_write_timeout(Some(Duration::from_secs(3)))
            .map_err(|err| PublishError::new(err.to_string()))?;
        let packet = connect_packet(client_id, will);
        stream
            .write_all(&packet)
            .map_err(|err| PublishError::new(err.to_string()))?;
        let mut ack = [0u8; 4];
        stream
            .read_exact(&mut ack)
            .map_err(|err| PublishError::new(format!("connack: {err}")))?;
        if ack[0] != 0x20 || ack[3] != 0x00 {
            return Err(PublishError::new(format!(
                "broker rejected connection: {:02x?}",
                ack
            )));
        }
        Ok(Self {
            stream,
            next_packet_id: 1,
        })
    }

    pub fn publish(&mut self, topic: &str, payload: &[u8]) -> Result<(), PublishError> {
        let packet_id = self.next_packet_id;
        self.next_packet_id = self.next_packet_id.wrapping_add(1).max(1);
        let packet = publish_packet(topic, payload, packet_id);
        self.stream
            .write_all(&packet)
            .map_err(|err| PublishError::new(err.to_string()))?;
        let mut ack = [0u8; 4];
        self.stream
            .read_exact(&mut ack)
            .map_err(|err| PublishError::new(format!("puback: {err}")))?;
        if ack[0] != 0x40 {
            return Err(PublishError::new(format!(
                "expected puback, got {:02x?}",
                ack
            )));
        }
        Ok(())
    }
}

pub fn connect_packet(client_id: &str, will: &LastWill) -> Vec<u8> {
    // Clean session, will flag, will QoS 1, will retain.
    let flags = 0b0010_1110u8;
    let mut variable = Vec::new();
    variable.extend(mqtt_string("MQTT"));
    variable.push(4);
    variable.push(flags);
    variable.extend_from_slice(&30u16.to_be_bytes());
    variable.extend(mqtt_string(client_id));
    variable.extend(mqtt_string(&will.topic));
    variable.extend(mqtt_bytes(&will.message));
    finish(0x10, variable)
}

pub fn publish_packet(topic: &str, payload: &[u8], packet_id: u16) -> Vec<u8> {
    let mut variable = Vec::new();
    variable.extend(mqtt_string(topic));
    variable.extend_from_slice(&packet_id.to_be_bytes());
    variable.extend_from_slice(payload);
    // QoS 1 publish, not retained.
    finish(0x32, variable)
}

fn finish(header: u8, variable: Vec<u8>) -> Vec<u8> {
    let mut packet = vec![header];
    packet.extend(encode_remaining_length(variable.len()));
    packet.extend(variable);
    packet
}

fn mqtt_string(text: &str) -> Vec<u8> {
    mqtt_bytes(text.as_bytes())
}

fn mqtt_bytes(bytes: &[u8]) -> Vec<u8> {
    let mut out = Vec::with_capacity(bytes.len() + 2);
    out.extend_from_slice(&(bytes.len() as u16).to_be_bytes());
    out.extend_from_slice(bytes);
    out
}

fn encode_remaining_length(mut length: usize) -> Vec<u8> {
    let mut encoded = Vec::new();
    loop {
        let mut byte = (length % 128) as u8;
        length /= 128;
        if length > 0 {
            byte |= 0x80;
        }
        encoded.push(byte);
        if length == 0 {
            break;
        }
    }
    encoded
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn connect_packet_carries_the_will() {
        let packet = connect_packet(
            "edge-1",
            &LastWill {
                topic: "sentra/v1/motor-lab-01/status".to_owned(),
                message: br#"{"online":false}"#.to_vec(),
            },
        );
        assert_eq!(packet[0], 0x10);
        let text = String::from_utf8_lossy(&packet);
        assert!(text.contains("MQTT"));
        assert!(text.contains("edge-1"));
        assert!(text.contains("sentra/v1/motor-lab-01/status"));
        assert!(text.contains("\"online\":false"));
    }

    #[test]
    fn publish_packet_is_qos1() {
        let packet = publish_packet("a/b", b"{}", 7);
        assert_eq!(packet[0], 0x32);
        assert!(packet.windows(2).any(|pair| pair == 7u16.to_be_bytes()));
        assert!(packet.windows(2).any(|pair| pair == b"{}"));
    }

    #[test]
    fn remaining_length_encodes_multi_byte_values() {
        assert_eq!(encode_remaining_length(127), vec![127]);
        assert_eq!(encode_remaining_length(128), vec![0x80, 0x01]);
    }
}
