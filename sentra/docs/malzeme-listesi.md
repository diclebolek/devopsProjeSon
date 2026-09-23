# Malzeme listesi

Fiyatlar **23 Eylül 2026** tarihinde ürün sayfasından okundu. Tutarlar **KDV dahildir**. Kargo, stok ve satıcının sepet alt limiti bu toplama girmez. Sayfa ertesi gün değişebilir; siparişten önce aynı bağlantıdaki tutar yeniden bakılır.

Bu liste tezgâhı kurmak içindir. Yazılım yığını ve oturum kuralları `../README.md` içindedir. Bağlantı şeması `circuit-diagram.md` içindedir.

## Alınacak parçalar

| Parça | Adet | Neden | KDV dahil | Sayfa |
| --- | ---: | --- | ---: | --- |
| ESP32 geliştirme kartı, 38 pin, CP2102 | 1 | Kenar düğüm. 38 pin breadboard’a oturur. | 271,03 TL | [Robotistan](https://www.robotistan.com/esp32-wifi-bluetooth-gelistirme-karti-38-pin) |
| MPU6050 GY-521 | 1 | Gövde titreşimi. ±2 g ölçeği yazılımda 16384 sayım/g. | 190,28 TL | [Robotistan](https://www.robotistan.com/mpu6050-6-eksen-ivme-ve-gyro-sensoru-6-dof-3-axis-accelerometer-and-gyros) |
| Su geçirmez DS18B20 | 1 | Gövde sıcaklığı. Veri hattına 4,7 kΩ gerekir. | 75,52 TL | [Robotistan](https://www.robotistan.com/su-gecirmez-ds18b20-dijital-isi-sensoru) |
| INA219 modülü, 0,1 Ω şönt, ±3,2 A | 1 | Hat akımı ve bara gerilimi. Ayrı şönt alınmaz. | 108,00 TL | [Robolink](https://www.robolinkmarket.com/ina219-i2c-cift-yonlu-akim-sensoru-modulu) |
| A3144, TO-92S | 1 | Devir darbesi. Açık kolektör; 10 kΩ pull-up ile GPIO27. | 9,97 TL | [Robotistan](https://www.robotistan.com/a-3144-dip-hall-effect-sensoru) |
| 1/4 W 4,7 kΩ | 1 | DS18B20 veri hattı pull-up. | 0,24 TL | [Robotistan](https://www.robotistan.com/14w-4k7-direnc-paketi-10-adet) |
| 1/2 W 10 kΩ | 1 | A3144 çıkış pull-up. | 0,54 TL | [Robotistan](https://www.robotistan.com/12w-10k-10-adet) |
| 12 V 25 mm 280 rpm redüktörlü DC motor | 1 | Tezgâh motoru. Boşta 60 mA, zorlanma 0,45 A. | 255,10 TL | [Robotistan](https://www.robotistan.com/12v-25mm-280-rpm-reduktorlu-dc-motor) |
| 12 V 2 A adaptör, 5,5×2,1 mm, 1 m | 1 | Motor beslemesi. ESP32 USB’den beslenir; motor akımı kart pininden geçmez. | 232,50 TL | [Robotistan](https://www.robotistan.com/12v-2a-dc-adaptor-55mm-x-21mm-1-metre-kablolu) |
| Orta boy breadboard | 1 | Lehimsiz ilk kurulum. | 35,56 TL | [Robotistan](https://www.robotistan.com/orta-boy-breadboard) |
| 40’lı erkek-dişi jumper, 20 cm | 1 | Modül pinleri. | 55,45 TL | [Robotistan](https://www.robotistan.com/40-pin-ayrilabilen-disi-erkek-m-f-jumper-kablo-200-mm) |
| 40’lı erkek-erkek jumper, 20 cm | 1 | Breadboard köprüleri ve ortak toprak. | 57,08 TL | [Robotistan](https://www.robotistan.com/40-pin-ayrilabilen-erkek-erkek-m-m-jumper-kablo-200-mm) |
| Neodyum disk 10×2 mm | 1 | Mile yapışır; her turda Hall’u bir kez tetikler. | 5,40 TL | [Mıknatıs Market](https://www.miknatismarket.com.tr/10x2-mm-yuvarlak-guclu-neodyum-miknatis-cap-10-mm-kalinlik-2-mm-kdv-dahil-fiyat-pmu363) |

Doğrulanan satırların toplamı **1.296,67 TL**.

Direnç sayfalarının başlığı paket, gövde metni tek direnç diyor. Görünen tutar 1 TL’nin altındadır; sepet adedi siparişte kontrol edilir, bütçeyi taşımaz. Mıknatıs satıcısının kargo alt limiti parçanın kendisinden büyük olabilir.

INA219 bu taramada Robotistan vitrininde yoktu. Tutar Robolink indirimli fiyatıdır (üstü çizili 151,50 TL, indirimli 108,00 TL). Modül 0,1 Ω ve ±3,2 A yazar; laboratuvar kalibrasyonu da budur (`Ina219Config`, yazmaç 4194).

## Neden bu motor

Zorlanma akımı 0,45 A. INA219 tam ölçeği 3,2 A, adaptör 2 A. Kilit akımı tam ölçeğin içinde kalır. Daha büyük gövde (örneğin 37 mm) zorlanma akımını sayfada açık yazmıyorsa bu listeye alınmaz; şönt ve kalibrasyon yeniden hesaplanır.

280 rpm, 1× hattını yaklaşık 4,7 Hz’e koyar. 1000 Hz örnekleme ve 256 örneklik pencere bu hattı görür. Amaç rulman zarfı (BPFI/BPFO) değildir.

Motor uzun süre kilitli bırakılmaz. Aşırı yük oturumu kısa tutulur, sıcaklık izlenir.

## Bu toplamda olmayanlar

- Mikro USB kablo. 38 pin kartın sayfası kabloyu paket içeriği diye yazmıyor. Elde yoksa ayrıca alınır; bu taramada fiyat satırı yok.
- 5,5×2,1 mm dişi jak veya klemens. Adaptör fişlidir, motor iki tellidir. Fiş motora doğrudan girmez. Bu taramada jak fiyatı kilitlenmedi.
- Multimetre, tahta veya mengene, mile kelepçelenecek somun, kısa fren için kayış. Atölyede varsa alınmaz.
- Hall modülü alternatifi: KY-003, sayfada 29,98 TL. [Robotistan](https://www.robotistan.com/hall-effect-manyetik-alan-sensoru). Devre şeması çıplak A3144 ve 10 kΩ ile çizildiği için toplam çıplak sensörü kullanır. Modül alınırsa pull-up’ın kartta olup olmadığı ölçülür; ikinci 10 kΩ paralel bağlanmaz.

## Yazılım

Rust, Go, Python, Mosquitto ve TimescaleDB imajları ücretli parça değildir. Sürümler `../README.md` kurulum bölümündedir.
