// Copyright 2026 Christophe Pallier <christophe@pallier.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"gostim2/engine"
	"gostim2/internal/version"

	"github.com/Zyko0/go-sdl3/bin/binimg"
	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/Zyko0/go-sdl3/bin/binttf"
	"github.com/Zyko0/go-sdl3/sdl"
	"github.com/Zyko0/go-sdl3/ttf"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	showVersion := flag.Bool("version", false, "Print version info and exit")
	flag.Parse()

	if *showVersion {
		fmt.Print(version.Info())
		os.Exit(0)
	}

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
		// Push a quit event to the SDL queue to trigger graceful shutdown
		sdl.PushEvent(&sdl.Event{Type: sdl.EVENT_QUIT})
	}()

	cfg := engine.DefaultConfig()
	cfg.LoadCache()
	cfg.SubjectID = "" // Always start with empty subject ID field

	for {
		if engine.RunGuiSetup(cfg) {
			_, err := engine.Run(cfg)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				engine.ShowErrorDialog(err.Error())
			}
		} else {
			break
		}
	}
}
