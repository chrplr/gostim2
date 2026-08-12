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

import "testing"

func TestStimulus_TotalDuration(t *testing.T) {
	tests := []struct {
		name     string
		stim     Stimulus
		expected uint64
	}{
		{
			name: "Simple Image",
			stim: Stimulus{
				Type:       StimImage,
				DurationMS: 1000,
			},
			expected: 1000,
		},
		{
			name: "Image Stream with gaps",
			stim: Stimulus{
				Type:           StimImageStream,
				FrameDurations: []uint64{200, 300},
				FrameGaps:      []uint64{50, 400},
			},
			expected: 950,
		},
		{
			name: "Text Stream",
			stim: Stimulus{
				Type:           StimTextStream,
				FrameDurations: []uint64{100, 200},
				FrameGaps:      []uint64{20, 0},
			},
			expected: 320,
		},
		{
			name: "Sound Stream",
			stim: Stimulus{
				Type:           StimSoundStream,
				FrameDurations: []uint64{500},
				FrameGaps:      []uint64{100},
			},
			expected: 600,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.stim.TotalDuration(); got != tt.expected {
				t.Errorf("Stimulus.TotalDuration() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
