# Cursor İçin Ana Proje Promptu — ESP32 Tabanlı Motor Öngörücü Bakım Sistemi

Bu dosyayı Cursor'da proje klasörünün köküne `PROJECT_BRIEF.md` olarak kaydet ve her yeni Cursor sohbetinin (chat/composer) başında Cursor'a "önce PROJECT_BRIEF.md dosyasını oku, sonra devam et" de. Böylece Cursor her seferinde projenin bütününü hatırlar, parça parça kod yazarken bağlamı kaybetmez.

Aşağıdaki metni olduğu gibi Cursor'a yapıştırabilirsin.

---

## KOPYALANACAK PROMPT (Cursor'a yapıştır)

Sen benim yazılım mimarım ve kıdemli pair-programming ortağımsın. Birlikte, bir ESP32 tabanlı endüstriyel motor öngörücü bakım (predictive maintenance) sistemi kuruyoruz. Bu proje benim için iki amaç taşıyor: (1) gerçek, çalışan bir mühendislik sistemi, (2) TÜBİTAK proje kalitesinde ve iş başvurularında portföy olarak gösterebileceğim bir GitHub deposu. Bu yüzden kod kalitesi, dokümantasyon ve mühendislik disiplini "hızlı prototip" değil, "profesyonel/araştırma kalitesi" seviyesinde olmalı.

### Projenin amacı

12V DC motor üzerinde, ESP32 aracılığıyla titreşim (MPU6050), sıcaklık (DS18B20), akım/gerilim (INA219) ve RPM (Hall effect sensör) verisini gerçek zamanlı topluyoruz. Bu veriyi MQTT ile bilgisayara aktarıp PostgreSQL'de saklıyoruz, sinyal işleme (FFT, RMS, kurtosis, peak) ile özellik çıkarıyoruz, makine öğrenmesiyle anormal çalışma durumlarını (overload, unbalanced, high vibration) tespit ediyoruz ve sonucu bir dashboard'da gösteriyoruz.

### Teknoloji yığını (kesinlikle bunları kullan, alternatif önerme)

