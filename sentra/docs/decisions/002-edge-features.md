# 002 — Ham pencere yerine kenar özeti

İlk prompt, ESP32'nin ham sensör örneklerini MQTT ile basmasını ve RMS ile FFT'nin Python'da hesaplanmasını istiyordu. Bu, laboratuvar defterinde kolaydır. 256 örnek, 1 kHz, sürekli yayın: her saniye yüzlerce sayı, yerçekimi ofseti henüz ayıklanmamış.

SENTRA pencereyi cihazda kapatır. Yayınlanan çerçeve onlarca bayttır. Model daha sonra bu özeti görür. Ham dalga biçimi gerekirse V2'de, yalnızca şüpheli pencerede ve ayrı bir konuda istenir. V1'de o konu yoktur; yarım bir kanca da bırakılmadı.

Bedel: ingest, RMS'i yeniden hesaplayamaz çünkü örnek ondadır. Bu yüzden Rust ve Go hesapları `testdata/dsp_case.json` üzerinde birebir kilitlendi. Simülatör Go hesabını kullanır. Kart, Rust hesabını kullanacak. İkisi saparsa sınama kırılır.
