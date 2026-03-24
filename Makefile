.PHONY: all frontend backend clean

# 最终产物
BINARY := crashpanel

all: frontend backend

# 1. 构建前端（输出到 backend/dist）
frontend:
	cd frontend && npm install && npm run build

# 2. 构建后端（embed 前端 dist）
backend:
	cd backend && go build -o ../$(BINARY) .

# 开发模式：前后端分离运行
dev-frontend:
	cd frontend && npm run dev

dev-backend:
	cd backend && go run .

clean:
	rm -rf frontend/node_modules
	rm -rf backend/dist
	rm -f $(BINARY)

# 交叉编译（Linux amd64）
build-linux:
	cd frontend && npm install && npm run build
	cd backend && GOOS=linux GOARCH=amd64 CGO_ENABLED=1 \
		go build -ldflags="-s -w" -o ../$(BINARY)-linux-amd64 .

# 交叉编译（Linux arm64）
build-arm64:
	cd frontend && npm install && npm run build
	cd backend && GOOS=linux GOARCH=arm64 CGO_ENABLED=1 \
		CC=aarch64-linux-gnu-gcc \
		go build -ldflags="-s -w" -o ../$(BINARY)-linux-arm64 .
