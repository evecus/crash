# CrashPanel

基于 Go + Vue 3 的透明代理管理面板，支持 Clash Meta (mihomo) 和 sing-box 双核心。

## 功能

- 核心启动/停止/重启控制
- 订阅管理（支持 subconverter 转换）
- DNS 配置（redir-host / fake-ip / mix）
- 防火墙规则（iptables / nftables）
- 计划任务管理
- 实时日志查看
- 全局设置

## 快速开始

### 部署

1. 下载对应架构的二进制文件（Release 页面）
2. 创建配置目录：

```bash
mkdir -p /etc/crashpanel
```

3. 可选：创建配置文件 `/etc/crashpanel/config.json`：

```json
{
  "port": 8080,
  "db_path": "/etc/crashpanel/crashpanel.db",
  "jwt_secret": "your-secret-key-here",
  "admin_password": "your-password",
  "debug": false
}
```

4. 放置核心二进制（mihomo 或 sing-box），默认路径 `/usr/local/bin/mihomo`

5. 启动面板：

```bash
chmod +x crashpanel-linux-amd64
./crashpanel-linux-amd64
```

6. 访问 `http://your-ip:8080`，使用配置的密码登录

### systemd 服务

```ini
# /etc/systemd/system/crashpanel.service
[Unit]
Description=CrashPanel
After=network.target

[Service]
ExecStart=/usr/local/bin/crashpanel
Restart=on-failure
Environment=ADMIN_PASSWORD=yourpassword

[Install]
WantedBy=multi-user.target
```

```bash
systemctl enable --now crashpanel
```

## 环境变量

| 变量               | 说明              | 默认值  |
|--------------------|-------------------|---------|
| `CRASHPANEL_CONFIG`| 配置文件路径      | `/etc/crashpanel/config.json` |
| `ADMIN_PASSWORD`   | 登录密码          | `admin` |
| `JWT_SECRET`       | JWT 签名密钥      | 内置默认值（请务必修改） |

## 从源码构建

```bash
# 需要：Node.js 20+，Go 1.22+，gcc

git clone https://github.com/yourname/crashpanel
cd crashpanel
make all
```

## 目录结构

```
crashpanel/
├── backend/          # Go 后端
│   ├── main.go
│   ├── api/          # 路由
│   ├── handlers/     # HTTP 处理器
│   ├── service/      # 业务逻辑
│   ├── models/       # 数据模型
│   ├── middleware/   # 中间件
│   ├── database/     # 数据库
│   └── config/       # 配置
├── frontend/         # Vue 3 前端
│   └── src/
│       ├── views/    # 页面
│       ├── api/      # API 封装
│       ├── stores/   # Pinia 状态
│       └── layouts/  # 布局
├── Makefile
└── .github/workflows/build.yml
```
