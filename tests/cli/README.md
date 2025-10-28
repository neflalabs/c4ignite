# CLI Tests

Simple Bats tests for the `c4ignite` script. Requires [bats-core](https://github.com/bats-core/bats-core) plus access to the public internet (tests download the AppStarter tarball once and reuse the cache).

Run tests locally:

```bash
bats tests/cli
```

Tips:
- Set `C4IGNITE_ALLOW_HOST_PYTHON=1` kalau mau pakai Python host tanpa kontainer.
- Set `C4IGNITE_SKIP_DOCKER_FIX=1` pas jalan di CI yang nggak punya Docker daemon.
