#!/usr/bin/env bats

APPSTARTER_TEST_TAG="v4.6.3"

setup() {
  export SCRIPT_DIR="$(pwd)/scripts"
  export C4IGNITE="${SCRIPT_DIR}/c4ignite"
  export C4IGNITE_ALLOW_HOST_PYTHON=1
  export C4IGNITE_SKIP_DOCKER_FIX=1
}

@test "c4ignite help shows commands" {
  run "${C4IGNITE}" --help
  [ "$status" -eq 0 ]
  [[ "$output" =~ "Usage: c4ignite" ]]
}

@test "c4ignite doctor detects docker command" {
  run "${C4IGNITE}" doctor
  [ "$status" -eq 0 ]
  [[ "$output" =~ "docker" ]]
}

make_workspace() {
  local dir
  dir="$(mktemp -d)"
  cp -a docker "${dir}/"
  cp -a scripts "${dir}/"
  cp -a templates "${dir}/"
  mkdir -p "${dir}/backups"
  echo "${dir}"
}

@test "c4ignite init caches downloaded tarball" {
  workspace="$(make_workspace)"
  pushd "${workspace}" >/dev/null

  run env C4IGNITE_ALLOW_HOST_PYTHON=1 C4IGNITE_SKIP_DOCKER_FIX=1 ./scripts/c4ignite init --force-download --version "${APPSTARTER_TEST_TAG}"
  [ "$status" -eq 0 ]
  [ -f "backups/cache/appstarter-${APPSTARTER_TEST_TAG}.tar.gz" ]

  rm -rf src

  run env C4IGNITE_ALLOW_HOST_PYTHON=1 C4IGNITE_SKIP_DOCKER_FIX=1 ./scripts/c4ignite init --version "${APPSTARTER_TEST_TAG}"
  [ "$status" -eq 0 ]
  [[ "$output" =~ "Using cached CodeIgniter appstarter ${APPSTARTER_TEST_TAG}" ]]

  popd >/dev/null
  rm -rf "${workspace}"
}

@test "c4ignite backup create produces metadata" {
  workspace="$(make_workspace)"
  pushd "${workspace}" >/dev/null

  run env C4IGNITE_ALLOW_HOST_PYTHON=1 C4IGNITE_SKIP_DOCKER_FIX=1 ./scripts/c4ignite init --force-download --version "${APPSTARTER_TEST_TAG}"
  [ "$status" -eq 0 ]

  run env C4IGNITE_ALLOW_HOST_PYTHON=1 C4IGNITE_SKIP_DOCKER_FIX=1 ./scripts/c4ignite backup create -y
  [ "$status" -eq 0 ]

  archive="$(ls backups/c4ignite-src-*.tar.gz 2>/dev/null | tail -n 1)"
  [ -f "${archive}" ]

  tar -tzf "${archive}" c4ignite-backup.json
  metadata="$(tar -xOf "${archive}" c4ignite-backup.json)"
  [[ "${metadata}" =~ "includes_vendor" ]]

  popd >/dev/null
  rm -rf "${workspace}"
}
