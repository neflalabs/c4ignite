# Troubleshooting Guide

### 1. Port 8000 or 33060 is already in use
If another service on your host is using port 8000 or 33060:
* Run `c4ignite doctor` to inspect container and port status.
* Stop the host service, or edit the exposed ports in `.c4ignite/docker-compose.yml`.

---

### 2. Missing Framework Dependencies (`Boot.php` error)
If your application shows `Failed opening required Boot.php`:
* Run `c4ignite composer install` inside your project directory to provision all vendor packages.

---

### 3. File Permissions on Linux
`c4ignite` automatically passes your host `HOST_UID` and `HOST_GID` to container processes.
If files were previously created with root permissions:
```bash
sudo chown -R $(id -u):$(id -g) .
```

---

### 4. Database Reset
To completely reset the database and start fresh:
```bash
c4ignite down -v
c4ignite up -d
c4ignite migrate
```

---

### 5. Xdebug Breakpoints
1. Run `c4ignite xdebug on`.
2. Ensure your IDE (VS Code / PhpStorm) is listening on port `9003`.
3. Check status anytime with `c4ignite xdebug status`.
