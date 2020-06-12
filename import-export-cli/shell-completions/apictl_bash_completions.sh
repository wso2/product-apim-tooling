# bash completion for apictl                               -*- shell-script -*-

__apictl_debug()
{
    if [[ -n ${BASH_COMP_DEBUG_FILE} ]]; then
        echo "$*" >> "${BASH_COMP_DEBUG_FILE}"
    fi
}

# Homebrew on Macs have version 1.3 of bash-completion which doesn't include
# _init_completion. This is a very minimal version of that function.
__apictl_init_completion()
{
    COMPREPLY=()
    _get_comp_words_by_ref "$@" cur prev words cword
}

__apictl_index_of_word()
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

__apictl_contains_word()
{
    local w word=$1; shift
    for w in "$@"; do
        [[ $w = "$word" ]] && return
    done
    return 1
}

__apictl_handle_reply()
{
    __apictl_debug "${FUNCNAME[0]}"
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
                flag="${cur%=*}"
                __apictl_index_of_word "${flag}" "${flags_with_completion[@]}"
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
    __apictl_index_of_word "${prev}" "${flags_with_completion[@]}"
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

    # If there is only 1 completion and it is a flag with an = it will be completed
    # but we don't want a space after the =
    if [[ "${#COMPREPLY[@]}" -eq "1" ]] && [[ $(type -t compopt) = "builtin" ]] && [[ "${COMPREPLY[0]}" == --*= ]]; then
       compopt -o nospace
    fi
}

# The arguments should be in the form "ext1|ext2|extn"
__apictl_handle_filename_extension_flag()
{
    local ext="$1"
    _filedir "@(${ext})"
}

__apictl_handle_subdirs_in_dir_flag()
{
    local dir="$1"
    pushd "${dir}" >/dev/null 2>&1 && _filedir -d && popd >/dev/null 2>&1
}

