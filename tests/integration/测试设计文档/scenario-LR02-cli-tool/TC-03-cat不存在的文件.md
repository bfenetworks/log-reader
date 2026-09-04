# TC-03 cat 不存在的文件

## 用例编号与名称

TC-03 cat 不存在的文件

## 所属场景

LR02 CLI 工具功能验证

## 版本声明

- `bfe-pblog-tool`：当前源码版本

## 测试目的

验证 `bfe-pblog-tool cat` 在指定不存在的文件路径时能正确报错并返回非零退出码。

## 运行模式

命令行模式：测试代码通过 `common.ProcessEnv.RunPblogTool` 调用 `bfe-pblog-tool cat` 二进制。

## 前置条件

1. 已编译 `bfe-pblog-tool` 可执行文件。
2. 确保 `/nonexistent/path/to/file.log` 路径不存在。

## 输入数据

| 项目 | 内容 |
|------|------|
| 日志文件 | `/nonexistent/path/to/file.log`（不存在） |
| CLI 命令 | `bfe-pblog-tool cat /nonexistent/path/to/file.log` |

## 操作步骤

1. 编译 `bfe-pblog-tool` 二进制。
2. 调用 `RunPblogTool("cat", "/nonexistent/path/to/file.log")` 执行 cat 命令。
3. 检查退出码和错误信息。

## 预期结果

- 退出码非零。
- `err` 非 nil，或 stderr 包含错误信息（如 "stat"、"no such file" 等）。

## 清理

删除临时目录。