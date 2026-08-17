# Bug 是什么

新增支出后，金额的 `amount_cents`、`amount` 以及按日汇总结果不一致：创建接口返回的 cents 被截断，列表和详情中的 amount 被按原始 cents 数值直接解释，并且列表与分类汇总中会混入一个空项。

# 如何触发

1. 启动服务并创建一个分类。
2. 新增一条 `amount=12.34` 的支出。
3. 调用支出列表、详情和 `GET /api/v1/summary/daily?date=...`。

# 错误信息

服务层集成测试中的关键失败：

```text
unexpected expense after get: &{ID:1 AmountCents:1234 Amount:1234 ...}
unexpected listed expense: ...
unexpected daily summary: ...
```
