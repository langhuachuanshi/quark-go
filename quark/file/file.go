// Package file 实现夸克网盘文件相关操作。
//
// 夸克文件列表接口：GET /file/sort
// 参数：pdir_fid（父目录，根目录"0"）、_page（页码从1开始）、_size（每页数，默认50）、
//      _fetch_total（是否返回总数，1）、_sort（排序，如 "file_type:asc,updated_at:desc"）
package file

import (
	"context"
	"fmt"

	"github.com/langhuachuanshi/quark-go/quark/invoker"
	"github.com/langhuachuanshi/quark-go/quark/types"
)

// Service 文件相关操作入口。
type Service struct {
	inv invoker.Invoker
}

// New 创建 file Service。
func New(inv invoker.Invoker) *Service { return &Service{inv: inv} }

// ListRequest 文件列表请求参数。
type ListRequest struct {
	PDirFID string // 父目录 fid，根目录用 "0"
	Page    int    // 页码，从 1 开始（0 视为 1）
	Size    int    // 每页数量，默认 50
	Sort    string // 排序，如 "file_type:asc,updated_at:desc"（空用默认）
}

// listResponse 文件列表响应（夸克外层结构：{metadata,data:{list,total}}）。
type listResponse struct {
	Code int `json:"code"`
	Msg  string `json:"message"`
	Data struct {
		List  []*types.File `json:"list"`
		Total int           `json:"total"`
	} `json:"data"`
	Status int `json:"status"`
}

// List 列出指定目录下的文件。自动分页，返回所有。
func (s *Service) List(ctx context.Context, req *ListRequest) ([]*types.File, error) {
	if req == nil {
		req = &ListRequest{}
	}
	pdir := req.PDirFID
	if pdir == "" {
		pdir = "0"
	}
	page := req.Page
	if page == 0 {
		page = 1
	}
	size := req.Size
	if size == 0 {
		size = 50
	}

	var all []*types.File
	for {
		params := map[string]any{
			"pdir_fid":      pdir,
			"_page":         page,
			"_size":         size,
			"_fetch_total":  1,
			"_sort":         req.Sort,
		}
		var resp listResponse
		if err := invoker.GetAndDecode(ctx, s.inv, "/file/sort", params, &resp); err != nil {
			return nil, err
		}
		if resp.Code != 0 {
			return nil, invoker.NewAPIError(resp.Code, resp.Msg)
		}
		if len(resp.Data.List) == 0 {
			break
		}
		all = append(all, resp.Data.List...)
		// 够了或没下一页。
		if len(resp.Data.List) < size || len(all) >= resp.Data.Total {
			break
		}
		page++
	}
	return all, nil
}

// ListPage 列出单页文件（手动分页）。返回 (文件, 总数)。
func (s *Service) ListPage(ctx context.Context, req *ListRequest) ([]*types.File, int, error) {
	if req == nil {
		req = &ListRequest{}
	}
	pdir := req.PDirFID
	if pdir == "" {
		pdir = "0"
	}
	page := req.Page
	if page == 0 {
		page = 1
	}
	size := req.Size
	if size == 0 {
		size = 50
	}
	params := map[string]any{
		"pdir_fid":     pdir,
		"_page":        page,
		"_size":        size,
		"_fetch_total": 1,
		"_sort":        req.Sort,
	}
	var resp listResponse
	if err := invoker.GetAndDecode(ctx, s.inv, "/file/sort", params, &resp); err != nil {
		return nil, 0, err
	}
	if resp.Code != 0 {
		return nil, 0, invoker.NewAPIError(resp.Code, resp.Msg)
	}
	return resp.Data.List, resp.Data.Total, nil
}

// Get 获取文件详情。
// 夸克接口：POST /file/info，body {fid, pdir_fid}。
func (s *Service) Get(ctx context.Context, fid string) (*types.File, error) {
	body := map[string]any{"fid": fid}
	var resp struct {
		Code int         `json:"code"`
		Msg  string      `json:"message"`
		Data types.File  `json:"data"`
	}
	if err := invoker.PostAndDecode(ctx, s.inv, "/file/info", body, nil, &resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, invoker.NewAPIError(resp.Code, resp.Msg)
	}
	f := resp.Data
	return &f, nil
}

// 确保编译期 fmt 被使用（错误格式化预留）。
var _ = fmt.Sprintf
