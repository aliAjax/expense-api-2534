# Bug 是什么

数据库访问层没有向底层查询传递调用方的 context，列表和汇总接口在 context 已取消时仍会继续查询并返回成功，导致取消/超时语义失效。

# 如何触发

1. 启动服务并准备数据。
2. 使用已经取消的 context 调用分类列表、支出列表或日汇总接口。

# 错误信息

```text
list categories with canceled context: expected error, got nil
```
