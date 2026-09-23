# Telemetri şeması, sürüm 1

`schema_version` 1 olmayan çerçeve reddedilir. Bilinmeyen JSON alanları yutulur; zorunlu alan eksikse sıfır değer doğrulamada düşer.

| Alan | Birim | Aralık |
| --- | --- | --- |
| device_id | slug | `^[a-z0-9](?:[a-z0-9-]{0,62}[a-z0-9])?$` |
| firmware_version | metin | 1–32, harf rakam `.` `_` `-` |
| captured_at_ms | unix ms | > 0 |
| sample_rate_hz | Hz | 100–8000 |
| window_samples | örnek | 16–2048, ikinin kuvveti |
| rpm | dev/dk | 0–20000 |
| current_a | A | 0–20 |
| voltage_v | V | 0–30 |
| temperature_c | °C | -40–150 |
| accel_rms_g | g | 0–50 |
| accel_peak_g | g | 0–50 |
| crest_factor | — | 0–100 |
| kurtosis | — | 0–100 |
| dominant_freq_hz | Hz | 0–örnek_hızı/2 |
| sensor_ok | bool | |
| wifi_rssi_dbm | dBm | yok sayılabilir, varsa -120–0 |

`condition` yalnızca etiket konusundadır: `NORMAL`, `OVERLOAD`, `UNBALANCED`, `HIGH_VIBRATION`.

Örnek: `testdata/frame_normal.json`.
