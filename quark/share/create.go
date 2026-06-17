package share

import (
	"context"

	"github.com/langhuachuanshi/quark-go/quark/invoker"
	"github.com/langhuachuanshi/quark-go/quark/types"
)

// 本文件实现夸克网盘创建分享。
//
// 接口：POST /share（创建分享链接）
// 夸克的文件分享要求文件在网盘里（和阿里云盘"必须资源盘"不同，夸克无此限制）。

// CreateRequest 创建分享的便捷请求（对外）。
type CreateRequest struct {
	FIDs        []string // 文件/文件夹 fid 列表（必填）
	Title       string   // 分享标题（可选）
	Forever     bool     // 是否永久有效（true=永久，false 用 ExpiredDays）
	ExpiredDays int      // 有效天数（Forever=false 时生效，如 1/7/30）
	WithPasscode bool    // 是否带提取码（私密分享）
	Passcode    string   // 提取码（WithPasscode=true 时用，空则自动生成）
}

// Create 创建分享链接，返回 https://pan.quark.cn/s/xxx。
//
// 支持永久/限时（1/7/30天）、公开/私密（带提取码）。
//
// 抓包确认：POST /share，body 仅含必要字段（公开分享不传 passcode）。
// 真实请求示例：{"fid_list":["xxx"],"title":"test","url_type":1,"expired_type":1}
func (s *Service) Create(ctx context.Context, req *CreateRequest) (*types.CreateShareResponse, error) {
	if req == nil || len(req.FIDs) == 0 {
		return nil, invoker.NewAPIError(0, "fids is required")
	}

	// 转成夸克 API 的 expired_type。
	expiredType := types.ExpiredForever
	if !req.Forever {
		switch req.ExpiredDays {
		case 1:
			expiredType = types.Expired1Day
		case 7:
			expiredType = types.Expired7Day
		case 30:
			expiredType = types.Expired30Day
		default:
			expiredType = types.Expired30Day
		}
	}

	// 按需构造 body，避免多余字段（抓包确认公开分享无 passcode）。
	body := map[string]any{
		"fid_list":     req.FIDs,
		"title":        req.Title,
		"url_type":     1,
		"expired_type": expiredType,
	}
	if req.WithPasscode {
		body["url_type"] = 2
		if req.Passcode != "" {
			body["passcode"] = req.Passcode
		}
	}

	// 夸克创建分享是异步任务：响应嵌套 data.task_resp.data.share_id。
	// 分享链接不在响应里，由 share_id 拼接：https://pan.quark.cn/s/<share_id>。
	var wrapper struct {
		Code int    `json:"code"`
		Msg  string `json:"message"`
		Data struct {
			TaskID   string `json:"task_id"`
			TaskResp struct {
				Code int    `json:"code"`
				Msg  string `json:"message"`
				Data struct {
					ShareID string `json:"share_id"`
					Status  int    `json:"status"` // 2=完成
				} `json:"data"`
			} `json:"task_resp"`
		} `json:"data"`
	}
	if err := invoker.PostAndDecode(ctx, s.inv, "/share", body, nil, &wrapper); err != nil {
		return nil, err
	}
	if wrapper.Code != 0 {
		return nil, invoker.NewAPIError(wrapper.Code, wrapper.Msg)
	}
	// task_resp 可能失败。
	if wrapper.Data.TaskResp.Code != 0 {
		return nil, invoker.NewAPIError(wrapper.Data.TaskResp.Code, wrapper.Data.TaskResp.Msg)
	}
	shareID := wrapper.Data.TaskResp.Data.ShareID
	if shareID == "" {
		return nil, invoker.NewAPIError(0, "创建分享成功但未返回 share_id")
	}
	return &types.CreateShareResponse{
		ShareID:  shareID,
		ShareURL: types.ShareURLBase + shareID,
	}, nil
}
