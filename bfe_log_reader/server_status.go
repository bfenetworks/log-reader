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

var COUNTER_KEYS = []string{
	// errors
	// for logReader
	"ERR_PB_OPEN",   // in open pb file (bfe_access_pb)
	"ERR_PB_SEEK",   // in seek in pb file
	"ERR_PB_STAT",   // in stat for pb file
	"ERR_PB_CLOSE",  // in close pb file
	"ERR_PB_READ",   // in read pb file
	"ERR_PB_DECODE", // in decode pb record

	// counters
	// for logReader
	"SUM_READ_DATA",          // total number of reading data
	"SUM_READ_DATA_EMPTY",    // total number of reading empty data
	"SUM_READ_RECORDS_EMPTY", // total number of reading empty records
	"PB_LOG_RELOCATE",        // relocate pb log file
	"SUM_PB_RECORD",          // total number of pb records decoded
	"SUM_PB_BATCH_COUNT",     // total number of pb batches
}

// get srv.srvState
func (br *BfeLogReader) srvStateGet(params map[string][]string) ([]byte, error) {
	// get data for srvState
	s := br.srvState.GetAll()
	return s.FormatOutput(params) // TBD: wait for modification of golang-lib
}

// get srv.srvStateDiff
func (br *BfeLogReader) srvStateDiffGet(params map[string][]string) ([]byte, error) {
	// get data for srvStateDiff
	diff := br.srvStateDiff.Get()
	return diff.FormatOutput(params) // TBD: wait for modification of golang-lib
}
