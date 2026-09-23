# DSP sözleşmesi

Rust `vibration_features` ve Go `dsp.VibrationFeatures` aynı adımları uygular. Sayılar `testdata/dsp_case.json` ile kilitlenir.

1. Pencere boşsa, örnek hızı sıfırsa, uzunluk ikinin kuvveti değilse veya örnek sonlu değilse hesap yapılmaz.
2. Ortalama çıkarılır. MPU6050 dururken yaklaşık 1 g okur; bu değer titreşim değildir.
3. RMS, ortalaması alınmış örneklerin karelerinin ortalamasının kareköküdür.
4. Tepe, ortalaması alınmış örneklerin mutlak maksimumudur.
5. Crest factor, RMS sıfıra çok yakınsa 0, değilse tepe bölü RMS'tir.
6. Kurtosis, dördüncü populasyon momentinin ikinci momentin karesine bölümüdür. Fazlalık (excess) değildir. Saf sinüs 1.5, normal dağılım yaklaşık 3'tür. Varyans yoksa kurtosis 0 yazılır; NaN yayınlanmaz.
7. Baskın frekans, ortalaması alınmış pencerenin FFT'sinde doğru akım kutusu hariç en büyük genlikli kutudur. Eşitlikte daha düşük frekans kalır. Çözünürlük `örnek_hızı / pencere` hertzdir.

Telemetri şeması pencereyi 16 ile 2048 arasında ve ikinin kuvveti olmaya zorlar. DSP fonksiyonu, oracle'daki 8 örnekli dalgayı da kabul eder; çerçeve doğrulaması onu yayınlatmaz.

Yerçekimi testi: 1 g ofsetli 8 Hz sinüs, 256 Hz örneklemede RMS olarak 0.7071067811865476 g vermelidir. Ofset çıkarılmazsa RMS 1 g'nin üzerinde kalır.
