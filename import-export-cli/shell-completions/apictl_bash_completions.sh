# bash completion for apictl                               -*- shell-script -*-

__debug()
{
    if [[ -n ${BASH_COMP_DEBUG_FILE} ]]; then
        echo "$*" >> "${BASH_COMP_DEBUG_FILE}"
    fi
}

# Homebrew on Macs have version 1.3 of bash-completion which doesn't include
# _init_completion. This is a very minimal version of that function.
__my_init_completion()
{
    COMPREPLY=()
    _get_comp_words_by_ref "$@" cur prev words cword
}

__index_of_word()
{
    local w word=$1
    shift
    index=0
    for w in "$@"; do
        [[ $w = "$word" ]] && return
        index=$((index+1))
    done
    index=-1
}

__contains_word()
{
    local w word=$1; shift
    for w in "$@"; do
        [[ $w = "$word" ]] && return
    done
    return 1
}

__handle_reply()
{
    __debug "${FUNCNAME[0]}"
    case $cur in
        -*)
            if [[ $(type -t compopt) = "builtin" ]]; then
                compopt -o nospace
            fi
            local allflags
            if [ ${#must_have_one_flag[@]} -ne 0 ]; then
                allflags=("${must_have_one_flag[@]}")
            else
                allflags=("${flags[*]} ${two_word_flags[*]}")
            fi
            COMPREPLY=( $(compgen -W "${allflags[*]}" -- "$cur") )
            if [[ $(type -t compopt) = "builtin" ]]; then
                [[ "${COMPREPLY[0]}" == *= ]] || compopt +o nospace
            fi

            # complete after --flag=abc
            if [[ $cur == *=* ]]; then
                if [[ $(type -t compopt) = "builtin" ]]; then
                    compopt +o nospace
                fi

                local index flag
                flag="${cur%%=*}"
                __index_of_word "${flag}" "${flags_with_completion[@]}"
                COMPREPLY=()
                if [[ ${index} -ge 0 ]]; then
                    PREFIX=""
                    cur="${cur#*=}"
                    ${flags_completion[${index}]}
                    if [ -n "${ZSH_VERSION}" ]; then
                        # zsh completion needs --flag= prefix
                        eval "COMPREPLY=( \"\${COMPREPLY[@]/#/${flag}=}\" )"
                    fi
                fi
            fi
            return 0;
            ;;
    esac

    # check if we are handling a flag with special work handling
    local index
    __index_of_word "${prev}" "${flags_with_completion[@]}"
    if [[ ${index} -ge 0 ]]; then
        ${flags_completion[${index}]}
        return
    fi

    # we are parsing a flag and don't have a special handler, no completion
    if [[ ${cur} != "${words[cword]}" ]]; then
        return
    fi

    local completions
    completions=("${commands[@]}")
    if [[ ${#must_have_one_noun[@]} -ne 0 ]]; then
        completions=("${must_have_one_noun[@]}")
    fi
    if [[ ${#must_have_one_flag[@]} -ne 0 ]]; then
        completions+=("${must_have_one_flag[@]}")
    fi
    COMPREPLY=( $(compgen -W "${completions[*]}" -- "$cur") )

    if [[ ${#COMPREPLY[@]} -eq 0 && ${#noun_aliases[@]} -gt 0 && ${#must_have_one_noun[@]} -ne 0 ]]; then
        COMPREPLY=( $(compgen -W "${noun_aliases[*]}" -- "$cur") )
    fi

    if [[ ${#COMPREPLY[@]} -eq 0 ]]; then
        declare -F __custom_func >/dev/null && __custom_func
    fi

    # available in bash-completion >= 2, not always present on macOS
    if declare -F __ltrim_colon_completions >/dev/null; then
        __ltrim_colon_completions "$cur"
    fi
}

# The arguments should be in the form "ext1|ext2|extn"
__handle_filename_extension_flag()
{
    local ext="$1"
    _filedir "@(${ext})"
}

__handle_subdirs_in_dir_flag()
{
    local dir="$1"
    pushd "${dir}" >/dev/null 2>&1 && _filedir -d && popd >/dev/null 2>&1
}

__handle_flag()
{
    __debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    # if a command required a flag, and we found it, unset must_have_one_flag()
    local flagname=${words[c]}
    local flagvalue
    # if the word contained an =
    if [[ ${words[c]} == *"="* ]]; then
        flagvalue=${flagname#*=} # take in as flagvalue after the =
        flagname=${flagname%%=*} # strip everything after the =
        flagname="${flagname}=" # but put the = back
    fi
    __debug "${FUNCNAME[0]}: looking for ${flagname}"
    if __contains_word "${flagname}" "${must_have_one_flag[@]}"; then
        must_have_one_flag=()
    fi

    # if you set a flag which only applies to this command, don't show subcommands
    if __contains_word "${flagname}" "${local_nonpersistent_flags[@]}"; then
      commands=()
    fi

    # keep flag value with flagname as flaghash
    if [ -n "${flagvalue}" ] ; then
        flaghash[${flagname}]=${flagvalue}
    elif [ -n "${words[ $((c+1)) ]}" ] ; then
        flaghash[${flagname}]=${words[ $((c+1)) ]}
    else
        flaghash[${flagname}]="true" # pad "true" for bool flag
    fi

    # skip the argument to a two word flag
    if __contains_word "${words[c]}" "${two_word_flags[@]}"; then
        c=$((c+1))
        # if we are looking for a flags value, don't show commands
        if [[ $c -eq $cword ]]; then
            commands=()
        fi
    fi

    c=$((c+1))

}

__handle_noun()
{
    __debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    if __contains_word "${words[c]}" "${must_have_one_noun[@]}"; then
        must_have_one_noun=()
    elif __contains_word "${words[c]}" "${noun_aliases[@]}"; then
        must_have_one_noun=()
    fi

    nouns+=("${words[c]}")
    c=$((c+1))
}

__handle_command()
{
    __debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    local next_command
    if [[ -n ${last_command} ]]; then
        next_command="_${last_command}_${words[c]//:/__}"
    else
        if [[ $c -eq 0 ]]; then
            next_command="_$(basename "${words[c]//:/__}")"
        else
            next_command="_${words[c]//:/__}"
        fi
    fi
    c=$((c+1))
    __debug "${FUNCNAME[0]}: looking for ${next_command}"
    declare -F "$next_command" >/dev/null && $next_command
}

__handle_word()
{
    if [[ $c -ge $cword ]]; then
        __handle_reply
        return
    fi
    __debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"
    if [[ "${words[c]}" == -* ]]; then
        __handle_flag
    elif __contains_word "${words[c]}" "${commands[@]}"; then
        __handle_command
    elif [[ $c -eq 0 ]] && __contains_word "$(basename "${words[c]}")" "${commands[@]}"; then
        __handle_command
    else
        __handle_noun
    fi
    __handle_word
}

_apictl_add_api()
{
    last_command="apictl_add_api"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--from-file=")
    two_word_flags+=("-f")
    local_nonpersistent_flags+=("--from-file=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--name=")
    two_word_flags+=("-n")
    local_nonpersistent_flags+=("--name=")
    flags+=("--namespace=")
    local_nonpersistent_flags+=("--namespace=")
    flags+=("--override")
    local_nonpersistent_flags+=("--override")
    flags+=("--replicas=")
    local_nonpersistent_flags+=("--replicas=")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_add()
{
    last_command="apictl_add"
    commands=()
    commands+=("api")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_add-env()
{
    last_command="apictl_add-env"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--admin=")
    local_nonpersistent_flags+=("--admin=")
    flags+=("--api_list=")
    local_nonpersistent_flags+=("--api_list=")
    flags+=("--apim=")
    local_nonpersistent_flags+=("--apim=")
    flags+=("--app_list=")
    local_nonpersistent_flags+=("--app_list=")
    flags+=("--environment=")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--registration=")
    local_nonpersistent_flags+=("--registration=")
    flags+=("--token=")
    local_nonpersistent_flags+=("--token=")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_export-api()
{
    last_command="apictl_export-api"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment=")
    flags+=("--format=")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--name=")
    two_word_flags+=("-n")
    local_nonpersistent_flags+=("--name=")
    flags+=("--preserveStatus")
    local_nonpersistent_flags+=("--preserveStatus")
    flags+=("--provider=")
    two_word_flags+=("-r")
    local_nonpersistent_flags+=("--provider=")
    flags+=("--version=")
    two_word_flags+=("-v")
    local_nonpersistent_flags+=("--version=")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_export-apis()
{
    last_command="apictl_export-apis"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment=")
    flags+=("--force")
    flags+=("--format=")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--preserveStatus")
    local_nonpersistent_flags+=("--preserveStatus")
    flags+=("--tenant=")
    two_word_flags+=("-t")
    local_nonpersistent_flags+=("--tenant=")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_export-app()
{
    last_command="apictl_export-app"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--name=")
    two_word_flags+=("-n")
    local_nonpersistent_flags+=("--name=")
    flags+=("--owner=")
    two_word_flags+=("-o")
    local_nonpersistent_flags+=("--owner=")
    flags+=("--withKeys")
    local_nonpersistent_flags+=("--withKeys")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_get-keys()
{
    last_command="apictl_get-keys"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--name=")
    two_word_flags+=("-n")
    local_nonpersistent_flags+=("--name=")
    flags+=("--provider=")
    two_word_flags+=("-r")
    local_nonpersistent_flags+=("--provider=")
    flags+=("--version=")
    two_word_flags+=("-v")
    local_nonpersistent_flags+=("--version=")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_import-api()
{
    last_command="apictl_import-api"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment=")
    flags+=("--file=")
    two_word_flags+=("-f")
    local_nonpersistent_flags+=("--file=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--params=")
    local_nonpersistent_flags+=("--params=")
    flags+=("--preserve-provider")
    local_nonpersistent_flags+=("--preserve-provider")
    flags+=("--skipCleanup")
    local_nonpersistent_flags+=("--skipCleanup")
    flags+=("--update")
    local_nonpersistent_flags+=("--update")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_flag+=("--file=")
    must_have_one_flag+=("-f")
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_import-app()
{
    last_command="apictl_import-app"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment=")
    flags+=("--file=")
    two_word_flags+=("-f")
    local_nonpersistent_flags+=("--file=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--owner=")
    two_word_flags+=("-o")
    local_nonpersistent_flags+=("--owner=")
    flags+=("--preserveOwner")
    local_nonpersistent_flags+=("--preserveOwner")
    flags+=("--skipKeys")
    local_nonpersistent_flags+=("--skipKeys")
    flags+=("--skipSubscriptions")
    flags+=("-s")
    local_nonpersistent_flags+=("--skipSubscriptions")
    flags+=("--update")
    local_nonpersistent_flags+=("--update")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_init()
{
    last_command="apictl_init"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--definition=")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--definition=")
    flags+=("--force")
    flags+=("-f")
    local_nonpersistent_flags+=("--force")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--oas=")
    local_nonpersistent_flags+=("--oas=")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_list_apis()
{
    last_command="apictl_list_apis"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment=")
    flags+=("--format=")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--query=")
    two_word_flags+=("-q")
    local_nonpersistent_flags+=("--query=")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_list_apps()
{
    last_command="apictl_list_apps"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment=")
    flags+=("--format=")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--owner=")
    two_word_flags+=("-o")
    local_nonpersistent_flags+=("--owner=")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_list_envs()
{
    last_command="apictl_list_envs"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--format=")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_list()
{
    last_command="apictl_list"
    commands=()
    commands+=("apis")
    commands+=("apps")
    commands+=("envs")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_login()
{
    last_command="apictl_login"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--password=")
    two_word_flags+=("-p")
    local_nonpersistent_flags+=("--password=")
    flags+=("--password-stdin")
    local_nonpersistent_flags+=("--password-stdin")
    flags+=("--username=")
    two_word_flags+=("-u")
    local_nonpersistent_flags+=("--username=")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_logout()
{
    last_command="apictl_logout"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_remove-env()
{
    last_command="apictl_remove-env"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_set()
{
    last_command="apictl_set"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--export-directory=")
    local_nonpersistent_flags+=("--export-directory=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--http-request-timeout=")
    local_nonpersistent_flags+=("--http-request-timeout=")
    flags+=("--mode=")
    two_word_flags+=("-m")
    local_nonpersistent_flags+=("--mode=")
    flags+=("--token-type=")
    two_word_flags+=("-t")
    local_nonpersistent_flags+=("--token-type=")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_update_api()
{
    last_command="apictl_update_api"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--from-file=")
    two_word_flags+=("-f")
    local_nonpersistent_flags+=("--from-file=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--name=")
    two_word_flags+=("-n")
    local_nonpersistent_flags+=("--name=")
    flags+=("--namespace=")
    local_nonpersistent_flags+=("--namespace=")
    flags+=("--replicas=")
    local_nonpersistent_flags+=("--replicas=")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_update()
{
    last_command="apictl_update"
    commands=()
    commands+=("api")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_version()
{
    last_command="apictl_version"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_apictl()
{
    last_command="apictl"
    commands=()
    commands+=("add")
    commands+=("add-env")
    commands+=("export-api")
    commands+=("export-apis")
    commands+=("export-app")
    commands+=("get-keys")
    commands+=("import-api")
    commands+=("import-app")
    commands+=("init")
    commands+=("list")
    commands+=("login")
    commands+=("logout")
    commands+=("remove-env")
    commands+=("set")
    commands+=("update")
    commands+=("version")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

__start_apictl()
{
    local cur prev words cword
    declare -A flaghash 2>/dev/null || :
    if declare -F _init_completion >/dev/null 2>&1; then
        _init_completion -s || return
    else
        __my_init_completion -n "=" || return
    fi

    local c=0
    local flags=()
    local two_word_flags=()
    local local_nonpersistent_flags=()
    local flags_with_completion=()
    local flags_completion=()
    local commands=("apictl")
    local must_have_one_flag=()
    local must_have_one_noun=()
    local last_command
    local nouns=()

    __handle_word
}

if [[ $(type -t compopt) = "builtin" ]]; then
    complete -o default -F __start_apictl apictl
else
    complete -o default -o nospace -F __start_apictl apictl
fi

# ex: ts=4 sw=4 et filetype=sh
