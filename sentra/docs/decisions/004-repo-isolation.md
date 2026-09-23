# 004 — Insucom ağacına dokunmamak

Bu git deposu canlı bir sigorta sitesini de taşır (`backend/`, `frontend/`, kök `docker-compose.yml`). SENTRA ayrı bir ürün olduğu için o ağacın içine backend olarak yazılmadı. Kök pipeline'a `sentra/` hariç tutuldu; aksi halde `main` birleşmesi klasörü canlı sunucuya rsync ederdi.

Yeni bir GitHub deposu açmak bu çalışma alanının yazma yetkisi dışında bırakıldı. Kod `sentra/` altındadır ve kendi `docker-compose.yml` dosyasını kullanır.
