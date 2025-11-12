# ME-Chain Docker 私有网络

Docker 相关文件已整理到 **`docker/`** 目录中。

## 📁 快速访问

- **快速开始**: [docker/README.md](docker/README.md)
- **目录索引**: [docker/INDEX.md](docker/INDEX.md)
- **完整教程**: [docker/QUICKSTART.md](docker/QUICKSTART.md)
- **快速参考**: [docker/QUICKREF.md](docker/QUICKREF.md)

## 🚀 三步启动

```bash
make docker-private-net        # 1. 构建镜像
make docker-private-net-start  # 2. 启动网络
make docker-private-net-test   # 3. 运行测试
```

## 📚 文档结构

```
docker/
├── README.md              # 快速入门
├── INDEX.md               # 目录索引
├── QUICKSTART.md          # 详细教程
├── QUICKREF.md            # 快速参考
├── SUMMARY.md             # 部署总结
├── BUILD_REPORT.md        # 构建报告
├── Dockerfile             # 镜像定义
├── docker-compose.yml     # Compose 配置
├── scripts/               # 相关脚本
│   ├── setup_local_docker.sh
│   ├── start_private_net.sh
│   ├── test_private_net.sh
│   ├── private_net_demo.sh
│   └── prepare_docker_build.sh
└── docs/                  # 详细文档
    ├── DOCKER_PRIVATE_NET.md
    └── SETUP_COMPARISON.md
```

## 🎯 主要特性

- ✅ 一键启动单节点私有网络
- ✅ 完全自动化，无需交互
- ✅ Docker 环境完全隔离
- ✅ 预配置账户和验证者
- ✅ 14项自动化测试

## 📡 访问端点

| 服务 | 地址 |
|------|------|
| RPC | http://localhost:36657 |
| REST API | http://localhost:1318 |
| JSON-RPC | http://localhost:9545 |
| gRPC | localhost:8090 |

## 💡 更多信息

详细文档请查看 [docker/](docker/) 目录。

---

**提示**: 首次构建约需 6 分钟，后续使用缓存仅需 1-2 分钟。
