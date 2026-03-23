# CLAUDE.md — Project Context for Claude Code

## What is this project?

Gostim2 is a multimedia stimulus delivery system for experimental psychology / cognitive neuroscience (fMRI, MEG, EEG). It presents images, sounds, text, and video according to a **fixed, pre-defined schedule** described in a CSV/TSV file. No programming is required to build an experiment.

There are two binaries sharing the same engine:
- `gostim2` — command-line interface
- `gostim2-gui` — graphical setup dialog (SDL3-based)

## Build & test commands

```bash
make build                # build CLI + GUI for the current platform
make build-multiplatform  # build for all platforms into dist/
make install              # install to $GOPATH/bin
make test                 # run all tests
make fmt                  # gofmt
make vet                  # go vet
make run ARGS="-csv ..."  # run CLI without installing
make run-gui              # run GUI without installing
make clean                # remove binaries and dist/
```

## Documentation commands

```bash
make docs-install   # pip install -r docs/requirements.txt (mkdocs + material)
make docs-html      # build static HTML site into site/
make docs-pdf       # build gostim2-userguide.pdf via pandoc + xelatex
make docs-serve     # live-reload preview at http://127.0.0.1:8000
```

The user guide lives in `docs/UserGuide-Gostim2.md`. The mkdocs config is `mkdocs.yml`. PDF generation uses `pandoc --pdf-engine=xelatex` (no WeasyPrint dependency).

## Repository layout

```
cmd/gostim2/          CLI entry point and flag definitions
cmd/gostim2-gui/      GUI entry point
engine/               Core engine (parsing, rendering, timing, audio, video)
  config.go           Config struct, TOML load/save, color parsing
  csv_parser.go       CSV/TSV parsing, stream syntax (~, :dur:gap)
  experiment.go       Main experiment loop, event logging
  stimuli.go          Stimulus type definitions
  validator.go        CSV validation (required columns, overlaps, missing files)
  resources.go        Pre-load cache (images, sounds, fonts)
  audio.go / audio_decode.go  Custom mixer + WAV/FLAC/MP3/OGG decoding
  video.go            Video playback
  gui_setup.go        SDL3-based GUI dialog
  dlp.go              DLP-IO8-G serial trigger support
internal/version/     Version metadata injected at build time via ldflags
examples/             Three complete working experiments
docs/                 User guide (Markdown) and mkdocs requirements
scripts/              Shell completion installer
```

## Key architectural points

- **Fixed-schedule model**: no real-time branching or feedback. All stimuli are pre-loaded before the experiment starts (`engine/resources.go`).
- **Timing**: VSYNC loop with predictive look-ahead (`laNS`). VRR mode uses 2 ms busy-wait for sub-millisecond precision.
- **Window modes** (stored as int in config): 0 = windowed, 1 = fullscreen-desktop, 2 = fullscreen-exclusive.
- **SDL3 coordinate note**: rendering uses physical pixels; always query `renderer.CurrentOutputSize()` for centering, not `cfg.ScreenWidth/ScreenHeight`.
- **Audio**: custom software mixer in `engine/audio_mixer.go`; audio is resampled at load time if sample rate mismatches.
- **DLP triggers**: lines 1/2/3 are hardcoded to IMAGE/VIDEO, SOUND, and TEXT/BOX onsets respectively.

## Config persistence

The GUI writes `gostim2_config.toml` next to itself after each session. The CLI accepts `-config <path>` to load the same format. `LastDir` in the TOML is used by the GUI to remember the last opened directory.

## Testing

Tests cover non-GUI logic only (parsing, color math, validation, event logging). All tests use `t.TempDir()` and are self-contained. Run with `make test`.

## Dependencies

- `github.com/Zyko0/go-sdl3` — SDL3 bindings (requires CGo + system SDL3 libs)
- `github.com/BurntSushi/toml` — TOML config parsing
- `github.com/funatsufumiya/go-gv-video` — video decoding
- `go.bug.st/serial` — serial port for DLP-IO8-G triggers (no CGo)

SDL3 system libraries must be installed separately:
- **Linux**: `libsdl3-0 libsdl3-image-0 libsdl3-ttf-0`
- **macOS**: `brew install sdl3 sdl3_image sdl3_ttf`
- **Windows**: DLLs bundled with releases
