<p align="center">
  <img src="docs/img/c4ignite.png" alt="c4ignite" width="260" height="260">
</p>

<h1 align="center">c4ignite</h1>
<p align="center">
  <strong>Docker-powered local development environment for CodeIgniter 4.</strong><br>
  Zero host PHP, Composer, or MySQL required. Single standalone binary.
</p>

---

## ⚡ Installation

Install `c4ignite` globally in seconds:

```bash
curl -fsSL https://raw.githubusercontent.com/neflalabs/c4ignite/main/install.sh | bash
```

*Supports Linux (amd64/arm64), macOS (Intel & Apple Silicon), and Windows (WSL).*

To uninstall:
```bash
curl -fsSL https://raw.githubusercontent.com/neflalabs/c4ignite/main/install.sh | bash -s -- --uninstall
```

---

## 🚀 3-Step Quickstart

### 1. Initialize a new project
```bash
c4ignite init MyApp
```
*(Or simply run `c4ignite init` to accept the default folder `Codeigniter4`)*

### 2. Enter the project directory
```bash
cd MyApp
```

### 3. Start the stack
```bash
c4ignite up
```
Open **[http://localhost:8000](http://localhost:8000)** in your browser.

---

## 🛠️ Daily Commands

All commands run inside isolated Docker containers automatically with your local user UID/GID permissions:

| Command | Description |
|---|---|
| `c4ignite up` | Start local development stack (PHP, Nginx, MariaDB) |
| `c4ignite down` | Stop container stack (`-v` to purge volumes) |
| `c4ignite restart [service]` | Restart all or specific service container |
| `c4ignite status` | View container status and health |
| `c4ignite logs -f` | Tail container logs |
| `c4ignite spark [command]` | Run CodeIgniter 4 Spark command |
| `c4ignite migrate` | Run database migrations (`spark migrate`) |
| `c4ignite seed [name]` | Seed database (`spark db:seed`) |
| `c4ignite composer [args]` | Run Composer in PHP container |
| `c4ignite php [script]` | Run PHP script in container |
| `c4ignite db` | Open interactive MySQL / MariaDB shell |
| `c4ignite shell [service]` | Drop into bash/sh shell inside container |
| `c4ignite xdebug on\|off` | Toggle Xdebug profiling dynamically |
| `c4ignite test` | Run PHPUnit test suite |
| `c4ignite lint` | Run PHP code linter |
| `c4ignite doctor` | Diagnose Docker & host environment prerequisites |

---

## ⚙️ Ports & Services

* **Application (Web):** `http://localhost:8000`
* **MySQL / MariaDB:** `localhost:33060` (user: `app`, pass: `secret`, db: `app`)

---

## 📄 License
MIT License. Crafted for modern CodeIgniter 4 developers.
