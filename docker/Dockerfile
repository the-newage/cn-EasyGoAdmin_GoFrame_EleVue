# ============================================================
# Dockerfile：EasyGoAdmin 后端镜像（多阶段构建）
# 构建命令（上下文为项目根目录，在项目根执行）：
#   docker build -f docker/Dockerfile -t easygoadmin:latest .
#
# 设计要点
#   - 纯 Go 项目（sqlite 驱动为 modernc 系，无 CGO），CGO_ENABLED=0 静态编译
#   - builder 只负责编译，runtime 仅含二进制 + 运行时资产，镜像最小化
#   - 运行时资产四件：manifest/config、manifest/data(ip2region.xdb)、
#     resource(字体/swagger)、templates(代码生成器模板)
#   - 敏感配置不改镜像内文件，运行时用 GF_ 前缀环境变量覆盖
#     （见 manifest/config/config.yaml 头部说明）
#   - 源码 COPY 放在依赖下载之后，保护 go mod 层缓存
#
# 构建参数
#   VAR_PORT       应用监听端口
# ============================================================


# ------------------------------------------------------------
# 全局构建参数
#   在 FROM 之前声明，可供后续各阶段引用
# ------------------------------------------------------------
ARG VAR_GOLANG_IMAGE=golang:1.27-alpine
ARG VAR_RUNTIME_IMAGE=alpine:3.22
ARG VAR_PORT=8000


# ============================================================
# 阶段一：builder —— 编译静态二进制
# ============================================================
FROM ${VAR_GOLANG_IMAGE} AS builder

LABEL MAINTAINER="easygoadmin@163.com"

# 国内代理：加速 go mod 下载（海外环境可改为 https://proxy.golang.org,direct）
ENV GOPROXY=https://goproxy.cn,direct \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /build

# 先拷依赖清单，利用层缓存：源码变更但 go.mod/go.sum 未变时不重复下载
COPY go.mod go.sum ./
RUN go mod download

# 拷源码并编译（-trimpath 去除本地路径，-s -w 去除符号表与调试信息）
COPY . .
RUN go build -trimpath -ldflags "-s -w" -o /build/bin/admin-server main.go


# ============================================================
# 阶段二：runtime —— 最终业务镜像
# ============================================================
FROM ${VAR_RUNTIME_IMAGE} AS runtime

LABEL MAINTAINER="easygoadmin@163.com"

# 时区（应用配置 app.timezone=Asia/Shanghai，日志与定时任务依赖系统时区）
ENV TZ=Asia/Shanghai

# tzdata 提供时区数据库；ca-certificates 供 HTTPS 外呼（无外呼场景可去掉）
RUN set -eux; \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories; \
    apk add --no-cache tzdata ca-certificates; \
    # 非 root 运行（低权限用户跑业务）
    addgroup -g 1000 app && adduser -u 1000 -G app -s /sbin/nologin -D app; \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /opt/easygoadmin

# 运行时资产（先于二进制 COPY，业务代码变动不破坏资产层缓存）
COPY --chown=app:app manifest/config ./manifest/config
COPY --chown=app:app manifest/data   ./manifest/data
COPY --chown=app:app resource        ./resource
COPY --chown=app:app templates       ./templates

# 二进制
COPY --chown=app:app --from=builder /build/bin/admin-server ./bin/admin-server

# 上传目录与日志目录：预建并授权，compose 挂载宿主卷时属主一致
RUN mkdir -p /data/upload /opt/easygoadmin/log && \
    chown -R app:app /data/upload /opt/easygoadmin/log

USER app

# ARG 转 ENV，运行期与 compose/HEALTHCHECK 可读
ENV VAR_PORT=${VAR_PORT}
EXPOSE ${VAR_PORT}

# 健康检查（/api/health 为公开白名单路径，不带 /v1）
HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
    CMD wget -qO- "http://127.0.0.1:${VAR_PORT}/api/health" || exit 1

# 前台运行；优雅重启由框架 graceful 配置处理
ENTRYPOINT ["./bin/admin-server"]
