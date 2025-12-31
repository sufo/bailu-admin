.PHONY: start build
SHELL = /bin/bash
APP = bailu-server
START_ARGS = --www public

ifeq ($(TAGS_OPT),)
TAGS_OPT = latest
else
endif

start:
	@go run -ldflags "-w -s" -a -o ${APP} main.go start ${START_ARGS}

build:
	@go build -ldflags "-w -s" -a -o ${APP}

# go install github.com/google/wire/cmd/wire@latest
wire:
	@wire gen ./app

# go install github.com/swaggo/swag/cmd/swag@latest
swagger:
	@swag init --parseDependency --generalInfo #./main.go --output ./docs/swagger

stop:
	./$(APP) stop


# 本地打包镜像为压缩文件
save-image-server:build-image-server
	docker save -o ${APP}.tar ${APP}:${TAGS_OPT}

#容器环境打包后端
build-server:
	docker run --name ${APP} --rm make build-server-local

#构建镜像
build-image-server:
	@docker build -t ${APP}:${TAGS_OPT} .

#本地环境打包后端
build-server-local:
	@cd if [ -f "bailu-server" ];then rm -rf ${APP}; else echo "OK!"; fi \
    && go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 && go env  && go mod tidy \
    && go build -ldflags "-w -s" -o ${APP}
