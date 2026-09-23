# V1 — Veri toplama

## Ne yapıldı

Kenar, bir titreşim penceresini ham örnek olarak taşımak yerine tek bir telemetri çerçevesine indirger. Çerçeve; devir, akım, gerilim, sıcaklık, RMS, tepe, crest factor, kurtosis ve baskın frekansı taşır. Go ingest bu çerçeveyi doğrular, TimescaleDB hazırsa hipertabloya yazar, değilse aynı şemayla düz PostgreSQL kullanır. Bellek deposu Docker olmadan panoyu açar.

Simülatör dört çalışma rejimini aynı sözleşmeyle üretir: `NORMAL`, `OVERLOAD`, `UNBALANCED`, `HIGH_VIBRATION`. Etiket ayrı konuda gider ve ekranda "model kararı" diye gösterilmez.

Rust host düğümü seri çıktı verir, MQTT kopunca gecikmeyi 1 saniyeden 30 saniyeye kadar iki katına çıkarır ve yeniden bağlanır. Bağlantı paketine çevrimdışı vasiyet mesajı konur.

## Neden böyle

Sürekli ham ivme yayınlamak hem bant genişliğini hem de sonraki modeli kirletir. Yerçekimi ofseti çıkarılmazsa RMS yaklaşık 1 g görünür ve motor titreşimi kaybolur. Formül bu yüzden ortalaması alınmış pencere üzerindedir. Ayrıntı `dsp.md` dosyasında.

Python ingest yerine Go seçildi çünkü bu katmanda öğrenme kütüphanesi yok; eşzamanlı MQTT ve HTTP var. Öğrenme V2'ye bırakıldı, atlanmadı.

## Nasıl test edilir

```bash
cd sentra
make test
```

Bellekli uçtan uca:

```bash
cd sentra/pipeline
SENTRA_STORE=memory SENTRA_HTTP_ADDR=127.0.0.1:8088 go run ./cmd/ingest
```

```bash
SENTRA_TRANSPORT=http SENTRA_INGEST_URL=http://127.0.0.1:8088 \
  SENTRA_MAX_FRAMES=4 SENTRA_REGIME=cycle SENTRA_REGIME_PERIOD=1 \
  SENTRA_INTERVAL_MS=50 go run ./cmd/simulator
curl -s http://127.0.0.1:8088/api/v1/telemetry
```

Kenar:

```bash
cd sentra/firmware
SENTRA_MAX_FRAMES=2 cargo run -p sentra-edge
```

Tam yığın: `cd sentra && docker compose up --build`, sonra [http://127.0.0.1:8088](http://127.0.0.1:8088).

## Karşılaşılan sınırlar

- Bu ortamda Docker yoktu. Compose dosyası yazıldı; çalışan doğrulama `make test`, host düğümün JSON çıktısı ve bellek modundaki HTTP turudur.
- ESP32 Rust alet zinciri (espup / Xtensa) kurulu değil. Kart kodu, derlenmeyen bir hedef olarak eklenmedi. Test edilen çekirdek ile pin tablosu ayrıldı.
- `rumqttc` güncel bağımlılıkları Rust 1.85 istiyordu. Kenar, laboratuvar broker'ı için kendi MQTT 3.1.1 kodlayıcısını kullanır. Paket baytları birim testlidir. TLS bu istemcide yoktur; laboratuvar broker'ı düz metindir.

## Bilinçli olarak yok

Sağlık skoru, Isolation Forest, xgboost, SHAP, bakım emri ve filo ekranı sonraki V aşamalarıdır.
