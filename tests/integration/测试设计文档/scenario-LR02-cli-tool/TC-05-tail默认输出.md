# TC-05 tail 默认输出

## 用例编号与名称

TC-05 tail 默认输出

## 所属场景

LR02 CLI 工具功能验证

## 版本声明

- `bfe-pblog-tool`：当前源码版本

## 测试目的

验证 `bfe-pblog-tool tail` 在不指定 `-n` 时默认输出最后 10 条记录（文件不足 10 条则全部输出）。

## 运行模式

命令行模式：测试代码通过 `common.ProcessEnv.RunPblogTool` 调用 `bfe-pblog-tool tail` 二进制。

## 前置条件

1. 已编译 `bfe-pblog-tool` 可执行文件。
2. 已通过 `LogGenerator` 生成包含 9 条记录的 pb 日志文件。

## 输入数据

| 项目 | 内容 |
|------|------|
| 日志文件 | 由 `common.MakeRequestLog` 生成的 9 条 `BfeLog` 请求日志 |
| logid 范围 | 10000 ~ 10008 |
| CLI 命令 | `bfe-pblog-tool tail <logfile>` |

## 操作步骤

1. 编译 `bfe-pblog-tool` 二进制。
2. 通过 `LogGenerator` 写入 9 条 b2log 记录。
3. 调用 `RunPblogTool("tail", logFile)` 执行 tail 命令。
4. 检查退出码和输出记录数。

## 预期结果

- 退出码为 0。
- stdout 包含 9 条记录（文件只有 9 条，默认 N=10 大于文件总记录数）。

## 清理

删除临时目录。