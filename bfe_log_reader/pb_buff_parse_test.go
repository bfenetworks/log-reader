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

package bfe_log_reader

import (
	"io/ioutil"
	"testing"
)

// test of pbBuffParse(), case 1
// It's the normal situation
func Test_pbBuffParse_1(t *testing.T) {
	// read testing data from file
	data, err := ioutil.ReadFile("test_data/pb_access_1.log")
	if err != nil {
		t.Error("fail to open file for testing data")
		return
	}

	// parse b2log record from data
	records, buffer := pbBuffParse(data, nil)

	if len(records) != 9 {
		t.Errorf("len(records) should be 9, but now it's %d", len(records))
	}
	if len(buffer) != 0 {
		t.Errorf("len(buffer) should be 0, but now it's %d", len(buffer))
	}
}

// test of pbBuffParse(), case 4
// try to parse only 32 bytes
func Test_pbBuffParse_2(t *testing.T) {
	// read testing data from file
	data, err := ioutil.ReadFile("test_data/pb_access_1.log")
	if err != nil {
		t.Error("fail to open file for testing data")
		return
	}

	// parse b2log record from data
	records, buffer := pbBuffParse(data[0:32], nil)

	if len(records) != 0 {
		t.Errorf("len(records) should be 0, but now it's %d", len(records))
	}
	if len(buffer) != 32 {
		t.Errorf("len(buffer) should be 32, but now it's %d", len(buffer))
	}
}
