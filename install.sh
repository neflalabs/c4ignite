#!/usr/bin/env bash
# ==============================================================================
#  c4ignite — Universal GitHub One-Line Installer
#  https://github.com/neflalabs/c4ignite
#
#  Usage:
#    curl -fsSL https://raw.githubusercontent.com/neflalabs/c4ignite/main/install.sh | bash
# ==============================================================================

set -euo pipefail

REPO="neflalabs/c4ignite"
BINARY_NAME="c4ignite"
INSTALL_DIR="${C4IGNITE_INSTALL_DIR:-/usr/local/bin}"

# ANSI Colors
BOLD="\033[1m"
GREEN="\033[0;32m"
BLUE="\033[0;34m"
YELLOW="\033[1;33m"
RED="\033[0;31m"
NC="\033[0m"

log_info() {
    printf "${BLUE}==>${NC} ${BOLD}%s${NC}\n" "$1"
}

log_success() {
    printf "${GREEN}==>${NC} ${BOLD}%s${NC}\n" "$1"
}

log_warn() {
    printf "${YELLOW}==> WARNING:${NC} %s\n" "$1"
}

log_error() {
    printf "${RED}==> ERROR:${NC} %s\n" "$1" >&2
}

# 1. Detect Operating System and Architecture
detect_platform() {
    local os arch
    os="$(uname -s | tr '[:upper:]' '[:lower:]')"
    arch="$(uname -m)"

    case "$os" in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        *)
            log_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac

    case "$arch" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            log_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac

    TARGET="${BINARY_NAME}-${OS}-${ARCH}"
}

# 2. Get latest release tag from GitHub API
get_latest_release() {
    local tag
    tag=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || true)
    if [ -z "$tag" ]; then
        # Fallback to main branch build or default tag
        tag="v2026.08.04"
    fi
    VERSION="$tag"
}

# 3. Uninstall Flow
uninstall_c4ignite() {
    printf "${BOLD}==============================================${NC}\n"
    printf "${BOLD}  c4ignite Uninstaller                       ${NC}\n"
    printf "${BOLD}==============================================${NC}\n\n"

    local found=false
    local targets=("/usr/local/bin/${BINARY_NAME}" "${HOME}/.local/bin/${BINARY_NAME}")

    for target in "${targets[@]}"; do
        if [ -f "${target}" ]; then
            found=true
            log_info "Removing binary at ${target}..."
            if [ -w "$(dirname "${target}")" ] || [ "$EUID" -eq 0 ]; then
                rm -f "${target}"
            else
                if command -v sudo >/dev/null 2>&1; then
                    sudo rm -f "${target}"
                else
                    log_error "Permission denied to delete ${target}. Please run with sudo."
                    exit 1
                fi
            fi
            log_success "Removed ${target}"
        fi
    done

    # 2. Clean up shell autocompletions
    log_info "Cleaning up shell autocompletion snippets..."
    if [ -f "${HOME}/.bashrc" ]; then
        sed -i '/c4ignite completion/d' "${HOME}/.bashrc" 2>/dev/null || true
    fi
    if [ -f "${HOME}/.zshrc" ]; then
        sed -i '/c4ignite completion/d' "${HOME}/.zshrc" 2>/dev/null || true
    fi
    rm -f "${HOME}/.config/fish/completions/c4ignite.fish" 2>/dev/null || true

    # 3. Clean up cache directory
    log_info "Cleaning up cache directory..."
    rm -rf "${HOME}/.cache/c4ignite" "${HOME}/.config/c4ignite" 2>/dev/null || true

    if [ "$found" = false ]; then
        log_warn "c4ignite binary not found in standard locations (/usr/local/bin, ~/.local/bin)."
    else
        log_success "c4ignite, its autocompletions, and cache have been completely uninstalled!"
    fi
}

