# c4ignite Handbook 🚀

## Requirements
- Docker & Docker Compose v2.

---

## 🛠️ Quickstart

1. **Install `c4ignite` globally** (or use `./bin/c4ignite`):
   ```bash
   curl -fsSL https://raw.githubusercontent.com/neflalabs/c4ignite/main/install.sh | bash
   ```
2. **Bootstrap CodeIgniter 4 AppStarter**:
   ```bash
   c4ignite init
   ```
   *(Pilih nama folder aplikasi Anda atau tekan Enter untuk default `src/`).*
3. **Start the Development Stack**:
   ```bash
   c4ignite up
   ```
   - Web App: [http://localhost:8000](http://localhost:8000)
   - Mailhog: [http://localhost:8025](http://localhost:8025)

---

## ⚡ CLI Command Reference

- `c4ignite up [--build] [-d]`: Start containerized stack.
- `c4ignite down [-v]`: Stop and remove containers.
- `c4ignite restart [service]`: Restart all or specific service.
- `c4ignite status`: Show live container health.
- `c4ignite logs [-f] [service]`: Tail service logs.
- `c4ignite spark <command>`: Run CodeIgniter 4 Spark command.
- `c4ignite migrate`: Shortcut for `spark migrate`.
- `c4ignite seed [seeder]`: Shortcut for `spark db:seed`.
- `c4ignite db`: Interactive MySQL/MariaDB shell with auto-credentials.
- `c4ignite composer <command>`: Execute Composer in PHP container.
- `c4ignite php <script>`: Execute PHP script in container.
- `c4ignite shell [service]`: Open bash/sh shell inside service container.
- `c4ignite xdebug [on|off|status]`: Toggle Xdebug dynamically.
- `c4ignite test [options]`: Run PHPUnit test suite.
- `c4ignite lint`: Run PHP code style linter.
- `c4ignite backup create/restore`: Fast native encrypted backup & restore.
- `c4ignite build [--tag=...]`: Build production multi-stage OCI container.
- `c4ignite release`: Run zero-downtime release & deployment pipeline.
- `c4ignite doctor`: Run host diagnostic checks.
- `c4ignite completion [bash|zsh]`: Generate shell auto-completions.

---

## 🔌 Default Services & Ports

- **PHP 8.4-FPM**: Internal App Engine
- **Nginx**: `http://localhost:8000`
- **MariaDB 10.11**: `127.0.0.1:33060` (user: `app`, pass: `secret`, db: `app`)
- **Redis 7**: `127.0.0.1:63790`
- **Mailhog UI**: `http://localhost:8025`
