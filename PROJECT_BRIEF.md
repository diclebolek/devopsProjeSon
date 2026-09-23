# SENTRA — güncel proje brifingi

Yeni bir Cursor sohbetine başlarken ilk mesaj şu olsun: **PROJECT_BRIEF.md dosyasını oku ve ona göre devam et.**

Bu dosya projenin yaşayan brifingidir. İlk yapıştırılan prompt, değiştirilmeden `sentra/docs/original-prompt.md` içinde durur. Aşağıdaki mimari onun yerine geçer.

## Ne değişti ve neden

İlk metin, ESP32 + MQTT + PostgreSQL + scikit-learn şeklinde yaygın bir öğrenci yığını öneriyordu ve "alternatif önerme" diyordu. Kullanıcı bunun yerine daha özgün bir sistem, Rust veya Go ile daha iyi verim ve portföy kalitesi istedi.

SENTRA bu isteğin karşılığıdır. Amaç aynıdır: 12 V DC motorda titreşim, sıcaklık, akım, gerilim ve devir toplayıp anormal çalışmayı görünür kılmak. Verinin taşınma biçimi değişti.

- **Rust (`sentra/firmware`)** kenarda durur. Pencere ortalamasını çıkarır; RMS, tepe, crest factor, kurtosis ve baskın frekansı cihazda hesaplar. Ham pencereyi sürekli yayınlamaz.
- **Go (`sentra/pipeline`)** MQTT'den çerçeveyi alır, doğrular, TimescaleDB'ye yazar ve panoyu sunar. Eşzamanlı bağlantı ve düşük gecikmeli yazma bu katmanın işidir.
- **Python** ancak V2'den itibaren, xgboost ve SHAP'ın yaşadığı yerde devreye girer. Makine öğrenmesini Rust veya Go'ya taşımak bu projede verim değil, ekosistem kaybı olurdu.

Tek dil seçilmedi. Her dil, çöp toplayıcısız ve deterministik kenar (Rust) ile eşzamanlı servis (Go) arasındaki gerçek ayrımı taşır. Karar notu: `sentra/docs/decisions/001-language-split.md`.

## Bu depoda nerede durur

Kod `sentra/` altındadır. Kökteki `backend/`, `frontend/` ve kök `docker-compose.yml` Insucom sitesine aittir; SENTRA işi onları değiştirmez. Azure pipeline `sentra/` klasörünü canlı site senkronunun dışında tutar.

## Katmanlar

```text
sentra/firmware/     Rust kenar sözleşmesi, DSP, sensör kod çözücüleri, host düğüm
sentra/pipeline/     Go ingest, simülatör, HTTP API, pano
sentra/docs/         aşama notları ve karar kayıtları
sentra/schema/       telemetri sözleşmesi
sentra/testdata/     Rust ve Go'nun paylaştığı sayısal oracle
```

## V1 durumu

V1 veri toplama bu dalda çalışır durumda teslim edildi. Ayrıntı: `sentra/docs/v1-data-collection.md` ve `sentra/README.md`.

İzleme kuralı, literatür seçimi ve ihtiyaç listesi `sentra/README.md` ile `sentra/docs/methodology.md` içindedir. Sağlam oturumun dışındaki ISO titreşim bölgeleri bu motora uygulanmaz. Model rakamı, oturum bölmeli testten önce yazılmaz.

Sonraki aşamaya, kullanıcı onaylamadan geçilmez.

- **V2** Python'da aynı DSP formüllerinin referans modülü ve özellik kaydı
- **V3** Isolation Forest
- **V4** xgboost ile NORMAL / OVERLOAD / UNBALANCED / HIGH_VIBRATION
- **V5** 0–100 sağlık skoru
- **V6** SHAP
- **V7** bakım geçmişi ve alarm
- **V8** filo operasyonu (şema bugünden birden fazla `device_id` kabul eder; ürünleşmiş çoklu motor yönetimi V8'dir)

## Çalışma kuralları

- Konfigürasyon ortam değişkenlerindedir. Örnek: `sentra/.env.example`.
- Testler `cd sentra && make test`.
- Donanım yokken Go simülatörü ve Rust host düğümü aynı JSON sözleşmesini üretir.
- Sessiz çökme yok: kenar seri porta yazar ve MQTT'yi yeniden dener; ingest reddettiği mesajı loglar.
- İlerideki iş için yarım fonksiyon bırakılmayacak. Bilinçli olarak V1 dışında kalan şey kodda `TODO` diye değil, bu brifingde aşama olarak durur.
