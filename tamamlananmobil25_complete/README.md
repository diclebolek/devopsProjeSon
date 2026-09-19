# tamamlananmobil25 — Orion Flutter App

Eksik Flutter proje dosyaları tamamlandı. Artık çalışan bir Flutter mobil uygulaması olarak açılabilir.

## Neler eklendi

- `pubspec.yaml` / `pubspec.lock`
- `android/` (Android Studio ile açılabilir)
- `web/` (web preview)
- `test/`, `analysis_options.yaml`, `.gitignore`, `.metadata`
- Güncellenmiş `lib/` (Orion markası, şifresiz giriş, bottom navbar, placeholder görseller)

## Çalıştırma

```bash
flutter pub get
flutter run
```

Android Studio: `android/` klasörünü veya proje kökünü açıp cihaz/emülatör seçerek Run.

Web:

```bash
flutter run -d web-server --web-hostname=0.0.0.0 --web-port=8080
```

## Notlar

- Marka: **Orion** (sektör-nötr)
- Şifresiz giriş: email/şifre boş bırakıp Giriş Yap
- Bottom navbar mobil/web viewport’ta görünür
