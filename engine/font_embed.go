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

package engine

import (
	_ "embed"

	"github.com/Zyko0/go-sdl3/sdl"
	"github.com/Zyko0/go-sdl3/ttf"
)

//go:embed Inconsolata.ttf
var inconsolataFont []byte

// OpenDefaultFont opens the embedded Inconsolata font at the given point size.
// If a user font path is set in cfg, that is tried first and this is the fallback.
func OpenDefaultFont(ptSize float32) (*ttf.Font, error) {
	stream, err := sdl.IOFromConstMem(inconsolataFont)
	if err != nil {
		return nil, err
	}
	return ttf.OpenFontIO(stream, true, ptSize)
}
