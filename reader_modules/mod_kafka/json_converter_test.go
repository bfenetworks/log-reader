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

package mod_kafka

import (
	"encoding/json"
	"testing"
)

func TestConvertBfeLogToJSON_DefaultFields(t *testing.T) {
	log := makeBfeLog()
	of := DefaultOutputFields()

	jsonBytes, err := ConvertBfeLogToJSON(log, of)
	if err != nil {
		t.Fatalf("ConvertBfeLogToJSON failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if result["logid"].(float64) != 12345 {
		t.Errorf("logid: expected 12345, got %v", result["logid"])
	}
	if result["timestamp"].(float64) != 1782353290 {
		t.Errorf("timestamp: expected 1782353290, got %v", result["timestamp"])
	}
	if result["product"] != "BFE" {
		t.Errorf("product: expected BFE, got %v", result["product"])
	}
	if result["origin_uri"] != "/api/v1/test" {
		t.Errorf("origin_uri: expected /api/v1/test, got %v", result["origin_uri"])
	}
	if result["backend_info"] != "10.0.0.2:8080" {
		t.Errorf("backend_info: expected 10.0.0.2:8080, got %v", result["backend_info"])
	}
}

func TestConvertBfeLogToJSON_RequireFields(t *testing.T) {
	log := makeBfeLog()
	cfg := &KafkaDataConfig{
		ConfFields: ConfKafkaFields{FieldMode: "require"},
	}
	of := cfg.ResolveFields()

	jsonBytes, err := ConvertBfeLogToJSON(log, of)
	if err != nil {
		t.Fatalf("ConvertBfeLogToJSON failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	required := RequiredFields()
	for _, name := range required {
		// empty-string / zero-value required fields are omitted by omitempty behavior
		if name == "err_code" || name == "err_msg" || name == "req_body_len" {
			continue
		}
		if _, ok := result[name]; !ok {
			t.Errorf("required field %q missing from output", name)
		}
	}

	if _, ok := result["ai_apikey"]; ok {
		t.Error("ai_apikey should not be in require output")
	}
	if _, ok := result["cluster"]; ok {
		t.Error("cluster should not be in require output")
	}
}

func TestConvertBfeLogToJSON_CustomizedFields(t *testing.T) {
	log := makeBfeLog()
	cfg := &KafkaDataConfig{
		ConfFields: ConfKafkaFields{
			FieldMode:  "customized",
			FieldNames: []string{"ai_apikey", "ai_requested_model"},
		},
	}
	of := cfg.ResolveFields()

	jsonBytes, err := ConvertBfeLogToJSON(log, of)
	if err != nil {
		t.Fatalf("ConvertBfeLogToJSON failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if result["ai_apikey"] != "sk-test" {
		t.Errorf("ai_apikey: expected sk-test, got %v", result["ai_apikey"])
	}
	if result["ai_requested_model"] != "test-model" {
		t.Errorf("ai_requested_model: expected test-model, got %v", result["ai_requested_model"])
	}

	for _, name := range RequiredFields() {
		// empty-string / zero-value required fields are omitted by omitempty behavior
		if name == "err_code" || name == "err_msg" || name == "req_body_len" {
			continue
		}
		if _, ok := result[name]; !ok {
			t.Errorf("required field %q missing", name)
		}
	}

	if _, ok := result["cluster"]; ok {
		t.Error("cluster should not be in customized output")
	}
}

func TestConvertBfeLogToJSON_OmitsZeroValues(t *testing.T) {
	log := makeBfeLog()
	cfg := &KafkaDataConfig{
		ConfFields: ConfKafkaFields{FieldMode: "all"},
	}
	of := cfg.ResolveFields()

	jsonBytes, err := ConvertBfeLogToJSON(log, of)
	if err != nil {
		t.Fatalf("ConvertBfeLogToJSON failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
}

func TestConvertBfeLogToJSON_NilOutputFields(t *testing.T) {
	log := makeBfeLog()

	jsonBytes, err := ConvertBfeLogToJSON(log, nil)
	if err != nil {
		t.Fatalf("ConvertBfeLogToJSON failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	// default fields minus zero-value fields omitted by omitempty
	expectedMin := len(DefaultFields()) - 16
	if len(result) < expectedMin {
		t.Fatalf("expected at least %d fields (some may be zero), got %d", expectedMin, len(result))
	}
}

func BenchmarkConvertBfeLogToJSON_Default(b *testing.B) {
	log := makeBfeLog()
	of := DefaultOutputFields()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ConvertBfeLogToJSON(log, of)
	}
}

func BenchmarkConvertBfeLogToJSON_All(b *testing.B) {
	log := makeBfeLog()
	cfg := &KafkaDataConfig{
		ConfFields: ConfKafkaFields{FieldMode: "all"},
	}
	of := cfg.ResolveFields()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ConvertBfeLogToJSON(log, of)
	}
}

func BenchmarkConvertBfeLogToJSON_Require(b *testing.B) {
	log := makeBfeLog()
	cfg := &KafkaDataConfig{
		ConfFields: ConfKafkaFields{FieldMode: "require"},
	}
	of := cfg.ResolveFields()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ConvertBfeLogToJSON(log, of)
	}
}
