# Rilis c4ignite

Panduan singkat buat rilis versi baru sekaligus publish image pra-build ke GitHub Container Registry (GHCR).

## 1. Siapkan rilis
1. Pastikan cabang `main` bersih dan semua perubahan sudah melewati tes lokal:
   ```bash
   python -m unittest discover -s tests/python
   bats tests/cli
   ```
2. Update versi / changelog bila perlu (contoh `CHANGELOG.md` atau catatan rilis GitHub).
3. Buat tag semantik, misal `v1.2.0`:
   ```bash
   git tag v1.2.0
   git push origin v1.2.0
   ```

## 2. Workflow build-image.yml
Setiap push tag `v*` akan menjalankan workflow:
1. Membangun image PHP dev stack dari `docker/dev/php/Dockerfile`.
2. Push dua tag ke GHCR (gunakan nama organisasi Anda):
   - `ghcr.io/<owner>/c4ignite-dev:<tag>` (misal `v1.2.0`)
   - `ghcr.io/<owner>/c4ignite-dev:latest`

> **Catatan:** Workflow memakai token GitHub standar dengan izin `packages:write`. Pastikan repository sudah mengizinkan akun/organisasi mengakses GHCR.

## 3. Konsumsi image
- Sambungkan `docker/dev/docker-compose.yml` ke image dengan men-set environment variable:
  ```bash
  export C4IGNITE_PHP_IMAGE=ghcr.io/<owner>/c4ignite-dev:latest
  ```
- Jalankan `docker compose pull php` sebelum `./scripts/c4ignite up` untuk menarik image terbaru.

## 4. Opsional: lampirkan tarball
Jika ingin distribusi offline:
1. Setelah workflow sukses, unduh image:
   ```bash
   docker pull ghcr.io/<owner>/c4ignite-dev:<tag>
   docker save ghcr.io/<owner>/c4ignite-dev:<tag> -o c4ignite-dev-<tag>.tar
   ```
2. Unggah sebagai Release asset di GitHub.

## 5. Housekeeping & produksi
- Hapus tag image lawas bila tidak dibutuhkan via GHCR UI atau `docker` CLI.
- Untuk image produksi berbasis AppStarter, jalankan manual:
  ```bash
  ./scripts/c4ignite init --force-download
  ./scripts/c4ignite build -t ghcr.io/<owner>/c4ignite-app:<tag>
  docker push ghcr.io/<owner>/c4ignite-app:<tag>
  ```
- Refresh dokumen bila ada perubahan besar pada tooling atau pipeline.
