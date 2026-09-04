# TC-09 tail 跟随新增数据

## 用例编号与名称

TC-09 tail 跟随新增数据

## 所属场景

LR02 CLI 工具功能验证

## 版本声明

- `bfe-pblog-tool`：当前源码版本

## 测试目的

验证 `bfe-pblog-tool tail -f` 在后台持续运行期间，能捕获到新追加到日志文件的记录。

## 运行模式

命令行模式：测试代码通过 `common.ProcessEnv.StartPblogTool` 在后台启动 `bfe-pblog-tool tail -f` 进程，然后通过 `LogGenerator` 追加新数据，读取 stdout 管道验证新数据被捕获。

## 前置条件

1. 已编译 `bfe-pblog-tool` 可执行文件。
2. 已通过 `LogGenerator` 生成包含 5 条记录的 pb 日志文件。

## 输入数据

| 项目 | 内容 |
|------|------|
| 初始日志 | 5 条 `BfeLog` 请求日志（logid 10000 ~ 10004） |
| 追加日志 | 3 条 `BfeLog` 请求日志（logid 10005 ~ 10007） |
| CLI 命令 | `bfe-pblog-tool tail -n 5 -f --interval 100 <logfile>` |

## 操作步骤

1. 编译 `bfe-pblog-tool` 二进制。
2. 写入 5 条初始记录。
3. 通过 `StartPblogTool` 在后台启动 `tail -n 5 -f --interval 100`。
4. 等待 500ms 让初始输出完成。
5. 通过 `LogGenerator` 追加 3 条新记录。
6. 在 5 秒超时内读取 stdout，验证新记录出现。
7. 发送 SIGINT 停止进程。

## 预期结果

- 后台进程正常启动。
- stdout 包含至少 8 条记录（5 条初始 + 3 条新增）。
- 进程能被 SIGINT 优雅停止。

## 清理

停止后台进程，删除临时目录。