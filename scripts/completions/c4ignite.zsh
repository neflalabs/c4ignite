#compdef c4ignite

_c4ignite() {
  local -a commands
  commands=(
    up down restart status pull shell spark composer php logs
    init doctor test lint audit build xdebug fresh migrate backup setup
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
        local -a shell_services=(php mysql nginx python)
        _describe 'service' shell_services
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
            local -a setup_env_sub=(list copy)
            _describe 'env command' setup_env_sub
            return
          fi
          if [[ "$words[4]" == copy && CURRENT == 5 ]]; then
            local -a env_templates=(dev staging prod docker)
            _describe 'template' env_templates
            return
          fi
          ;;
      esac
      ;;
  esac
}

compdef _c4ignite c4ignite ./scripts/c4ignite