# 4. Main Flow
main() {
    if [ "${1:-}" = "--uninstall" ] || [ "${1:-}" = "uninstall" ]; then
        uninstall_c4ignite
        exit 0
    fi

    printf "${BOLD}==============================================${NC}\n"
    printf "${BOLD}  c4ignite Installer — CodeIgniter 4 CLI     ${NC}\n"
    printf "${BOLD}==============================================${NC}\n\n"

    detect_platform
    log_info "Detected platform: ${OS}/${ARCH}"

    get_latest_release
    log_info "Resolved target version: ${VERSION}"

    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${TARGET}"
    TMP_DIR=$(mktemp -d)
    TMP_BIN="${TMP_DIR}/${BINARY_NAME}"

    trap 'rm -rf "${TMP_DIR}"' EXIT

    log_info "Downloading ${TARGET} from GitHub..."
    if ! curl -fSL "${DOWNLOAD_URL}" -o "${TMP_BIN}" 2>/dev/null; then
        # If release asset not yet available, try direct raw binary fallback
        RAW_URL="https://raw.githubusercontent.com/${REPO}/main/bin/${BINARY_NAME}"
        log_warn "Asset ${TARGET} not found in release. Trying repository binary..."
        if ! curl -fsSL "${RAW_URL}" -o "${TMP_BIN}"; then
            log_error "Failed to download c4ignite binary from ${DOWNLOAD_URL}"
            exit 1
        fi
    fi

    chmod +x "${TMP_BIN}"

    # Determine installation target folder (allow non-root fallback to ~/.local/bin)
    # Use install command or mv instead of cp to avoid 'Text file busy' when updating self
    if [ ! -w "${INSTALL_DIR}" ] && [ "$EUID" -ne 0 ]; then
        if command -v sudo >/dev/null 2>&1; then
            log_info "Installing to ${INSTALL_DIR}/${BINARY_NAME} using sudo..."
            if command -v install >/dev/null 2>&1; then
                sudo install -m 755 "${TMP_BIN}" "${INSTALL_DIR}/${BINARY_NAME}"
            else
                sudo rm -f "${INSTALL_DIR}/${BINARY_NAME}" 2>/dev/null || true
                sudo cp "${TMP_BIN}" "${INSTALL_DIR}/${BINARY_NAME}"
            fi
        else
            INSTALL_DIR="${HOME}/.local/bin"
            mkdir -p "${INSTALL_DIR}"
            log_warn "Cannot write to /usr/local/bin. Installing to ${INSTALL_DIR} instead."
            if command -v install >/dev/null 2>&1; then
                install -m 755 "${TMP_BIN}" "${INSTALL_DIR}/${BINARY_NAME}"
            else
                rm -f "${INSTALL_DIR}/${BINARY_NAME}" 2>/dev/null || true
                cp "${TMP_BIN}" "${INSTALL_DIR}/${BINARY_NAME}"
            fi
        fi
    else
        mkdir -p "${INSTALL_DIR}"
        if command -v install >/dev/null 2>&1; then
            install -m 755 "${TMP_BIN}" "${INSTALL_DIR}/${BINARY_NAME}"
        else
            rm -f "${INSTALL_DIR}/${BINARY_NAME}" 2>/dev/null || true
            cp "${TMP_BIN}" "${INSTALL_DIR}/${BINARY_NAME}"
        fi
    fi

    log_success "c4ignite installed successfully at ${INSTALL_DIR}/${BINARY_NAME}!"

    # 4. Setup Shell Autocompletion
    log_info "Configuring shell auto-completion..."
    "${INSTALL_DIR}/${BINARY_NAME}" completion install || true

    printf "\n"
    printf "${GREEN}To get started, simply run:${NC}\n"
    printf "  ${BOLD}c4ignite init${NC}   (to bootstrap a fresh CodeIgniter 4 project)\n"
    printf "  ${BOLD}c4ignite up${NC}     (to start the container stack)\n"
    printf "  ${BOLD}c4ignite --help${NC} (to view all available commands)\n\n"
}

main "$@"
