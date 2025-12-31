# Bailu-backend

> A lightweight, flexible, elegant and full-featured web management console based on GIN + JWT + GORM 2.0 + Casbin 2.0 + Wire DI
<div align=center>
<img src="https://img.shields.io/badge/golang-1.23-blue"/>
<img src="https://img.shields.io/badge/gin-1.10.0-lightBlue"/>
<img src="https://img.shields.io/badge/casbin-v2.99.0-9cf"/>
<img src="https://img.shields.io/badge/gorm-1.25.11-red"/>
<img src="https://img.shields.io/badge/wire-0.6.0-green"/>
</div> 

[中文](./README.zh-CN.md) | English

## Feature

- :scroll: Follow the `RESTful API` design specification & interface-based programming specifications
- :house: A simpler project structure and modular design improve the readability and maintainability of the code
- :rocket: Based on the `GIN` framework, it provides rich middleware support (JWTAuth, CORS, RequestLogger, RequestRateLimiter, Casbin, Recovery, OperationRecord, Locale, SSE, GZIP, StaticWebsite)
- :closed_lock_with_key: RBAC access control model based on `Casbin`
- :page_facing_up: Database access layer based on `GORM 2.0`
- :electric_plug: Dependency injection based on `WIRE` -- The role of dependency injection itself is to solve the cumbersome initialization process of hierarchical dependencies between modules
- :memo: Log output based on `Zap`
- :key: User authentication based on `JWT`
- :microscope: Automatically generate `Swagger` documents based on `Swaggo` - [Preview](https://sufo.me:8081/swagger/index.html)

## Front End

- [Front-end project based on Vue.js]() - [Preview](https://sufo.me:3333): sufo/admin123

## Install dependent tools

- [Go](https://golang.org/) 1.19+
- [Wire](github.com/google/wire) `go install github.com/google/wire/cmd/wire@latest`
- [Swag](github.com/swaggo/swag) `go install github.com/swaggo/swag/cmd/swag@latest`

## Download and deploy
1. Clone the project from git 
```shell
    git clone https://github.com/sufo/bailu-admin.git 
```
2. Install MySQL database, create db, and run the init.mysql.sql script under scripts
3. Modify config.yml
```yaml
    datasource:
    dbType: 'mysql'
    mysql:
     driver: mysql
     host : 127.0.0.1
     username: test #Change to your own database user name
     password: 123456 #Change the password to your own database user
 ```
4. Generate Dependency Injection Code
```shell
   wire gen ./app 
```
5. Generate Swagger documentation
```shell
   swag init 
```
6. start
```shell
   go run main.go
```
7. Visit [https://sufo.me:8081/swagger/index.html](https://sufo.me:8081/swagger/index.html) to see the api doc page


## Generate Docker image

```shell
sudo docker build -f ./Dockerfile -t bailu-admin:v1.0.0 .
```


## Project structure overview

```text
├── app
│   ├── api                         (API)
│   │   ├── admin                   (Controller)
│   │   ├── home                    (Client)
│   │   └── api.go                  (API wire)
│   ├── config                            
│   │   └── config.go               (Configuration file Struct)
│   ├── core                              
│   │   ├── appctx                  (App context)
│   │   │   ├── context.go          (context)
│   │   ├── engine.go               (Gin Router)
│   │   └── viper.go                (Viper)
│   ├── domain                              
│   │   ├── dto                     (DTO)
│   │   ├── entity                  (Entity)
│   │   ├── repo                    (Repository)
│   │   ├── resp                    (Response)
│   │   └── vo                      (VO)
│   ├── locales                           
│   │   └── lang                    (LOCALE)
│   ├── middleware                  (Middleware)
│   ├── router                      (Router)
│   ├── service                     (Services)
│   │   ├── base                    (Base service)
│   │   ├── cron                    (Scheduled tasks)
│   │   ├── message                 (Message)
│   │   ├── sys                     (System)
│   │   └── service.go              (service wire)
│   ├── app.go                      (Application startup entry)
│   ├── casbin.go                   (RBAC)
│   ├── injector.go                 (Dependency Injection)
│   ├── wire.go                     (Dependency Injection)
│   └── wire_gen.go                 (Dependency Injection)
├── assets                          (Static resource)
├── cmd                             (CMD)
│   ├── cli                         (cli)
│   ├── admin
│   │   └── api.go                  (API Start)
│   ├── version
│   │   └── version.go              (Version)
├── config
│   ├── config.yml                  (Configuration yaml)
│   ├── menu.yml                    (menu yaml)
│   └── rbac_model.conf             (Casbin RBAC Model)
│── docs                            (swagger Directory)
│── global                          (Global)
│── log                             (Log)
│── pkg                             (Package)
│── utils                           (utilities)
├── Dockerfile
├── Makefile
├── README.md
├── go.mod
├── go.sum
└── main.go                        
```
