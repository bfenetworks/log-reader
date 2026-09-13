# bfe-pblog-tool 使用指南

本文档介绍如何使用 `bfe-pblog-tool` 命令行工具离线查看 BFE protobuf 格式访问日志（`pb_access3.log`）。

---

## 1. 概述

`bfe-pblog-tool` 是 `log-reader` 仓库的第二个构建产物，与 daemon 进程 `log-reader` 并列。它复用 `bfe_log_reader` 的同一套解析内核，提供两个子命令：

```
┌──────────────────┐         ┌──────────────────────┐
│ pb_access3.log   │ ──────> │ bfe-pblog-tool       │
│ (b2log 二进制)    │  读取    │                      │
└──────────────────┘         │  cat  → 全量输出       │
                             │  tail → 末尾N条+跟随   │
                             └──────────────────────┘
                                      │
                                      ▼ stdout
```

| 子命令 | 用途 |
|--------|------|
| `cat` | 全量输出日志文件中的所有记录 |
| `tail` | 输出末尾 N 条记录，支持 `-f` 持续跟随新增数据 |

> 支持平台：macOS / Linux / Windows。

---

## 2. 前提条件

| 条件 | 说明 |
|------|------|
| Go | 1.22+ |
| 日志源 | BFE 开启 PB 访问日志后产出的 `pb_access3.log`（[bfe-access-pb](https://github.com/bfenetworks/bfe-access-pb) 二进制格式） |

---

## 3. 编译

```bash
cd log-reader
make bfe-pblog-tool
```

编译产物输出到：

```
output/bin/bfe-pblog-tool
```

---

## 4. cat 子命令

### 4.1. 用途

从文件头开始读取，全量输出 pb 日志文件中的所有记录。

### 4.2. 用法

```bash
bfe-pblog-tool cat [options] <logfile>
```

### 4.3. 参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `-n` | bool | `false` | 显示行号 |

### 4.4. 示例

**全量输出：**

```bash
./bfe-pblog-tool cat /path/to/pb_access3.log
```

**带行号输出：**

```bash
./bfe-pblog-tool cat -n /path/to/pb_access3.log
```

输出示例：

```
1 logid:10001 timestamp:1782353290 product:BFE ...
2 logid:10002 timestamp:1782353291 product:BFE ...
3 logid:10003 timestamp:1782353292 product:BFE ...
Time taken: 12.345ms
```

### 4.5. 说明

- 输出格式为 `bfe_access_pb.BfeLog` 的 `String()` 表示。
- 输出末尾附带 `Time taken: Xms` 显示耗时。
- 大文件可能产生大量输出，建议配合 `head` 或重定向使用。

---

## 5. tail 子命令

### 5.1. 用途

输出日志文件末尾 N 条记录。可选开启持续跟随模式，实时输出新增记录。

### 5.2. 用法

```bash
bfe-pblog-tool tail [options] <logfile>
```

### 5.3. 参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `-n` | int | `10` | 显示末尾 N 条记录 |
| `-f` | bool | `false` | 开启持续跟随模式 |
| `--interval` | int | `500` | 跟随模式下的轮询间隔（毫秒） |

### 5.4. 示例

**默认末尾 10 条：**

```bash
./bfe-pblog-tool tail /path/to/pb_access3.log
```

**末尾 3 条：**

```bash
./bfe-pblog-tool tail -n 3 /path/to/pb_access3.log
```

**持续跟随（Ctrl+C 退出）：**

```bash
./bfe-pblog-tool tail -n 20 -f /path/to/pb_access3.log
```

**自定义轮询间隔 200ms：**

```bash
./bfe-pblog-tool tail -n 20 -f --interval 200 /path/to/pb_access3.log
```

### 5.5. 跟随模式说明

| 特性 | 说明 |
|------|------|
| 轮询机制 | 按 `--interval` 间隔检查文件大小变化 |
| 日志轮转检测 | 检测 inode 变更（unix）或文件大小缩小（truncate），自动重新打开文件 |
| 退出方式 | 按一次 Ctrl+C 优雅停止；连按两次 Ctrl+C 强制退出 |
| N 上限 | `-n` 最大值为 100000（防止 OOM），超过自动截断 |

---

## 6. 常见使用场景

### 6.1. 排查线上问题

快速查看最近的请求日志：

```bash
bfe-pblog-tool tail -n 5 /home/work/bfe/log/pb_access3.log
```

### 6.2. 验证日志完整性

用 cat 全量输出确认日志文件可正常解析：

```bash
bfe-pblog-tool cat /path/to/pb_access3.log | wc -l
```

### 6.3. 持续监控

实时观察新日志写入：

```bash
bfe-pblog-tool tail -f /home/work/bfe/log/pb_access3.log
```

### 6.4. 配合管道过滤

搜索特定关键字：

```bash
bfe-pblog-tool cat /path/to/pb_access3.log | grep "error"
bfe-pblog-tool cat /path/to/pb_access3.log | grep "ai.example.org"
```

---

## 7. 注意事项

| 项目 | 说明 |
|------|------|
| 输入格式 | 仅支持 b2log 二进制格式的 `pb_access3.log`，无法读取普通文本日志 |
| 大文件 | `cat` 全量输出可能产生大量内容，建议配合 `head`、`grep` 或重定向 |
| `-n` 上限 | `tail -n` 最大值为 100000，超过后自动截断以防止 OOM |
| 耗时显示 | 仅 `cat` 子命令在末尾显示 `Time taken:`，`tail` 不显示 |
| 跨平台 | unix 平台通过 inode 检测日志轮转；windows 平台不支持 inode 检测，仅通过文件大小变化检测 truncate |

---

## 8. 附录：验证清单

| 序号 | 验证项 | 命令 | 预期 |
|:---:|--------|------|------|
| 1 | 编译成功 | `make bfe-pblog-tool` | 无报错，产物在 `output/bin/bfe-pblog-tool` |
| 2 | cat 全量输出 | `bfe-pblog-tool cat <logfile>` | 输出全部记录 + "Time taken:" |
| 3 | cat 带行号 | `bfe-pblog-tool cat -n <logfile>` | 每行带 "1 " "2 " 前缀 |
| 4 | tail 默认输出 | `bfe-pblog-tool tail <logfile>` | 输出末尾 10 条（或全部，如不足 10 条） |
| 5 | tail 指定条数 | `bfe-pblog-tool tail -n 3 <logfile>` | 仅输出末尾 3 条 |
| 6 | tail 跟随模式 | `bfe-pblog-tool tail -f <logfile>` | 持续运行，新日志写入时输出 |
| 7 | 错误处理 | `bfe-pblog-tool cat /nonexistent` | 非零退出码 + 错误信息 |
