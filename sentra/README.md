# SENTRA

12 V doğru akım motorunun gövde titreşimini, sıcaklığını, akımını, gerilimini ve devrini izleyen öngörücü bakım tezgahı. Amaç, motor bozulduktan sonra müdahale etmek değil; aynı motorun kendi sağlam kaydından sapmayı erken görmek.

Her saniye cihaz bir titreşim penceresi toplar, ortalamayı (yerçekimi ofsetini) çıkarır ve yalnızca özeti yayınlar: RMS, tepe, crest factor, kurtosis, baskın frekans, artı yavaş kanallar (devir, akım, gerilim, sıcaklık). Ham örnek hatta kalmaz. Pano bu çerçeveyi canlı gösterir. Model katmanı sonraki sürümlerdedir; aşağıdaki metodoloji o katmanın hangi yayına dayandığını ve hangi veriyle eğitileceğini kilitler.

Kod `sentra/` altındadır. Bu deponun `backend/` ve `frontend/` ağacı ayrı bir sigorta sitesidir. SENTRA onu değiştirmez.

## Neyi çözer

On iki voltluk bir motor tezgâhında dört işi aynı anda yapmak: ham titreşimi sürekli hatta basmamak, mekanik sapmayı yük sapmasından ayırmak, başka tezgâhın doğruluk yüzdesini ve ISO 20816 bölgesini bu gövdeye yapıştırmamak, donanım yokken aynı sözleşmeyle panoyu doldurmak.

Çıktı bir durum adıdır. Sağlam kayda göre sapma var mı; varsa yük mü, dengesizlik mi, gevşek montaj mı. “Kaç saat sonra müdahale” sorusu bu sürümün dışındadır. Kalan ömür gelecek çalışmadır (`docs/methodology.md`).

## Benzeri var

Var. Üç komşu aile duruyor. SENTRA onların yerine yeni bir fizik iddiası koymaz.

| Aile | Ne yapıyor | Bu tezgâhta durmayan yeri |
| --- | --- | --- |
| ESP32 öğrenci kitleri (JAAFR, JETIR, UPC tez özeti, IJEAST 2026) | MPU6050 veya ADXL345, INA219 veya ACS712, DS18B20, MQTT. | Sensör çantası ortaktır. Çoğu ham veya seyrek örnek basar. “%100 doğruluk” cümlesi oturum bölmeli test değildir. |
| Fabrika kenar izleme (Vermesan 2022) | Titreşim, sıcaklık, ses, akım. ISO 20816-1:2016 RMS hız bölgeleri. | Bölge tablosu makine sınıfına bağlıdır. 12 V laboratuvar motoru o sınıflara girmez. |
| CWRU özellik makaleleri (Branco 2024, Alonso-González 2023) | 12 kHz ve üzeri, rulman zarfı, SHAP veya kurtogram. | MPU6050 ile 1 kHz, BPFI/BPFO ayrımı hedeflemez. |

Ayrım şurada: özellik formülü Rust, Go ve sonra Python’da aynı oracle’a kilitlenir; erken uyarı yalnız sağlam veriyle Isolation Forest’tır, dört sınıf adı XGBoost’tur; rapor oturum bölmesi ve sınıf başına recall ister.

## Literatürden ne çıktı

Taranan işler iki aileye ayrılıyor. Birinci aile, kontrollü tezgâhta titreşim toplayıp sınıflandırıcı kıyaslıyor. İkinci aile, ESP32 ve birkaç sensörle MQTT’ye veri basan öğrenci uygulaması. Birinci ailenin ölçüm dili bu projeye girer. İkinci ailenin doğruluk cümleleri, oturum sızıntısı ve endüstriyel ivmeölçer olmadığı için rakam olarak devralınmaz.

