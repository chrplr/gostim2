package engine

import (
	"fmt"
	"path/filepath"
	"strconv"
	"unicode/utf8"

	"gostim2/internal/version"

	"github.com/Zyko0/go-sdl3/sdl"
	"github.com/Zyko0/go-sdl3/ttf"
)


func getTargetField(focusBox int, cfg *Config, configFilePath, displayStr, fontSizeStr *string) *string {
	switch focusBox {
	case 0: return configFilePath
	case 1: return &cfg.SubjectID
	case 2: return &cfg.CSVFile
	case 3: return &cfg.StimuliDir
	case 4: return &cfg.StartSplash
	case 5: return &cfg.FontFile
	case 6: return &cfg.ResultsDir
	case 7: return &cfg.DLPDevice
	case 8: return displayStr
	case 9: return fontSizeStr
	default: return nil
	}
}

func renderText(renderer *sdl.Renderer, font *ttf.Font, text string, x, y float32, color sdl.Color) {
	if text == "" {
		return
	}
	surf, err := font.RenderTextBlended(text, color)
	if err != nil || surf == nil {
		return
	}
	defer surf.Destroy()
	tex, err := renderer.CreateTextureFromSurface(surf)
	if err != nil {
		return
	}
	defer tex.Destroy()
	r := sdl.FRect{X: x, Y: y, W: float32(surf.W), H: float32(surf.H)}
	renderer.RenderTexture(tex, nil, &r)
}

func renderInputBox(renderer *sdl.Renderer, font *ttf.Font, label, text string, x, y, w, h float32, isFocused bool, showBrowse bool) {
	black := sdl.Color{R: 0, G: 0, B: 0, A: 255}
	renderText(renderer, font, label, x, y-25, black)

	renderer.SetDrawColor(255, 255, 255, 255)
	box := sdl.FRect{X: x, Y: y, W: w, H: h}
	renderer.RenderFillRect(&box)
	if isFocused {
		renderer.SetDrawColor(0, 120, 255, 255)
	} else {
		renderer.SetDrawColor(180, 180, 180, 255)
	}
	renderer.RenderRect(&box)

	displayText := text
	if len(text) > 40 {
		displayText = "..." + text[len(text)-37:]
	}
	renderText(renderer, font, displayText, x+5, y+5, black)

	if showBrowse {
		renderer.SetDrawColor(200, 200, 200, 255)
		btn := sdl.FRect{X: x + w + 10, Y: y, W: 35, H: h}
		renderer.RenderFillRect(&btn)
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.RenderRect(&btn)
		renderText(renderer, font, "...", x+w+20, y+5, black)
	}
}

func renderCheckbox(renderer *sdl.Renderer, font *ttf.Font, label string, x, y float32, checked bool) {
	black := sdl.Color{R: 0, G: 0, B: 0, A: 255}
	const checkSize = 20

	renderer.SetDrawColor(255, 255, 255, 255)
	box := sdl.FRect{X: x, Y: y, W: checkSize, H: checkSize}
	renderer.RenderFillRect(&box)
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.RenderRect(&box)
	if checked {
		mark := sdl.FRect{X: x + 4, Y: y + 4, W: checkSize - 8, H: checkSize - 8}
		renderer.SetDrawColor(0, 150, 0, 255)
		renderer.RenderFillRect(&mark)
	}
	renderText(renderer, font, label, x+30, y, black)
}

func renderCenteredText(renderer *sdl.Renderer, font *ttf.Font, text string, rect sdl.FRect, color sdl.Color) {
	if text == "" {
		return
	}
	tw, th, err := font.StringSize(text)
	if err != nil {
		return
	}
	x := rect.X + (rect.W-float32(tw))/2
	y := rect.Y + (rect.H-float32(th))/2
	renderText(renderer, font, text, x, y, color)
}

