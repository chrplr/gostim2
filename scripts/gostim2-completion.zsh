#compdef gostim2

_gostim2() {
    local -a opts
    opts=(
        "-assets[Directory containing stimuli (alias for -stimuli-dir)]:directory:_files -/"
        "-bg-color[Background color (R,G,B,A)]:color:"
        "-config[Load parameters from a TOML config file]:file:_files -g '*.toml'"
        "-csv[Stimulus CSV/TSV file]:file:_files -g '*.(csv|tsv)'"
        "-display[Display index]:index:"
        "-dlp[DLP-IO8-G device]:device:_files"
        "-end-splash[End splash image]:file:_files -g '*.(jpg|jpeg|png|bmp|gif)'"
        "-fixation-color[Fixation color (R,G,B,A)]:color:"
        "-font[TTF font file]:file:_files -g '*.ttf'"
        "-font-size[Font size]:size:"
        "-fullscreen[Enable fullscreen]"
        "-height[Screen height]:pixels:"
        "-no-fixation[Disable fixation cross]"
        "-no-vsync[Disable VSync]"
        "-res[Screen resolution (e.g. 1920x1080 or Autodetect)]:res:(Autodetect 1920x1080 1440x900 1280x720 1024x768 800x600)"
        "-results-dir[Directory where result files are saved]:directory:_files -/"
        "-scale[Scale factor for stimuli]:factor:"
        "-skip-wait[Skip 'Press any key to start' message]"
        "-start-splash[Start splash image]:file:_files -g '*.(jpg|jpeg|png|bmp|gif)'"
        "-stimuli-dir[Directory containing stimuli]:directory:_files -/"
        "-subject[Subject ID]:id:"
        "-text-color[Text color (R,G,B,A)]:color:"
        "-tsv[Stimulus TSV file (alias for -csv)]:file:_files -g '*.(csv|tsv)'"
        "-version[Print version info and exit]"
        "-vrr[Enable Variable Refresh Rate mode (disables VSync)]"
        "-width[Screen width]:pixels:"
    )

    _arguments -s : $opts
}

_gostim2 "$@"
