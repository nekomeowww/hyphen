# syntax=docker/dockerfile:1

# 设定构建步骤所使用的来源镜像为基于 Alpine 发行版的 Go 1.18 版本镜像
FROM golang:1.19-alpine as builder

ARG VERSION

# 设定 Go 使用 模块化依赖 管理方式：GO111MODULE
RUN GO111MODULE=on

# 创建路径 /app
RUN mkdir /app

# 复制当前目录下 hyphen 到 /app/hyphen
COPY . /app/hyphen

# 切换到 /app/hyphen 目录
WORKDIR /app/hyphen

RUN go env
RUN go env -w CGO_ENABLED=0
RUN go build -a -o "release/hyphen" -ldflags " -X 'github.com/nekomeowww/hyphen/pkg/meta.Version=${VERSION}'" "github.com/nekomeowww/hyphen/cmd/hyphen"

# 设定运行步骤所使用的镜像为基于 Alpine 发行版的 Go 1.18 版本镜像
FROM alpine as runner

# 创建路径 /app
RUN mkdir /app
# 创建路径 /app/hyphen/bin
RUN mkdir -p /app/hyphen/bin
# 创建路径 /var/lib/hyphen
RUN mkdir -p /var/lib/hyphen

COPY --from=builder /app/hyphen/release/hyphen /app/hyphen/bin/

# 映射日志文件路径
VOLUME [ "/app/hyphen/logs" ]

# 入点是编译好的 hyphen 应用程序
ENTRYPOINT [ "/app/hyphen/bin/hyphen", "run", "--data", "/var/lib/hyphen/bbolt_data.db", "-l", "0.0.0.0:9010" ]

# 暴露端口 9010
EXPOSE 9010
