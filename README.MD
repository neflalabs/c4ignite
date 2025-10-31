<p align="center">
  <img src="docs/img/c4ignite.png" alt="C4ignite" width="400" height="400">
</p>

---

## 🚀 C4ignite — Toolkit santai buat ngerjain CodeIgniter 4 dengan rasa Laravel Sail

C4ignite bantu kamu bikin lingkungan pengembangan CodeIgniter 4 yang ringan, cepat, dan siap pakai — full Dockerized seperti Laravel Sail.

---

### ⚙️ Fitur & Perintah Utama

- **Bootstrap CI4 AppStarter**  
  Jalankan:  
  `./scripts/c4ignite init`

  Buat skeleton CI4 (AppStarter) muncul di folder `src/` — folder ini sengaja di-ignore supaya repo tetap bersih.

  - Tarball AppStarter otomatis di-cache di `backups/cache/`, jadi init berikutnya tinggal ekstrak ulang.
  - Gunakan `--force-download` kalau mau refresh AppStarter.
  - Setelah bootstrap, `composer install` akan otomatis dijalankan (bisa dilewati dengan `--no-install`).

- **Jalankan stack dev lengkap**  
  Hidupkan PHP, Nginx, MariaDB, Redis, dan Mailhog:

  ```
  ./scripts/c4ignite up
  ```

  Tambah `--build` kalau mau rebuild image lokal.

- **Shortcut siap pakai**  
  Beberapa perintah berguna:

  ```
  ./scripts/c4ignite lint
  ./scripts/c4ignite backup
  ./scripts/c4ignite build
  ```

  Jalankan `./scripts/c4ignite lint --setup` sekali biar lint siap dipakai.  
  Mau build dengan panduan? Pakai `./scripts/c4ignite build --interactive` untuk pilih tag, opsi push, dan build-arg secara interaktif.

- **Auto-complete & alias shell**  
  Aktifkan dengan:

  ```
  ./scripts/c4ignite setup shell
  source ~/.bashrc   # atau ~/.zshrc
  ```

  Pengin refresh tanpa buka terminal baru? Jalankan `eval "$(./scripts/c4ignite setup shell --refresh)"`.

- **Backup & restore skeleton**  
  Buat arsip `src/` (bisa terenkripsi):
  ```
  ./scripts/c4ignite backup create
  ```
  Temanmu bisa restore dengan:
  ```
  ./scripts/c4ignite backup restore <arsip>
  ```
  Atau pakai mode interaktif:
  ```
  ./scripts/c4ignite backup restore --interactive
  ```

---

### 📘 Dokumentasi

- Panduan umum: [`docs/c4ignite.md`](docs/c4ignite.md)
- Troubleshooting: [`docs/troubleshooting.md`](docs/troubleshooting.md)
- Workflow & rilis: [`docs/RELEASE.md`](docs/RELEASE.md)

---

<p align="center">
  <sub>Made with ☕ & 🧠 — by developers who love CodeIgniter 4</sub>
</p>
