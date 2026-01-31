# 在线书店系统 (Bookstore Management System)

一个基于 Go 语言的完整电商平台，支持图书展示、购物车、订单管理、支付集成（支付宝）等功能，包含前后台管理系统。

## 🏗️ 项目结构

```
E:.
├── air.toml                 # 热重载配置
├── go.mod                   # Go 模块定义
├── go.sum                   # 依赖校验
├── main.go                  # 程序入口
├── project_tree.txt         # 项目结构说明
├── .idea/                   # IDE 配置 (JetBrains)
├── common/                  # 公共模块
├── config/                  # 配置文件与数据库、支付等配置
├── controller/              # 控制器层
│   ├── admin/               # 后台管理接口
│   └── frontend/            # 前端用户接口
├── cron/                    # 定时任务
├── dao/                     # 数据访问层
├── model/                   # 数据模型定义
├── router/                  # 路由与中间件
├── service/                 # 业务逻辑服务
├── test/                    # 单元测试
├── tmp/                     # 临时文件与构建输出
└── utils/                   # 工具包（JWT、日志、邮件、哈希等）
```

## ✨ 功能特性

### 后台管理 (`/controller/admin`)
- 管理员身份验证
- 图书管理
- 分类管理
- 订单管理
- 用户管理
- 数据看板

### 前端用户 (`/controller/frontend`)
- 用户注册/登录
- 图书浏览与搜索
- 购物车管理
- 订单流程
- 支付宝支付集成
- 个人中心

### 核心模块
- **JWT 身份验证** (`/utils/jwt.go`)
- **Redis 缓存与库存同步** (`/dao/redisdao.go`, `/service/stock_sync_service.go`)
- **支付宝支付回调处理** (`/config/alipay.go`, `/controller/frontend/alipay_notify.go`)
- **分布式ID生成** (`/utils/snowflake.go`)
- **日志与邮件工具** (`/utils/logger.go`, `/utils/email.go`)

## 🛠️ 技术栈

- **语言**: Go 1.19+
- **Web框架**: Gin / Echo (基于 `router/` 推断)
- **数据库**: MySQL + Redis
- **配置管理**: YAML (`config/config.yaml`)
- **支付集成**: 支付宝
- **工具**: JWT、雪花ID、邮件验证、日志系统

## 🚀 快速开始

### 1. 克隆项目
```bash
git clone <仓库地址>
cd bookstore
```

### 2. 配置环境
复制并编辑配置文件：
```bash
cp config/config.yaml.example config/config.yaml
# 修改数据库、Redis、支付宝等配置
```

### 3. 安装依赖
```bash
go mod download
```

### 4. 运行项目
```bash
go run main.go
```
或使用 `air` 热重载：
```bash
air
```

### 5. 访问系统
- 前台地址：`http://localhost:8080`
- 后台地址：`http://localhost:8080/admin`

## 📦 数据库

项目使用 MySQL 存储业务数据，Redis 用于缓存和库存同步。数据库表结构请参考 `/model/` 中的定义。

## 🔧 测试

运行测试用例：
```bash
go test ./test/...
```

## 📄 许可证

MIT License

---

你可以将这个 `README.md` 放置在项目根目录，并根据实际情况调整链接、端口、示例配置等内容。如果需要补充构建、部署或 API 文档，也可以进一步扩展。