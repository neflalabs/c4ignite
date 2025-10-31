<p align="center">
  <img src="docs/img/c4ignite.png" alt="C4ignite" width="400" height="400">
</p>

---

## 🚀 C4ignite — Chill toolkit for CodeIgniter 4 with Laravel Sail energy

C4ignite spins up a lightweight, fully Dockerized CodeIgniter 4 workspace so you can ship faster with zero local-env drama.

---

### ⚙️ Core commands, nice and breezy

- **Bootstrap CI4 AppStarter**  
  Run:  
  `./scripts/c4ignite init`

  Drops the AppStarter skeleton into `src/` (git-ignored on purpose).

  - Cache lands in `backups/cache/`, so re-runs just extract.
  - Add `--force-download` when you want a fresh tarball.
  - Skip `composer install` with `--no-install` if you’re offline or already set.

- **Spin up the full dev stack**  
  Brings up PHP, Nginx, MariaDB, Redis, and Mailhog:

  ```
  ./scripts/c4ignite up
  ```

  Toss in `--build` when you want to rebuild the local image.

- **Handy shortcuts**  
  Your daily drivers:

  ```
  ./scripts/c4ignite lint
  ./scripts/c4ignite backup
  ./scripts/c4ignite build
  ```

  Prime linting once with `./scripts/c4ignite lint --setup`.  
  Want a guided image build? Hit `./scripts/c4ignite build --interactive` to pick tags, push, and build args from a prompt.

- **Shell autocomplete & alias**  
  Install the snippet:

  ```
  ./scripts/c4ignite setup shell
  source ~/.bashrc   # or ~/.zshrc
  ```

  Need a refresh without a new tab? Run `eval "$(./scripts/c4ignite setup shell --refresh)"`.

- **Backup & restore the app skeleton**  
  Archive `src/` (encryption optional):
  ```
  ./scripts/c4ignite backup create
  ```
  Teammates can restore with:
  ```
  ./scripts/c4ignite backup restore <arsip>
  ```
  Or let them pick interactively:
  ```
  ./scripts/c4ignite backup restore --interactive
  ```

---

### 📘 Docs

- Quick guide: [`docs/c4ignite.md`](docs/c4ignite.md)
- Troubleshooting: [`docs/troubleshooting.md`](docs/troubleshooting.md)

---

<p align="center">
  <sub>Built with ☕ & enough serotonin — by devs who vibe with CodeIgniter 4</sub>
</p>
