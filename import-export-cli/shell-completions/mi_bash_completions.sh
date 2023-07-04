# bash completion for mi                                   -*- shell-script -*-

__mi_debug()
{
    if [[ -n ${BASH_COMP_DEBUG_FILE:-} ]]; then
        echo "$*" >> "${BASH_COMP_DEBUG_FILE}"
    fi
}

# Homebrew on Macs have version 1.3 of bash-completion which doesn't include
# _init_completion. This is a very minimal version of that function.
__mi_init_completion()
{
    COMPREPLY=()
    _get_comp_words_by_ref "$@" cur prev words cword
}

__mi_index_of_word()
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

__mi_contains_word()
{
    local w word=$1; shift
    for w in "$@"; do
        [[ $w = "$word" ]] && return
    done
    return 1
}

__mi_handle_go_custom_completion()
{
    __mi_debug "${FUNCNAME[0]}: cur is ${cur}, words[*] is ${words[*]}, #words[@] is ${#words[@]}"

    local shellCompDirectiveError=1
    local shellCompDirectiveNoSpace=2
    local shellCompDirectiveNoFileComp=4
    local shellCompDirectiveFilterFileExt=8
    local shellCompDirectiveFilterDirs=16

    local out requestComp lastParam lastChar comp directive args

    # Prepare the command to request completions for the program.
    # Calling ${words[0]} instead of directly mi allows to handle aliases
    args=("${words[@]:1}")
    # Disable ActiveHelp which is not supported for bash completion v1
    requestComp="MI_ACTIVE_HELP=0 ${words[0]} __completeNoDesc ${args[*]}"

    lastParam=${words[$((${#words[@]}-1))]}
    lastChar=${lastParam:$((${#lastParam}-1)):1}
    __mi_debug "${FUNCNAME[0]}: lastParam ${lastParam}, lastChar ${lastChar}"

    if [ -z "${cur}" ] && [ "${lastChar}" != "=" ]; then
        # If the last parameter is complete (there is a space following it)
        # We add an extra empty parameter so we can indicate this to the go method.
        __mi_debug "${FUNCNAME[0]}: Adding extra empty parameter"
        requestComp="${requestComp} \"\""
    fi

    __mi_debug "${FUNCNAME[0]}: calling ${requestComp}"
    # Use eval to handle any environment variables and such
    out=$(eval "${requestComp}" 2>/dev/null)

    # Extract the directive integer at the very end of the output following a colon (:)
    directive=${out##*:}
    # Remove the directive
    out=${out%:*}
    if [ "${directive}" = "${out}" ]; then
        # There is not directive specified
        directive=0
    fi
    __mi_debug "${FUNCNAME[0]}: the completion directive is: ${directive}"
    __mi_debug "${FUNCNAME[0]}: the completions are: ${out}"

    if [ $((directive & shellCompDirectiveError)) -ne 0 ]; then
        # Error code.  No completion.
        __mi_debug "${FUNCNAME[0]}: received error from custom completion go code"
        return
    else
        if [ $((directive & shellCompDirectiveNoSpace)) -ne 0 ]; then
            if [[ $(type -t compopt) = "builtin" ]]; then
                __mi_debug "${FUNCNAME[0]}: activating no space"
                compopt -o nospace
            fi
        fi
        if [ $((directive & shellCompDirectiveNoFileComp)) -ne 0 ]; then
            if [[ $(type -t compopt) = "builtin" ]]; then
                __mi_debug "${FUNCNAME[0]}: activating no file completion"
                compopt +o default
            fi
        fi
    fi

    if [ $((directive & shellCompDirectiveFilterFileExt)) -ne 0 ]; then
        # File extension filtering
        local fullFilter filter filteringCmd
        # Do not use quotes around the $out variable or else newline
        # characters will be kept.
        for filter in ${out}; do
            fullFilter+="$filter|"
        done

        filteringCmd="_filedir $fullFilter"
        __mi_debug "File filtering command: $filteringCmd"
        $filteringCmd
    elif [ $((directive & shellCompDirectiveFilterDirs)) -ne 0 ]; then
        # File completion for directories only
        local subdir
        # Use printf to strip any trailing newline
        subdir=$(printf "%s" "${out}")
        if [ -n "$subdir" ]; then
            __mi_debug "Listing directories in $subdir"
            __mi_handle_subdirs_in_dir_flag "$subdir"
        else
            __mi_debug "Listing directories in ."
            _filedir -d
        fi
    else
        while IFS='' read -r comp; do
            COMPREPLY+=("$comp")
        done < <(compgen -W "${out}" -- "$cur")
    fi
}

__mi_handle_reply()
{
    __mi_debug "${FUNCNAME[0]}"
    local comp
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
            while IFS='' read -r comp; do
                COMPREPLY+=("$comp")
            done < <(compgen -W "${allflags[*]}" -- "$cur")
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
                __mi_index_of_word "${flag}" "${flags_with_completion[@]}"
                COMPREPLY=()
                if [[ ${index} -ge 0 ]]; then
                    PREFIX=""
                    cur="${cur#*=}"
                    ${flags_completion[${index}]}
                    if [ -n "${ZSH_VERSION:-}" ]; then
                        # zsh completion needs --flag= prefix
                        eval "COMPREPLY=( \"\${COMPREPLY[@]/#/${flag}=}\" )"
                    fi
                fi
            fi

            if [[ -z "${flag_parsing_disabled}" ]]; then
                # If flag parsing is enabled, we have completed the flags and can return.
                # If flag parsing is disabled, we may not know all (or any) of the flags, so we fallthrough
                # to possibly call handle_go_custom_completion.
                return 0;
            fi
            ;;
    esac

    # check if we are handling a flag with special work handling
    local index
    __mi_index_of_word "${prev}" "${flags_with_completion[@]}"
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
        completions+=("${must_have_one_noun[@]}")
    elif [[ -n "${has_completion_function}" ]]; then
        # if a go completion function is provided, defer to that function
        __mi_handle_go_custom_completion
    fi
    if [[ ${#must_have_one_flag[@]} -ne 0 ]]; then
        completions+=("${must_have_one_flag[@]}")
    fi
    while IFS='' read -r comp; do
        COMPREPLY+=("$comp")
    done < <(compgen -W "${completions[*]}" -- "$cur")

    if [[ ${#COMPREPLY[@]} -eq 0 && ${#noun_aliases[@]} -gt 0 && ${#must_have_one_noun[@]} -ne 0 ]]; then
        while IFS='' read -r comp; do
            COMPREPLY+=("$comp")
        done < <(compgen -W "${noun_aliases[*]}" -- "$cur")
    fi

    if [[ ${#COMPREPLY[@]} -eq 0 ]]; then
        if declare -F __mi_custom_func >/dev/null; then
            # try command name qualified custom func
            __mi_custom_func
        else
            # otherwise fall back to unqualified for compatibility
            declare -F __custom_func >/dev/null && __custom_func
        fi
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
__mi_handle_filename_extension_flag()
{
    local ext="$1"
    _filedir "@(${ext})"
}

__mi_handle_subdirs_in_dir_flag()
{
    local dir="$1"
    pushd "${dir}" >/dev/null 2>&1 && _filedir -d && popd >/dev/null 2>&1 || return
}

__mi_handle_flag()
{
    __mi_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    # if a command required a flag, and we found it, unset must_have_one_flag()
    local flagname=${words[c]}
    local flagvalue=""
    # if the word contained an =
    if [[ ${words[c]} == *"="* ]]; then
        flagvalue=${flagname#*=} # take in as flagvalue after the =
        flagname=${flagname%=*} # strip everything after the =
        flagname="${flagname}=" # but put the = back
    fi
    __mi_debug "${FUNCNAME[0]}: looking for ${flagname}"
    if __mi_contains_word "${flagname}" "${must_have_one_flag[@]}"; then
        must_have_one_flag=()
    fi

    # if you set a flag which only applies to this command, don't show subcommands
    if __mi_contains_word "${flagname}" "${local_nonpersistent_flags[@]}"; then
      commands=()
    fi

    # keep flag value with flagname as flaghash
    # flaghash variable is an associative array which is only supported in bash > 3.
    if [[ -z "${BASH_VERSION:-}" || "${BASH_VERSINFO[0]:-}" -gt 3 ]]; then
        if [ -n "${flagvalue}" ] ; then
            flaghash[${flagname}]=${flagvalue}
        elif [ -n "${words[ $((c+1)) ]}" ] ; then
            flaghash[${flagname}]=${words[ $((c+1)) ]}
        else
            flaghash[${flagname}]="true" # pad "true" for bool flag
        fi
    fi

    # skip the argument to a two word flag
    if [[ ${words[c]} != *"="* ]] && __mi_contains_word "${words[c]}" "${two_word_flags[@]}"; then
        __mi_debug "${FUNCNAME[0]}: found a flag ${words[c]}, skip the next argument"
        c=$((c+1))
        # if we are looking for a flags value, don't show commands
        if [[ $c -eq $cword ]]; then
            commands=()
        fi
    fi

    c=$((c+1))

}

__mi_handle_noun()
{
    __mi_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    if __mi_contains_word "${words[c]}" "${must_have_one_noun[@]}"; then
        must_have_one_noun=()
    elif __mi_contains_word "${words[c]}" "${noun_aliases[@]}"; then
        must_have_one_noun=()
    fi

    nouns+=("${words[c]}")
    c=$((c+1))
}

__mi_handle_command()
{
    __mi_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    local next_command
    if [[ -n ${last_command} ]]; then
        next_command="_${last_command}_${words[c]//:/__}"
    else
        if [[ $c -eq 0 ]]; then
            next_command="_mi_root_command"
        else
            next_command="_${words[c]//:/__}"
        fi
    fi
    c=$((c+1))
    __mi_debug "${FUNCNAME[0]}: looking for ${next_command}"
    declare -F "$next_command" >/dev/null && $next_command
}

__mi_handle_word()
{
    if [[ $c -ge $cword ]]; then
        __mi_handle_reply
        return
    fi
    __mi_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"
    if [[ "${words[c]}" == -* ]]; then
        __mi_handle_flag
    elif __mi_contains_word "${words[c]}" "${commands[@]}"; then
        __mi_handle_command
    elif [[ $c -eq 0 ]]; then
        __mi_handle_command
    elif __mi_contains_word "${words[c]}" "${command_aliases[@]}"; then
        # aliashash variable is an associative array which is only supported in bash > 3.
        if [[ -z "${BASH_VERSION:-}" || "${BASH_VERSINFO[0]:-}" -gt 3 ]]; then
            words[c]=${aliashash[${words[c]}]}
            __mi_handle_command
        else
            __mi_handle_noun
        fi
    else
        __mi_handle_noun
    fi
    __mi_handle_word
}

_mi_activate_endpoint()
{
    last_command="mi_activate_endpoint"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_activate_help()
{
    last_command="mi_activate_help"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    has_completion_function=1
    noun_aliases=()
}

_mi_activate_message-processor()
{
    last_command="mi_activate_message-processor"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_activate_proxy-service()
{
    last_command="mi_activate_proxy-service"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_activate()
{
    last_command="mi_activate"

    command_aliases=()

    commands=()
    commands+=("endpoint")
    commands+=("help")
    commands+=("message-processor")
    commands+=("proxy-service")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_mi_add_env()
{
    last_command="mi_add_env"

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
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_mi_add_help()
{
    last_command="mi_add_help"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    has_completion_function=1
    noun_aliases=()
}

_mi_add_log-level()
{
    last_command="mi_add_log-level"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_add_role()
{
    last_command="mi_add_role"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_add_user()
{
    last_command="mi_add_user"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_add()
{
    last_command="mi_add"

    command_aliases=()

    commands=()
    commands+=("env")
    commands+=("help")
    commands+=("log-level")
    commands+=("role")
    commands+=("user")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_mi_deactivate_endpoint()
{
    last_command="mi_deactivate_endpoint"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_deactivate_help()
{
    last_command="mi_deactivate_help"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    has_completion_function=1
    noun_aliases=()
}

_mi_deactivate_message-processor()
{
    last_command="mi_deactivate_message-processor"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_deactivate_proxy-service()
{
    last_command="mi_deactivate_proxy-service"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_deactivate()
{
    last_command="mi_deactivate"

    command_aliases=()

    commands=()
    commands+=("endpoint")
    commands+=("help")
    commands+=("message-processor")
    commands+=("proxy-service")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_mi_delete_env()
{
    last_command="mi_delete_env"

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
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_mi_delete_help()
{
    last_command="mi_delete_help"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    has_completion_function=1
    noun_aliases=()
}

_mi_delete_role()
{
    last_command="mi_delete_role"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--domain=")
    two_word_flags+=("--domain")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--domain")
    local_nonpersistent_flags+=("--domain=")
    local_nonpersistent_flags+=("-d")
    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_delete_user()
{
    last_command="mi_delete_user"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--domain=")
    two_word_flags+=("--domain")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--domain")
    local_nonpersistent_flags+=("--domain=")
    local_nonpersistent_flags+=("-d")
    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_delete()
{
    last_command="mi_delete"

    command_aliases=()

    commands=()
    commands+=("env")
    commands+=("help")
    commands+=("role")
    commands+=("user")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_mi_get_apis()
{
    last_command="mi_get_apis"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--format=")
    two_word_flags+=("--format")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_get_composite-apps()
{
    last_command="mi_get_composite-apps"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--format=")
    two_word_flags+=("--format")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_get_connectors()
{
    last_command="mi_get_connectors"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--format=")
    two_word_flags+=("--format")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_get_data-services()
{
    last_command="mi_get_data-services"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--format=")
    two_word_flags+=("--format")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_get_endpoints()
{
    last_command="mi_get_endpoints"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--format=")
    two_word_flags+=("--format")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_get_envs()
{
    last_command="mi_get_envs"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--format=")
    two_word_flags+=("--format")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_mi_get_help()
{
    last_command="mi_get_help"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    has_completion_function=1
    noun_aliases=()
}

_mi_get_inbound-endpoints()
{
    last_command="mi_get_inbound-endpoints"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--format=")
    two_word_flags+=("--format")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_get_local-entries()
{
    last_command="mi_get_local-entries"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--format=")
    two_word_flags+=("--format")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_get_log-levels()
{
    last_command="mi_get_log-levels"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--format=")
    two_word_flags+=("--format")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_get_logs()
{
    last_command="mi_get_logs"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--format=")
    two_word_flags+=("--format")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--path=")
    two_word_flags+=("--path")
    two_word_flags+=("-p")
    local_nonpersistent_flags+=("--path")
    local_nonpersistent_flags+=("--path=")
    local_nonpersistent_flags+=("-p")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_get_message-processors()
{
    last_command="mi_get_message-processors"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--format=")
    two_word_flags+=("--format")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_get_message-stores()
{
    last_command="mi_get_message-stores"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--format=")
    two_word_flags+=("--format")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_get_proxy-services()
{
    last_command="mi_get_proxy-services"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--format=")
    two_word_flags+=("--format")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_get_roles()
{
    last_command="mi_get_roles"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--domain=")
    two_word_flags+=("--domain")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--domain")
    local_nonpersistent_flags+=("--domain=")
    local_nonpersistent_flags+=("-d")
    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--format=")
    two_word_flags+=("--format")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_get_sequences()
{
    last_command="mi_get_sequences"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--format=")
    two_word_flags+=("--format")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_get_tasks()
{
    last_command="mi_get_tasks"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--format=")
    two_word_flags+=("--format")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_get_templates()
{
    last_command="mi_get_templates"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--format=")
    two_word_flags+=("--format")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_get_transaction-counts()
{
    last_command="mi_get_transaction-counts"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--format=")
    two_word_flags+=("--format")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_get_transaction-reports()
{
    last_command="mi_get_transaction-reports"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--path=")
    two_word_flags+=("--path")
    two_word_flags+=("-p")
    local_nonpersistent_flags+=("--path")
    local_nonpersistent_flags+=("--path=")
    local_nonpersistent_flags+=("-p")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_get_users()
{
    last_command="mi_get_users"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--domain=")
    two_word_flags+=("--domain")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--domain")
    local_nonpersistent_flags+=("--domain=")
    local_nonpersistent_flags+=("-d")
    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--format=")
    two_word_flags+=("--format")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--pattern=")
    two_word_flags+=("--pattern")
    two_word_flags+=("-p")
    local_nonpersistent_flags+=("--pattern")
    local_nonpersistent_flags+=("--pattern=")
    local_nonpersistent_flags+=("-p")
    flags+=("--role=")
    two_word_flags+=("--role")
    two_word_flags+=("-r")
    local_nonpersistent_flags+=("--role")
    local_nonpersistent_flags+=("--role=")
    local_nonpersistent_flags+=("-r")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_get()
{
    last_command="mi_get"

    command_aliases=()

    commands=()
    commands+=("apis")
    commands+=("composite-apps")
    commands+=("connectors")
    commands+=("data-services")
    commands+=("endpoints")
    commands+=("envs")
    commands+=("help")
    commands+=("inbound-endpoints")
    commands+=("local-entries")
    commands+=("log-levels")
    commands+=("logs")
    commands+=("message-processors")
    commands+=("message-stores")
    commands+=("proxy-services")
    commands+=("roles")
    commands+=("sequences")
    commands+=("tasks")
    commands+=("templates")
    commands+=("transaction-counts")
    commands+=("transaction-reports")
    commands+=("users")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_mi_help()
{
    last_command="mi_help"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    has_completion_function=1
    noun_aliases=()
}

_mi_login()
{
    last_command="mi_login"

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
    local_nonpersistent_flags+=("-h")
    flags+=("--password=")
    two_word_flags+=("--password")
    two_word_flags+=("-p")
    local_nonpersistent_flags+=("--password")
    local_nonpersistent_flags+=("--password=")
    local_nonpersistent_flags+=("-p")
    flags+=("--password-stdin")
    local_nonpersistent_flags+=("--password-stdin")
    flags+=("--username=")
    two_word_flags+=("--username")
    two_word_flags+=("-u")
    local_nonpersistent_flags+=("--username")
    local_nonpersistent_flags+=("--username=")
    local_nonpersistent_flags+=("-u")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_mi_logout()
{
    last_command="mi_logout"

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
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_mi_update_hashicorp-secret()
{
    last_command="mi_update_hashicorp-secret"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_update_help()
{
    last_command="mi_update_help"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    has_completion_function=1
    noun_aliases=()
}

_mi_update_log-level()
{
    last_command="mi_update_log-level"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_update_user()
{
    last_command="mi_update_user"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--environment=")
    two_word_flags+=("--environment")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--environment")
    local_nonpersistent_flags+=("--environment=")
    local_nonpersistent_flags+=("-e")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_flag+=("--environment=")
    must_have_one_flag+=("-e")
    must_have_one_noun=()
    noun_aliases=()
}

_mi_update()
{
    last_command="mi_update"

    command_aliases=()

    commands=()
    commands+=("hashicorp-secret")
    commands+=("help")
    commands+=("log-level")
    commands+=("user")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_mi_version()
{
    last_command="mi_version"

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
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_mi_root_command()
{
    last_command="mi"

    command_aliases=()

    commands=()
    commands+=("activate")
    commands+=("add")
    commands+=("deactivate")
    commands+=("delete")
    commands+=("get")
    commands+=("help")
    commands+=("login")
    commands+=("logout")
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
    local_nonpersistent_flags+=("-h")
    flags+=("--insecure")
    flags+=("-k")
    flags+=("--verbose")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

__start_mi()
{
    local cur prev words cword split
    declare -A flaghash 2>/dev/null || :
    declare -A aliashash 2>/dev/null || :
    if declare -F _init_completion >/dev/null 2>&1; then
        _init_completion -s || return
    else
        __mi_init_completion -n "=" || return
    fi

    local c=0
    local flag_parsing_disabled=
    local flags=()
    local two_word_flags=()
    local local_nonpersistent_flags=()
    local flags_with_completion=()
    local flags_completion=()
    local commands=("mi")
    local command_aliases=()
    local must_have_one_flag=()
    local must_have_one_noun=()
    local has_completion_function=""
    local last_command=""
    local nouns=()
    local noun_aliases=()

    __mi_handle_word
}

if [[ $(type -t compopt) = "builtin" ]]; then
    complete -o default -F __start_mi mi
else
    complete -o default -o nospace -F __start_mi mi
fi

# ex: ts=4 sw=4 et filetype=sh