// wrapText splits text into lines of at most maxChars characters, breaking at word boundaries.
func wrapText(text string, maxChars int) []string {
	words := []string{}
	current := ""
	for _, ch := range text {
		if ch == ' ' {
			words = append(words, current)
			current = ""
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		words = append(words, current)
	}

	var lines []string
	line := ""
	for _, word := range words {
		if line == "" {
			line = word
		} else if len(line)+1+len(word) <= maxChars {
			line += " " + word
		} else {
			lines = append(lines, line)
			line = word
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}

// ShowAboutDialog opens a modal window with program version and authorship information.
func ShowAboutDialog() {
	const (
		winW    = 540
		winH    = 300
		padding = float32(30)
		btnW    = float32(100)
		btnH    = float32(40)
		lineH   = float32(26)
	)

	window, renderer, err := sdl.CreateWindowAndRenderer("About Gostim2", winW, winH, sdl.WINDOW_ALWAYS_ON_TOP)
	if err != nil {
		return
	}
	defer window.Destroy()
	defer renderer.Destroy()

	fontPath := GetDefaultFontPath()
	if fontPath == "" {
		return
	}
	font, err := ttf.OpenFont(fontPath, 18)
	if err != nil {
		return
	}
	defer font.Close()

	boldFont, err := ttf.OpenFont(fontPath, 20)
	if err != nil {
		boldFont = font
	} else {
		boldFont.SetStyle(ttf.FontStyleFlags(1)) // TTF_STYLE_BOLD
		defer boldFont.Close()
	}

	lines := []string{
		"Gostim2 " + version.Version,
		"",
		"\u00a9 2026 Christophe Pallier <christophe@pallier.org>",
		"License: GNU General Public License v3",
		"Source:  https://github.com/chrplr/gostim2",
	}

	closeBtn := sdl.FRect{X: float32(winW)/2 - btnW/2, Y: float32(winH) - padding - btnH, W: btnW, H: btnH}
	white := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	dark := sdl.Color{R: 30, G: 30, B: 30, A: 255}
	blue := sdl.Color{R: 0, G: 80, B: 180, A: 255}

	for {
		var e sdl.Event
		for sdl.PollEvent(&e) {
			switch e.Type {
			case sdl.EVENT_QUIT:
				return
			case sdl.EVENT_KEY_DOWN:
				ke := e.KeyboardEvent()
				if ke.Key == sdl.K_ESCAPE || ke.Key == sdl.K_RETURN {
					return
				}
			case sdl.EVENT_MOUSE_BUTTON_DOWN:
				me := e.MouseButtonEvent()
				if me.X >= closeBtn.X && me.X <= closeBtn.X+closeBtn.W &&
					me.Y >= closeBtn.Y && me.Y <= closeBtn.Y+closeBtn.H {
					return
				}
			}
		}

		renderer.SetDrawColor(245, 245, 255, 255)
		renderer.Clear()

		y := padding
		for i, line := range lines {
			if line == "" {
				y += lineH / 2
				continue
			}
			f := font
			c := dark
			if i == 0 {
				f = boldFont
				c = blue
			}
			renderText(renderer, f, line, padding, y, c)
			y += lineH
		}

		renderer.SetDrawColor(0, 100, 200, 255)
		renderer.RenderFillRect(&closeBtn)
		renderCenteredText(renderer, font, "Close", closeBtn, white)

		renderer.Present()
		sdl.Delay(10)
	}
}

// ShowWarningDialog opens a modal window displaying a warning message with an OK button.
func ShowWarningDialog(message string) {
	const (
		winW    = 520
		winH    = 200
		padding = float32(30)
		btnW    = float32(80)
		btnH    = float32(40)
	)

	window, renderer, err := sdl.CreateWindowAndRenderer("Warning", winW, winH, sdl.WINDOW_ALWAYS_ON_TOP)
	if err != nil {
		fmt.Printf("Warning: %v\n", message)
		return
	}
	defer window.Destroy()
	defer renderer.Destroy()

	fontPath := GetDefaultFontPath()
	if fontPath == "" {
		fmt.Printf("Warning: %v\n", message)
		return
	}
	font, err := ttf.OpenFont(fontPath, 18)
	if err != nil {
		fmt.Printf("Warning: %v\n", message)
		return
	}
	defer font.Close()

	lines := wrapText(message, 55)
	okBtn := sdl.FRect{X: float32(winW)/2 - btnW/2, Y: float32(winH) - padding - btnH, W: btnW, H: btnH}
	white := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	amber := sdl.Color{R: 180, G: 100, B: 0, A: 255}

	for {
		var e sdl.Event
		for sdl.PollEvent(&e) {
			switch e.Type {
			case sdl.EVENT_QUIT:
				return
			case sdl.EVENT_KEY_DOWN:
				ke := e.KeyboardEvent()
				if ke.Key == sdl.K_ESCAPE || ke.Key == sdl.K_RETURN {
					return
				}
			case sdl.EVENT_MOUSE_BUTTON_DOWN:
				me := e.MouseButtonEvent()
				if me.X >= okBtn.X && me.X <= okBtn.X+okBtn.W &&
					me.Y >= okBtn.Y && me.Y <= okBtn.Y+okBtn.H {
					return
				}
			}
		}

		renderer.SetDrawColor(255, 248, 220, 255) // cornsilk background
		renderer.Clear()

		renderText(renderer, font, "Warning", padding, padding, amber)
		for i, line := range lines {
			renderText(renderer, font, line, padding, padding+30+float32(i)*26, sdl.Color{R: 30, G: 30, B: 30, A: 255})
		}

		renderer.SetDrawColor(amber.R, amber.G, amber.B, amber.A)
		renderer.RenderFillRect(&okBtn)
		renderCenteredText(renderer, font, "OK", okBtn, white)

		renderer.Present()
		sdl.Delay(10)
	}
}

// ShowErrorDialog opens a modal-style window displaying the error message with a Quit button.
func ShowErrorDialog(message string) {
	const (
		winW    = 620
		winH    = 280
		padding = float32(30)
		btnW    = float32(100)
		btnH    = float32(40)
	)

	window, renderer, err := sdl.CreateWindowAndRenderer("Error", winW, winH, sdl.WINDOW_ALWAYS_ON_TOP)
	if err != nil {
		fmt.Printf("Error (could not open dialog): %v\n", message)
		return
	}
	defer window.Destroy()
	defer renderer.Destroy()

	fontPath := GetDefaultFontPath()
	if fontPath == "" {
		fmt.Printf("Error: %v\n", message)
		return
	}
	font, err := ttf.OpenFont(fontPath, 18)
	if err != nil {
		fmt.Printf("Error: %v\n", message)
		return
	}
	defer font.Close()

	lines := wrapText(message, 60)
	quitBtn := sdl.FRect{X: float32(winW)/2 - btnW/2, Y: float32(winH) - padding - btnH, W: btnW, H: btnH}
	white := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	red := sdl.Color{R: 180, G: 0, B: 0, A: 255}

	for {
		var e sdl.Event
		for sdl.PollEvent(&e) {
			switch e.Type {
			case sdl.EVENT_QUIT:
				return
			case sdl.EVENT_KEY_DOWN:
				ke := e.KeyboardEvent()
				if ke.Key == sdl.K_ESCAPE || ke.Key == sdl.K_RETURN {
					return
				}
			case sdl.EVENT_MOUSE_BUTTON_DOWN:
				me := e.MouseButtonEvent()
				if me.X >= quitBtn.X && me.X <= quitBtn.X+quitBtn.W &&
					me.Y >= quitBtn.Y && me.Y <= quitBtn.Y+quitBtn.H {
					return
				}
			}
		}

		renderer.SetDrawColor(255, 240, 240, 255)
		renderer.Clear()

		// Title
		renderText(renderer, font, "Error", padding, padding, sdl.Color{R: 180, G: 0, B: 0, A: 255})

		// Message lines
		for i, line := range lines {
			renderText(renderer, font, line, padding, padding+30+float32(i)*26, sdl.Color{R: 30, G: 30, B: 30, A: 255})
		}

		// Quit button
		renderer.SetDrawColor(red.R, red.G, red.B, red.A)
		renderer.RenderFillRect(&quitBtn)
		renderCenteredText(renderer, font, "Quit", quitBtn, white)

		renderer.Present()
		sdl.Delay(10)
	}
}

func RunGuiSetup(cfg *Config) bool {
	window, renderer, err := sdl.CreateWindowAndRenderer("Gostim2", 800, 800, 0)
	if err != nil {
		fmt.Printf("CreateWindowAndRenderer Error: %v\n", err)
		return false
	}
	defer window.Destroy()
	defer renderer.Destroy()

	fontPath := GetDefaultFontPath()
	if fontPath == "" {
		fmt.Println("Error: No default font found for GUI setup")
		return false
	}
	guiFont, err := ttf.OpenFont(fontPath, 18)
	if err != nil {
		fmt.Printf("Failed to load GUI font: %v\n", err)
		return false
	}
	defer guiFont.Close()

	btnFont, err := ttf.OpenFont(fontPath, 18)
	if err != nil {
		btnFont = guiFont // fallback
	} else {
		btnFont.SetStyle(ttf.FontStyleFlags(1)) // TTF_STYLE_BOLD
		defer btnFont.Close()
	}

	configFilePath := ""
	displayStr := strconv.Itoa(cfg.DisplayIndex)
	fontSizeStr := strconv.Itoa(cfg.FontSize)

	focusBox := -1

	type ResOption struct {
		W, H  int
		Label string
	}
	resOptions := []ResOption{
		{800, 600, "800x600 (SVGA)"},
		{1024, 768, "1024x768 (XGA)"},
		{1366, 1024, "1366x1024 (SXGA-)"},
		{1920, 1080, "1920x1080 (FHD)"},
		{2560, 1440, "2560x1440 (QHD)"},
		{3840, 2160, "3840x2160 (4K UHD)"},
	}
	selectedRes := 3
	if cfg.AutodetectRes {
		selectedRes = -1
	} else {
		for i, res := range resOptions {
			if cfg.ScreenWidth == res.W && cfg.ScreenHeight == res.H {
				selectedRes = i
				break
			}
		}
	}

	window.StartTextInput()
	defer window.StopTextInput()

	const (
		C1X           = 30
		C2X           = 450
		BoxW          = 350
		BoxH          = 30
		RowSpacing    = 60
		BrowseX       = 390
		BrowseW       = 35
		ResStartYS    = 240
		CheckSize     = 20
		CheckSpacing  = 30
		OptionsY      = 460
		StartBtnY     = 720
	)

	for {
		var e sdl.Event
		for sdl.PollEvent(&e) {
			switch e.Type {
			case sdl.EVENT_QUIT:
				return false
			case sdl.EVENT_MOUSE_BUTTON_DOWN:
				me := e.MouseButtonEvent()
				mx, my := me.X, me.Y

				focusBox = -1
				for i := 0; i < 7; i++ {
					by := float32(40 + i*RowSpacing)
					if mx >= C1X && mx <= C1X+BoxW && my >= by && my <= by+BoxH {
						focusBox = i
						break
					}
				}
				if focusBox == -1 {
					for i := 0; i < 3; i++ {
						by := float32(40 + i*RowSpacing)
						if mx >= C2X && mx <= C2X+BoxW && my >= by && my <= by+BoxH {
							focusBox = 7 + i
							break
						}
					}
				}

				if mx >= 20 && mx <= 100 && my >= StartBtnY && my <= StartBtnY+40 {
					ShowAboutDialog()
				}

				if mx >= BrowseX && mx <= BrowseX+BrowseW {
					for i := 0; i < 7; i++ {
						by := float32(40 + i*RowSpacing)
						if my >= by && my <= by+BoxH {
							switch i {
							case 0:
								filters0 := []sdl.DialogFileFilter{{Name: "Config Files", Pattern: "toml"}}
								cb0 := sdl.NewDialogFileCallback(func(fileList []string, filter int32) {
									if len(fileList) > 0 {
										configFilePath = fileList[0]
										cfg.LastDir = filepath.Dir(configFilePath)
										if err := cfg.LoadFromFile(configFilePath); err == nil {
											displayStr = strconv.Itoa(cfg.DisplayIndex)
											fontSizeStr = strconv.Itoa(cfg.FontSize)
											if cfg.AutodetectRes {
												selectedRes = -1
											} else {
												selectedRes = 3
												for j, res := range resOptions {
													if cfg.ScreenWidth == res.W && cfg.ScreenHeight == res.H {
														selectedRes = j
														break
													}
												}
											}
										}
									}
								})
								sdl.ShowOpenFileDialog(cb0, window, filters0, cfg.LastDir, false)
							case 2:
								filters := []sdl.DialogFileFilter{{Name: "Experiment Files (CSV/TSV)", Pattern: "csv;tsv"}}
								cb := sdl.NewDialogFileCallback(func(fileList []string, filter int32) {
									if len(fileList) > 0 {
										cfg.CSVFile = fileList[0]
										cfg.LastDir = filepath.Dir(cfg.CSVFile)
									}
								})
								sdl.ShowOpenFileDialog(cb, window, filters, cfg.LastDir, false)
							case 3:
								cb := sdl.NewDialogFileCallback(func(fileList []string, filter int32) {
									if len(fileList) > 0 {
										cfg.StimuliDir = fileList[0]
										cfg.LastDir = cfg.StimuliDir
									}
								})
								sdl.ShowOpenFolderDialog(cb, window, cfg.LastDir, false)
							case 4:
								filters := []sdl.DialogFileFilter{{Name: "Images", Pattern: "png;jpg;jpeg;bmp"}}
								cb := sdl.NewDialogFileCallback(func(fileList []string, filter int32) {
									if len(fileList) > 0 {
										cfg.StartSplash = fileList[0]
										cfg.LastDir = filepath.Dir(cfg.StartSplash)
									}
								})
								sdl.ShowOpenFileDialog(cb, window, filters, cfg.LastDir, false)
							case 5:
								filters := []sdl.DialogFileFilter{{Name: "TTF Fonts", Pattern: "ttf;ttc"}}
								cb := sdl.NewDialogFileCallback(func(fileList []string, filter int32) {
									if len(fileList) > 0 {
										cfg.FontFile = fileList[0]
										cfg.LastDir = filepath.Dir(cfg.FontFile)
									}
								})
								sdl.ShowOpenFileDialog(cb, window, filters, cfg.LastDir, false)
							case 6:
								cb := sdl.NewDialogFileCallback(func(fileList []string, filter int32) {
									if len(fileList) > 0 {
										cfg.ResultsDir = fileList[0]
										cfg.LastDir = cfg.ResultsDir
									}
								})
								sdl.ShowOpenFolderDialog(cb, window, cfg.LastDir, false)
							}
						}
					}
				}

				for i := range resOptions {
					ry := float32(ResStartYS + i*CheckSpacing)
					if mx >= C2X && mx <= C2X+200 && my >= ry && my <= ry+CheckSize {
						selectedRes = i
						cfg.AutodetectRes = false
					}
				}
				// Autodetect checkbox
				ryAuto := float32(ResStartYS + len(resOptions)*CheckSpacing)
				if mx >= C2X && mx <= C2X+300 && my >= ryAuto && my <= ryAuto+CheckSize {
					cfg.AutodetectRes = !cfg.AutodetectRes
					if cfg.AutodetectRes {
						selectedRes = -1
					} else {
						selectedRes = 3 // Default back to 1080p if unchecking autodetect
					}
				}

				if mx >= C2X && mx <= C2X+200 && my >= OptionsY && my <= OptionsY+CheckSize {
					cfg.UseFixation = !cfg.UseFixation
				}
				if mx >= C2X && mx <= C2X+200 && my >= OptionsY+CheckSpacing && my <= OptionsY+CheckSpacing+CheckSize {
					cfg.Fullscreen = !cfg.Fullscreen
				}
				if mx >= C2X && mx <= C2X+200 && my >= OptionsY+2*CheckSpacing && my <= OptionsY+2*CheckSpacing+CheckSize {
					cfg.SkipWait = !cfg.SkipWait
				}
				if mx >= C2X && mx <= C2X+200 && my >= OptionsY+3*CheckSpacing && my <= OptionsY+3*CheckSpacing+CheckSize {
					cfg.VRR = !cfg.VRR
				}

				if mx >= 350 && mx <= 450 && my >= StartBtnY && my <= StartBtnY+40 {
					if cfg.CSVFile == "" {
						ShowWarningDialog("Select a config file or an Experiment CSV file")
					} else {
						if selectedRes >= 0 && selectedRes < len(resOptions) {
							cfg.ScreenWidth = resOptions[selectedRes].W
							cfg.ScreenHeight = resOptions[selectedRes].H
						}
						if v, err := strconv.Atoi(displayStr); err == nil {
							cfg.DisplayIndex = v
						}
						if v, err := strconv.Atoi(fontSizeStr); err == nil {
							cfg.FontSize = v
						}
						cfg.SaveCache()
						return true
					}
				}

				if mx >= 690 && mx <= 790 && my >= StartBtnY && my <= StartBtnY+40 {
					return false
				}
			case sdl.EVENT_TEXT_INPUT:
				ti := e.TextInputEvent()
				if target := getTargetField(focusBox, cfg, &configFilePath, &displayStr, &fontSizeStr); target != nil {
					*target += ti.Text
				}
			case sdl.EVENT_KEY_DOWN:
				ke := e.KeyboardEvent()
				if focusBox != -1 && ke.Key == sdl.K_BACKSPACE {
					if target := getTargetField(focusBox, cfg, &configFilePath, &displayStr, &fontSizeStr); target != nil {
						if len(*target) > 0 {
							_, size := utf8.DecodeLastRuneInString(*target)
							*target = (*target)[:len(*target)-size]
						}
					}
				}
			}
		}

		renderer.SetDrawColor(240, 240, 240, 255)
		renderer.Clear()
		black := sdl.Color{R: 0, G: 0, B: 0, A: 255}

		col1Labels := []string{"Config File:", "Subject ID:", "Experiment CSV:", "Stimuli Directory:", "Start Splash Image:", "TTF Font File:", "Results Directory:"}
		col1ShowBrowse := []bool{true, false, true, true, true, true, true}
		for i, label := range col1Labels {
			text := ""
			switch i {
			case 0: text = configFilePath
			case 1: text = cfg.SubjectID
			case 2: text = cfg.CSVFile
			case 3: text = cfg.StimuliDir
			case 4: text = cfg.StartSplash
			case 5: text = cfg.FontFile
			case 6: text = cfg.ResultsDir
			}
			renderInputBox(renderer, guiFont, label, text, C1X, float32(40+i*RowSpacing), BoxW, BoxH, focusBox == i, col1ShowBrowse[i])
		}

		col2Labels := []string{"DLP Device:", "Display Index:", "Font Size:"}
		for i, label := range col2Labels {
			text := ""
			switch i {
			case 0: text = cfg.DLPDevice
			case 1: text = displayStr
			case 2: text = fontSizeStr
			}
			renderInputBox(renderer, guiFont, label, text, C2X, float32(40+i*RowSpacing), BoxW, BoxH, focusBox == 7+i, false)
		}

		renderText(renderer, guiFont, "Resolution:", C2X, ResStartYS-25, black)
		for i, opt := range resOptions {
			renderCheckbox(renderer, guiFont, opt.Label, C2X, float32(ResStartYS+i*CheckSpacing), selectedRes == i)
		}
		renderCheckbox(renderer, guiFont, "Autodetect (Exclusive Fullscreen)", C2X, float32(ResStartYS+len(resOptions)*CheckSpacing), cfg.AutodetectRes)

		renderCheckbox(renderer, guiFont, "Show fixation cross", C2X, OptionsY, cfg.UseFixation)
		renderCheckbox(renderer, guiFont, "Fullscreen mode", C2X, OptionsY+CheckSpacing, cfg.Fullscreen)
		renderCheckbox(renderer, guiFont, "Skip 'Press any key' screen", C2X, OptionsY+2*CheckSpacing, cfg.SkipWait)
		renderCheckbox(renderer, guiFont, "Variable Refresh Rate (VRR)", C2X, OptionsY+3*CheckSpacing, cfg.VRR)

		renderer.SetDrawColor(0, 120, 255, 255)
		aboutBtn := sdl.FRect{X: 20, Y: StartBtnY, W: 80, H: 40}
		renderer.RenderFillRect(&aboutBtn)
		white := sdl.Color{R: 255, G: 255, B: 255, A: 255}
		renderCenteredText(renderer, btnFont, "ABOUT", aboutBtn, white)

		renderer.SetDrawColor(0, 150, 0, 255)
		startBtn := sdl.FRect{X: 350, Y: StartBtnY, W: 100, H: 40}
		renderer.RenderFillRect(&startBtn)
		renderCenteredText(renderer, btnFont, "START", startBtn, white)

		renderer.SetDrawColor(180, 0, 0, 255)
		quitBtn := sdl.FRect{X: 690, Y: StartBtnY, W: 100, H: 40}
		renderer.RenderFillRect(&quitBtn)
		renderCenteredText(renderer, btnFont, "QUIT", quitBtn, white)

		renderer.Present()
		sdl.Delay(10)
	}
}