__apictl_handle_flag()
{
    __apictl_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    # if a command required a flag, and we found it, unset must_have_one_flag()
    local flagname=${words[c]}
    local flagvalue
    # if the word contained an =
    if [[ ${words[c]} == *"="* ]]; then
        flagvalue=${flagname#*=} # take in as flagvalue after the =
        flagname=${flagname%=*} # strip everything after the =
        flagname="${flagname}=" # but put the = back
    fi
    __apictl_debug "${FUNCNAME[0]}: looking for ${flagname}"
    if __apictl_contains_word "${flagname}" "${must_have_one_flag[@]}"; then
        must_have_one_flag=()
    fi

    # if you set a flag which only applies to this command, don't show subcommands
    if __apictl_contains_word "${flagname}" "${local_nonpersistent_flags[@]}"; then
      commands=()
    fi

    # keep flag value with flagname as flaghash
    # flaghash variable is an associative array which is only supported in bash > 3.
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        if [ -n "${flagvalue}" ] ; then
            flaghash[${flagname}]=${flagvalue}
        elif [ -n "${words[ $((c+1)) ]}" ] ; then
            flaghash[${flagname}]=${words[ $((c+1)) ]}
        else
            flaghash[${flagname}]="true" # pad "true" for bool flag
        fi
    fi

    # skip the argument to a two word flag
    if __apictl_contains_word "${words[c]}" "${two_word_flags[@]}"; then
        c=$((c+1))
        # if we are looking for a flags value, don't show commands
        if [[ $c -eq $cword ]]; then
            commands=()
        fi
    fi

    c=$((c+1))

}

__apictl_handle_noun()
{
    __apictl_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    if __apictl_contains_word "${words[c]}" "${must_have_one_noun[@]}"; then
        must_have_one_noun=()
    elif __apictl_contains_word "${words[c]}" "${noun_aliases[@]}"; then
        must_have_one_noun=()
    fi

    nouns+=("${words[c]}")
    c=$((c+1))
}

__apictl_handle_command()
{
    __apictl_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    local next_command
    if [[ -n ${last_command} ]]; then
        next_command="_${last_command}_${words[c]//:/__}"
    else
        if [[ $c -eq 0 ]]; then
            next_command="_apictl_root_command"
        else
            next_command="_${words[c]//:/__}"
        fi
    fi
    c=$((c+1))
    __apictl_debug "${FUNCNAME[0]}: looking for ${next_command}"
    declare -F "$next_command" >/dev/null && $next_command
}

__apictl_handle_word()
{
    if [[ $c -ge $cword ]]; then
        __apictl_handle_reply
        return
    fi
    __apictl_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"
    if [[ "${words[c]}" == -* ]]; then
        __apictl_handle_flag
    elif __apictl_contains_word "${words[c]}" "${commands[@]}"; then
        __apictl_handle_command
    elif [[ $c -eq 0 ]]; then
        __apictl_handle_command
    elif __apictl_contains_word "${words[c]}" "${command_aliases[@]}"; then
        # aliashash variable is an associative array which is only supported in bash > 3.
        if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
            words[c]=${aliashash[${words[c]}]}
            __apictl_handle_command
        else
            __apictl_handle_noun
        fi
    else
        __apictl_handle_noun
    fi
    __apictl_handle_word
}

_apictl_add_api()
{
    last_command="apictl_add_api"

    command_aliases=()

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
    flags+=("--mode=")
    two_word_flags+=("-m")
    local_nonpersistent_flags+=("--mode=")
    flags+=("--name=")
    two_word_flags+=("-n")
    local_nonpersistent_flags+=("--name=")
    flags+=("--namespace=")
    local_nonpersistent_flags+=("--namespace=")
    flags+=("--override")
    local_nonpersistent_flags+=("--override")
    flags+=("--replicas=")
    local_nonpersistent_flags+=("--replicas=")
    flags+=("--version=")
    two_word_flags+=("-v")
    local_nonpersistent_flags+=("--version=")
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

    command_aliases=()

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

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--admin=")
    local_nonpersistent_flags+=("--admin=")
    flags+=("--apim=")
    local_nonpersistent_flags+=("--apim=")
    flags+=("--devportal=")
    local_nonpersistent_flags+=("--devportal=")
    flags+=("--environment=")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--publisher=")
    local_nonpersistent_flags+=("--publisher=")
    flags+=("--registration=")
    local_nonpersistent_flags+=("--registration=")
    flags+=("--token=")
    local_nonpersistent_flags+=("--token=")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_flag+=("--token=")
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_change_registry()
{
    last_command="apictl_change_registry"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--key-file=")
    two_word_flags+=("-c")
    local_nonpersistent_flags+=("--key-file=")
    flags+=("--password=")
    two_word_flags+=("-p")
    local_nonpersistent_flags+=("--password=")
    flags+=("--password-stdin")
    local_nonpersistent_flags+=("--password-stdin")
    flags+=("--registry-type=")
    two_word_flags+=("-R")
    local_nonpersistent_flags+=("--registry-type=")
    flags+=("--repository=")
    two_word_flags+=("-r")
    local_nonpersistent_flags+=("--repository=")
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

_apictl_change()
{
    last_command="apictl_change"

    command_aliases=()

    commands=()
    commands+=("registry")

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

_apictl_change-status_api()
{
    last_command="apictl_change-status_api"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--action=")
    two_word_flags+=("-a")
    local_nonpersistent_flags+=("--action=")
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
    must_have_one_flag+=("--action=")
    must_have_one_flag+=("-a")
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_flag+=("--name=")
    must_have_one_flag+=("-n")
    must_have_one_flag+=("--version=")
    must_have_one_flag+=("-v")
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_change-status()
{
    last_command="apictl_change-status"

    command_aliases=()

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

_apictl_delete_api()
{
    last_command="apictl_delete_api"

    command_aliases=()

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
    must_have_one_flag+=("--name=")
    must_have_one_flag+=("-n")
    must_have_one_flag+=("--version=")
    must_have_one_flag+=("-v")
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_delete_api-product()
{
    last_command="apictl_delete_api-product"

    command_aliases=()

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
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_flag+=("--name=")
    must_have_one_flag+=("-n")
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_delete_app()
{
    last_command="apictl_delete_app"

    command_aliases=()

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
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_flag+=("--name=")
    must_have_one_flag+=("-n")
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_delete()
{
    last_command="apictl_delete"

    command_aliases=()

    commands=()
    commands+=("api")
    commands+=("api-product")
    commands+=("app")

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

_apictl_export_api-product()
{
    last_command="apictl_export_api-product"

    command_aliases=()

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
    must_have_one_flag+=("--name=")
    must_have_one_flag+=("-n")
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_export()
{
    last_command="apictl_export"

    command_aliases=()

    commands=()
    commands+=("api-product")

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

_apictl_export-api()
{
    last_command="apictl_export-api"

    command_aliases=()

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
    must_have_one_flag+=("--name=")
    must_have_one_flag+=("-n")
    must_have_one_flag+=("--version=")
    must_have_one_flag+=("-v")
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_export-apis()
{
    last_command="apictl_export-apis"

    command_aliases=()

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

    command_aliases=()

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
    must_have_one_flag+=("--name=")
    must_have_one_flag+=("-n")
    must_have_one_flag+=("--owner=")
    must_have_one_flag+=("-o")
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_get-keys()
{
    last_command="apictl_get-keys"

    command_aliases=()

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
    must_have_one_flag+=("--name=")
    must_have_one_flag+=("-n")
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_import_api-product()
{
    last_command="apictl_import_api-product"

    command_aliases=()

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
    flags+=("--import-apis")
    local_nonpersistent_flags+=("--import-apis")
    flags+=("--preserve-provider")
    local_nonpersistent_flags+=("--preserve-provider")
    flags+=("--skipCleanup")
    local_nonpersistent_flags+=("--skipCleanup")
    flags+=("--update-api-product")
    local_nonpersistent_flags+=("--update-api-product")
    flags+=("--update-apis")
    local_nonpersistent_flags+=("--update-apis")
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

_apictl_import()
{
    last_command="apictl_import"

    command_aliases=()

    commands=()
    commands+=("api-product")

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

_apictl_import-api()
{
    last_command="apictl_import-api"

    command_aliases=()

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

    command_aliases=()

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
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_flag+=("--file=")
    must_have_one_flag+=("-f")
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_init()
{
    last_command="apictl_init"

    command_aliases=()

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
    flags+=("--initial-state=")
    local_nonpersistent_flags+=("--initial-state=")
    flags+=("--oas=")
    local_nonpersistent_flags+=("--oas=")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_install_api-operator()
{
    last_command="apictl_install_api-operator"

    command_aliases=()

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
    flags+=("--key-file=")
    two_word_flags+=("-c")
    local_nonpersistent_flags+=("--key-file=")
    flags+=("--password=")
    two_word_flags+=("-p")
    local_nonpersistent_flags+=("--password=")
    flags+=("--password-stdin")
    local_nonpersistent_flags+=("--password-stdin")
    flags+=("--registry-type=")
    two_word_flags+=("-R")
    local_nonpersistent_flags+=("--registry-type=")
    flags+=("--repository=")
    two_word_flags+=("-r")
    local_nonpersistent_flags+=("--repository=")
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

_apictl_install_wso2am-operator()
{
    last_command="apictl_install_wso2am-operator"

    command_aliases=()

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
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_apictl_install()
{
    last_command="apictl_install"

    command_aliases=()

    commands=()
    commands+=("api-operator")
    commands+=("wso2am-operator")

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

_apictl_list_api-products()
{
    last_command="apictl_list_api-products"

    command_aliases=()

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
    flags+=("--limit=")
    two_word_flags+=("-l")
    local_nonpersistent_flags+=("--limit=")
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

_apictl_list_apis()
{
    last_command="apictl_list_apis"

    command_aliases=()

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
    flags+=("--limit=")
    two_word_flags+=("-l")
    local_nonpersistent_flags+=("--limit=")
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

    command_aliases=()

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
    flags+=("--limit=")
    two_word_flags+=("-l")
    local_nonpersistent_flags+=("--limit=")
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

    command_aliases=()

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

    command_aliases=()

    commands=()
    commands+=("api-products")
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

    command_aliases=()

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

    command_aliases=()

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

_apictl_remove_env()
{
    last_command="apictl_remove_env"

    command_aliases=()

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

_apictl_remove()
{
    last_command="apictl_remove"

    command_aliases=()

    commands=()
    commands+=("env")

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

_apictl_set()
{
    last_command="apictl_set"

    command_aliases=()

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

_apictl_uninstall_api-operator()
{
    last_command="apictl_uninstall_api-operator"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--force")
    local_nonpersistent_flags+=("--force")
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

_apictl_uninstall_wso2am-operator()
{
    last_command="apictl_uninstall_wso2am-operator"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--force")
    local_nonpersistent_flags+=("--force")
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

_apictl_uninstall()
{
    last_command="apictl_uninstall"

    command_aliases=()

    commands=()
    commands+=("api-operator")
    commands+=("wso2am-operator")

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

_apictl_update_api()
{
    last_command="apictl_update_api"

    command_aliases=()

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
    flags+=("--mode=")
    two_word_flags+=("-m")
    local_nonpersistent_flags+=("--mode=")
    flags+=("--name=")
    two_word_flags+=("-n")
    local_nonpersistent_flags+=("--name=")
    flags+=("--namespace=")
    local_nonpersistent_flags+=("--namespace=")
    flags+=("--replicas=")
    local_nonpersistent_flags+=("--replicas=")
    flags+=("--version=")
    two_word_flags+=("-v")
    local_nonpersistent_flags+=("--version=")
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

    command_aliases=()

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

    command_aliases=()

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

_apictl_root_command()
{
    last_command="apictl"

    command_aliases=()

    commands=()
    commands+=("add")
    commands+=("add-env")
    commands+=("change")
    commands+=("change-status")
    commands+=("delete")
    commands+=("export")
    commands+=("export-api")
    commands+=("export-apis")
    commands+=("export-app")
    commands+=("get-keys")
    commands+=("import")
    commands+=("import-api")
    commands+=("import-app")
    commands+=("init")
    commands+=("install")
    commands+=("list")
    commands+=("login")
    commands+=("logout")
    commands+=("remove")
    commands+=("set")
    commands+=("uninstall")
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
    declare -A aliashash 2>/dev/null || :
    if declare -F _init_completion >/dev/null 2>&1; then
        _init_completion -s || return
    else
        __apictl_init_completion -n "=" || return
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

    __apictl_handle_word
}

if [[ $(type -t compopt) = "builtin" ]]; then
    complete -o default -F __start_apictl apictl
else
    complete -o default -o nospace -F __start_apictl apictl
fi

# ex: ts=4 sw=4 et filetype=sh
