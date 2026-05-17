package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/studio-b12/gowebdav"

	"super-supply-chain/utils"
)

// 从当前测试工作目录向上回溯，找到 configs/.env，找不到返回空串。
func findRepoDotEnv() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for i := 0; i < 6; i++ {
		cand := filepath.Join(dir, "configs", ".env")
		if _, err := os.Stat(cand); err == nil {
			return cand
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
	return ""
}

// 加载仓库根的 configs/.env；找不到或字段不全时跳过测试，CI 安全。
func loadDotEnvOrSkip(t *testing.T) (server, user, pass string) {
	t.Helper()
	envPath := findRepoDotEnv()
	if envPath == "" {
		t.Skip("configs/.env not found; skipping real-WebDAV integration test")
	}
	// 使用 Load（非 Overload）：进程已有的环境变量优先。
	if err := godotenv.Load(envPath); err != nil {
		t.Fatalf("godotenv.Load(%s): %v", envPath, err)
	}
	server = os.Getenv("UPLOAD_SERVER")
	user = os.Getenv("UPLOAD_USER")
	pass = os.Getenv("UPLOAD_PASSWORD")
	if server == "" || user == "" || pass == "" {
		t.Skipf("UPLOAD_SERVER / UPLOAD_USER / UPLOAD_PASSWORD 任一为空，跳过 (%s)", envPath)
	}
	return server, user, pass
}

// 用 .env 中真实的 WebDAV 凭据做一次端到端 PUT + 读回 + DELETE。
// 注意：会在真实 NAS 上短暂留下一个以 _ssc_it_ 前缀命名的探测文件，t.Cleanup 会删除。
func TestUploadToNasEndToEndAgainstDotEnvWebDAV(t *testing.T) {
	server, user, pass := loadDotEnvOrSkip(t)
	withWebDAVConfig(t, server, user, pass)

	// 模拟 controller 的命名约定：local basename == newFileName。
	// upload-to-nas.go 内部用的是 fileInfo.Name()，只有这样上传路径才与
	// UploadToNas 返回的 URL 真正一致。
	newFileName := fmt.Sprintf("_ssc_it_%d.txt", time.Now().UnixNano())
	content := fmt.Sprintf("ssc-integration-payload-%d", time.Now().UnixNano())
	filePath := writeTempFile(t, newFileName, content)

	// 确保无论后续断言成败，远端文件都会被清理。
	t.Cleanup(func() {
		c := gowebdav.NewClient(server, user, pass)
		if err := c.Remove("/" + newFileName); err != nil {
			t.Logf("cleanup remote /%s: %v", newFileName, err)
		}
	})

	gotURL, err := utils.UploadToNas(filePath, newFileName)
	if err != nil {
		t.Fatalf("UploadToNas: %v", err)
	}
	wantURL := server + "/" + newFileName
	if gotURL != wantURL {
		t.Errorf("UploadToNas() URL = %q, want %q", gotURL, wantURL)
	}

	// 用 gowebdav 把刚上传的对象读回来，比对内容。
	c := gowebdav.NewClient(server, user, pass)
	got, err := c.Read("/" + newFileName)
	if err != nil {
		t.Fatalf("read back /%s: %v", newFileName, err)
	}
	if string(got) != content {
		t.Errorf("remote body = %q, want %q", string(got), content)
	}
}
