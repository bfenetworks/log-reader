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
	"fmt"

	gcfg "gopkg.in/gcfg.v1"
)

var openDebug bool

// ConfModKafka Kafka 模块配置
type ConfModKafka struct {
	Basic ModKafkaBasic
	Kafka ConfKafka
}

// ConfKafka Kafka 连接和发送配置
type ConfKafka struct {
	Brokers         string // Kafka broker 地址列表，逗号分隔
	Topic           string // 目标 Topic
	DeadLetterTopic string // 死信 Topic
	Compression     string // 压缩方式：none / snappy / gzip / lz4 /zstd
	BatchSize       int    // 批量发送大小
	LingerMs        int    // 批量发送等待时间（毫秒）
	MaxRetries      int    // 最大重试次数
}

// ModKafkaBasic Kafka 数据配置文件路径
type ModKafkaBasic struct {
	DataPath  string // 数据配置文件路径，相对于 mod_kafka.conf 所在目录
	OpenDebug bool   // 是否开启 Debug 模式
}

// LoadConfig 加载 mod_kafka 配置文件
func LoadConfig(filePath string) (*ConfModKafka, error) {
	var cfg ConfModKafka

	if err := gcfg.ReadFileInto(&cfg, filePath); err != nil {
		return &cfg, err
	}

	if err := ConfModKafkaCheck(&cfg); err != nil {
		return &cfg, err
	}

	openDebug = cfg.Basic.OpenDebug

	return &cfg, nil
}

// ConfModKafkaCheck 校验配置
func ConfModKafkaCheck(cfg *ConfModKafka) error {
	if cfg.Kafka.Brokers == "" {
		return fmt.Errorf("Kafka.Brokers is empty")
	}

	if cfg.Kafka.Topic == "" {
		return fmt.Errorf("Kafka.Topic is empty")
	}

	if cfg.Kafka.BatchSize <= 0 {
		cfg.Kafka.BatchSize = 1000
	}

	if cfg.Kafka.LingerMs <= 0 {
		cfg.Kafka.LingerMs = 100
	}

	if cfg.Kafka.MaxRetries <= 0 {
		cfg.Kafka.MaxRetries = 3
	}

	switch cfg.Kafka.Compression {
	case "", "none", "snappy", "gzip", "lz4", "zstd":
	default:
		return fmt.Errorf("Kafka.Compression invalid: %s", cfg.Kafka.Compression)
	}

	return nil
}
