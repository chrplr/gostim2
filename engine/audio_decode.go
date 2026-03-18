package engine

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	mp3dec "github.com/hajimehoshi/go-mp3"
	"github.com/jfreymuth/oggvorbis"
	"github.com/mewkiz/flac"
	"github.com/Zyko0/go-sdl3/sdl"
)

// loadAudioFile reads an audio file and returns its PCM data as S16LE bytes
// together with the source AudioSpec (sample rate, channels).
// Supported formats: .wav, .flac, .mp3, .ogg
func loadAudioFile(fullPath string) ([]byte, sdl.AudioSpec, error) {
	ext := strings.ToLower(filepath.Ext(fullPath))
	switch ext {
	case ".wav":
		spec := &sdl.AudioSpec{}
		data, err := sdl.LoadWAV(fullPath, spec)
		if err != nil {
			return nil, sdl.AudioSpec{}, err
		}
		return data, *spec, nil
	case ".flac":
		return loadFLAC(fullPath)
	case ".mp3":
		return loadMP3(fullPath)
	case ".ogg":
		return loadOGG(fullPath)
	default:
		return nil, sdl.AudioSpec{}, fmt.Errorf("unsupported audio format %q (supported: .wav, .flac, .mp3, .ogg)", ext)
	}
}

func loadFLAC(path string) ([]byte, sdl.AudioSpec, error) {
	stream, err := flac.Open(path)
	if err != nil {
		return nil, sdl.AudioSpec{}, fmt.Errorf("flac open: %w", err)
	}
	defer stream.Close()

	info := stream.Info
	spec := sdl.AudioSpec{
		Format:   sdl.AUDIO_S16,
		Channels: int32(info.NChannels),
		Freq:     int32(info.SampleRate),
	}

	// Shift amount to map source bit depth to 16-bit.
	// Positive = right-shift (e.g. 24-bit → shift 8), negative = left-shift (e.g. 8-bit → shift -8).
	shift := int(info.BitsPerSample) - 16

	var samples []int16
	for {
		f, err := stream.ParseNext()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, sdl.AudioSpec{}, fmt.Errorf("flac decode: %w", err)
		}
		nSamples := len(f.Subframes[0].Samples)
		for i := 0; i < nSamples; i++ {
			for ch := 0; ch < int(info.NChannels); ch++ {
				s := f.Subframes[ch].Samples[i]
				var s16 int16
				if shift > 0 {
					s16 = int16(s >> shift)
				} else if shift < 0 {
					s16 = int16(s << uint(-shift))
				} else {
					s16 = int16(s)
				}
				samples = append(samples, s16)
			}
		}
	}

	return int16sToBytes(samples), spec, nil
}

func loadMP3(path string) ([]byte, sdl.AudioSpec, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, sdl.AudioSpec{}, err
	}
	defer f.Close()

	dec, err := mp3dec.NewDecoder(f)
	if err != nil {
		return nil, sdl.AudioSpec{}, fmt.Errorf("mp3 decode: %w", err)
	}

	// go-mp3 always outputs stereo S16LE regardless of source channel count.
	spec := sdl.AudioSpec{
		Format:   sdl.AUDIO_S16,
		Channels: 2,
		Freq:     int32(dec.SampleRate()),
	}

	data, err := io.ReadAll(dec)
	if err != nil {
		return nil, sdl.AudioSpec{}, fmt.Errorf("mp3 read: %w", err)
	}

	return data, spec, nil
}

func loadOGG(path string) ([]byte, sdl.AudioSpec, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, sdl.AudioSpec{}, err
	}
	defer f.Close()

	r, err := oggvorbis.NewReader(f)
	if err != nil {
		return nil, sdl.AudioSpec{}, fmt.Errorf("ogg decode: %w", err)
	}

	spec := sdl.AudioSpec{
		Format:   sdl.AUDIO_S16,
		Channels: int32(r.Channels()),
		Freq:     int32(r.SampleRate()),
	}

	buf := make([]float32, 4096)
	var samples []int16
	for {
		n, err := r.Read(buf)
		for i := 0; i < n; i++ {
			s := buf[i]
			if s > 1.0 {
				s = 1.0
			} else if s < -1.0 {
				s = -1.0
			}
			samples = append(samples, int16(s*32767))
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, sdl.AudioSpec{}, fmt.Errorf("ogg read: %w", err)
		}
	}

	return int16sToBytes(samples), spec, nil
}

func int16sToBytes(samples []int16) []byte {
	data := make([]byte, len(samples)*2)
	for i, s := range samples {
		binary.LittleEndian.PutUint16(data[i*2:], uint16(s))
	}
	return data
}
