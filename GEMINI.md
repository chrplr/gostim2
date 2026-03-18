# GEMINI.md - Technical Context & Architecture

This document provides technical context, architectural decisions, and hidden dependencies for Gostim2, intended for developers and AI assistants.

## 🏗️ Architecture Overview

Gostim2 is a high-precision stimulus delivery system. Unlike many experimental platforms, it follows a **Fixed Schedule** model: the entire experiment is deterministic and known before execution.

### Core Components
- **`engine/run.go`**: The main execution loop. It orchestrates the transition from setup to experiment.
- **`engine/experiment.go`**: Contains the `experimentState` and the primary `RunExperiment` loop.
- **`engine/csv_parser.go`**: Handles the flexible parsing of CSV/TSV files, including the complex "Stream" syntax (`~` and `:` delimiters).
- **`engine/audio_mixer.go`**: A custom software mixer designed for low-latency playback without thread-blocking issues common in higher-level libraries.

## ⏱️ Timing & Precision

Precision is the primary technical goal of Gostim2.

### The Timing Loop
- **VSYNC Synchronization**: When `-vsync` is enabled, the loop synchronizes with the monitor's refresh rate.
- **Predictive Look-ahead (`laNS`)**: To ensure frame-perfect onsets, the engine uses a look-ahead (typically half a frame duration). If the next stimulus onset falls within this window, it is prepared for the *next* VBLANK.
- **VRR Mode**: When `-vrr` is enabled, VSYNC is disabled, and the engine uses a busy-wait loop (`vrrWaitNS` = 2ms) to hit timestamps with sub-millisecond accuracy.

### Resolution & Scaling
- **Autodetect Resolution**: Sets `cfg.ScreenWidth/ScreenHeight` from `SDL_GetDesktopDisplayMode()` (screen coordinates, not physical pixels). This is independent of the window mode.
- **Scale Factor**: All visual stimuli are scaled by `cfg.ScaleFactor`. This allows running the same experiment on different monitors while maintaining relative stimulus size.

### Window Modes (`cfg.WindowMode`)
Three modes are available, stored as an integer in config (`window_mode`):
- **0 — Windowed**: A standard resizable window at the configured resolution.
- **1 — Fullscreen Desktop** (`SDL_WINDOW_FULLSCREEN` + `SetFullscreenMode(nil)`): Borderless fullscreen at the desktop resolution. The compositor stays active; no mode switch. Safest option for multi-monitor setups.
- **2 — Fullscreen Exclusive** (`SDL_WINDOW_FULLSCREEN` + `SetFullscreenMode(&desktopMode)`): Exclusive ownership of the display. Bypasses the OS compositor, which minimises latency and prevents frame drops. Most relevant for EEG/MEG/fMRI where timing is critical. Can crash on some systems if the display mode negotiation fails.

**Important SDL3 coordinate system note**: SDL3 renders in *physical pixel space*. `DesktopDisplayMode()` returns screen coordinates (logical points), which differ from physical pixels on HiDPI/Retina displays. All centering calculations therefore query `renderer.CurrentOutputSize()` at draw time rather than using `cfg.ScreenWidth/ScreenHeight`, ensuring correct centering in all modes and display densities.

## 📦 Dependencies & Hidden Requirements

### SDL3 (The Foundation)
Gostim2 relies heavily on **SDL3**. Most of the "magic" in timing and cross-platform compatibility comes from SDL3's high-level abstractions over DRM (Linux), CoreAudio (macOS), and WASAPI/DirectX (Windows).
- **Bindings**: Uses `github.com/Zyko0/go-sdl3`.
- **CGo**: While the core logic is Go, the SDL3 bindings require CGo for linking against system libraries.

### Audio Decoding
The system supports `.wav`, `.flac`, `.mp3`, and `.ogg`.
- **Custom Loaders**: Found in `engine/audio_decode.go`.
- **Resampling**: If a sound's sample rate doesn't match the engine's target (usually 44100Hz), it is automatically resampled during the resource loading phase to prevent real-time performance hits.

### Serial Triggers (DLP-IO8-G)
- **Library**: `go.bug.st/serial`.
- **Logic**: Lines 1, 2, and 3 are hardcoded for specific trigger types:
    - **Line 1**: IMAGE/VIDEO onset.
    - **Line 2**: SOUND onset (5ms pulse).
    - **Line 3**: TEXT/BOX onset.

## 🛠️ Developer Notes

### Config Persistence
- **`gostim2_config.toml`**: The local cache file for GUI settings. 
- **`LastDir`**: A critical field for UX; it stores the directory of the last opened file to avoid repetitive navigation.

### Resource Caching
Resources are pre-loaded in `engine/resources.go` before the experiment starts. This prevents I/O jitters during execution. The `ResourceCache` deduplicates assets used multiple times in a CSV.

### Testing Strategy
- **Unit Tests**: Focus on non-GUI logic (parsing, color math, log saving).
- **Self-Contained Tests**: Use `t.TempDir()` to avoid polluting the workspace.
- **Pre-existing Failures**: `engine/stream_test.go` was previously broken by a missing file; it has been refactored to be self-contained.

## ⚠️ Known Constraints
- **No Real-time Feedback**: The architecture does not support branching or feedback based on user input.
- **Memory Usage**: Since all resources are pre-loaded, experiments with many large videos or images may require significant RAM.
