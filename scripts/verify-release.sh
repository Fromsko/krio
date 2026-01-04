#!/bin/bash

# 验证发布配置脚本

echo "========================================="
echo "  验证发布配置"
echo "========================================="
echo ""

# 1. 检查 .version 文件
echo "1. 检查版本文件..."
if [ -f ".version" ]; then
    VERSION=$(cat .version)
    echo "   ✅ 版本号: $VERSION"
else
    echo "   ❌ .version 文件不存在"
    exit 1
fi

# 2. 检查 GitHub Actions 工作流
echo ""
echo "2. 检查 GitHub Actions 工作流..."
if [ -f ".github/workflows/release.yml" ]; then
    echo "   ✅ release.yml 存在"
else
    echo "   ❌ release.yml 不存在"
    exit 1
fi

# 3. 检查 app 包
echo ""
echo "3. 检查 app 包..."
if [ -f "app/version.go" ]; then
    echo "   ✅ app/version.go 存在"
else
    echo "   ❌ app/version.go 不存在"
    exit 1
fi

# 4. 检查敏感信息
echo ""
echo "4. 检查敏感信息..."
if grep -r "44a695c982014" config.example.yaml 2>/dev/null; then
    echo "   ❌ config.example.yaml 中包含敏感信息"
    exit 1
else
    echo "   ✅ config.example.yaml 无敏感信息"
fi

if grep -r "D:/notes/Fromsko" config.example.yaml 2>/dev/null; then
    echo "   ❌ config.example.yaml 中包含本地路径"
    exit 1
else
    echo "   ✅ config.example.yaml 无本地路径"
fi

# 5. 检查 config.yaml 在 gitignore 中
echo ""
echo "5. 检查 .gitignore..."
if grep -q "^config.yaml$" .gitignore; then
    echo "   ✅ config.yaml 在 .gitignore 中"
else
    echo "   ❌ config.yaml 不在 .gitignore 中"
    exit 1
fi

# 6. 测试构建
echo ""
echo "6. 测试构建..."
if command -v go &> /dev/null; then
    echo "   正在构建..."
    if go build -o /tmp/krio-test . 2>/dev/null; then
        echo "   ✅ 构建成功"
        rm -f /tmp/krio-test
    else
        echo "   ❌ 构建失败"
        exit 1
    fi
else
    echo "   ⚠️  Go 未安装，跳过构建测试"
fi

echo ""
echo "========================================="
echo "  ✅ 所有检查通过！"
echo "========================================="
echo ""
echo "发布新版本："
echo "  echo '1.1.0' > .version"
echo "  git add .version"
echo "  git commit -m 'bump: release v1.1.0'"
echo "  git push origin main"
echo ""
