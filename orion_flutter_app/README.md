# Orion Gym & Fitness (Flutter)

Çok sektörlü randevu / işletme yönetim uygulamasının **spor salonu (gym)** varsayılanlı mobil & web istemcisi.

> **Önemli — veri kaynağı:** Bu proje **normalde Supabase’e bağlıdır** (Auth, işletme, randevu, içerik tabloları).  
> Aşağıda anlatılan **SQLite / yerel depolama yalnızca geçici** bir katmandır (demo, offline, “Beni hatırla”). Üretim verisi ve kalıcı iş kuralları için **Supabase** kullanılır.

---

## Özellikler

- **Müşteri:** Ana sayfa (galeri, hizmetler, neden bizi seçmelisiniz), randevu, profil
- **Admin:** Randevu yönetimi, çalışan yönetimi, performans, işletme / profil ayarları
- **Şifresiz demo giriş:** E-posta + boş şifre ile müşteri veya admin paneline giriş
- **Beni hatırla:** E-posta ve rol yerel olarak saklanır (yeniden açılışta doldurulur)
- **Çok dil:** TR / EN (`LanguageProvider`)
- **Tema:** Açık / koyu; aksan rengi `#228AE6`
- **UI:** Dialog’lar saydam + mavi border; hizmet kartlarında hover/lift; why-us kartlarında mavi border

---

## Mimari

| Katman | Teknoloji | Not |
|--------|-----------|-----|
| UI | Flutter 3.x (Material) | `orion_flutter_app/` |
| Durum | Provider | Auth, tema, dil |
| **Asıl backend** | **Supabase** | Auth, Postgres, Storage |
| Geçici yerel | SQLite (`sqflite`) / Web’de SharedPreferences | Sadece demo & remember-me |
| Demo veri | `DemoData` | Supabase yanıt vermezse UI dolu kalsın diye |

```
lib/
  main.dart                 # Supabase init (async) + LocalDb init
  providers/                # AuthProvider, ThemeProvider, LanguageProvider
  screens/                  # login, home, admin, profile, appointment, …
  services/
    db_service.dart         # Supabase sorguları
    auth_service.dart       # Supabase Auth
    demo_data.dart          # Yerel mock içerik
    local_db_service.dart   # GEÇİCİ SQLite / prefs
  widgets/                  # Ortak navbar (Icons.fitness_center marka ikonu)
  constants/colors.dart     # SiriusColors / Orion tema + dialog stili
```

---

## Supabase (asıl bağlantı)

1. Supabase projesi oluşturun; `isletme`, randevu, menü/hizmet, galeri vb. şemayı bağlayın.
2. `lib/main.dart` içindeki URL / anon key veya compile-time env:

```bash
flutter run --dart-define=SUPABASE_URL=https://xxx.supabase.co \
            --dart-define=SUPABASE_ANON_KEY=eyJ...
```

3. Auth: e-posta/şifre ve (opsiyonel) Google Sign-In.
4. RLS politikalarını işletme (`isletme_id`) bazlı tutun.

Uygulama şu an demo modda Supabase’i **bloklamadan** başlatır; bağlantı yoksa `DemoData` ile ekranlar dolu kalır. Bu, geliştirme kolaylığı içindir — **hedef mimari yine Supabase’tir.**

---

## Geçici SQLite / yerel depolama

`LocalDbService` (`lib/services/local_db_service.dart`):

- **Mobil:** `sqflite` ile `orion_local_temp.db` (kv + remembered_user)
- **Web:** Tarayıcıda native SQLite olmadığı için **SharedPreferences** fallback
- Kullanım: “Beni hatırla”, isteğe bağlı küçük yerel cache

**Supabase’in yerine geçmez.** Canlıya alırken yerel katmanı sınırlı tutun veya kaldırın.

---

## Kurulum & çalıştırma

```bash
cd orion_flutter_app
flutter pub get

# Web (geliştirme)
flutter run -d web-server --web-hostname=0.0.0.0 --web-port=8080

# Release web
flutter run -d web-server --web-hostname=0.0.0.0 --web-port=8080 --release

# Android / iOS
flutter run
```

Demo giriş:

| Rol | Nasıl |
|-----|--------|
| Müşteri | Rol: Müşteri → e-posta (opsiyonel) → şifre boş → Giriş Yap |
| Admin | Rol: Admin → Giriş Yap (şifresiz) → `/admin` |

“Beni hatırla” işaretliyse e-posta ve rol bir sonraki açılışta gelir.

---

## Marka ikonu

Navbar (`CommonAppBar`) ve giriş ekranındaki dönen yuvarlak **aynı** sembolü kullanır: `Icons.fitness_center` + tema mavisi.  
Eski akışta girişte işletme `logo_url` / placeholder görseli vardı; navbar ise fitness ikonu kullanıyordu — bu yüzden farklı görünüyordu. Artık marka ikonu hizalandı.

---

## UI notları

- Dialog / AlertDialog: `ThemeData.dialogTheme` + birçok yerde `SiriusColors.accent` border, hafif saydam arka plan
- “Neden bizi seçmelisiniz” kartları: her zaman mavi border
- “Hizmetlerimiz” kartları: gölge + hover/press ile yükselme

---

## Çeviri

`LanguageProvider` TR/EN sözlükleri. Gym metinleri (why-us, hizmetler) demo + i18n ile uyumlu tutulur. Yeni anahtar eklerken **hem `en` hem `tr`** map’ine yazın.

---

## Bilinen geliştirme notları

- Cloud agent ortamında commit’ler şu an `devopsProjeSon` remote’una gidebilir; `GNP_App` reposuna push için ayrı yetki gerekir.
- Supabase erişilemezse bazı admin işlemleri snackbar gösterebilir; demo ID (`demo-isletme`) ile devam edilir.

---

## Lisans / paket adı

`pubspec.yaml` paket adı tarihsel olarak `hairsalon_flutter` kalabilir; ürün adı **Orion Gym & Fitness**.
