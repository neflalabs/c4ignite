<p align="center">
  <img src="img/c4ignite.png" alt="c4ignite" width="260" height="260">
</p>

<h1 align="center">c4ignite</h1>
<p align="center">
  <strong>Environment development lokal berbasis Docker untuk CodeIgniter 4.</strong><br>
  Tanpa perlu install PHP, Composer, atau MySQL di komputer host. Berbentuk single binary mandiri.
</p>

---

## ⚡ Instalasi

Pasang `c4ignite` secara global ke sistem Anda:

```bash
curl -fsSL https://raw.githubusercontent.com/neflalabs/c4ignite/main/install.sh | bash
```

Untuk uninstall:
```bash
curl -fsSL https://raw.githubusercontent.com/neflalabs/c4ignite/main/install.sh | bash -s -- --uninstall
```

---

## 🚀 3 Langkah Mulai

### 1. Inisialisasi proyek baru
```bash
c4ignite init Aplikasiku
```
*(Atau ketik `c4ignite init` untuk menggunakan nama folder default `Codeigniter4`)*

### 2. Masuk ke folder proyek
```bash
cd Aplikasiku
```

### 3. Nyalakan container stack
```bash
c4ignite up
```
Buka browser di **[http://localhost:8000](http://localhost:8000)**.

---

## 🛠️ Perintah Harian

| Perintah | Deskripsi |
|---|---|
| `c4ignite up` | Menyalakan stack development (PHP, Nginx, MariaDB) |
| `c4ignite down` | Menghentikan container stack (`-v` untuk reset database volume) |
| `c4ignite restart [service]` | Restart container service |
| `c4ignite status` | Melihat status dan health container |
| `c4ignite logs -f` | Melihat log real-time container |
| `c4ignite spark [command]` | Menjalankan perintah CodeIgniter 4 Spark |
| `c4ignite migrate` | Menjalankan migrasi database (`spark migrate`) |
| `c4ignite seed [name]` | Menjalankan seeder database (`spark db:seed`) |
| `c4ignite composer [args]` | Menjalankan Composer di dalam container |
| `c4ignite php [script]` | Menjalankan script PHP di dalam container |
| `c4ignite db` | Membuka interactive console MySQL / MariaDB |
| `c4ignite shell [service]` | Masuk ke terminal shell container |
| `c4ignite xdebug on\|off` | Mengaktifkan/menonaktifkan Xdebug |
| `c4ignite test` | Menjalankan unit test PHPUnit |
| `c4ignite lint` | Menjalankan PHP code style linter |
| `c4ignite doctor` | Memeriksa kesiapan Docker & environment host |

---

## ⚙️ Informasi Port Bawaan

* **Aplikasi Web:** `http://localhost:8000`
* **MySQL / MariaDB:** `localhost:33060` (user: `app`, pass: `secret`, db: `app`)
