package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadStreamExperiment(t *testing.T) {
	content := `onset_time	duration	type	stimuli
1000	500	IMAGE_STREAM	img1.png:200:50~img2.png:300:400
2000	200	TEXT_STREAM	txt1:100:20~txt2
3000	500	SOUND_STREAM	snd1:100:10~snd2
`
	dir := t.TempDir()
	path := filepath.Join(dir, "experiment_stream_test.tsv")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	exp, err := LoadExperiment(path)
	if err != nil {
		t.Fatalf("Failed to load experiment: %v", err)
	}

	if len(exp.Stimuli) != 3 {
		t.Fatalf("Expected 3 stimuli, got %d", len(exp.Stimuli))
	}

	// Test IMAGE_STREAM
	s0 := exp.Stimuli[0]
	if s0.Type != StimImageStream {
		t.Errorf("Expected IMAGE_STREAM, got %v", s0.Type)
	}
	if len(s0.FilePaths) != 2 {
		t.Errorf("Expected 2 file paths, got %d", len(s0.FilePaths))
	}
	if s0.FrameDurations[0] != 200 || s0.FrameGaps[0] != 50 {
		t.Errorf("Expected frame 0 duration 200, gap 50, got %d, %d", s0.FrameDurations[0], s0.FrameGaps[0])
	}
	if s0.FrameDurations[1] != 300 || s0.FrameGaps[1] != 400 {
		t.Errorf("Expected frame 1 duration 300, gap 400, got %d, %d", s0.FrameDurations[1], s0.FrameGaps[1])
	}

	// Test TEXT_STREAM
	s1 := exp.Stimuli[1]
	if s1.Type != StimTextStream {
		t.Errorf("Expected TEXT_STREAM, got %v", s1.Type)
	}
	if s1.FrameDurations[0] != 100 || s1.FrameGaps[0] != 20 {
		t.Errorf("Expected frame 0 duration 100, gap 20, got %d, %d", s1.FrameDurations[0], s1.FrameGaps[0])
	}
	if s1.FrameDurations[1] != 200 || s1.FrameGaps[1] != 0 {
		t.Errorf("Expected frame 1 duration 200 (explicit), gap 0, got %d, %d", s1.FrameDurations[1], s1.FrameGaps[1])
	}

	// Test SOUND_STREAM
	s2 := exp.Stimuli[2]
	if s2.Type != StimSoundStream {
		t.Errorf("Expected SOUND_STREAM, got %v", s2.Type)
	}
	if s2.FrameDurations[0] != 100 || s2.FrameGaps[0] != 10 {
		t.Errorf("Expected frame 0 duration 100, gap 10, got %d, %d", s2.FrameDurations[0], s2.FrameGaps[0])
	}
	if s2.FrameDurations[1] != 500 || s2.FrameGaps[1] != 0 {
		t.Errorf("Expected frame 1 duration 500 (default), gap 0, got %d, %d", s2.FrameDurations[1], s2.FrameGaps[1])
	}
}
