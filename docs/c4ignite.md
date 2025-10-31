# c4ignite Handbook

## Syarat dasar
- Punya Docker + Docker Compose v2.
- Shell yang ada `curl` dan `tar` (Python disiapkan lewat container `python`).
- Bonus kalau ada `rsync`, bikin proses init makin cepat.

## Cara bootstrap project
1. Pastikan folder `src/app/Config/App.php` belum ada (repo ini emang sengaja kosongin `src/`).
2. Jalanin `./scripts/c4ignite init` biar AppStarter terbaru diambil ke folder `src/` (folder ini nggak ikut ke git).
   - Tarball yang sudah kebawa bakal di-cache di `backups/cache/`, jadi init berikutnya cukup ekstrak ulang.
   - Mau paksa download ulang? Tambahin `--force-download`.
3. Edit `src/.env` atau `src/.env.docker` sesuai kebutuhan. Kalau belum ada, init bakal bikinin otomatis.
4. Nyalain stack via `./scripts/c4ignite up` (tambah `--build` kalau mau rebuild image lokal); `./scripts/c4ignite init` sudah otomatis menjalankan `composer install` (lewati dengan `--no-install`).
   - Kalau kena limit API GitHub, set `GH_TOKEN` (PAT read-only) sebelum panggil `c4ignite init`.
   - Punya image pra-build? Set `C4IGNITE_PHP_IMAGE=ghcr.io/neflalabs/c4ignite-dev:latest` terus `docker compose pull php` supaya container pakai image siap pakai.
5. Opsional: jalankan `./scripts/c4ignite setup env copy dev` (atau `staging`/`prod`/`docker`) buat nurunin template `.env`.
6. Opsional: jalankan `./scripts/c4ignite setup shell` biar alias dan auto-complete langsung aktif.

## Cheat sheet perintah CLI
- `up | down | restart | status`: ngatur hidup matinya service Docker.
- `shell [service]`: masuk ke shell container tertentu (default `php`).
- `php …`: jalankan command PHP CLI di container (otomatis di `/var/www/html`), misal `./scripts/c4ignite php -v`.
- `spark …`: akses `php spark` di dalam container PHP.
- `composer …`: larikan Composer di container PHP (langsung di `/var/www/html`).
- `logs [service]`: intip log semua service atau pilih salah satu.
- `lint`: jalanin linting; pakai `./scripts/c4ignite lint --setup` sekali untuk tambahin script + PHPCS.
- `build [opsi]`: build image produksi (`-t/--tag`, `--push`, `--build-arg`, `--no-cache` tersedia).
- `fresh [--reinit]`: matiin stack + bersihin volume; kalau tambah `--reinit` sekalian reset `src/`.
- `tinker …`: buka `php spark shell` buat uji cepat.
- `backup create/list/restore/info`: utility backup/restore `src/` (support enkripsi, autosave).
- `setup env list/copy`: manage template `.env` (development/staging/production/docker).
- `init [tag] [--force-download]`: download AppStarter (default ke versi terbaru; `--force-download` buat abaikan cache).
- `doctor`: cek kesiapan environment lokal dan status port.
- `xdebug [on|off|status]`: toggle Xdebug tanpa rebuild container.
- `setup shell [opsi]`: wizard interaktif buat pasang alias + auto-complete (Bash/Zsh), dukung `--alias`, `--yes`, dan `--uninstall`.

### Default shell tiap service
- **php**: `/bin/bash` (biar nyaman buat develop di container app)
- **mysql**: `/bin/sh` (bawaan image MariaDB)
- **nginx**: `/bin/sh` (bawaan image Nginx resmi)
- **python**: `/bin/sh` (image Python Alpine)

## Service & port default
- Web app: http://localhost:8000 (Nginx serve `public/`).
- MySQL: `127.0.0.1:33060` user/pass `app/secret`.
- Redis: `127.0.0.1:63790`.
- Mailhog UI: http://localhost:8025.

## Tips ringan
- Semua source CodeIgniter kita taruh di `src/`; volume container juga ke sana.
- Ekspor `HOST_UID`/`HOST_GID` sebelum `c4ignite up` biar permission aman:  
  `export HOST_UID=$(id -u)` dan `export HOST_GID=$(id -g)`.
- Mau debug pakai Xdebug? Set `PHP_IDE_CONFIG` plus host/port di `src/.env`, bisa override pakai env `XDEBUG_CLIENT_HOST/PORT`, lalu toggle `./scripts/c4ignite xdebug on`.
- Sebelum pull perubahan gede, mending `./scripts/c4ignite down` dulu supaya nggak ke-lock volume.
- Pengguna VSCode bisa langsung pakai setup di `.devcontainer/` dan `.vscode/`.
- Pake CI/CD? Workflow contoh ada di `.github/workflows/` dan panduan rilis di `docs/RELEASE.md`.
- Kalau ada masalah, langsung cek `docs/troubleshooting.md`.
- `./scripts/c4ignite init` otomatis ngaktipin `CI_ENVIRONMENT = development` dan bikin `.env.docker` default kalau belum ada.
- Mau alias & auto-complete? Jalankan wizard `./scripts/c4ignite setup shell` (bisa install/uninstall buat Bash/Zsh), habis itu `source ~/.bashrc` atau `source ~/.zshrc` supaya langsung aktif.

## Backup & bagi-bagi skeleton
- Gunakan `./scripts/c4ignite backup create [opsi]` buat bikin arsip `src/` (default exclude vendor/writable). Tambahkan `--encrypt` kalau mau pakai passphrase.
- Mau lihat daftar backup? `./scripts/c4ignite backup list`. Info detail: `./scripts/c4ignite backup info <file>`.
- Restore dari dev lain: `./scripts/c4ignite backup restore --file path/to/arsip` (bakal minta konfirmasi, bisa auto-backup kondisi lama). Tambahkan `--interactive` buat pilih arsip langsung dari daftar di folder `backups/`.
- Perlu image dev pra-build? `docker compose -f docker/dev/docker-compose.yml build php` lalu simpan dengan `docker save` sebelum dibagikan.

## Build buat produksi
- `./scripts/c4ignite build -t registry.example.com/ci4-app:latest` bakal bikin image produksi berdasarkan `docker/prod/Dockerfile`.
- Mau langsung kirim ke registry? Tambahin `--push` (jangan lupa `docker login` dulu).
- Perlu override ARG/ENV? Pakai `--build-arg KEY=VALUE`, contoh `--build-arg APP_BASE_URL=https://app.example.com`.
