# c4ignite Roadmap (versi santai)

## Phase 1 – Skeleton & Stack
- Rapihin struktur repo: folder `docker/dev/`, config Nginx, Dockerfile PHP, template env, sampe skrip `scripts/c4ignite`.
- Tulis `docker-compose.yml` buat PHP-FPM, Nginx, MariaDB, Redis, Mailhog, lengkap sama volume dan network.
- Pasang command dasar CLI: `up`, `down`, `restart`, `status`, `shell`, `spark`, `composer`, `logs`.

## Phase 2 – Onboarding & Init
- Tambah `c4ignite init` buat bootstrap CI4 (download skeleton, salin `.env`, benerin permission).
- Update dokumentasi di `docs/c4ignite.md` plus README.
- Bikin `c4ignite doctor` buat ngecek dependency dev dan konflik port.

## Phase 3 – Developer Experience
- Tambahin perintah QA: `c4ignite test`, `lint`, `audit` (optional coverage).
- Sediakan toggle `c4ignite xdebug on/off` dan alias `c4ignite php/composer` biar eksekusi via container makin enak.
- Kalau perlu, siapin konfigurasi override (profile) dan autocompletion bash/zsh.

## Phase 4 – Production Tooling
- Bangun Dockerfile multi-stage untuk produksi (`docker/prod/`) plus profile compose.
- Implement `c4ignite build` (dan opsi push) buat bikin image siap deploy.
- Tambah helper `c4ignite deploy/release` buat migrasi & health check setelah rilis.

## Phase 5 – Enhancements & Integration
- Kasih DevContainer/IDE helper, panduan debugging, dan sample pipeline CI.
- Eksplore fitur ekstra kayak preset worker/scheduler.
- Tambah test suite buat skrip `c4ignite` dan terus iterasi dari feedback dev.
- Tambah utilitas produktivitas: wizard lint (`lint --setup`), manajemen template `.env`, serta sistem backup/restore terenkripsi + mode interaktif.
