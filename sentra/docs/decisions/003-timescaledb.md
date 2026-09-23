# 003 — TimescaleDB

Düz PostgreSQL satırı tutar. Motor çerçevesi zamana göre büyür, son N pencere ve cihaz başına son durum okunur. TimescaleDB, aynı SQL'i hipertablo ile böler. Uzantı yoksa ingest tabloyu yine oluşturur ve uyarı yazar; geliştirme Postgres'i bu yüzden özel bir çatala kilitlenmez.

Birincil anahtar `(device_id, captured_at)` şeklindedir. Hipertablo, bölüm sütununun benzersiz anahtarda olmasını ister. `captured_at` bu yüzden anahtardadır.

Port 5433'tür. Bu depodaki Insucom Postgres'i 5432'yi kullanır.
