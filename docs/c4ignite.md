# c4ignite Handbook

## Prerequisites

- Docker plus Docker Compose v2.
- A shell with `curl` and `tar` (Python tooling ships via the `python` container).
- `rsync` is optional but speeds up the first init.

## Bootstrapping your project

1. Make sure `src/app/Config/App.php` does not exist (the repo intentionally keeps `src/` empty).
2. Run `./scripts/c4ignite init` to pull the latest AppStarter into `src/` (still ignored by git).
   - The downloaded tarball is cached in `backups/cache/`, so subsequent runs just extract.
   - Add `--force-download` if you want to refresh the tarball.
3. Tweak `src/.env` or `src/.env.docker` to taste. If they’re missing, the init step generates them.
4. Bring the stack online with `./scripts/c4ignite up` (add `--build` to rebuild the local image). The init already runs `composer install` unless you pass `--no-install`.
   - Smash the GitHub API rate limit? Set `GH_TOKEN` (read-only PAT) before running `c4ignite init`.
5. Optional: `./scripts/c4ignite setup env copy dev` (or `staging` / `prod` / `docker`) to pull in a template `.env`.
6. Optional: `./scripts/c4ignite setup shell` to drop aliases and autocompletion straight into your shell.

## CLI cheat sheet

- `up | down | restart | status`: manage Docker services.
- `shell [service]`: hop into a container shell (defaults to `php`).
- `php …`: run PHP CLI inside the container (cwd is `/var/www/html`).
- `spark …`: proxy `php spark`.
- `composer …`: run Composer within the app container.
- `logs [service]`: tail logs for all services or a specific one.
- `lint`: run linting; first bootstrap via `./scripts/c4ignite lint --setup`.
- `build [options]`: build the production image (`-t/--tag`, `--push`, `--build-arg`, `--no-cache`, `--interactive`). Use `--interactive` for a guided prompt (tag, push, build args).
- `fresh [--reinit]`: stop the stack, wipe volumes, and optionally reset `src/`.
- `backup create/list/restore/info`: manage encrypted/unencrypted `src/` backups.
- `setup env list/copy`: list or copy `.env` templates (dev/staging/prod/docker).
- `init [tag] [--force-download]`: grab AppStarter (latest by default).
- `doctor`: check local prerequisites and port collisions.
- `xdebug [on|off|status]`: toggle Xdebug without rebuilding.
- `setup shell [options]`: interactive wizard for aliases + completions (Bash/Zsh). Supports `--alias`, `--yes`, and `--uninstall`.

### Default shell per service

- **php**: `/bin/bash`
- **mysql**: `/bin/sh`
- **nginx**: `/bin/sh`
- **python**: `/bin/sh`

## Default services & ports

- Web app: http://localhost:8000 (Nginx serves `public/`).
- MySQL: `127.0.0.1:33060` with credentials `app/secret`.
- Redis: `127.0.0.1:63790`.
- Mailhog UI: http://localhost:8025.

## Quick tips

- All project code lives in `src/`; containers mount the same directory.
- Export `HOST_UID` / `HOST_GID` before `c4ignite up` to keep permissions chill:
  `export HOST_UID=$(id -u)` and `export HOST_GID=$(id -g)`.
- Need Xdebug? Configure `PHP_IDE_CONFIG` + host/port in `src/.env`, override with `XDEBUG_CLIENT_HOST/PORT` if needed, then run `./scripts/c4ignite xdebug on`.
- Pulling a big change? Drop the stack with `./scripts/c4ignite down` so volumes don’t get cranky.
- VS Code users can lean on `.devcontainer/` and `.vscode/`.
- Looking at CI/CD? Example workflows live in `.github/workflows/`.
- Trouble? Hit `docs/troubleshooting.md`.
- `./scripts/c4ignite init` defaults `CI_ENVIRONMENT` to `development` and generates `.env.docker` if it’s missing.
- Run `./scripts/c4ignite setup shell` to install/uninstall aliases + completions; source your shell rc afterward, or shortcut with `eval "$(./scripts/c4ignite setup shell --refresh)"`.

## Backup & share the skeleton

- `./scripts/c4ignite backup create [options]` archives `src/` (vendor/writable excluded by default). Add `--encrypt` to use a passphrase.
- See what’s available via `./scripts/c4ignite backup list`. Inspect details with `./scripts/c4ignite backup info <file>`.
- Restore from a teammate using `./scripts/c4ignite backup restore --file path/to/archive` (auto-backup prompt included). Add `--interactive` to pick from the list in `backups/`.
- Have a pre-built dev image? `docker compose -f docker/dev/docker-compose.yml build php` then `docker save` before sharing.

## Production build flow

- `./scripts/c4ignite build -t registry.example.com/ci4-app:latest` builds the production image using `docker/prod/Dockerfile`.
- Want to push right away? Slip in `--push` (and make sure you’re logged in).
- Need build overrides? Use `--build-arg KEY=VALUE`, e.g. `--build-arg APP_BASE_URL=https://app.example.com`.
