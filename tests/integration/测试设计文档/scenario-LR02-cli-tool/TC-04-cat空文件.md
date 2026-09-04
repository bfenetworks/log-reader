# TC-04 cat 空文件

## 用例编号与名称

TC-04 cat 空文件

## 所属场景

LR02 CLI 工具功能验证

## 版本声明

- `bfe-pblog-tool`：当前源码版本

## 测试目的

验证 `bfe-pblog-tool cat` 在处理空 pb 日志文件时不会崩溃，返回零退出码，且无记录输出。

## 运行模式

命令行模式：测试代码通过 `common.ProcessEnv.RunPblogTool` 调用 `bfe-pblog-tool cat` 二进制。

## 前置条件

1. 已编译 `bfe-pblog-tool` 可执行文件。
2. 已创建一个空的 pb 日志文件。

## 输入数据

| 项目 | 内容 |
|------|------|
| 日志文件 | 空文件（0 字节） |
| CLI 命令 | `bfe-pblog-tool cat <emptyfile>` |

## 操作步骤

1. 编译 `bfe-pblog-tool` 二进制。
2. 创建一个空文件。
3. 调用 `RunPblogTool("cat", emptyFile)` 执行 cat 命令。
4. 检查退出码和 stdout 内容。

## 预期结果

- 退出码为 0。
- stdout 中无记录输出（可能包含 "Time taken:" 行）。
- 无崩溃或异常退出。

## 清理

删除临时目录。