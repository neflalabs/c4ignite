# Bash completion untuk c4ignite
_c4ignite_complete() {
    local cur prev
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    local commands="up down restart status shell spark composer php logs init doctor test lint audit build xdebug fresh migrate tinker backup setup"

    if [[ ${COMP_CWORD} -eq 1 ]]; then
        COMPREPLY=( $(compgen -W "$commands" -- "$cur") )
        return 0
    fi

    local first="${COMP_WORDS[1]}"
    case "${first}" in
        backup)
            local backup_sub="create list restore info"
            if [[ ${COMP_CWORD} -eq 2 ]]; then
                COMPREPLY=( $(compgen -W "$backup_sub" -- "$cur") )
                return 0
            fi
            case "${COMP_WORDS[2]}" in
                create)
                    local opts="--output --include-vendor --include-writable --exclude-env --encrypt --yes -y"
                    COMPREPLY=( $(compgen -W "$opts" -- "$cur") )
                    return 0
                    ;;
                restore)
                    local opts="--file --auto-backup --interactive --yes -y"
                    COMPREPLY=( $(compgen -W "$opts" -- "$cur") )
                    return 0
                    ;;
            esac
            ;;
        setup)
            local setup_sub="shell env"
            if [[ ${COMP_CWORD} -eq 2 ]]; then
                COMPREPLY=( $(compgen -W "$setup_sub" -- "$cur") )
                return 0
            fi
            case "${COMP_WORDS[2]}" in
                shell)
                    local opts="--shell --alias --uninstall --yes -y"
                    COMPREPLY=( $(compgen -W "$opts" -- "$cur") )
                    return 0
                    ;;
                env)
                    local env_sub="list copy"
                    if [[ ${COMP_CWORD} -eq 3 ]]; then
                        COMPREPLY=( $(compgen -W "$env_sub" -- "$cur") )
                        return 0
                    fi
                    if [[ "${COMP_WORDS[3]}" == "copy" && ${COMP_CWORD} -eq 4 ]]; then
                        local templates="dev staging prod docker"
                        COMPREPLY=( $(compgen -W "$templates" -- "$cur") )
                        return 0
                    fi
                    ;;
            esac
            ;;
    esac
}

complete -F _c4ignite_complete c4ignite
complete -F _c4ignite_complete ./scripts/c4ignite
