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
	"testing"
)

func TestValidateExperiment(t *testing.T) {
	dir := t.TempDir()
	dummyFile := filepath.Join(dir, "image.png")
	os.WriteFile(dummyFile, []byte("fake"), 0644)

	exp := &Experiment{
		Stimuli: []Stimulus{
			{TimestampMS: 0, DurationMS: 1000, Type: StimImage, FilePaths: []string{"image.png"}},
			{TimestampMS: 1000, DurationMS: 500, Type: StimImage, FilePaths: []string{"missing.png"}},
			{TimestampMS: 500, DurationMS: 500, Type: StimText, FilePaths: []string{"Text"}},
		},
	}

	errs := ValidateExperiment(exp, dir)

	if len(errs) != 3 {
		t.Fatalf("expected 3 validation errors, got %d: %v", len(errs), errs)
	}

	expectedMissing := "stimulus 2 (at 1000ms): resource file not found: " + filepath.Join(dir, "missing.png")
	if errs[0].Error() != expectedMissing {
		t.Errorf("unexpected error 1: %v", errs[0])
	}

	if errs[1].Error() != "stimulus 3 (at 500ms) is out of order (previous at 1000ms)" {
		t.Errorf("unexpected error 2: %v", errs[1])
	}

	if errs[2].Error() != "stimulus 3 (visual, at 500ms) overlaps with previous visual stimulus 2 (ends at 1500ms)" {
		t.Errorf("unexpected error 3: %v", errs[2])
	}
}
