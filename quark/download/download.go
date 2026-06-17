// Package download 实现夸克网盘文件下载（占位，后续实现）。
package download

import (
	"github.com/langhuachuanshi/quark-go/quark/invoker"
)

// Service 下载操作入口。
type Service struct {
	inv invoker.Invoker
}

// New 创建 download Service。
func New(inv invoker.Invoker) *Service { return &Service{inv: inv} }
