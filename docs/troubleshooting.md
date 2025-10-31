# Troubleshooting santai

## Masalah yang sering muncul

### Permission denied di `src/`
Kalau `./scripts/c4ignite init` komplain soal izin, artinya folder `src/` lagi dimiliki root. Skrip sekarang sudah coba betulin otomatis pake container Alpine, tapi kalau masih ngeyel tinggal jalanin `sudo chown -R $(id -u):$(id -g) src` atau sekalian `./scripts/c4ignite fresh --reinit` buat reset total.
Lari di environment yang nggak punya Docker? Set `C4IGNITE_SKIP_DOCKER_FIX=1` sebelum jalanin perintahnya biar langkah perbaikan otomatis dilewati.

### Port udah dipakai service lain
Gas `./scripts/c4ignite doctor` buat cek port mana yang bentrok, terus matiin service yang ngeribetin.

### Build pertama lama banget
Wajar banget karena lagi download base image sama dependency. Bisa akalin dengan sharing cache/image ke temen tim atau sediain image pre-built di registry internal.
Kalau udah punya image sendiri, tinggal push ke registry tim, `docker compose pull php`, terus jalankan `./scripts/c4ignite up` supaya stack langsung pakai image yang udah siap pakai.

### Composer cache rewel
Pastikan `COMPOSER_HOME` bisa ditulis. Stack dev udah sediain volume cache, tapi kalau masih berulah, coba `./scripts/c4ignite composer clear-cache`.

### Xdebug nggak nendang
Pastikan udah `./scripts/c4ignite xdebug on`, IDE listen di port 9003, dan setting host di `src/.env` cocok (apalagi kalau kamu pakai WSL/VM).

## FAQ singkat

- **Bisa re-bootstrap framework?** Bisa banget. Kosongin aja `src/` (atau hapus sekalian) terus jalankan `./scripts/c4ignite init`. Folder ini memang di-ignore git.
- **Gimana build image produksi?** Tinggal `./scripts/c4ignite build -t your/image:tag`.
- **Alias & auto-complete gimana?** Jalankan `./scripts/c4ignite setup shell` (ada opsi `--uninstall` kalau mau bersih-bersih).
- **Harus pakai Docker Desktop?** Nggak wajib. Semua Docker Engine 20.x+ oke, yang penting Compose v2 tersedia.
