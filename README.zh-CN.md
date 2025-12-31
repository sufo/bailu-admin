# Bailu-backend

> 基于 GIN + JWT + GORM 2.0 + Casbin 2.0 + Wire DI 的轻量级、灵活、优雅且功能齐全的web管理台
<div align=center>
<img src="https://img.shields.io/badge/golang-1.23-blue"/>
<img src="https://img.shields.io/badge/gin-1.10.0-lightBlue"/>
<img src="https://img.shields.io/badge/casbin-v2.99.0-9cf"/>
<img src="https://img.shields.io/badge/gorm-1.25.11-red"/>
<img src="https://img.shields.io/badge/wire-0.6.0-green"/>
</div>

[English](./README.md) | 中文

## 功能特性

- :scroll: 遵循 `RESTful API` 设计规范 & 基于接口的编程规范
- :house: 更加简洁的项目结构，模块化的设计，提高代码的可读性和可维护性
- :rocket: 基于 `GIN` 框架，提供了丰富的中间件支持（JWTAuth, CORS, RequestLogger, RequestRateLimiter, Casbin, Recovery, OperationRecord, Locale, SSE, GZIP, StaticWebsite）
- :closed_lock_with_key: 基于 `Casbin` 的 RBAC 访问控制模型
- :page_facing_up: 基于 `GORM 2.0` 的数据库访问层
- :electric_plug: 基于 `WIRE` 的依赖注入 -- 依赖注入本身的作用是解决了各个模块间层级依赖繁琐的初始化过程
- :memo: 基于 `Zap` 实现了日志输出
- :key: 基于 `JWT` 的用户认证
- :microscope: 基于 `Swaggo` 自动生成 `Swagger` 文档 - [预览](https://sufo.me:8081/swagger/index.html)

## 前端项目

- [基于 Vue.js 实现的前端项目]() - [预览](https://sufo.me:3000/): sufo/admin123

## 安装依赖工具

- [Go](https://golang.org/) 1.19+
- [Wire](github.com/google/wire) `go install github.com/google/wire/cmd/wire@latest`
- [Swag](github.com/swaggo/swag) `go install github.com/swaggo/swag/cmd/swag@latest`

## 下载部署
1. 从git下载项目
```shell
    git clone https://github.com/sufo/bailu-admin.git 
```
2. 安装mysql数据库，创建db，运行scripts下init.mysql.sql脚本
3. 修改config.yml
```yaml
   datasource:
   dbType: 'mysql'
   mysql:
     driver: mysql
     host : 127.0.0.1
     username: test #修改为自己数据库用户名
     password: 123456 #修改为自己数据库用户密码
 ```
4. 生成依赖注入代码
```shell
   wire gen ./app
```
5. 生成 Swagger 文档
```shell
   swag init 
```
6. 启动项目
```shell
   go run main.go
```
7. 访问[http://localhost:8081/swagger/index.html](http://localhost:8081/swagger/index.html)即可看到接口页面


## 生成 Docker 镜像

```shell
sudo docker build -f ./Dockerfile -t bailu-admin:v1.0.0 .
```

## 项目结构概览

```text
├── app
│   ├── api                         (API层)
│   │   ├── admin                   (管理台控制器)
│   │   ├── home                    (客户端)
│   │   └── api.go                  (API wire)
│   ├── config                            
│   │   └── config.go               (配置文件结构体)
│   ├── core                              
│   │   ├── appctx                  (app context)
│   │   │   ├── context.go          (context)
│   │   ├── engine.go               (gin路由)
│   │   └── viper.go                (viper)
│   ├── domain                              
│   │   ├── dto                     (数据传输对象)
│   │   ├── entity                  (数据库实体模型)
│   │   ├── repo                    (持久话层)
│   │   ├── resp                    (response响应对象)
│   │   └── vo                      (视图对象)
│   ├── locales                           
│   │   └── lang                    (语言文件)
│   ├── middleware                  (中间件)
│   ├── router                      (路由)
│   ├── service                     (服务层)
│   │   ├── base                    (基础服务)
│   │   ├── cron                    (定时任务)
│   │   ├── message                 (消息)
│   │   ├── sys                     (系统服务)
│   │   └── service.go              (service wire)
│   ├── app.go                      (应用启动入口)
│   ├── casbin.go                   (RBAC 模块)
│   ├── injector.go                 (依赖注入)
│   ├── wire.go                     (依赖注入)
│   └── wire_gen.go                 (依赖注入)
├── assets                          (静态资源文件)
├── cmd                             (命令行定义目录)
│   ├── cli                         (启动入口)
│   ├── admin
│   │   └── api.go                  (API启动命令)
│   ├── version
│   │   └── version.go              (版本命令)
├── config
│   ├── config.yml                  (系统配置文件)
│   ├── menu.yml                    (初始化菜单文件)
│   └── rbac_model.conf             (Casbin RBAC 模型配置文件)
│── docs                            (swagger文档目录)
│── global                          (全局常量和函数)
│── log                             (日志目录)
│── pkg                             (扩展功能包)
│── utils                           (工具类)
├── Dockerfile
├── Makefile
├── README.md
├── go.mod
├── go.sum
└── main.go                         (入口文件)
```
