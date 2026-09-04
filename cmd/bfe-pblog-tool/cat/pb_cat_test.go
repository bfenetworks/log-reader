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

package cat

import "testing"

func TestPblogCat(t *testing.T) {
	var allRecords []string
	err := PblogCat("test_data/pb_access_1.log", func(res []string) {
		allRecords = append(allRecords, res...)
	})
	if err != nil {
		t.Fatalf("PblogCat failed: %v", err)
	}
	if len(allRecords) != 9 {
		t.Errorf("expected 9 records, got %d", len(allRecords))
	}
}
