// Copyright (c) 2026 The BFE Authors.
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

package tail

import (
	"testing"
	"time"
)

// TestValidateTailOptions covers defaulting and validation of tail options.
func TestValidateTailOptions(t *testing.T) {
	tests := []struct {
		name     string
		opts     TailOptions
		wantOpts TailOptions
		wantErr  bool
	}{
		{
			name:    "negative N errors",
			opts:    TailOptions{N: -1},
			wantErr: true,
		},
		{
			name:     "N zero defaults to 10",
			opts:     TailOptions{N: 0},
			wantOpts: TailOptions{N: defaultTailN},
		},
		{
			name:     "N over cap is capped",
			opts:     TailOptions{N: MaxTailRecords + 100},
			wantOpts: TailOptions{N: MaxTailRecords},
		},
		{
			name:    "negative interval errors when following",
			opts:    TailOptions{Follow: true, Interval: -1 * time.Second},
			wantErr: true,
		},
		{
			name:     "follow with zero interval defaults to 500ms",
			opts:     TailOptions{Follow: true, Interval: 0},
			wantOpts: TailOptions{N: defaultTailN, Follow: true, Interval: defaultFollowInterval},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateTailOptions(tt.opts)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.wantOpts {
				t.Fatalf("validateTailOptions() = %+v, want %+v", got, tt.wantOpts)
			}
		})
	}
}

// TestPblogTail_NonFollow verifies the non-follow path emits the last N records.
func TestPblogTail_NonFollow(t *testing.T) {
	const fp = "test_data/pb_access_1.log"

	t.Run("N=3 emits exactly 3 records", func(t *testing.T) {
		var emitted []string
		err := PblogTail(fp, TailOptions{N: 3}, func(s string) error {
			emitted = append(emitted, s)
			return nil
		})
		if err != nil {
			t.Fatalf("PblogTail() error: %v", err)
		}
		if len(emitted) != 3 {
			t.Fatalf("expected 3 records, got %d", len(emitted))
		}
	})

	t.Run("N=0 defaults and emits all 9 records", func(t *testing.T) {
		var emitted []string
		err := PblogTail(fp, TailOptions{}, func(s string) error {
			emitted = append(emitted, s)
			return nil
		})
		if err != nil {
			t.Fatalf("PblogTail() error: %v", err)
		}
		if len(emitted) != 9 {
			t.Fatalf("expected 9 records, got %d", len(emitted))
		}
	})
}
