#compdef c4ignite

_c4ignite() {
  local -a commands
  commands=(
    up down restart status shell spark composer php logs
    init doctor test lint audit build xdebug fresh migrate tinker backup setup
  )

  if (( CURRENT == 2 )); then
    _describe 'command' commands
    return
  fi

  case "$words[2]" in
    backup)
      local -a backup_sub=(create list restore info)
      if (( CURRENT == 3 )); then
        _describe 'backup command' backup_sub
        return
      fi
      case "$words[3]" in
        create)
          _arguments '*:options:(--output --include-vendor --include-writable --exclude-env --encrypt --yes -y)'
          return
          ;;
        restore)
          _arguments '*:options:(--file --auto-backup --interactive --yes -y)'
          return
          ;;
      esac
      ;;
    init)
      _arguments '*:options:(--version --force-download --force --no-install --help -h)'
      return
      ;;
    fresh)
      _arguments '*:options:(--reinit --help -h)'
      return
      ;;
    shell)
      if (( CURRENT == 3 )); then
        _describe 'service' (php mysql nginx python)
        return
      fi
      ;;
    setup)
      local -a setup_sub=(shell env)
      if (( CURRENT == 3 )); then
        _describe 'setup command' setup_sub
        return
      fi
      case "$words[3]" in
        shell)
          _arguments '*:options:(--shell --alias --uninstall --yes -y)'
          return
          ;;
        env)
          if (( CURRENT == 4 )); then
            _describe 'env command' (list copy)
            return
          fi
          if [[ "$words[4]" == copy && CURRENT == 5 ]]; then
            _describe 'template' (dev staging prod docker)
            return
          fi
          ;;
      esac
      ;;
  esac
}

compdef _c4ignite c4ignite ./scripts/c4ignite
