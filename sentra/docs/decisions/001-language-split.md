# 001 — Rust kenarda, Go hatta, Python öğrenmede

Tek dil daha az görünür sürtünme yaratır. Bu sistemde sürtünme dil sayısından değil, işin şeklinden gelir.

ESP32 üzerinde örnekleme deterministik olmalı. Çöp toplayıcı duraklaması bir titreşim penceresini kaydırır. Rust, `sentra-core` sayesinde aynı RMS ve FFT kodunu dizüstünde test edip kartta çalıştırır. TinyGo ESP32'de derlenir ama ESP-IDF sürücüleri ve kesme modeli Rust ekosisteminde daha olgundur. Bu yüzden aygıt yazılımı Go değil.

Ingest tarafında kazanç, dilin mikro kıyaslamasından çok eşzamanlı MQTT, HTTP ve veritabanı havuzunun sade yazılmasıdır. Go bunu tek çalışma zamanında verir. Aynı hizmeti Rust ile yazmak mümkündür; V1'de kazanç, axum ve sqlx öğrenmek değil, sözleşmeyi ve panoyu bitirmektir.

xgboost ve SHAP Python'da kalır. Onları Rust veya Go'ya taşımak modelin açıklanabilirliğini düşürür. V2 bu yüzden yeni bir dil yarışı değil, formülün üçüncü kez — bu sefer NumPy ile — kilitlenmesidir.

Üç süreç birbirinin kütüphanesini çağırmaz. Aralarındaki tek bağ JSON sözleşmesidir.
