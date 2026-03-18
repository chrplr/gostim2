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
