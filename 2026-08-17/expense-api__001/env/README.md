# 个人支出记录后端 API

这是一个从 0 实现的 Go 个人支出记录后端服务，用于记录每天消费。数据存储在 SQLite，使用 `modernc.org/sqlite` 纯 Go 驱动，不依赖 cgo。

## 项目简介

服务支持：

- 支出新增、查询、修改、删除
- 按分类 ID 或分类名称筛选支出
- 按天、按月汇总支出，并返回分类明细
- 分类新增、查询、修改、删除
- 统一 JSON 响应与错误码
- SQLite 自动建表/迁移

## 目录结构

```text
.
├── cmd/server/main.go
├── internal
│   ├── config
│   ├── handler
│   │   ├── category_handler.go
│   │   ├── expense_handler.go
│   │   └── summary_handler.go
│   ├── model
│   ├── repository
│   ├── service
│   └── server
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

`handler`、`service`、`repository`、`model` 四层分离，支出、分类、汇总统计各自独立。

## 启动方式

### Docker 构建和运行

在项目目录执行：

```bash
docker build -t expense-api .
docker run -d --name expense-api -p 18090:8080 -e PORT=8080 -e DB_PATH=/app/data/expense.db expense-api
```

查看日志：

```bash
docker logs expense-api
```

停止并清理：

```bash
docker rm -f expense-api
```

### 本机 Go 运行（可选）

需要 Go 1.23 或更高版本：

```bash
go mod download
go run ./cmd/server
```

## 环境变量

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `PORT` | `8080` | HTTP 服务监听端口 |
| `DB_PATH` | `./data/expense.db` | SQLite 数据库文件路径，容器内默认为 `/app/data/expense.db` |

## API 路径

### 健康检查

`GET /healthz`

### 分类

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `POST` | `/api/v1/categories` | 新增分类 |
| `GET` | `/api/v1/categories` | 查询分类列表 |
| `PUT` | `/api/v1/categories/{id}` | 修改分类 |
| `DELETE` | `/api/v1/categories/{id}` | 删除分类 |

### 支出

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `POST` | `/api/v1/expenses` | 新增支出 |
| `GET` | `/api/v1/expenses` | 查询支出列表 |
| `GET` | `/api/v1/expenses/{id}` | 查询单条支出 |
| `PUT` | `/api/v1/expenses/{id}` | 修改支出 |
| `DELETE` | `/api/v1/expenses/{id}` | 删除支出 |

支出筛选参数：

- `category_id`：按分类 ID 精确筛选
- `category`：按分类名称筛选

### 汇总统计

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `GET` | `/api/v1/summary/daily?date=YYYY-MM-DD` | 按天汇总 |
| `GET` | `/api/v1/summary/monthly?month=YYYY-MM` | 按月汇总 |

## curl 示例

先新增分类：

```bash
curl -sS -X POST http://localhost:18090/api/v1/categories \
  -H 'Content-Type: application/json' \
  -d '{"name":"餐饮"}'
```

新增支出：

```bash
curl -sS -X POST http://localhost:18090/api/v1/expenses \
  -H 'Content-Type: application/json' \
  -d '{"amount":28.50,"category_id":1,"date":"2026-08-16","payment_method":"微信","note":"午餐"}'
```

按分类筛选：

```bash
curl -sS 'http://localhost:18090/api/v1/expenses?category_id=1'
curl -sS 'http://localhost:18090/api/v1/expenses?category=餐饮'
```

按天汇总：

```bash
curl -sS 'http://localhost:18090/api/v1/summary/daily?date=2026-08-16'
```

按月汇总：

```bash
curl -sS 'http://localhost:18090/api/v1/summary/monthly?month=2026-08'
```
