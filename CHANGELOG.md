# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-07-20

### Added

- Core log reader with protobuf-encoded BFE access log parsing and tailing support.
- Module framework (`reader_module`) for extensible log processing pipelines.
- Config loading system (`reader_conf`) with support for basic and access-pb config types.
- Built-in Kafka output module (`mod_kafka`) for forwarding parsed access logs to Kafka.


