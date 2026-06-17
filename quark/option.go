package quark

// options 是 New 的配置容器。
type options struct {
	cookie     string
	cookieFile string
}

func defaultOptions() *options { return &options{} }

// Option 是 New 的配置函数。
type Option func(*options)

// WithCookie 直接传入完整 cookie 字符串（从浏览器 F12 复制）。
// 格式："__puus=xxx; __pus=yyy; ..."，必须含 __puus。
func WithCookie(cookie string) Option { return func(o *options) { o.cookie = cookie } }

// WithCookieFile 指定 cookie 文件路径。
// 文件可以是 {"cookie":"..."} 格式或纯 cookie 文本。
func WithCookieFile(path string) Option { return func(o *options) { o.cookieFile = path } }