- **Firmware:** ESP32 (Arduino framework veya PlatformIO), C++
- **Mesajlaşma:** MQTT (Mosquitto broker, Docker içinde)
- **Backend:** Python 3.11+, FastAPI, paho-mqtt, SQLAlchemy
- **Veritabanı:** PostgreSQL (Docker içinde)
- **Sinyal işleme / ML:** NumPy, SciPy, pandas, scikit-learn, xgboost, (ileride SHAP)
- **Dashboard:** React (veya proje hızlanana kadar basit FastAPI + Jinja2 template, sonra React'e geçiş)
- **Altyapı:** Docker + docker-compose
- **Versiyon kontrol:** Git, anlamlı commit mesajlarıyla

### Mimari prensipler

1. **Katmanlı mimari** kur: firmware / veri toplama / veri işleme / ML / sunum katmanları birbirinden net şekilde ayrılmalı. Her katman kendi klasöründe, kendi sorumluluğuyla dursun (`firmware/`, `backend/`, `ml/`, `dashboard/`, `docs/`, `data/`).
2. **Konfigürasyon kodun içine gömülmesin.** WiFi şifresi, MQTT broker adresi, DB bağlantı bilgisi gibi değerler `.env` dosyasında veya `config.yaml`'da tutulmalı, `.gitignore`'a eklenmeli. Bunun yerine `.env.example` dosyası ile hangi değişkenlerin gerektiğini göster.
3. **Her modül test edilebilir olmalı.** Backend ve ML tarafında `pytest` ile birim testleri yaz. En azından: MQTT mesaj parse fonksiyonları, sinyal işleme fonksiyonları (RMS/FFT/kurtosis hesaplamaları), veritabanı modelleri için temel testler olsun.
4. **Hata yönetimi ciddi olsun.** ESP32 firmware'i WiFi/MQTT bağlantısı koptuğunda yeniden bağlanmayı denemeli, hata durumlarını seri port'a loglamalı. Backend tarafında da exception handling ve loglama (`logging` modülü) olsun, sessizce çökmesin.
5. **Kod stiline dikkat et:** Python tarafında PEP8, tip belirteçleri (type hints) kullan, fonksiyonlara docstring yaz. C++ tarafında anlamlı değişken isimleri ve yorum satırları kullan. Kısa değişken isimleri (`x`, `tmp`, `data2`) kullanma, açıklayıcı isimler kullan (`vibration_rms`, `motor_temp_celsius`).
6. **Her aşamayı bağımsız çalışır halde teslim et.** Yani örneğin ESP32 firmware'i tek başına, backend hazır olmadan bile derlenip seri port'a veri basabilmeli. Backend, ESP32 gerçekten bağlı olmasa bile mock/test verisiyle çalışabilmeli. Bu, hem geliştirmeyi kolaylaştırır hem de "bu sistemi parça parça test edebilmiş" izlenimi verir.

### Aşamalı geliştirme planı (V1'den V8'e — sırayla ilerle, atlama yapma)

- **V1 — Veri toplama:** ESP32 firmware, sensör okuma, MQTT'ye ham veri yayınlama. Backend'in bu veriyi dinleyip PostgreSQL'e yazması.
- **V2 — Özellik mühendisliği:** Ham titreşim verisinden RMS, peak, kurtosis, FFT hesaplama. Bu hesaplamaları ayrı, test edilebilir bir Python modülünde tut (`ml/features.py` gibi).
- **V3 — Anomali tespiti:** Denetimsiz (unsupervised) bir anomali tespit modeli (örn. Isolation Forest) ile "normal dışı" durumları yakalama.
- **V4 — Arıza sınıflandırma:** Etiketlenmiş verilerle (NORMAL / OVERLOAD / UNBALANCED / HIGH VIBRATION) denetimli (supervised) sınıflandırma modeli (xgboost).
- **V5 — Sağlık skoru:** Sensör verilerinden 0-100 arası bir "motor sağlık skoru" üretme mantığı.
- **V6 — Açıklanabilir AI:** SHAP ile modelin hangi özelliğe göre karar verdiğini gösterme.
- **V7 — Bakım yönetimi:** Anomali geçmişini ve öneri/alarm mekanizmasını dashboard'a ekleme.
- **V8 — Çoklu motor desteği:** Sistemi birden fazla cihazı aynı anda izleyecek şekilde genişletme.

Her V aşaması bittiğinde bana şunu sun: (a) o aşamada ne yapıldığının kısa özeti, (b) nasıl test edileceği, (c) `docs/` klasörüne eklenecek kısa bir Markdown notu. Bir sonraki aşamaya geçmeden önce bana onaylat.

### Dokümantasyon beklentisi (TÜBİTAK/iş başvurusu kalitesi için kritik)

- Kök dizinde kapsamlı bir `README.md` olsun: proje amacı, mimari diyagram (ASCII veya Mermaid), kurulum adımları, kullanılan teknolojiler, ekran görüntüleri/örnek çıktılar için yer tutucular.
- `docs/` klasöründe her V aşaması için ayrı bir not dosyası (`docs/v1-data-collection.md` gibi) — ne yapıldığı, neden o teknik tercih yapıldığı, karşılaşılan sorunlar ve çözümleri.
- Devre şeması ve bağlantı diyagramını `docs/circuit-diagram.md` içinde metin/ASCII ve mümkünse Fritzing/KiCad dosyası olarak tut.
- Her önemli mimari karar için kısa bir "neden bunu seçtim" notu (örneğin "neden xgboost, neden Isolation Forest").

### Kod üretirken izlemeni istediğim çalışma şekli

- Büyük bir özelliği tek seferde yazmak yerine, önce planını bana kısaca özetle, onaylarsam kodu yaz.
- Bir dosyayı değiştirirken ilgisiz kısımlara dokunma, değişikliği minimal ve odaklı tut.
- Her yeni bağımlılık eklediğinde `requirements.txt` (Python) veya `platformio.ini`/kütüphane listesini güncelle.
- Kodun çalıştığını varsaymak yerine, mümkün olduğunca benimle birlikte terminal üzerinden test et (docker compose up, pytest çalıştırma, ESP32 seri port çıktısı kontrolü gibi).
- Ben aksini söylemedikçe kısayol/geçici çözüm (hardcoded değer, TODO bırakılmış eksik fonksiyon) üretme; üretmek zorunda kalırsan açıkça `# TODO:` yorumuyla işaretle ve bana söyle.

Şimdi önce projenin klasör yapısını oluştur, `README.md` iskeletini yaz ve V1 aşamasının ilk adımı olan ESP32 firmware'ine başlamak için bana kısa bir plan sun. Kod yazmadan önce planı onaylamamı bekle.

---

## Kullanım notları

1. Bu dosyayı `PROJECT_BRIEF.md` olarak proje kök dizinine kaydet.
2. Cursor'da yeni bir chat açtığında ilk mesaj olarak: *"PROJECT_BRIEF.md dosyasını oku ve ona göre devam et"* yaz.
3. Uzun bir çalışma oturumundan sonra context dolarsa, yeni sohbette tekrar aynı dosyayı okutarak devam et — böylece proje tutarlılığı bozulmaz.
4. Her V aşaması bittiğinde `docs/` klasörüne not eklemeyi unutma; bu notlar hem senin öğrenme sürecini kayıt altına alır hem de mülakatlarda "projeni anlat" dendiğinde kullanacağın hazır malzeme olur.
