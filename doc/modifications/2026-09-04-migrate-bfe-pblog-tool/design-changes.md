# log-reader 新增 bfe-pblog-tool CLI 工具

## 背景

`log-reader` daemon 负责实时 tailing BFE 的 protobuf 访问日志（`pb_access3.log`）并分发到下游（如 Kafka）。在日常运维和问题排查中，经常需要离线查看 PB 日志文件的内容，但此前只能通过 `bfe-access-pb` 仓库中的独立工具 `bfe-pblog-tool` 来完成。

将 `bfe-pblog-tool` 纳入 `log-reader` 仓库作为第二个构建产物有以下优势：

- 共享 `bfe_log_reader` 的解析内核，避免跨仓库维护重复代码；
- 统一版本发布和 CI/CD 流程；
- protobuf 协议升级时只需在一处修改解析逻辑。

## 修改目标

1. 新增 `bfe-pblog-tool` 命令行工具，提供 `cat` 和 `tail` 两个子命令，用于离线查看 PB 格式访问日志。
2. 复用 `bfe_log_reader` 已导出的文件操作和 PB 解析 API（`LogFileOpen`、`FileRead`、`PbBuffParse` 等）。
3. 支持跨平台运行（macOS / Linux / Windows），tail 跟随模式能检测日志轮转。
4. 新增集成测试场景 LR02，覆盖 CLI 全参数组合与边界条件。
5. 新增用户使用指南文档。

---

## 功能设计

### 1. 子命令设计

#### cat 子命令

全量读取 PB 日志文件，输出所有记录的文本表示。

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `-n` | bool | `false` | 在每行前显示递增行号 |

用法：`bfe-pblog-tool cat [-n] <logfile>`

实现方式：从头开始读取文件，以 8KB 块为单位通过 `FileRead` + `DataBufferParse` 解析，每解析出一批记录即通过回调输出。输出末尾附带 `Time taken: Xms` 耗时信息。

#### tail 子命令

输出日志文件末尾 N 条记录，可选开启持续跟随模式。

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `-n` | int | `10` | 显示末尾 N 条记录 |
| `-f` | bool | `false` | 开启持续跟随模式 |
| `--interval` | int | `500` | 跟随模式轮询间隔（毫秒） |

用法：`bfe-pblog-tool tail [-n N] [-f] [--interval MS] <logfile>`

实现方式分两阶段：

- **初始输出**（`tailLastN`）：从文件末尾反向读取数据块（初始 8MB，指数增长至 64MB 上限），解析出记录后保留最后 N 条，按时间正序输出。
- **跟随模式**（`followFrom`）：按 `--interval` 间隔轮询文件大小变化，增量读取并解析新数据。通过 inode 检测（unix）或文件大小缩小检测（truncate/轮转），自动重新打开文件。

`-n` 上限为 100000 条（`MaxTailRecords`），超过自动截断以防止 OOM。

### 2. 代码组织

```
cmd/bfe-pblog-tool/
├── main.go              # CLI 入口，注册子命令并启动 urfave/cli
├── cat/                 # cat 子命令
│   ├── cat.go           # CLI Action：参数解析 + 调用 PblogCat
│   ├── pb_cat.go        # PblogCat 实现
│   └── pb_cat_test.go   # 单元测试
├── tail/                # tail 子命令
│   ├── tail.go          # CLI Action：参数解析 + 调用 PblogTail
│   ├── pb_tail.go       # PblogTail / TailOptions / tailLastN / followFrom
│   ├── pb_tail_unix.go  # unix inode 获取（//go:build !windows）
│   ├── pb_tail_windows.go # windows inode 获取（//go:build windows）
│   └── pb_tail_test.go  # 单元测试
└── common/              # 共享基础设施
    └── system.go        # RegistCmd / Cmds 命令注册机制
```

各子命令包在 `init()` 中调用 `common.RegistCmd()` 注册自身命令定义。`main.go` 通过空白导入 `_ "cat"` 和 `_ "tail"` 触发注册。

