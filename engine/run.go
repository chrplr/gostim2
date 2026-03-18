package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Zyko0/go-sdl3/sdl"
	"github.com/Zyko0/go-sdl3/ttf"
)

func Run(cfg *Config) (string, error) {
	if cfg.CSVFile == "" {
		return "", fmt.Errorf("CSV file is required")
	}

	// Resolve screen resolution first (autodetect reads display info)
	if cfg.AutodetectRes {
		displays, err := sdl.GetDisplays()
		if err == nil && cfg.DisplayIndex >= 0 && cfg.DisplayIndex < len(displays) {
			if dm, err := displays[cfg.DisplayIndex].DesktopDisplayMode(); err == nil {
				cfg.ScreenWidth = int(dm.W)
				cfg.ScreenHeight = int(dm.H)
			}
		}
	}

	windowFlags := sdl.WINDOW_RESIZABLE
	if cfg.WindowMode > 0 {
		windowFlags |= sdl.WINDOW_FULLSCREEN
	}

	window, renderer, err := sdl.CreateWindowAndRenderer("gostim2 (Go)", cfg.ScreenWidth, cfg.ScreenHeight, windowFlags)
	if err != nil {
		return "", fmt.Errorf("CreateWindowAndRenderer Error: %v", err)
	}
	defer window.Destroy()
	defer renderer.Destroy()

	switch cfg.WindowMode {
	case 1: // Fullscreen desktop (borderless)
		if err := window.SetFullscreenMode(nil); err != nil {
			fmt.Printf("Warning: Failed to set fullscreen desktop mode: %v\n", err)
		}
		if (window.Flags() & sdl.WINDOW_FULLSCREEN) == 0 {
			if err := window.SetFullscreen(true); err != nil {
				fmt.Printf("Warning: Failed to set fullscreen: %v\n", err)
			}
		}
	case 2: // Fullscreen exclusive
		display := sdl.GetDisplayForWindow(window)
		if desktopMode, err := display.DesktopDisplayMode(); err == nil {
			if err := window.SetFullscreenMode(desktopMode); err != nil {
				fmt.Printf("Warning: Failed to set exclusive fullscreen mode: %v\n", err)
			}
		}
		if (window.Flags() & sdl.WINDOW_FULLSCREEN) == 0 {
			if err := window.SetFullscreen(true); err != nil {
				fmt.Printf("Warning: Failed to set fullscreen: %v\n", err)
			}
		}
	}

	if cfg.VSync {
		renderer.SetVSync(1)
	} else {
		renderer.SetVSync(0)
	}

	var font *ttf.Font
	if cfg.FontFile != "" {
		font, err = ttf.OpenFont(cfg.FontFile, float32(cfg.FontSize))
		if err != nil {
			fmt.Printf("Failed to load font: %s (%v)\n", cfg.FontFile, err)
			cfg.FontFile = "" // Clear it if it failed to load
		}
	}

	// If no font loaded yet (either none specified or loading failed), use embedded default
	if font == nil {
		font, err = OpenDefaultFont(float32(cfg.FontSize))
		if err != nil {
			fmt.Printf("Failed to load embedded default font: %v\n", err)
		}
	}
	defer func() {
		if font != nil {
			font.Close()
		}
	}()

	exp, err := LoadExperiment(cfg.CSVFile)
	if err != nil {
		return "", fmt.Errorf("failed to load experiment: %v", err)
	}

	if cfg.StimuliDir == "" {
		csvDir := filepath.Dir(cfg.CSVFile)
		for _, candidate := range []string{"stimuli", "assets"} {
			if info, err := os.Stat(filepath.Join(csvDir, candidate)); err == nil && info.IsDir() {
				cfg.StimuliDir = filepath.Join(csvDir, candidate)
				break
			}
		}
	}

	validationErrs := ValidateExperiment(exp, cfg.StimuliDir)
	if len(validationErrs) > 0 {
		return "", fmt.Errorf("experiment configuration contains %d errors (first: %v)", len(validationErrs), validationErrs[0])
	}

	if len(exp.Stimuli) > 0 {
		lastStim := exp.Stimuli[len(exp.Stimuli)-1]
		cfg.TotalDuration = lastStim.TimestampMS + lastStim.TotalDuration() + uint64(gracePeriod)
	}

	cache := NewResourceCache()
	defer cache.Destroy()

	resources, warnings, err := cache.Load(renderer, exp, font, cfg.TextColor, cfg.StimuliDir)
	if err != nil {
		return "", fmt.Errorf("failed to load resources: %v", err)
	}

	mixer := NewAudioMixer()
	spec := DefaultAudioSpec()
	cb := sdl.NewAudioStreamCallback(mixer.Callback)
	stream := sdl.AUDIO_DEVICE_DEFAULT_PLAYBACK.OpenAudioDeviceStream(&spec, cb)
	if stream == nil {
		return "", fmt.Errorf("failed to open audio stream")
	}
	defer stream.Destroy()
	stream.ResumeDevice()

	var dlp *DLPIO8G
	if cfg.DLPDevice != "" {
		dlp, err = NewDLPIO8G(cfg.DLPDevice, 9600)
		if err != nil {
			fmt.Printf("Failed to initialize DLP device: %v\n", err)
		} else {
			defer dlp.Close()
		}
	}

	subjID := cfg.SubjectID
	if subjID == "" {
		subjID = "unknown"
	}

	hostname, _ := os.Hostname()
	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("LOGNAME")
	}
	if username == "" {
		username = os.Getenv("USERNAME")
	}
	if username == "" {
		username = "unknown"
	}

	sdlVer := sdl.GetVersion()
	sdlVersionStr := fmt.Sprintf("%d.%d.%d", sdlVer/1000000, (sdlVer/1000)%1000, sdlVer%1000)

	displayModeStr := ""
	if window != nil {
		display := sdl.GetDisplayForWindow(window)
		if dm, err := display.CurrentDisplayMode(); err == nil {
			displayModeStr = fmt.Sprintf("%dx%d @ %.2fHz (Physical)", dm.W, dm.H, dm.RefreshRate)
		}
	}

	rendererName, _ := renderer.Name()

	log := &EventLog{
		SubjectID:         subjID,
		CSVHeader:         exp.Header,
		Entries:           make([]EventLogEntry, 0, len(exp.Stimuli)*4+100),
		SDLVersion:        sdlVersionStr,
		Platform:          runtime.GOOS,
		Hostname:          hostname,
		Username:          username,
		VideoDriver:       sdl.GetCurrentVideoDriver(),
		AudioDriver:       sdl.GetCurrentAudioDriver(),
		AudioSampleRate:   int(DefaultAudioSpec().Freq),
		Renderer:          rendererName,
		Warnings:          warnings,
		DisplayMode:       displayModeStr,
		LogicalResolution: fmt.Sprintf("%dx%d", cfg.ScreenWidth, cfg.ScreenHeight),
		Font:              cfg.FontFile,
		FontSize:          cfg.FontSize,
		CommandLine:       strings.Join(os.Args, " "),
	}

	if !DisplaySplash(renderer, cfg.StartSplash, cfg.ScaleFactor, cfg.BGColor) {
		return "", nil
	}

	if !cfg.SkipWait {
		if !WaitForKeyPress(renderer, font, cfg.TextColor, cfg.BGColor) {
			return "", nil
		}
	}

	log.StartTime = time.Now().Format("2006-01-02 15:04:05.000")

	success := RunExperiment(cfg, exp, resources, renderer, mixer, log, dlp, font)

	if success {
		log.Completed = true
		DisplaySplash(renderer, cfg.EndSplash, cfg.ScaleFactor, cfg.BGColor)
	}

	log.EndTime = time.Now().Format("2006-01-02 15:04:05.000")

	baseName := filepath.Base(cfg.CSVFile)
	baseName = strings.TrimSuffix(baseName, filepath.Ext(baseName))
	timestamp := time.Now().Format("20060102-150405")
	outputName := fmt.Sprintf("%s_sub-%s_%s.tsv", baseName, subjID, timestamp)

	// Determine final results directory relative to CSV if not absolute
	finalResultsDir := cfg.ResultsDir
	if !filepath.IsAbs(finalResultsDir) && cfg.CSVFile != "" {
		finalResultsDir = filepath.Join(filepath.Dir(cfg.CSVFile), finalResultsDir)
	}

	if finalResultsDir != "" {
		if err := os.MkdirAll(finalResultsDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create results directory %s: %v", finalResultsDir, err)
		}
		outputName = filepath.Join(finalResultsDir, outputName)
	} else if cfg.CSVFile != "" {
		// Fallback to CSV directory if ResultsDir is empty
		outputName = filepath.Join(filepath.Dir(cfg.CSVFile), outputName)
	}

	if err := log.Save(outputName); err != nil {
		return outputName, fmt.Errorf("failed to save event log: %v", err)
	}
	fmt.Printf("\nResults saved to %s\n", outputName)
	return outputName, nil
}
