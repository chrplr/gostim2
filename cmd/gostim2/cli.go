package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"gostim2/engine"
	"gostim2/internal/version"

	"github.com/Zyko0/go-sdl3/sdl"
	"github.com/Zyko0/go-sdl3/ttf"
	"github.com/Zyko0/go-sdl3/bin/binimg"
	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/Zyko0/go-sdl3/bin/binttf"
)

func init() {
	runtime.LockOSThread()
}

// prescanConfig extracts the value of -config / --config from os.Args before
// flag.Parse() runs, so we can load defaults from that file first.
func prescanConfig(args []string) string {
	for i, arg := range args {
		for _, prefix := range []string{"-config=", "--config="} {
			if strings.HasPrefix(arg, prefix) {
				return strings.TrimPrefix(arg, prefix)
			}
		}
		if (arg == "-config" || arg == "--config") && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

func main() {
	// Load config file early so its values can serve as flag defaults.
	cfg := engine.DefaultConfig()
	if path := prescanConfig(os.Args[1:]); path != "" {
		if err := cfg.LoadFromFile(path); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config file %s: %v\n", path, err)
			os.Exit(1)
		}
	}

	showVersion := flag.Bool("version", false, "Print version info and exit")
	configFile  := flag.String("config", "", "Load parameters from a TOML config file (e.g. gostim2_config.toml)")
	csvFile     := flag.String("csv", cfg.CSVFile, "Stimulus CSV/TSV file")
	tsvFile     := flag.String("tsv", "", "Stimulus TSV file (alias for -csv)")
	subjectID   := flag.String("subject", cfg.SubjectID, "Subject ID")
	resultsDir  := flag.String("results-dir", cfg.ResultsDir, "Directory where result files are saved")
	stimuliDir  := flag.String("stimuli-dir", cfg.StimuliDir, "Directory containing stimuli")
	assetsDir   := flag.String("assets", "", "Directory containing stimuli (alias for -stimuli-dir)")
	startSplash := flag.String("start-splash", cfg.StartSplash, "Start splash image")
	endSplash   := flag.String("end-splash", cfg.EndSplash, "End splash image")
	fontFile    := flag.String("font", cfg.FontFile, "TTF font file")
	fontSize    := flag.Int("font-size", cfg.FontSize, "Font size")
	dlpDevice   := flag.String("dlp", cfg.DLPDevice, "DLP-IO8-G device")
	res         := flag.String("res", "", "Screen resolution (e.g. 1920x1080 or Autodetect)")
	screenW     := flag.Int("width", cfg.ScreenWidth, "Screen width")
	screenH     := flag.Int("height", cfg.ScreenHeight, "Screen height")
	displayIdx  := flag.Int("display", cfg.DisplayIndex, "Display index")
	scaleFactor := flag.Float64("scale", float64(cfg.ScaleFactor), "Scale factor for stimuli")
	noVSync     := flag.Bool("no-vsync", false, "Disable VSync")
	noFixation      := flag.Bool("no-fixation", cfg.FixationMode == 0, "Never show fixation cross")
	fixationAlways  := flag.Bool("fixation-always", cfg.FixationMode == 2, "Show fixation cross superimposed on stimuli")
	fullscreen        := flag.Bool("fullscreen", cfg.WindowMode == 2, "Enable exclusive fullscreen")
	fullscreenDesktop := flag.Bool("fullscreen-desktop", cfg.WindowMode == 1, "Enable fullscreen desktop (borderless)")
	vrr         := flag.Bool("vrr", cfg.VRR, "Enable Variable Refresh Rate mode (disables VSync)")
	bgColorStr  := flag.String("bg-color", fmt.Sprintf("%d,%d,%d,%d", cfg.BGColor.R, cfg.BGColor.G, cfg.BGColor.B, cfg.BGColor.A), "Background color (R,G,B,A)")
	textColorStr := flag.String("text-color", fmt.Sprintf("%d,%d,%d,%d", cfg.TextColor.R, cfg.TextColor.G, cfg.TextColor.B, cfg.TextColor.A), "Text color (R,G,B,A)")
	fixColorStr  := flag.String("fixation-color", fmt.Sprintf("%d,%d,%d,%d", cfg.FixationColor.R, cfg.FixationColor.G, cfg.FixationColor.B, cfg.FixationColor.A), "Fixation color (R,G,B,A)")
	skipWait    := flag.Bool("skip-wait", cfg.SkipWait, "Skip 'Press any key to start' message")

	flag.Parse()

	if *showVersion {
		fmt.Print(version.Info())
		os.Exit(0)
	}

	// If -config was given but not pre-scanned (shouldn't happen), load now.
	// This also handles the case where the user wants to see config applied.
	_ = configFile

	defer binsdl.Load().Unload()
	defer binimg.Load().Unload()
	defer binttf.Load().Unload()

	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO | sdl.INIT_EVENTS); err != nil {
		fmt.Printf("SDL_Init Error: %v\n", err)
		os.Exit(1)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		fmt.Printf("TTF_Init Error: %v\n", err)
		os.Exit(1)
	}
	defer ttf.Quit()

	// Handle Ctrl-C (SIGINT)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt, exiting...")
		sdl.PushEvent(&sdl.Event{Type: sdl.EVENT_QUIT})
	}()

	cfg.CSVFile = *csvFile
	if cfg.CSVFile == "" {
		cfg.CSVFile = *tsvFile
	}
	if cfg.CSVFile == "" && flag.NArg() > 0 {
		cfg.CSVFile = flag.Arg(0)
	}
	cfg.SubjectID = *subjectID
	cfg.ResultsDir = *resultsDir
	cfg.StimuliDir = *stimuliDir
	if cfg.StimuliDir == "" {
		cfg.StimuliDir = *assetsDir
	}
	cfg.StartSplash = *startSplash
	cfg.EndSplash = *endSplash
	cfg.FontFile = *fontFile
	cfg.FontSize = *fontSize
	cfg.DLPDevice = *dlpDevice

	if *res != "" {
		if strings.ToLower(*res) == "autodetect" {
			cfg.AutodetectRes = true
		} else {
			var w, h int
			if n, _ := fmt.Sscanf(*res, "%dx%d", &w, &h); n == 2 {
				cfg.ScreenWidth = w
				cfg.ScreenHeight = h
			} else {
				fmt.Printf("Warning: Invalid resolution format '%s', using defaults.\n", *res)
				cfg.ScreenWidth = *screenW
				cfg.ScreenHeight = *screenH
			}
		}
	} else {
		cfg.ScreenWidth = *screenW
		cfg.ScreenHeight = *screenH
	}

	cfg.DisplayIndex = *displayIdx
	cfg.ScaleFactor = float32(*scaleFactor)
	cfg.VSync = !*noVSync
	if *vrr {
		cfg.VSync = false
		cfg.VRR = true
	}
	switch {
	case *noFixation:
		cfg.FixationMode = 0
	case *fixationAlways:
		cfg.FixationMode = 2
	default:
		cfg.FixationMode = 1
	}
	switch {
	case *fullscreen:
		cfg.WindowMode = 2
	case *fullscreenDesktop:
		cfg.WindowMode = 1
	default:
		cfg.WindowMode = 0
	}
	cfg.BGColor = engine.ParseColor(*bgColorStr)
	cfg.TextColor = engine.ParseColor(*textColorStr)
	cfg.FixationColor = engine.ParseColor(*fixColorStr)
	cfg.SkipWait = *skipWait

	_, err := engine.Run(cfg)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
