package engine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEventLog_Save(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "test_log.tsv")

	log := &EventLog{
		SubjectID: "sub001",
		CSVHeader: []string{"onset", "duration", "type", "stimuli"},
		Entries: []EventLogEntry{
			{
				IntendedMS:  1000,
				TimestampMS: 1005,
				Type:        "IMAGE_ONSET",
				Label:       "img1.png",
				StimulusRow: []string{"1000", "500", "image", "img1.png"},
			},
		},
		SDLVersion: "3.0.0",
		Platform:   "test-platform",
		Completed:  true,
	}

	err := log.Save(logPath)
	if err != nil {
		t.Fatalf("failed to save log: %v", err)
	}

	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read log: %v", err)
	}

	sContent := string(content)
	if !strings.Contains(sContent, "subject_id\tintended_ms\tactual_ms\ttype\tlabel\tonset\tduration\ttype\tstimuli") {
		t.Errorf("header missing or incorrect: %s", sContent)
	}
	if !strings.Contains(sContent, "sub001\t1000\t1005\tIMAGE_ONSET\timg1.png\t1000\t500\timage\timg1.png") {
		t.Errorf("entry missing or incorrect: %s", sContent)
	}
	if !strings.Contains(sContent, "# SDL Version: 3.0.0") {
		t.Errorf("metadata missing: %s", sContent)
	}
}
