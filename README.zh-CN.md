# Bailu Admin
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/sufo/bailu-admin)

基于 Gin + Vue3 + Naive UI 的前后端分离通用管理台系统。


<div align="center">
  <strong>中文</strong> | <strong><a href="./README.md">English</a></strong>
</div>

---

## ✨ 功能特性

-   **高性能 API**: 基于 Gin 框架构建，提供快速高效的路由。
-   **RESTful API 设计**: 遵循 RESTful 原则，实现清晰、标准化和可扩展的 API 端点。
-   **灵活的数据访问**: 利用 GORM 实现强大且对开发者友好的数据库交互。
-   **基于角色的访问控制 (RBAC)**: 集成 Casbin 进行细粒度的权限管理。
-   **JWT 认证**: 使用 JWT 进行安全的 API 认证，并结合 Redis 进行令牌管理。
-   **简洁的架构**: 采用 Google Wire 实现编译时依赖注入，促进代码的模块化和可维护性。
-   **标准化的 JSON 响应**: 提供辅助函数以实现一致且可预测的 API 响应。
-   **API 文档**: 使用 Swagger (OpenAPI) 生成交互式 API 文档。
-   **动态查询生成**: 从请求 DTO 自动构建数据库查询，减少样板代码。
-   **结构化日志**: 使用 Zap 实现高性能的结构化日志记录。
-   **灵活的配置**: 通过 Viper 管理应用配置，支持多种文件格式。


## ✨ 使用示例

Bailu 旨在提高开发效率。以下是一些便捷功能的使用示例：

### 1. 轻松分页

只需在控制器中添加一行代码，即可为任何查询启用分页。

```go
// 在您的 API 控制器函数中：
func (a *UserAPI) GetUserList(c *gin.Context) {
    // 从查询参数自动应用 limit 和 offset (例如, ?page=1&pageSize=10)
    page.StartPage(c) 
    
    // 您的数据检索逻辑
    users, err := a.userService.ListByBuilder(c)
    if err != nil {
        resp.FailWithError(c, err)
        return
    }
    
    // 响应被自动包装在分页结构中
    resp.OKWithData(c, users)
}
```

### 2. 统一响应包装

使用简单、富有表现力的辅助函数来标准化您的 API 响应。

```go
// 成功响应：
resp.OK(c) // 返回标准成功消息
resp.OKWithData(c, data) // 返回带数据负载的成功响应

// 失败响应：
resp.Fail(c) // 返回标准失败消息
resp.FailWithError(c, someError) // 返回带特定错误的失败响应

// 对于更复杂的场景，您可以 panic 一个响应错误。
// 全局恢复中间件将捕获它并格式化 JSON 响应。
if user == nil {
    panic(resp.ErrNotFound)
}
if err != nil {
    panic(resp.InternalServerErrorWithError(err))
}
```

### 3. 自动查询构建

使用结构体标签直接从您的请求 DTO 构建复杂的 GORM 查询。这消除了冗长的 `db.Where()` 子句。

```go
// 1. 在您的 DTO 中使用 `query` 标签定义查询参数。
//    格式: `query:"[column_name],[operator]"`
//    支持的操作符: eq, neq, gt, gte, lt, lte, like, in
type UserQueryParams struct {
    dto.Pagination
    Username string `form:"username" query:"username,like"`
    Email    string `form:"email" query:"email,eq"`
    Status   int    `form:"status" query:"status,eq"`
}

// 2. 在您的仓库层使用 QueryBuilder。
func (r *UserRepo) FindByParams(ctx context.Context, params *dto.UserQueryParams) ([]*entity.User, error) {
    // 构建器会自动构造 WHERE 子句。
    // 例如, WHERE username LIKE '%...%' AND status = ?
    builder := base.NewQueryBuilder().WithWhereStruct(params)
    
    var users []*entity.User
    err := r.FindByBuilder(ctx, builder).Find(&users).Error
    return users, err
}
```

## 前端项目

- **Bailu Admin (Vue)**: 配套的前端项目正在开发中。（链接待添加）
- **在线演示**: （链接待添加）
- **默认凭据**: `sufo` / `admin123`

## 🚀 快速开始

按照以下步骤在本地启动和运行开发环境。

### 环境准备

- [Go](https://golang.org/dl/) 1.21+
- [MySQL](https://www.mysql.com/downloads/) 5.7+
- [Make](https://www.gnu.org/software/make/)
- [Wire](https://github.com/google/wire): `go install github.com/google/wire/cmd/wire@latest`
- [Swag](https://github.com/swaggo/swag): `go install github.com/swaggo/swag/cmd/swag@latest`

### 安装与运行

1.  **克隆仓库:**
    ```shell
    git clone https://github.com/sufo/bailu-admin.git
    cd bailu-admin
    ```

2.  **配置应用:**
    -   复制开发配置文件：`cp config/config.dev.yml config/config.yml`。
    -   编辑 `config/config.yml` 并更新 `mysql` 部分的数据库凭据。

3.  **初始化数据库:**
    -   在 MySQL 中创建一个新数据库 (例如, `bailu`)。
    -   从 `sql/init_mysql.sql` 导入初始结构和数据。

4.  **生成依赖注入代码:**
    ```shell
    make wire
    ```

5.  **生成 API 文档:**
    ```shell
    make swagger
    ```

6.  **运行服务:**
    ```shell
    make start
    ```
    服务将在您配置的端口上启动 (默认为 `8081`)。

7.  **访问 API 文档:**
    访问 `http://localhost:8081/swagger/index.html` 查看可交互的 API 文档。

## 🐳 Docker 快速启动

1.  **构建 Docker 镜像:**
    ```shell
    make build-image-server TAGS_OPT=latest
    ```

2.  **运行容器:**
    确保您的 `config/config.docker.yml` 已正确配置以连接到您的数据库。
    ```shell
    docker run -d -p 8081:8081 --name bailu-server bailu-server:latest
    ```

## 🧰 Makefile 命令

本项目使用 `make` 来简化常用任务。

- `make start`: 在开发模式下启动应用。
- `make build`: 构建应用二进制文件。
- `make wire`: 生成依赖注入代码。
- `make swagger`: 生成 Swagger API 文档。
- `make stop`: 停止正在运行的应用。
- `make build-image-server`: 构建 Docker 镜像。

## 📂 项目结构

项目遵循模块化、分层的架构，以促进关注点分离和可维护性。

```
/
├── app/                # 核心应用代码
│   ├── api/            # API 控制器和路由
│   ├── config/         # 配置结构体
│   ├── core/           # 核心组件 (服务引擎, 依赖注入)
│   ├── domain/         # 领域模型 (实体, DTO, 仓库)
│   ├── middleware/     # Gin 中间件
│   ├── service/        # 业务逻辑层
│   └── ...
├── config/             # 配置文件 (YAML, 等)
├── global/             # 全局变量和常量
├── pkg/                # 共享的工具包
├── sql/                # SQL 初始化脚本
├── utils/              # 通用工具函数
├── main.go             # 应用入口点
├── go.mod              # Go 模块定义
├── Makefile            # Makefile 用于简化常用任务
└── Dockerfile          # Docker 构建定义
```

## 📄 许可证

本项目采用 [MIT](./LICENSE) 许可证。