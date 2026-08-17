# Bug 是什么

分类更新时，找不到记录的错误被统一归类为内部错误，HTTP 状态码和业务错误码都从“未找到”变成了内部错误。

# 如何触发

1. 启动服务。
2. 调用分类更新接口，传入一个不存在的分类 ID。

# 错误信息

```text
update missing category error = *model.APIError failed to update category, want not found
```
