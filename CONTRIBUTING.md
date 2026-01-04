# 贡献指南

感谢你对 Krio 项目的关注！

## 如何贡献

### 报告问题

请在 [Issues](https://github.com/fromsko/krio/issues) 页面报告 bug 或提出功能建议。

### 提交代码

1. Fork 本仓库
2. 创建你的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启一个 Pull Request

### 开发流程

```bash
# 1. 克隆仓库
git clone https://github.com/fromsko/krio.git
cd krio

# 2. 安装依赖
make deps

# 3. 开发模式 (热重载)
make dev

# 4. 运行测试
make test

# 5. 构建 CLI 工具
make cli

# 6. 代码检查
make lint
```

## 发布流程

Krio 使用自动化发布系统。修改 `.version` 文件即可触发新版本发布：

```bash
echo "1.1.0" > .version
git add .version
git commit -m "bump: release v1.1.0"
git push origin main
```

详细的发布流程请参考 [RELEASE.md](RELEASE.md)。

## 代码规范

- 遵循 Go 语言官方代码风格
- 运行 `make fmt` 格式化代码
- 运行 `make lint` 进行代码检查
- 确保所有测试通过 (`make test`)

## 版本号规范

遵循 [语义化版本 2.0.0](https://semver.org/lang/zh-CN/)：

- **主版本号 (MAJOR)**: 不兼容的 API 变更
- **次版本号 (MINOR)**: 向下兼容的功能新增
- **修订号 (PATCH)**: 向下兼容的问题修复

## 许可证

通过提交代码，你同意你的代码将根据项目的许可证进行授权。
