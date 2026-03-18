#!/bin/bash
# gostim2 bash completion

_gostim2_completions() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    
    # List of all options (single dash)
    opts="-assets -bg-color -config -csv -display -dlp -end-splash -fixation-color -font -font-size -fullscreen -height -no-fixation -no-vsync -res -results-dir -scale -skip-wait -start-splash -stimuli-dir -subject -text-color -tsv -version -vrr -width"

    # Handle specific options that expect arguments
    case "${prev}" in
        -assets|-stimuli-dir|-results-dir)
            if type _filedir &>/dev/null; then
                _filedir -d
            else
                COMPREPLY=( $(compgen -d -- "${cur}") )
            fi
            return 0
            ;;
        -config)
            if type _filedir &>/dev/null; then
                _filedir toml
            else
                COMPREPLY=( $(compgen -f -X '!*.toml' -- "${cur}") )
            fi
            return 0
            ;;
        -csv|-tsv)
            if type _filedir &>/dev/null; then
                _filedir '?(t|c)sv'
            else
                COMPREPLY=( $(compgen -f -X '!*.@(csv|tsv)' -- "${cur}") )
            fi
            return 0
            ;;
        -end-splash|-start-splash)
            if type _filedir &>/dev/null; then
                _filedir '@(jpg|jpeg|png|bmp|gif)'
            else
                COMPREPLY=( $(compgen -f -X '!*.@(jpg|jpeg|png|bmp|gif)' -- "${cur}") )
            fi
            return 0
            ;;
        -font)
            if type _filedir &>/dev/null; then
                _filedir ttf
            else
                COMPREPLY=( $(compgen -f -X '!*.ttf' -- "${cur}") )
            fi
            return 0
            ;;
        -res)
            local resolutions="Autodetect 1920x1080 1440x900 1280x720 1024x768 800x600"
            COMPREPLY=( $(compgen -W "${resolutions}" -- "${cur}") )
            return 0
            ;;
        -dlp)
            # Typically a serial device like /dev/ttyUSB0
            if type _filedir &>/dev/null; then
                _filedir
            else
                COMPREPLY=( $(compgen -f -- "${cur}") )
            fi
            return 0
            ;;
    esac

    # Complete the option itself
    if [[ ${cur} == -* ]] ; then
        COMPREPLY=( $(compgen -W "${opts}" -- "${cur}") )
        return 0
    fi
}

# Complete for 'gostim2' binary
complete -F _gostim2_completions gostim2
# Complete for './gostim2' if running from current directory
complete -F _gostim2_completions ./gostim2
