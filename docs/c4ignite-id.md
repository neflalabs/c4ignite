# Panduan Lengkap c4ignite 🚀

## Syarat Utama
- Docker & Docker Compose v2.

---

## 🛠️ Langkah Cepat

1. **Install `c4ignite` secara global** (atau gunakan `./bin/c4ignite`):
   ```bash
   curl -fsSL https://raw.githubusercontent.com/neflalabs/c4ignite/main/install.sh | bash
   ```
2. **Inisialisasi CodeIgniter 4 AppStarter**:
   ```bash
   c4ignite init
   ```
   *(Pilih nama folder aplikasi Anda atau tekan Enter untuk default `src/`).*
3. **Nyalakan Development Stack**:
   ```bash
   c4ignite up
   ```
   - Web App: [http://localhost:8000](http://localhost:8000)
   - Mailhog: [http://localhost:8025](http://localhost:8025)

---

## ⚡ Referensi Perintah CLI

- `c4ignite up [--build] [-d]`: Menyalakan stack container.
- `c4ignite down [-v]`: Mematikan dan menghapus container.
- `c4ignite restart [service]`: Restart seluruh atau service tertentu.
- `c4ignite status`: Memeriksa kesehatan container aktif.
- `c4ignite logs [-f] [service]`: Memantau live output logs container.
- `c4ignite spark <command>`: Menjalankan command CodeIgniter 4 Spark.
- `c4ignite migrate`: Shortcut untuk `spark migrate`.
- `c4ignite seed [seeder]`: Shortcut untuk `spark db:seed`.
- `c4ignite db`: Membuka console MySQL/MariaDB interaktif.
- `c4ignite composer <command>`: Menjalankan Composer di dalam container PHP.
- `c4ignite php <script>`: Menjalankan script PHP di dalam container.
- `c4ignite shell [service]`: Masuk ke terminal bash/sh container.
- `c4ignite xdebug [on|off|status]`: Toggle Xdebug secara dinamis.
- `c4ignite test [options]`: Menjalankan test suite PHPUnit.
- `c4ignite lint`: Menjalankan code style linter PHP.
- `c4ignite backup create/restore`: Backup & restore terenkripsi cepat.
- `c4ignite build [--tag=...]`: Build image OCI multi-stage untuk server produksi.
- `c4ignite release`: Menjalankan pipeline rilis, migrasi DB, dan healthcheck.
- `c4ignite doctor`: Diagnostik host dan kesiapan sistem.
- `c4ignite completion [bash|zsh]`: Generate auto-complete terminal.

---

## 🔌 Service & Port Bawaan

- **PHP 8.4-FPM**: Engine Aplikasi
- **Nginx**: `http://localhost:8000`
- **MariaDB 10.11**: `127.0.0.1:33060` (user: `app`, pass: `secret`, db: `app`)
- **Redis 7**: `127.0.0.1:63790`
- **Mailhog UI**: `http://localhost:8025`
