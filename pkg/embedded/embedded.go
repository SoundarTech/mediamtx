// Package embedded 提供供宿主进程安全嵌入 MediaMTX 的最小生命周期接口。
// 它不暴露上游 internal 实现，也不接管宿主进程的信号、退出或升级策略。
package embedded

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/bluenviron/mediamtx/internal/core"
)

// ErrNilContext 表示 Wait 的调用方没有提供可取消的上下文。
// 宿主应始终使用带 deadline 或取消路径的上下文，避免等待逻辑无限阻塞。
var ErrNilContext = errors.New("embedded MediaMTX wait context is nil")

// Config 定义嵌入式 MediaMTX 实例的启动输入。
// ConfigPath 必须指向由宿主控制的完整 MediaMTX 配置文件；空值会被拒绝，避免从工作目录或系统路径隐式发现配置。
type Config struct {
	ConfigPath string
}

// Server 是一个已启动的嵌入式 MediaMTX 实例。
// Server 可被多个 goroutine 并发等待或关闭；Close 幂等，且不会退出宿主进程或处理进程级信号。
type Server struct {
	core      *core.Core
	closeOnce sync.Once
}

// Start 根据显式配置文件启动一个嵌入式 MediaMTX 实例。
// 启动成功后调用方必须调用 Close 释放监听器和后台 goroutine；配置或资源初始化错误会原样保留在错误链中返回。
func Start(config Config) (*Server, error) {
	if config.ConfigPath == "" {
		return nil, fmt.Errorf("start embedded MediaMTX: configuration path is required")
	}

	instance, err := core.NewForEmbedding(config.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("start embedded MediaMTX: %w", err)
	}

	return &Server{core: instance}, nil
}

// Close 停止实例并等待其后台 goroutine 退出。
// 该方法幂等且可并发调用；对 nil 接收者无副作用，便于调用方在延迟清理中安全使用。
func (s *Server) Close() {
	if s == nil || s.core == nil {
		return
	}

	s.closeOnce.Do(s.core.Close)
}

// Wait 等待实例结束或上下文取消。
// 实例正常结束返回 nil；上下文结束时返回其错误，调用方仍应随后调用 Close 释放实例资源。
func (s *Server) Wait(ctx context.Context) error {
	if ctx == nil {
		return ErrNilContext
	}
	if s == nil || s.core == nil {
		return fmt.Errorf("wait embedded MediaMTX: server is nil")
	}

	select {
	case <-s.core.Done():
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
