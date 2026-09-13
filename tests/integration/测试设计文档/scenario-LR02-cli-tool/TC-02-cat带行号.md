# TC-02 cat 带行号

## 用例编号与名称

TC-02 cat 带行号

## 所属场景

LR02 CLI 工具功能验证

## 版本声明

- `bfe-pblog-tool`：当前源码版本

## 测试目的

验证 `bfe-pblog-tool cat -n` 能在每条记录前显示递增的行号。

## 运行模式

命令行模式：测试代码通过 `common.ProcessEnv.RunPblogTool` 调用 `bfe-pblog-tool cat -n` 二进制。

## 前置条件

1. 已编译 `bfe-pblog-tool` 可执行文件。
2. 已通过 `LogGenerator` 生成包含 9 条记录的 pb 日志文件。

## 输入数据

| 项目 | 内容 |
|------|------|
| 日志文件 | 由 `common.MakeRequestLog` 生成的 9 条 `BfeLog` 请求日志 |
| logid 范围 | 10000 ~ 10008 |
| CLI 命令 | `bfe-pblog-tool cat -n <logfile>` |

## 操作步骤

1. 编译 `bfe-pblog-tool` 二进制。
2. 通过 `LogGenerator` 写入 9 条 b2log 记录。
3. 调用 `RunPblogTool("cat", "-n", logFile)` 执行 cat -n 命令。
4. 逐行检查 stdout，验证每条数据行以数字+空格开头。

## 预期结果

- 退出码为 0。
- 共 9 条数据行（排除 "Time taken:" 行）。
- 每条数据行以 "1 "、"2 "、... "9 " 开头的行号前缀。

## 清理

删除临时目录。