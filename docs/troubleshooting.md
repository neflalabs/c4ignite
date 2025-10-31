# Troubleshooting (keep it calm)

## Common hiccups

### Permission denied under `src/`
If `./scripts/c4ignite init` gripes about permissions, the `src/` directory is probably owned by root. The script attempts to patch it using a tiny Alpine container, but if it still acts up, run `sudo chown -R $(id -u):$(id -g) src` or nuke it with `./scripts/c4ignite fresh --reinit`.
Running somewhere without Docker? Export `C4IGNITE_SKIP_DOCKER_FIX=1` before the command so the auto-fix step is skipped.

### Ports already claimed
Run `./scripts/c4ignite doctor` to see what’s hogging your ports, then shut down the culprit.

### First build takes forever
Totally normal—base images and dependencies are hefty. Speed it up by sharing caches/images with the squad or hosting a pre-built variant in your registry.
If you already have an image, push it to the team registry, run `docker compose pull php`, then launch `./scripts/c4ignite up` so the stack starts with the pre-baked image.

### Composer cache throwing shade
Make sure `COMPOSER_HOME` is writable. The dev stack mounts a cache volume; if it still misbehaves, try `./scripts/c4ignite composer clear-cache`.

### Xdebug not hitting breakpoints
Toggle it on with `./scripts/c4ignite xdebug on`, ensure your IDE listens on port 9003, and verify host settings inside `src/.env` (especially for WSL/VM setups).

## Quick FAQ

- **Can I re-bootstrap the framework?** Yup. Empty out `src/` (or delete it) and run `./scripts/c4ignite init`. The directory is git-ignored by design.
- **How do I build a production image?** Run `./scripts/c4ignite build -t your/image:tag`.
- **Where do alias & autocomplete live?** Use `./scripts/c4ignite setup shell` (with `--uninstall` for cleanup).
- **Do I need Docker Desktop?** Nope. Any Docker Engine 20.x+ works—as long as Compose v2 is available.
