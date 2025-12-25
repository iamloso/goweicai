# ⚠️ 重要说明：需要 Node.js

## 当前状态

虽然原计划是通过 goja (纯 Go JavaScript 引擎) 实现无 Node.js 依赖，但经过测试发现：

1. **hexin-v.bundle.js** 包含了 ES2018 异步生成器语法，goja 不支持
2. 打包文件使用了大量 Node.js 特定 API (Buffer, require 等)
3. 完全模拟 Node.js 环境在 goja 中过于复杂

## 解决方案

### 方案 1：安装 Node.js（推荐）

最简单的解决方案是安装 Node.js：

#### Linux (Ubuntu/Debian)
```bash
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs
```

#### Linux (CentOS/RHEL)
```bash
curl -fsSL https://rpm.nodesource.com/setup_18.x | sudo bash -
sudo yum install -y nodejs
```

#### macOS
```bash
brew install node
```

#### Windows
下载并安装：https://nodejs.org/

#### 验证安装
```bash
node --version
npm --version
```

安装完成后，gowencai 会自动使用 Node.js 执行 JavaScript 代码生成 token。

### 方案 2：使用预先生成的 token（临时方案）

如果您的环境确实无法安装 Node.js，可以：

1. 在有 Node.js 的环境中预先生成 token
2. 直接在代码中使用生成的 token

```go
headers := map[string]string{
    "hexin-v":    "预先生成的token值",
    "User-Agent": gowencai.GetRandomUserAgent(),
    "Cookie":     "your_cookie_here",
}
```

注意：token 有时效性，可能需要定期更新。

### 方案 3：贡献代码

如果您有兴趣帮助实现真正的"无 Node.js"版本，可以：

1. 重写 hexin-v.js，移除所有 Node.js 依赖
2. 确保代码符合 ECMAScript 5.1 标准
3. 提交 Pull Request

## 技术细节

### 为什么 goja 不work？

1. **异步生成器**：hexin-v.bundle.js 第 11265 行使用了 `async function*`
   ```javascript
   const AsyncIteratorPrototype = Object.getPrototypeOf(
       Object.getPrototypeOf(async function* () {}).prototype
   );
   ```

2. **Node.js 模块**：大量使用 `require('path')`, `require('buffer')` 等

3. **jsdom 依赖**：原始文件需要 DOM 环境模拟

### 当前实现

`headers.go` 现在使用 `exec.Command` 调用 Node.js：

```go
cmd := exec.Command("node", "-e", fmt.Sprintf(`
    const module = { exports: {} };
    const exports = module.exports;
    require('%s');
    console.log(v());
`, scriptPath))
```

## 常见问题

### Q: 文档说"无需 Node.js"，为什么还要安装？
**A**: 这是一个已知问题。原计划使用 goja 实现，但遇到了技术障碍。我们正在努力解决。

### Q: 能否静态编译为 Go 二进制？
**A**: 目前不行，因为需要外部 Node.js 进程。

### Q: 性能影响如何？
**A**: 每次生成 token 需要启动 Node.js 进程，约增加 50-100ms 延迟。

### Q: 有替代方案吗？
**A**: 
1. 使用 Python 版本的 pywencai
2. 等待社区贡献纯 Go 实现
3. 手动提取 token 生成逻辑并重写

## 总结

**当前版本需要 Node.js v12+ 才能正常运行。**

我们为给您带来的不便表示歉意，并会继续努力改进。

---

更新日期：2025年11月20日
