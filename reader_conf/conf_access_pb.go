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
	"errors"
	"strings"
)

type PbAccessLogConf struct {
	LogFile         string   // the path of pb access log
	Modules         []string // modules for logReader to load, format: "mod_name" or "mod_name:true" or "mod_name:false"
	MaxSizePerBatch int      // max element size per batch, <=0, unlimited

	FinalModules []string
}

func (cfg *PbAccessLogConf) Check() error {
	if len(cfg.LogFile) > 0 && len(cfg.Modules) <= 0 {
		return errors.New("PbAccessLogConf: modules must be set, when logFile are set")
	}

	if len(cfg.LogFile) <= 0 && len(cfg.Modules) > 0 {
		return errors.New("PbAccessLogConf: logFile must be set, when modules are set")
	}

	if len(cfg.LogFile) > 0 && len(cfg.Modules) > 0 {
		for _, module := range cfg.Modules {
			moduleName, enable := ParseModule(module)
			if _, ok := modulesSupport[moduleName]; !ok {
				return errors.New("PbAccessLogConf: module[" + module + "] is not support")
			}
			if enable {
				cfg.FinalModules = append(cfg.FinalModules, moduleName)
			}
		}
	}

	if cfg.MaxSizePerBatch <= 0 {
		cfg.MaxSizePerBatch = -1
	}

	return nil
}

// ParseModule extracts from "mod_name:true" or "mod_name:false"  or "mod_name" format
func ParseModule(module string) (string, bool) {
	parts := strings.SplitN(module, ":", 2)
	enable := true
	mn := strings.TrimSpace(parts[0])
	if len(parts) >= 2 {
		tmp := strings.ToLower(strings.TrimSpace(parts[1]))
		if len(tmp) <= 0 {
			enable = true
		} else {
			enable = (tmp == "true")
		}
	}
	return mn, enable
}
