# copy this file to /etc/bash_completion.d/ and source it to enable bash completion
_apimcli()
{
    local current previous options base
    COMPREPLY=()
    current="${COMP_WORDS[COMP_CWORD]}"
    previous="${COMP_WORDS[COMP_CWORD-1]}"

    options="export-api import-api list add-env remove-env reset-user version author"

    case "${previous}" in
        export-api)
            local flags="--name -n --version -v --environment -e --help -h"
            COMPREPLY=( $(compgen -W "${flags}" -- ${current}) )
            return 0
            ;;
        import-api)
            local flags="--name -n --environment -e --help -h"
            COMPREPLY=( $(compgen -W "${flags}" -- ${current}) )
            return 0
            ;;
        list)
            local flags="apis envs"
            COMPREPLY=( $(compgen -W "${flags}" -- ${current}) )
            return 0
            ;;
        add-env)
            local flags="--name -n --publisher --registration --token --help -h"
            COMPREPLY=( $(compgen -W "${flags}" -- ${current}) )
            return 0
            ;;
        remove-env)
            local flags="--name -n --help -h"
            COMPREPLY=( $(compgen -W "${flags}" -- ${current}) )
            return 0
            ;;
        reset-user)
            local flags="--environment -e --help -h"
            COMPREPLY=( $(compgen -W "${flags}" -- ${current}) )
            return 0
            ;;
        set)
            local flags="--export-directory --http-request-timeout"
            COMPREPLY=( $(compgen -W "${flags}" -- ${current}) )
            return 0
            ;;
        version)
            return 0
            ;;
        author)
            return 0
            ;;
        *)
        ;;

    esac

    COMPREPLY=($(compgen -W "${options}" -- ${current}))  
    return 0
}

complete -F _apimcli apimcli
complete -F _apimcli ./apimcli
complete -F _apimcli import-export-cli
complete -F _apimcli ./import-export-cli
