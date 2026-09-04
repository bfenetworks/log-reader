# TC-01 cat 基本输出

## 用例编号与名称

TC-01 cat 基本输出

## 所属场景

LR02 CLI 工具功能验证

## 版本声明

- `bfe-pblog-tool`：当前源码版本

## 测试目的

验证 `bfe-pblog-tool cat` 能正确读取 pb 日志文件并全量输出所有记录，输出格式正确，并以 "Time taken:" 行结尾。

## 运行模式

命令行模式：测试代码通过 `common.ProcessEnv.RunPblogTool` 调用 `bfe-pblog-tool cat` 二进制。

## 前置条件

1. 已编译 `bfe-pblog-tool` 可执行文件。
2. 已通过 `LogGenerator` 生成包含 9 条记录的 pb 日志文件。

## 输入数据

| 项目 | 内容 |
|------|------|
| 日志文件 | 由 `common.MakeRequestLog` 生成的 9 条 `BfeLog` 请求日志 |
| logid 范围 | 10000 ~ 10008 |
| CLI 命令 | `bfe-pblog-tool cat <logfile>` |

## 操作步骤

1. 编译 `bfe-pblog-tool` 二进制。
2. 通过 `LogGenerator` 写入 9 条 b2log 记录。
3. 调用 `RunPblogTool("cat", logFile)` 执行 cat 命令。
4. 检查退出码、stdout 内容。

## 预期结果

- 退出码为 0。
- stdout 非空，包含 9 条记录的格式化输出。
- stdout 以 "Time taken:" 行结尾。
- stderr 为空。

## 清理

删除临时目录。