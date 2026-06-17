// Package upload 实现夸克网盘文件上传（占位，后续实现）。
package upload

import (
	"github.com/langhuachuanshi/quark-go/quark/invoker"
)

// Service 上传操作入口。
type Service struct {
	inv invoker.Invoker
}

// New 创建 upload Service。
func New(inv invoker.Invoker) *Service { return &Service{inv: inv} }
