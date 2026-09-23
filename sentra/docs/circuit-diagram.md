# Devre

Laboratuvar yükü küçük bir 12 V DC motordur. Motor akımı ESP32 pininden geçmez. Motorun kendi 12 V kaynağı vardır. ESP32, sensörler ve motor kaynağı yalnızca toprak hattında birleşir.

```text
 12 V kaynak (+) ----+---- motor (+) ----+
                     |                   |
                     |              INA219 şönt (0.1 ohm)
                     |                   |
                     +---- INA219 Vin- --+---- motor (-) ---- 12 V kaynak (-)
                                                          \
                                                           +---- ortak GND ---- ESP32 GND

 ESP32 3V3 ---- MPU6050 VCC (0x68)
 ESP32 3V3 ---- INA219 VCC (0x40)
 ESP32 3V3 ---- DS18B20 VDD
 ESP32 3V3 ---- Hall A3144 VCC

 ESP32 GPIO21 (SDA) ---- MPU6050 SDA ---- INA219 SDA
 ESP32 GPIO22 (SCL) ---- MPU6050 SCL ---- INA219 SCL
 ESP32 GPIO4  ---- DS18B20 DQ ---- 4.7 kΩ ---- 3V3
 ESP32 GPIO27 ---- Hall çıkışı ---- 10 kΩ ---- 3V3
```

| Sinyal | ESP32 | Not |
| --- | --- | --- |
| I2C SDA | GPIO21 | MPU6050 ve INA219 aynı hatta |
| I2C SCL | GPIO22 | |
| DS18B20 | GPIO4 | 4.7 kΩ pull-up 3V3'e |
| Hall | GPIO27 | kesme pini, 10 kΩ pull-up |
| INA219 şönt | motor hattı | 0.1 Ω, beklenen tam ölçek 3.2 A |

INA219 kalibrasyon yazmacı bu şönt ve 3.2 A için 4194'tür. Hesap TI veri sayfasındaki `trunc(0.04096 / (current_lsb * Rshunt))` formülüdür ve `sentra-core` testi bu sayıyı kilitler. Akım LSB değeri `3.2 / 32768` amperdir.

DS18B20 12 bit scratchpad örneği `0x0191` = 25.0625 °C. Hall testi: 1 saniyede 25 darbe ve devir başına 1 darbe = 1500 rpm. MPU6050 ±2 g ölçeğinde 16384 sayım = 1 g.

Fritzing veya KiCad sayfası kart basılacağı zaman eklenir. Bu dosyadaki bağlantı, o sayfanın kaynağıdır; ölçülmemiş bir PCB dosyası üretilmedi.

Güvenlik: motoru ESP32'nin 5 V veya 3V3 pininden beslemeyin. İlk çalıştırmada şönt yönünü ve ortak toprağı multimetre ile doğrulayın.