| Kaynak | Düzen | Ne işe yaradı | SENTRA’ya etkisi |
| --- | --- | --- | --- |
| Kim, Lee, Wang, Lee. *Sensors* 2023, 23(5), 2585. [10.3390/s23052585](https://doi.org/10.3390/s23052585) | Asenkron motor, endüstriyel ivmeölçer, 1240 pencere × 1024 örnek, normal / rotor / rulman. SVM, çok katmanlı ağ, CNN, GBM, XGBoost. Katmanlı K-kat. | XGBoost (700 ağaç): normal %99,11, rotor %99,76, rulman %100. Test süresi ortalama 0,050 s. SVM 1,65 s. Yazarlar en yüksek doğruluğu SVM ve CNN’de, en yüksek hızı XGBoost’ta bulmuş. | Ham dalga üzerinde CNN bu tezgâhta kazanabiliyor. Biz ham dalgayı sürekli taşımıyoruz. Özellik tablosunda XGBoost hem hızlı hem SHAP ile açıklanabilir. CNN, ancak özellik modeli yetersiz kalırsa ve şüpheli pencerenin ham hali ayrıca kaydedilirse ikinci deney olur. |
| Ali ve diğerleri, *IEEE Trans. Energy Conversion*, 2024. [10.1109/TEC.2024.3405897](https://doi.org/10.1109/TEC.2024.3405897) | Asenkron motor, stator elektrik ölçümü + titreşim, altı durum (sağlam, kırık rotor çubuğu, dengesizlik ve sargı kısa devresi bileşimleri). | Ayarlanmış XGBoost için bildirilen doğruluk %96,06. Ayarsız modele göre yaklaşık %15 iyileşme, kıyaslanan tanıma göre %57 daha kısa çalışma. | Elektrik ve titreşim birlikte kullanılınca bileşik arıza ayrılıyor. SENTRA’da akım yükü, titreşim mekaniği taşır. |
| *FIT* 2024. [10.1109/FIT63703.2024.10838415](https://doi.org/10.1109/FIT63703.2024.10838415) | Değişken hızda rulman titreşimi. | XGBoost ortalama %98,06, SVM %98. | Devir değişince tek bir frekans eşiği yetmez. Baskın frekans, devirle birlikte okunur (1× hattı). |
| Titreşim ve akım karşılaştırması, *Structural Health Monitoring*, 2024. [10.1177/14759217241289874](https://doi.org/10.1177/14759217241289874) | Sağlam, rulman, hizasızlık. Ham veride CNN; FFT ve dalgacık özelliğinde rastgele orman ve XGBoost. | Mekanik arızada titreşim neredeyse %100. Aynı arızada, işlenmiş akım sinyali %87,41. | Titreşim birincil mekanik kanal. Akım onun yedeği değil, yük ve elektriksel sapma kanalı. |
| Branco ve diğerleri, *Mach. Learn. Knowl. Extr.* 2024, 6(1), 16. [10.3390/make6010016](https://doi.org/10.3390/make6010016) | CWRU rulman seti, 12 kHz, 1 s pencerede tepe-tepe, crest factor, RMS, uç değerler, standart sapma, kurtosis. SVM + SHAP. | SHAP, hangi zaman özelliğinin tespitte, sınıfta ve şiddette taşındığını seçmek için kullanılmış. | V6’daki SHAP bu işin açıklama adımıdır. Özellik listesi baştan bu zaman düzlemi ailesindedir. |
| Alonso-González ve diğerleri, *IEEE Access* 2023. [10.1109/ACCESS.2023.3283466](https://doi.org/10.1109/ACCESS.2023.3283466) | CWRU, zarf analizi ve kurtogram. | Zarf genlikleriyle karar ağacı ve KNN %100, kernel naive Bayes %94,4. | Rulman frekansları (BPFI, BPFO) ancak örnekleme ve mekanik model yeterse anlamlı. MPU6050 ile 1 kHz’de bu ayrım hedeflenmez. |
| Zacharia ve diğerleri, *Sensors* 2022, 22(24), 9658. [10.3390/s22249658](https://doi.org/10.3390/s22249658) | Fırçasız DC motor, ses, “sağlam / kırık / ağır yük”, model mikrodenetleyicide. | Kenarda çıkarım, buluta gidiş gecikmesini ve kopan ağı tanımanın dışına iter. | Özeti kenarda hesaplamak bu sonucun hafif hali. Sınıflandırıcıyı ESP32’ye gömmek V1 işi değil. |
| Vermesan ve diğerleri, *Front. Chem. Eng.* 2022. [10.3389/fceng.2022.900096](https://doi.org/10.3389/fceng.2022.900096) | Fabrika motorları, titreşim + sıcaklık + ses + akım/gerilim, kenar işleme. | ISO 20816-1:2016, gövde titreşimini RMS hız ile sınıflar. Sınıf, makine boyuna ve montaja bağlıdır. | Gösterge ailesi (RMS) alınır. Bölge tablosu 12 V motora yapıştırılmaz. |
| ESP32 öğrenci raporları (JAAFR, JETIR, UPC tez özeti, IJEAST 2026) | MPU6050 veya ADXL345, INA219 veya ACS712, DS18B20, MQTT, bazen ThingSpeak. | Aynı sensör çantası laboratuvarda kurulabilir. Raporların “%100 doğruluk” cümleleri kontrollü kıyas değil. | Donanım listesi buradan gelir. Model iddiası birincil yayınlardan gelir. |

Allegro uygulama notu AN296276, fırçalı DC motorda hat akımının yük (tork) ve sargı arızası için kullanıldığını yazar. Asenkron motordaki yan bant akım imzası (MCSA) bu motora taşınmaz.

## Bu projede kullanılacaklar

Bunlar seçilmiş listedir. Yerine başka yığın önerilmez.

| Katman | Seçim | Yayınla bağı |
| --- | --- | --- |
| Kenar | ESP32, Rust, `sentra-core` | Özet kenarda kalsın (Zacharia 2022’nin gerekçesi). Çöp toplayıcı örneklemeyi kaydırmasın. |
| Titreşim | MPU6050, pencere ortalaması alınmış RMS, tepe, crest factor, kurtosis, baskın FFT kutusu | Kim 2023, Branco 2024, CWRU zaman özellikleri. 1 kHz, 256 örnek. |
| Sıcaklık | DS18B20 | Isıl sapma, titreşimin görmediği yavaş arıza. |
| Akım ve gerilim | INA219, 0,1 Ω, 3,2 A tam ölçek, kalibrasyon yazmacı 4194 | Yük kanalı (Ali 2024, Allegro). |
| Devir | Hall A3144, telemetride ayrı `rpm` alanı | 1× hattı `rpm/60`. RMS’i bölmek için kullanılmaz (FIT 2024). |
| Taşıma | MQTT 3.1.1, Mosquitto, konu `sentra/v1/{cihaz}/{telemetry,label,status}` | Öğrenci sistemlerinin ortak taşıması. Yük, özellik çerçevesi olduğu için küçük. |
| Kayıt | TimescaleDB (uzantı yoksa düz PostgreSQL) | Zaman serisi. InfluxDB alternatifleri görüldü; SQL ve mevcut testler için Timescale seçildi. |
| Toplama ve pano | Go | Eşzamanlı abone ve HTTP. Model dili değil. |
| Erken uyarı | Isolation Forest, yalnız sağlam oturum | “Sağlam bulutun dışında mı?” Sınıf adı vermez. Etiket birikmeden de çalışır. |
| Sınıflandırma | XGBoost: `NORMAL`, `OVERLOAD`, `UNBALANCED`, `HIGH_VIBRATION` | “Gösterilen dört rejimden hangisi?” (Kim 2023, Ali 2024). Gösterilmemiş arızaya isim koyamaz. |
| Açıklama | SHAP, sınıf dengesi kurulduktan sonra | Branco 2024. Sağlam çoğunluğun açıklaması diye kullanılmaz. |
| Sağlık skoru | 0–100, sağlam merkeze uzaklık ve sınıf olasılığı | ISO bölge rakamı değil. |

Python 3.11 bu listedeki Isolation Forest, XGBoost ve SHAP için V2’den itibaren girer. Rust ve Go onların yerine geçmez.

## Nasıl izlenir

Pano her çerçevede şunları gösterir: rejim etiketi (simülatör veya teknisyen; model kararı değildir), devir, sıcaklık, akım, gerilim, RMS, tepe, kurtosis, baskın frekans. MQTT yoksa bile HTTP ile aynı çerçeve yazılır.

İzlemenin kuralı mutlak eşik değil, sağlam oturumdur.

| Gözlem | Sağlam kayda göre anlamı |
| --- | --- |
| Akım artar, devir düşer, gerilim oturur veya biraz düşer, sıcaklık dakikalar içinde yükselir | Aşırı yük |
| RMS ve 1× (devir/60) büyür, kurtosis sinüse yakın (~1,5) kalır | Dengesiz kütle |
| RMS ve kurtosis birlikte büyür, baskın frekans tek bir hatta oturmaz | Gevşek montaj veya vuruşlu titreşim |
| Sıcaklık, akım sakinken yükselir | Sürtünme veya ortam; titreşim tek başına yetmez |
| `sensor_ok` yanlış veya çerçeve kesilir | Kablo, I2C veya besleme; motor arızası sayılmaz |

Sağlık skoru V5’tedir. O gelene kadar pano skoru uydurmaz. Alarm eşiği, sağlam oturumun RMS ve akım dağılımından hesaplanacak; ISO 20816 bölgesi yapıştırılmayacak.

Canlı konular:

- `sentra/v1/{device_id}/telemetry` ölçüm
- `sentra/v1/{device_id}/label` oturum etiketi
- `sentra/v1/{device_id}/status` çevrimiçi, vasiyet mesajıyla çevrimdışı

## Metodoloji

Sıra atlanmaz. Her adım bir öncekinin kaydı durunca başlar.

1. Tezgâh. Motor 12 V kaynağı ESP32 pininden geçmez. Ortak toprak vardır. Sensör gövdeye sabitlenir, yön not edilir.
2. Sağlam temel. En az on dakika, yük yok, kapak sıkı. Bu kayıt eşiklerin ve Isolation Forest’ın tek eğitimidir.
3. Kontrollü rejim. Her rejim ayrı oturum: kısa süreli fren (aşırı yük), mile küçük kütle (dengesizlik), gevşetilmiş ayak (yüksek titreşim). Etiket oturum adından gelir.
4. Kenar özellik. Aynı formül Rust’ta cihazda, Go simülatöründe ve sonra Python’da. Oracle `testdata/dsp_case.json`.
5. Bölme. Test oturumu eğitime girmez. Pencereyi rastgele karıştırmak yok.
6. Isolation Forest yalnız sağlam pencerede öğrenir. Bu erken uyarıdır. Arıza penceresi “görülmemiş” ise şüpheli sayılır. Sınıf adı bu adımda yazılmaz.
7. XGBoost’tan önce sınıf başına pencere sayısı yazılır. Sağlam oturum uzun, arıza oturumu kısa kalırsa ya sağlam pencereler eşit aralıkla alt örneklenir ya da sınıf başına örnek ağırlığı konur. Tek doğruluk yüzdesi bu dengesizliği gizler. Ayrıntı `docs/methodology.md`.
8. XGBoost dört sınıfı oturum bölmesiyle öğrenir. Rapor sınıf başına recall ve karışıklık matrisidir. Isolation Forest sapma deyip XGBoost `NORMAL` derse sonuç bilinmeyen sapmadır.
9. SHAP, 7. adımdan sonra, yüksek çıkan sınıf için hangi alanın ittiğini gösterir.
10. Sağlık skoru ve alarm, ancak 6 ve 8 kendi test oturumunda ayrımı bozmadan durursa panoya yazılır. Kalan ömür bu adımların arasında yoktur.

```mermaid
flowchart TD
  bench["Tezgah ve guvenlik"] --> baseline["Saglam oturum"]
  baseline --> faults["Yuk, dengesizlik, gevsek montaj"]
  faults --> edge["Kenar ozellik"]
  edge --> store["MQTT ve TimescaleDB"]
  store --> split["Oturum bazli bolme"]
  split --> iforest["Erken uyari: Isolation Forest"]
  split --> balance["Sinif penceresini esitle"]
  balance --> xgb["Siniflandirma: XGBoost"]
  iforest --> gate["Sapma var mi"]
  xgb --> gate
  gate --> score["Saglik skoru"]
  xgb --> shap["SHAP"]
  score --> board["Pano"]
  shap --> board
```

Aynı akışın bugünkü yazılım karşılığı:

```mermaid
flowchart LR
  sensors["MPU6050, DS18B20, INA219, Hall"] --> edge["ESP32 / Rust"]
  edge --> broker["Mosquitto"]
  sim["Go simulator"] --> broker
  broker --> ingest["Go ingest"]
  ingest --> db["TimescaleDB"]
  ingest --> ui["Pano"]
  db --> py["Python V2+"]
```

Ayrıntı: `docs/methodology.md`, formüller `docs/dsp.md`, kararlar `docs/decisions/`.

## Gerçeklemek için ihtiyaç listesi

Parça, adet, KDV dahil fiyat ve satıcı bağlantısı `docs/malzeme-listesi.md` içindedir. 23 Eylül 2026 vitrinine göre doğrulanan satırların toplamı 1.296,67 TL’dir. Kargo ve stok o dosyada toplama dahil değildir.

Motor, boşta 60 mA ve zorlanma 0,45 A çeken 12 V 280 rpm redüktörlü gövdedir. Zorlanma akımı INA219 tam ölçeğinin (3,2 A) ve 12 V 2 A adaptörün içindedir. Motor akımı ESP32 pininden geçmez.

Bağlantı: `docs/circuit-diagram.md`. Elde yoksa ayrıca duranlar: mikro USB kablo, 5,5×2,1 mm dişi jak, multimetre, mengene, mile kelepçelenecek somun. Jak fiyatı malzeme dosyasında kilitlenmedi.

### Alet ve yazılım

- Rust 1.83 veya üzeri, Go 1.22, Docker Compose
- Python 3.11, V2 ile: NumPy, SciPy, scikit-learn, xgboost, shap
- Mosquitto ve TimescaleDB imajları Compose içindedir
- Seri port için USB sürücüsü (CP2102 veya CH340, kartın köprüsüne göre)

### Veri

- En az bir sağlam oturum, on dakika
- Her arıza rejiminden ayrı oturum. Süre, ısınma motora zarar vermeyecek kadar kısa, pencereler istatistik için yeterince çok. Kısa arıza oturumu, eğitimde sağlam sınıfın pencere sayısıyla eşitlenir veya ağırlıklanır (`docs/methodology.md`)
- Her oturumun etiketi `NORMAL`, `OVERLOAD`, `UNBALANCED` veya `HIGH_VIBRATION`
- Teste ayrılmış ve eğitime hiç girmemiş ikinci bir sağlam oturum ve ikinci bir arıza oturumu

Simülatör bu oturumlar gelene kadar panoyu doldurur. Simülatör kaydı, tezgâh doğruluğu diye raporlanmaz.

## Kurulum

```bash
cd sentra
make test
```

Donanımsız pano:

```bash
cd sentra/pipeline
SENTRA_STORE=memory SENTRA_HTTP_ADDR=127.0.0.1:8088 go run ./cmd/ingest
```

```bash
cd sentra/pipeline
SENTRA_TRANSPORT=http SENTRA_INGEST_URL=http://127.0.0.1:8088 \
  SENTRA_REGIME=cycle SENTRA_REGIME_PERIOD=1 SENTRA_INTERVAL_MS=1000 \
  go run ./cmd/simulator
```

Pano: [http://127.0.0.1:8088](http://127.0.0.1:8088)

Kenar düğümünün dizüstü karşılığı:

```bash
cd sentra/firmware
SENTRA_MAX_FRAMES=3 cargo run -p sentra-edge
```

`SENTRA_MQTT_BROKER=mqtt://127.0.0.1:1883` aynı çerçeveyi QoS 1 ile basar.

Tam yığın, bu depodaki diğer Postgres ile çarpışmasın diye 5433 ve 8088 kullanır:

```bash
cd sentra
docker compose up --build
```

Değişkenler `.env.example` içindedir. Compose parolası `sentra` yalnızca yerel laboratuvardır.

## Düzen

```text
firmware/sentra-core   DSP, sensör kod çözücüleri, çerçeve doğrulama
firmware/sentra-edge   host düğüm ve MQTT 3.1.1 istemcisi
firmware/esp32         kart pinleri
pipeline/cmd/ingest    MQTT, HTTP, pano
pipeline/cmd/simulator dört rejim
docs/methodology.md    izleme kuralı, Hall kanalı, sınıf dengesi
docs/malzeme-listesi.md parça, fiyat, bağlantı
testdata/              Rust ve Go oracle
```

ESP32, `sentra-core` fonksiyonlarını çağıracak. Register dönüşümleri kart olmadan test edilir. Flash adımı `firmware/esp32/README.md` içindedir.

Sözleşme `schema/telemetry.md`. Örnek çerçeve `testdata/frame_normal.json`.

## Sınırlar

- MPU6050, Kim 2023’teki endüstriyel ivmeölçer değildir. 1 kHz örnekleme, 12 V motorun dönüş frekansını ve birkaç harmoniği görür. CWRU’daki 12–48 kHz rulman zarfı bu kartın iddiası değildir.
- Yayımlanmış %96–%100 bandı başka tezgâhın, başka arızasının sonucudur. Bu motor için rakam, oturum bölmeli testten önce yazılmaz.
- ISO 20816 bölgeleri bu gövdeye uygulanmaz.
- Hall `rpm` alanı RMS’i normalize etmez. 1× yorumu `rpm/60` iledir.
- Isolation Forest erken uyarı, XGBoost dört sınıf adıdır. İkisi aynı rozet değildir.
- XGBoost, sınıf başına pencere eşitlenmeden veya ağırlık yazılmadan raporlanmaz. SHAP ondan sonradır.
- Kalan ömür yoktur. Sistem durum tanır. Müdahale zamanı gelecek çalışmadır.
- V1’de sağlık skoru, Isolation Forest, XGBoost ve SHAP yoktur. Sıra `docs/v1-data-collection.md` ile başlar. V2, `docs/dsp.md` formüllerinin Python karşılığıdır ve onay beklemektedir.

## Kaynakça

Rakamlar kaynak tezgâhın sonucudur. Bu motor için doğruluk, oturum bölmeli testten önce yazılmaz.

1. Kim, M.-C.; Lee, J.-H.; Wang, D.-H.; Lee, I.-S. “Induction Motor Fault Diagnosis Using Support Vector Machine, Neural Networks, and Boosting Methods.” *Sensors* 2023, 23, 2585. https://doi.org/10.3390/s23052585 — asenkron motor titreşiminde XGBoost hızı ve tablo özelliği.
2. Ali ve diğerleri. *IEEE Transactions on Energy Conversion*, 2024. https://doi.org/10.1109/TEC.2024.3405897 — titreşim ile elektrik ölçümünün birlikte kullanımı, ayarlanmış XGBoost.
3. Değişken hızda rulman titreşimi. *FIT* 2024. https://doi.org/10.1109/FIT63703.2024.10838415 — sabit hertz eşiğinin yetmemesi.
4. Titreşim ve akım karşılaştırması. *Structural Health Monitoring*, 2024. https://doi.org/10.1177/14759217241289874 — mekanik arızada titreşim, yükte akım.
5. Branco ve diğerleri. *Machine Learning and Knowledge Extraction* 2024, 6(1), 16. https://doi.org/10.3390/make6010016 — zaman özellikleri ve SHAP.
6. Alonso-González ve diğerleri. *IEEE Access* 2023. https://doi.org/10.1109/ACCESS.2023.3283466 — zarf ve kurtogram. 1 kHz MPU6050 bu ayrımı hedeflemez.
7. Zacharia ve diğerleri. *Sensors* 2022, 22(24), 9658. https://doi.org/10.3390/s22249658 — çıkarımın kenarda durması.
8. Vermesan ve diğerleri. *Frontiers in Chemical Engineering* 2022. https://doi.org/10.3389/fceng.2022.900096 — ISO 20816-1:2016. Bölge tablosu bu gövdeye uygulanmaz.
9. Allegro AN296276. Fırçalı DC motorda hat akımı, yük ve sargı içindir. https://www.allegromicro.com/-/media/files/application-notes/an296276-current-sensing-in-motor-drives.pdf
10. ESP32 sensör çantası raporları, doğruluk iddiası olarak değil: JAAFR https://www.rjwave.org/jaafr/papers/JAAFRTH00034.pdf — JETIR https://www.jetir.org/papers/JETIR2511506.pdf — UPC tezi https://hdl.handle.net/2117/328861

## Test

`make test` şunları çalıştırır: DSP oracle’ı, INA219 kalibrasyonu 4194, Hall RPM, DS18B20, çerçeve doğrulama, MQTT paket baytları, konu ayrıştırma, bellek deposu, HTTP ingest, rejim ayrımı.

Postgres turu `SENTRA_TEST_DATABASE_URL` ile açılır. Değişken yokken atlanır.
