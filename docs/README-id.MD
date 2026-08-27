<p align="center">
  <img src="img/c4ignite.png" alt="C4ignite" width="300" height="300">
</p>

# C4ignite — Toolkit CodeIgniter 4 Cepat & Mandiri Berbasis Native Go 🚀

> **Lingkungan Pengembangan CodeIgniter 4 Instan & Full Dockerized** — Tanpa perlu install PHP, Composer, atau MySQL di komputer host.

---

## ⚡ Instalasi Cepat (One-Line Installer dari GitHub)

Pasang `c4ignite` secara global ke sistem Anda dalam hitungan detik:

```bash
curl -fsSL https://raw.githubusercontent.com/neflalabs/c4ignite/main/install.sh | bash
```

---

## 🛠️ Panduan Memulai Cepat

### 1. Inisialisasi Proyek CodeIgniter 4
```bash
c4ignite init
```
*CLI akan menanyakan nama folder aplikasi yang Anda inginkan (misal `src`, `app`, atau nama kustom lainnya).*

### 2. Nyalakan Development Stack
```bash
c4ignite up
```
* Buka browser di: [http://localhost:8000](http://localhost:8000)
* Mailhog (Web UI email testing): [http://localhost:8025](http://localhost:8025)

---

## ⚡ Perintah Utama Sehari-hari

| Aksi | Perintah |
|---|---|
| **Menjalankan Spark CI4** | `c4ignite spark make:controller Home` |
| **Migrasi Database** | `c4ignite migrate` |
| **Seeder Database** | `c4ignite seed UserSeeder` |
| **Interactive MySQL Console** | `c4ignite db` |
| **Menjalankan Composer** | `c4ignite composer require <package>` |
| **Menjalankan Unit Test** | `c4ignite test` |
| **Toggle Xdebug Dinamis** | `c4ignite xdebug on` / `c4ignite xdebug off` |
| **Backup Native Terenkripsi** | `c4ignite backup create --key=secret` |
| **Build Multi-Stage Produksi** | `c4ignite build --tag=myapp:v1.0.0` |
| **Deploy & Release Pipeline** | `c4ignite release` |
| **Cek Kesiapan Host** | `c4ignite doctor` |

---

## 📘 Dokumentasi Tambahan
- [Troubleshooting & Solusi Masalah](troubleshooting-id.md)
- [Histori Arsitektur & TaskDev](../TaskDev/README.md)
