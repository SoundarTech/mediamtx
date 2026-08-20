package embedded

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestStartRequiresExplicitConfig 验证嵌入模式不会从宿主工作目录隐式发现配置。
func TestStartRequiresExplicitConfig(t *testing.T) {
	t.Parallel()

	_, err := Start(Config{})
	if err == nil {
		t.Fatal("Start() error = nil, want configuration-path error")
	}
}

// TestServerLifecycle 验证关闭幂等，且 Wait 可由调用方上下文限制等待时间。
func TestServerLifecycle(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "mediamtx.yml")
	config := []byte("logLevel: error\n" +
		"api: false\n" +
		"metrics: false\n" +
		"pprof: false\n" +
		"playback: false\n" +
		"rtsp: false\n" +
		"rtmp: false\n" +
		"hls: false\n" +
		"webrtc: false\n" +
		"srt: false\n" +
		"moq: false\n")
	if err := os.WriteFile(configPath, config, 0o600); err != nil {
		t.Fatalf("write configuration: %v", err)
	}

	server, err := Start(Config{ConfigPath: configPath})
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if err := server.Wait(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Wait() error = %v, want context deadline exceeded", err)
	}

	server.Close()
	server.Close()

	completedCtx, completedCancel := context.WithTimeout(context.Background(), time.Second)
	defer completedCancel()
	if err := server.Wait(completedCtx); err != nil {
		t.Fatalf("Wait() after Close error = %v", err)
	}
}
