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

package reader_conf

import (
	"testing"

	gcfg "gopkg.in/gcfg.v1"
)

/* load config from config file    */
func confBasicLoad(filePath string) (ReaderConfig, error) {
	var cfg ReaderConfig
	var err error

	// read config from file
	err = gcfg.ReadFileInto(&cfg, filePath)
	if err != nil {
		return cfg, err
	}

	// check basic conf
	err = ConfBasicCheck(&cfg.Main)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

// normal case
func Test_confBasic_case1(t *testing.T) {
	config, err := confBasicLoad("./test_data/conf_basic/config_1.conf")
	if err != nil {
		t.Errorf("err in ConfBasicCheck():%s", err.Error())
		return
	}

	if config.Main.HttpPort != 8992 {
		t.Errorf("config.HttpPort should be 8992, but now it's %d", config.Main.HttpPort)
		return
	}

	if config.Main.MaxCpus != 8 {
		t.Errorf("config.MaxCpus should be 8, but now it's %d", config.Main.MaxCpus)
		return
	}

	if config.Main.MonitorInterval != 20 {
		t.Errorf("config.monitorInterval should be 20, but now it's %d",
			config.Main.MonitorInterval)
		return
	}
}
