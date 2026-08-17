# Bug 是什么

汇总金额的换算方向错误，日汇总和分类明细中的 `total` 被放大 100 倍；月度查询还因日期格式校验错误而返回校验失败。

# 如何触发

1. 启动服务并创建一个分类。
2. 新增若干支出。
3. 调用日汇总和月汇总接口。

# 错误信息

```text
unexpected daily total: &{Date:2026-08-16 TotalCents:1975 Total:197500 ...}
month must use YYYY-MM format
```
