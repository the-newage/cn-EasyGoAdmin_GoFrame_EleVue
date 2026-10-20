# ======================================================================
# EasyGoAdmin_GoFrame_EleVue 常用命令
# 前置：go install github.com/gogf/gf/cmd/gf/v2@latest
# ======================================================================
GF    ?= gf
APP   ?= admin-server
BIN   ?= bin/$(APP).exe

.PHONY: run build test vet fmt gen gen-dao gen-service clean

## 本地运行
run:
	go run main.go

## 编译输出
build:
	go build -o $(BIN) main.go

## 单元测试（全部）；单测指定包：go test ./internal/logic/...
test:
	go test ./...

## 静态检查
vet:
	go vet ./...

## 格式化
fmt:
	gofmt -w .

## 重新生成 dao + service（表结构变更后执行）
gen: gen-dao gen-service

gen-dao:
	$(GF) gen dao -c hack/config.yaml

gen-service:
	$(GF) gen service -c hack/config.yaml

clean:
	if exist $(BIN) del /f $(BIN)