### 3. 与 bfe_log_reader 的依赖关系

CLI 工具复用 `bfe_log_reader` 包已导出的 API：

| 导出符号 | 用途 |
|----------|------|
| `LogFileOpen()` | 打开 PB 日志文件 |
| `CloseFileAndInit()` | 关闭文件并重置状态 |
| `FileRead()` | 按块读取文件数据 |
| `DataBufferParse()` | 解析数据缓冲区中的 PB 记录（cat 使用） |
| `PbBuffParse()` | 解析原始 PB buffer（tail 使用） |
| `MAX_BUFF_SIZE` | 读取缓冲区大小上限 |
| `NewPbLogReader()` | 创建 PbLogReader 实例 |

---

## 详细改动

### 1. 新增 CLI 工具代码

**文件：**

- `cmd/bfe-pblog-tool/main.go`：CLI 入口，使用 `urfave/cli` 框架。
- `cmd/bfe-pblog-tool/cat/cat.go`：cat 子命令定义（名称、用法、flags）。
- `cmd/bfe-pblog-tool/cat/pb_cat.go`：`PblogCat` 函数实现。
- `cmd/bfe-pblog-tool/cat/pb_cat_test.go`：cat 单元测试。
- `cmd/bfe-pblog-tool/tail/tail.go`：tail 子命令定义。
- `cmd/bfe-pblog-tool/tail/pb_tail.go`：`PblogTail`、`TailOptions`、`ErrTailStop`、`tailLastN`、`followFrom` 实现。
- `cmd/bfe-pblog-tool/tail/pb_tail_unix.go`：unix 平台 `getInode()` 实现（`syscall.Stat_t.Ino`）。
- `cmd/bfe-pblog-tool/tail/pb_tail_windows.go`：windows 平台 `getInode()` 实现（返回 0）。
- `cmd/bfe-pblog-tool/tail/pb_tail_test.go`：tail 单元测试（`validateTailOptions` 验证 + `PblogTail` 非跟随模式测试）。
- `cmd/bfe-pblog-tool/common/system.go`：`RegistCmd` / `Cmds` 命令注册基础设施。

### 2. 导出 bfe_log_reader 符号

**文件：** `bfe_log_reader/` 下多个文件

将原本未导出的方法名首字母大写，使其可被 `cmd/bfe-pblog-tool/` 下的包调用：

| 原名称 | 导出后名称 |
|--------|-----------|
| `logFileOpen` | `LogFileOpen` |
| `closeFileAndInit` | `CloseFileAndInit` |
| `fileRead` | `FileRead` |
| `dataBufferParse` | `DataBufferParse` |
| `pbBuffParse` | `PbBuffParse` |

同时导出相关结构体字段 `LogFd`、`DataBuffer` 和常量 `MAX_BUFF_SIZE`。

### 3. Makefile 构建目标

**文件：** `Makefile`

新增 `bfe-pblog-tool` 构建目标：

```makefile
bfe-pblog-tool:
	go build -o output/bin/bfe-pblog-tool ./cmd/bfe-pblog-tool
```

`release` 目标中增加 `bfe-pblog-tool` 的交叉编译和打包。

### 4. 集成测试

**文件：**

- `tests/integration/common/process_env.go`：扩展 `ProcessEnv`，新增 `BuildPblogTool()`、`RunPblogTool()`、`StartPblogTool()` 方法。
- `tests/integration/implementation/scenario-LR02-cli-tool/lr02_cli_tool_test.go`：10 个集成测试用例。
- `tests/integration/测试设计文档/scenario-LR02-cli-tool/`：场景说明 + 9 个 TC 设计文档。
- `tests/integration/README.md`：新增 LR02 场景覆盖说明和参考文档链接。
- `tests/integration/测试设计文档/测试场景总体说明.md`：场景清单新增 LR02。

### 5. 文档

