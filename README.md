# quark-go

夸克网盘的 Go SDK。基于 cookie 鉴权（无 token、无签名、无扫码登录 API）。

> 开发中。当前已完成：骨架、cookie 鉴权、文件列表/详情。后续：转存分享（核心）、文件管理、上传、创建分享、下载。

## 安装

```sh
go get github.com/langhuachuanshi/quark-go
```

## 获取 Cookie

夸克只能用 cookie 鉴权（无官方扫码 API）：

1. 浏览器登录 https://pan.quark.cn
2. 按 F12 → Network 标签
3. 刷新页面，点击任意一个请求
4. 在 Request Headers 里找到 `Cookie:`，复制完整值（`__puus=xxx; __pus=yyy; ...`）
5. 保存到 `~/.quark/cookie.json`，内容 `{"cookie":"复制的完整cookie"}`

## 快速入门

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/langhuachuanshi/quark-go/quark"
    "github.com/langhuachuanshi/quark-go/quark/file"
)

func main() {
    ctx := context.Background()
    // 默认从 ~/.quark/cookie.json 读 cookie；也可 quark.WithCookie("...")
    c, err := quark.New(ctx)
    if err != nil { log.Fatal(err) }

    files, _ := c.Files().List(ctx, &file.ListRequest{PDirFid: "0"})
    for _, f := range files {
        fmt.Println(f.FileName, f.Size)
    }
}
```

## 目录结构

```
quark-go/
├── quark/          主包 Client（持 cookie，实现 invoker）
├── quark/auth/     cookie 管理 + 持久化
├── quark/file/     文件列表/详情/管理
├── quark/share/    转存他人分享（核心）+ 创建分享
├── quark/upload/   上传（秒传+分片+断点续传）
├── quark/download/ 下载
├── quark/types/    数据模型
└── quark/invoker/  共享调用接口
```

## 夸克 vs 阿里云盘

| 维度 | 阿里云盘 | 夸克网盘 |
|---|---|---|
| 鉴权 | token（扫码） | cookie（浏览器复制） |
| 签名 | 需 x-signature | 不需要 |
| 核心场景 | 上传+分享 | **转存他人分享** |

## License

仅供学习交流。
