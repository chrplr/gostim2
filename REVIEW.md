# Project Review: gostim2

## Overview
`gostim2` is a high-precision multimedia stimulus delivery system written in Go, utilizing SDL3 for low-latency graphics and audio. It is designed for experimental psychology and neuroscience, where timing accuracy is paramount.

## Technical Architecture

### 1. High-Precision Timing Loop (`engine/experiment.go`)
- **Mechanism**: Uses `sdl.TicksNS()` for nanosecond-level precision.
- **Strategies**:
    - **VSync Mode**: Predictive onset look-ahead (`laNS`) to hit the next vertical refresh.
    - **VRR Mode**: Busy-wait loop for immediate onset on variable refresh rate displays.
- **Optimization**: Disables Garbage Collection (`debug.SetGCPercent(-1)`) during the experiment to prevent unpredictable latency spikes.

### 2. Audio Engine (`engine/audio.go`)
- **Implementation**: A custom software mixer that manages multiple concurrent audio slots (`MaxActiveSounds = 16`).
- **Safety**: Thread-safe access via `sync.Mutex`.
- **Latency**: Direct interaction with SDL3 audio streams using a callback mechanism.

### 3. Resource Management (`engine/resources.go`)
- **Strategy**: Pre-loads all assets (textures, sounds, fonts, videos) before the experiment starts.
- **Caching**: Implements a `ResourceCache` to avoid redundant loading and disk I/O during critical timing sections.
- **Conversion**: Automatically converts audio samples to the target format (44.1kHz, 16-bit, Stereo) during loading.

### 4. Data Ingestion & Validation (`engine/csv_parser.go`, `engine/validator.go`)
- **Flexibility**: Supports various delimiters (comma, semicolon, tab) and auto-detects them.
- **Robustness**: Comprehensive validation of stimulus order, existence of resources, and visual overlaps.
- **Features**: Supports complex RSVP (Rapid Serial Visual Presentation) streams for images, text, and sound.

## Code Quality & Best Practices

### Strengths
- **Surgical Go Usage**: Efficient use of `unsafe.Slice` for audio mixing performance.
- **Platform Agnostic**: Robust font discovery across Windows, macOS, and Linux.
- **Error Handling**: Detailed error messages for malformed CSVs and missing resources.
- **Clean Separation**: Logical split between the CLI/GUI entry points and the core engine.

### Observations
- **Concurrency**: The audio mixer uses a global mutex. While safe, high-contention scenarios (though unlikely in this use case) could be optimized with lock-free structures.
- **Global State**: The project manages state well within structs, avoiding excessive global variables.
- **Testing**: Good unit tests for the CSV parser and streams, though more integration tests for the `RunExperiment` loop could be beneficial.

## User Experience

### 1. Dual Interfaces
- **CLI (`cmd/gostim2`)**: Powerful for automated runs and headless environments (e.g., Linux TTY/DRM).
- **GUI (`cmd/gostim2-gui`)**: Friendly for configuration and quick experimentation.

### 2. Results & Metadata
- Generates detailed TSV results with extensive metadata (system info, versions, drivers), ensuring experiment reproducibility.

## Conclusion
`gostim2` is a well-engineered tool that strikes a balance between Go's high-level productivity and the low-level control required for millisecond-accurate stimulus delivery. The architecture is modular, performance-oriented, and highly suitable for its intended scientific domain.
