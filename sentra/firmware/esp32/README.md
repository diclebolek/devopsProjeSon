# ESP32 taşıması

V1'de kart alet zinciri bu ortamda yoktur. Derlenmeyen bir `esp-idf` hedefi eklenmedi. Kartta çalışacak matematik ve register dönüşümleri `sentra-core` içindedir ve host testlerinden geçer.

Taşıma şunları çağırır:

- `accel_g_from_raw` — MPU6050, ±2 g, 16384 sayım/g, GPIO21/22
- `ds18b20_celsius_from_raw` — GPIO4, 4.7 kΩ
- `ina219_calibration_register(Ina219Config::LAB_MOTOR)` — kalibrasyon 4194
- `ina219_bus_voltage_volts` ve `ina219_current_amps`
- `rpm_from_pulses` — GPIO27 Hall kesmesi
- `Frame::from_window` — yayınlanacak JSON
- `Backoff` — Wi-Fi ve MQTT kopunca 1 s, 2 s, 4 s … 30 s

Host düğüm (`cargo run -p sentra-edge`) aynı JSON'u seri porta basar. Kart bağlanınca fark, örneklerin mock bitki yerine bu fonksiyonlardan gelmesidir.

Kurulum, kartın başına geçildiğinde:

```bash
cargo install espup
espup install
```

Ardından `esp-idf-hal` ile I2C ve GPIO kesmesi bu fonksiyonlara bağlanır. Pin tablosu `docs/circuit-diagram.md` dosyasındadır.
