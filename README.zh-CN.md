# Bailu 后端

<p align="center">
  <strong>一个轻量级、生产就绪、功能丰富的后端样板项目。</strong>
</p>
<p align="center">
  Bailu 使用 Go、Gin、GORM 和 Wire 构建，旨在帮助您快速启动安全且可扩展的管理面板、RESTful API 或微服务。
</p>

<p align="center">
  <img src="https://img.shields.io/badge/go-1.21+-blue.svg" alt="Go Version">
  <img src="https://img.shields.io/badge/gin-v1.10.0-blue.svg" alt="Gin Version">
  <img src="https://img.shields.io/badge/gorm-v1.25.11-orange.svg" alt="Gorm Version">
  <img src="https://img.shields.io/badge/casbin-v2.99.0-green.svg" alt="Casbin Version">
  <img src="https://img.shields.io/badge/wire-v0.6.0-purple.svg" alt="Wire Version">
  <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License">
</p>

<div align="center">
  <strong><a href="./README.zh-CN.md">中文</a></strong> | <strong>English</strong>
</div>

---

## ⚙️ 技术栈

-   **Web 框架**: Gin
-   **ORM**: GORM
-   **数据库**: MySQL
-   **依赖注入**: Google Wire
-   **权限管理**: Casbin
-   **认证**: JWT
-   **配置管理**: Viper
-   **日志**: Zap

## ✨ 用法示例

Bailu 旨在提高开发人员的生产力。以下是一些其便捷功能的示例：

### 1. 轻松分页

在您的控制器中添加一行代码即可实现任何查询的分页功能。

```go
// 在您的 API 控制器函数中：
func (a *UserAPI) GetUserList(c *gin.Context) {
    // 自动应用查询参数中的 limit 和 offset (例如：?page=1&pageSize=10)
    page.StartPage(c) 
    
    // 您的数据检索逻辑
    users, total, err := a.userService.ListByBuilder(c)
    if err != nil {
        resp.FailWithError(c, err)
        return
    }
    
    // 响应会自动封装为分页结构
    resp.OKWithData(c, page.New(users, total))
}
```

### 2. 统一响应封装

使用简单、富有表现力的辅助函数标准化您的 API 响应。

```go
// 成功响应：
resp.OK(c) // 返回标准成功消息
resp.OKWithData(c, data) // 返回成功消息和数据载荷

// 错误响应：
resp.Fail(c) // 返回标准失败消息
resp.FailWithError(c, someError) // 返回带有特定错误的失败消息

// 对于更复杂的场景，您可以使用 panic 抛出响应错误。
// 全局恢复中间件将捕获它并格式化 JSON 响应。
if user == nil {
    panic(resp.ErrNotFound)
}
if err != nil {
    panic(resp.InternalServerErrorWithError(err))
}
```

### 3. 自动查询构建器

直接使用结构体标签从请求 DTO 构建复杂的 GORM 查询。这消除了样板式的 `db.Where()` 子句。

```go
// 1. 在您的 DTO 中使用 `query` 标签定义查询参数。
//    格式：`query:"[列名],[运算符]"`
//    支持的运算符：eq, neq, gt, gte, lt, lte, like, in
type UserQueryParams struct {
    dto.Pagination
    Username string `form:"username" query:"username,like"`
    Email    string `form:"email" query:"email,eq"`
    Status   int    `form:"status" query:"status,eq"`
}

// 2. 在您的仓库层中使用 QueryBuilder。
func (r *UserRepo) FindByParams(ctx context.Context, params *dto.UserQueryParams) ([]*entity.User, error) {
    // 构建器自动构建 WHERE 子句。
    // 例如：WHERE username LIKE '%...%' AND status = ?
    builder := base.NewQueryBuilder().WithWhereStruct(params)
    
    var users []*entity.User
    err := r.FindByBuilder(ctx, builder).Find(&users).Error
    return users, err
}
```

## 前端项目

- **Bailu Admin (Vue)**: 配套的前端项目正在开发中。（链接待添加）
- **在线演示**: （链接待添加）
- **默认凭证**: `sufo` / `admin123`

## 🚀 快速开始

按照以下步骤在本地开发环境中启动并运行。

### 前置条件

- [Go](https://golang.org/dl/) 1.21+
- [MySQL](https://www.mysql.com/downloads/) 5.7+
- [Make](https://www.gnu.org/software/make/)
- [Wire](https://github.com/google/wire): `go install github.com/google/wire/cmd/wire@latest`
- [Swag](https://github.com/swaggo/swag): `go install github.com/swaggo/swag/cmd/swag@latest`

### 安装与运行

1.  **克隆仓库：**
    ```shell
    git clone https://github.com/sufo/bailu-backend.git
    cd bailu-backend
    ```

2.  **配置应用程序：**
    -   复制开发配置文件：`cp config/config.dev.yml config/config.yml`。
    -   编辑 `config/config.yml` 并使用您的数据库凭据更新 `mysql` 部分。

3.  **初始化数据库：**
    -   在 MySQL 中创建一个新数据库（例如 `bailu`）。
    -   从 `sql/init_mysql.sql` 导入初始架构和数据。

4.  **生成依赖注入代码：**
    ```shell
    make wire
    ```

5.  **生成 API 文档：**
    ```shell
    make swagger
    ```

6.  **运行服务器：**
    ```shell
    make start
    ```
    服务器将在您的配置中指定的端口上启动（默认：`8081`）。

7.  **访问 API 文档：**
    访问 `http://localhost:8081/swagger/index.html` 查看交互式 API 文档。

## 🐳 Docker 快速启动

1.  **构建 Docker 镜像：**
    ```shell
    make build-image-server TAGS_OPT=latest
    ```

2.  **运行容器：**
    请确保您的 `config/config.docker.yml` 已正确配置以连接到您的数据库。
    ```shell
    docker run -d -p 8081:8081 --name bailu-server bailu-server:latest
    ```

## 🧰 Makefile 命令

本项目使用 `make` 来简化常见任务。

- `make start`: 以开发模式启动应用程序。
- `make build`: 构建应用程序二进制文件。
- `make wire`: 生成依赖注入代码。
- `make swagger`: 生成 Swagger API 文档。
- `make stop`: 停止正在运行的应用程序。
- `make build-image-server`: 构建 Docker 镜像。

## 📂 项目结构

项目遵循模块化、分层的架构，以促进职责分离和可维护性。

```
/
├── app/                # 核心应用程序代码
│   ├── api/            # API 控制器和路由
│   ├── config/         # 配置结构体
│   ├── core/           # 核心组件 (服务器引擎, DI)
│   ├── domain/         # 领域模型 (实体, DTOs, 仓库)
│   ├── middleware/     # Gin 中间件
│   ├── service/        # 业务逻辑层
│   └── ...
├── config/             # 配置文件 (YAML 等)
├── global/             # 全局变量和常量
├── pkg/                # 共享工具包
├── sql/                # SQL 初始化脚本
├── utils/              # 通用工具函数
├── main.go             # 应用程序入口
├── go.mod              # Go 模块定义
├── Makefile            # 常用任务的 Makefile
└── Dockerfile          # Docker 构建定义
```

## 📄 许可证

本项目采用 [MIT](./LICENSE) 许可证。
