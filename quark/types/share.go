package types

// 本文件定义夸克分享相关数据模型。

// CreateShareResponse 创建分享响应。
// 夸克创建分享是异步任务：响应嵌套在 data.task_resp.data 里，
// 关键字段是 share_id，分享链接 = https://pan.quark.cn/s/<share_id>。
type CreateShareResponse struct {
	ShareID string `json:"share_id"`
	// ShareURL 由 SDK 根据 share_id 拼接（响应里不直接返回 URL）。
	ShareURL string `json:"-"`
}

// ExpiredType 常量。
// 抓包确认：1=永久。
// 2/3/4 暂按社区常见约定（1天/7天/30天），如不准以实际为准。
const (
	ExpiredForever int = 1 // 永久有效
	Expired1Day    int = 2 // 1天
	Expired7Day    int = 3 // 7天
	Expired30Day   int = 4 // 30天
)

// ShareURLBase 夸克分享链接前缀。
const ShareURLBase = "https://pan.quark.cn/s/"
