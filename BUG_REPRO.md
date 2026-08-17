# Bug 是什么

更新支出成功后，接口返回 nil 对象；再次读取支出详情也会得到 nil 结果，导致调用方解引用时发生空指针 panic。

# 如何触发

1. 启动服务并创建一个分类。
2. 新增一条支出。
3. 调用支出更新接口，然后读取更新后的支出详情。

# 错误信息

```text
panic: runtime error: invalid memory address or nil pointer dereference
```