| 文件 | 说明 |
|------|------|
| `README.md` | 新增 `bfe-pblog-tool` 编译与使用说明、目录结构展开 |
| `doc/howto/02-bfe-pblog-tool-使用指南.md` | 用户使用指南：子命令参数、示例、常见场景、注意事项 |
| `AGENTS.md` | 目录结构、WHERE TO LOOK、CODE MAP 更新 |
| `bfe_log_reader/AGENTS.md` | EXPORTED API 说明更新 |
| `doc/modifications/2026-09-04-migrate-bfe-pblog-tool/design-changes.md` | 本设计变更文档 |

---

## 代码变更说明

- `cmd/bfe-pblog-tool/`（新增）：完整的 CLI 工具实现，含 cat/tail 子命令、单元测试、命令注册框架。
- `bfe_log_reader/`：导出文件操作和 PB 解析方法供 CLI 复用。
- `Makefile`：新增 `bfe-pblog-tool` 构建目标和 release 打包。
- `tests/integration/common/process_env.go`：扩展测试框架支持 CLI 二进制的编译与执行。
- `tests/integration/implementation/scenario-LR02-cli-tool/`（新增）：LR02 集成测试。
- `tests/integration/测试设计文档/scenario-LR02-cli-tool/`（新增）：LR02 测试设计文档。
- `tests/integration/README.md`：LR02 场景覆盖说明。
- `tests/integration/测试设计文档/测试场景总体说明.md`：场景清单更新。
- `README.md`：新增 CLI 工具说明。
- `doc/howto/02-bfe-pblog-tool-使用指南.md`（新增）：用户使用指南。
- `AGENTS.md` / `bfe_log_reader/AGENTS.md`：项目知识库更新。

---

## 验证步骤

1. 编译 CLI 工具：

```bash
make bfe-pblog-tool
```

2. 单元测试：

```bash
go test ./cmd/bfe-pblog-tool/cat/... -v
go test ./cmd/bfe-pblog-tool/tail/... -v
```

3. 集成测试：

```bash
go test ./tests/integration/implementation/scenario-LR02-cli-tool/... -v
```

4. 全量测试：

```bash
make test
```

5. 手工验证：

```bash
./output/bin/bfe-pblog-tool cat -n <pb_access3.log>
./output/bin/bfe-pblog-tool tail -n 3 <pb_access3.log>
./output/bin/bfe-pblog-tool tail -f <pb_access3.log>
```

---

## 影响范围

| 文件 | 影响 |
|------|------|
| `cmd/bfe-pblog-tool/` | 新增：完整 CLI 工具（main + cat + tail + common） |
| `bfe_log_reader/` | 修改：导出 5 个方法 + 2 个字段 + 1 个常量 |
| `Makefile` | 修改：新增 `bfe-pblog-tool` 构建目标 + release 打包 |
| `tests/integration/common/process_env.go` | 修改：新增 `BuildPblogTool` / `RunPblogTool` / `StartPblogTool` |
| `tests/integration/implementation/scenario-LR02-cli-tool/` | 新增：10 个集成测试用例 |
| `tests/integration/测试设计文档/scenario-LR02-cli-tool/` | 新增：场景说明 + 9 个 TC 文档 |
| `tests/integration/README.md` | 修改：新增 LR02 场景说明 |
| `tests/integration/测试设计文档/测试场景总体说明.md` | 修改：场景清单新增 LR02 |
| `README.md` | 修改：新增 CLI 工具说明 |
| `doc/howto/02-bfe-pblog-tool-使用指南.md` | 新增：用户使用指南 |
| `AGENTS.md` | 修改：目录结构、代码地图更新 |
| `bfe_log_reader/AGENTS.md` | 修改：导出 API 说明更新 |

---

## 依赖与兼容性

- 新增 `github.com/urfave/cli` 依赖（`go.mod` 已包含）。
- CLI 工具为独立二进制，不影响 daemon 的运行。
- `bfe_log_reader` 导出的符号为新增导出，原有 daemon 内部调用不受影响。
- 集成测试新增 LR02 场景，不影响现有 LR01 场景。

---

*文档生成日期：2026-09-04*
