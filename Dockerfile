# 指定构建镜像的基础镜像
FROM golang:alpine AS builder

LABEL MAINTAINER=ouamour@gmail.com

ARG APP=bailu-server
ARG VERSION=v0.1.0

# ENV GOPROXY="https://goproxy.cn"

# 声明工作目录
WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

#把当前目录的文件拷过去，编译代码
COPY . .

#RUN go generate

# build the application
# [-ldflags](https://pkg.go.dev/cmd/link)
# RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -o bailu-admin main.go


RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s -X main.Version=1.0.0" -o ${APP} main.go


FROM alpine
ARG APP=bailu-server
WORKDIR /bailu
COPY --from=builder /build/${APP} ./
COPY --from=builder /build/config ./config
COPY --from=builder /build/app/locales ./app/locales
COPY --from=builder /build/sql  ./sql
COPY --from=builder /build/public  ./public
COPY --from=builder /build/assets  ./assets

ENTRYPOINT ["./bailu-server", "start" ,"--www", "public", "-c", "config/config.docker.yml"]
EXPOSE 8081