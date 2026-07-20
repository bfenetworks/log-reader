# log-reader

Log reader for BFE access logs, reads protobuf-encoded access logs and sends them to downstreams.

## Build

```sh
make
```

Binary output: `output/bin/log_reader`

## Test

```sh
make test
```

## Usage

```sh
# ./log_reader -h
  -a    automatically get name of bfe cluster
  -b    is read from begin
  -c string
        root path of config file (default "../conf")
  -d    to show debug log (otherwise >= info)
  -h    to show help
  -l string
        dir path of log (default "../log")
  -s    to show log in stdout

# ./log_reader -c ../conf/

```

## Project Structure

```
├── main/               # Entry point
├── bfe_log_reader/     # Core log reader (pb parsing, log file tailing)
├── reader_conf/        # Configuration loading
├── reader_module/      # Module framework
├── reader_modules/     # Built-in modules
│   └── mod_kafka/      # Kafka output module
├── reader_util/        # Utility functions
└── conf/               # Default config files
```
