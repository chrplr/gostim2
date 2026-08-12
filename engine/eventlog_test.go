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
