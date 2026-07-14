// Copyright (c) 2026 The BFE Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Required fields: bfeClusterName, idc, isSmallCluster, clusters have been removed as per requirements.
// This test file has been updated to reflect the changes.

package reader_conf

import (
	"testing"
)

// normal case
func Test_confAll_case1(t *testing.T) {
	config, err := ReaderConfigLoad("./test_data/conf_all/config_1.conf")
	if err != nil {
		t.Errorf("err in ReaderConfigLoad():%s", err.Error())
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

	if config.PbAccessLogConf.LogFile != "/home/work/bfe/log/pb_access3.log" {
		t.Errorf("config.PbAccessLogConf.LogFile should be /home/work/bfe/log/pb_access3.log , but it's %s",
			config.PbAccessLogConf.LogFile)
		return
	}

	if len(config.PbAccessLogConf.Modules) != 1 {
		t.Errorf("len(config.PbAccessLogConf.Modules) should be 1, but it's %d",
			len(config.PbAccessLogConf.Modules))
		return
	}

	if config.PbAccessLogConf.Modules[0] != "mod_kafka" {
		t.Errorf("config.PbAccessLogConf.Modules[0] should be mod_kafka, but it's %s",
			config.PbAccessLogConf.Modules[0])
		return
	}
}
