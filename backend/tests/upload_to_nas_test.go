package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"super-supply-chain/configs"
	"super-supply-chain/utils"
)

// 用临时配置覆写 WebDAV 全局变量，返回还原函数。
func withWebDAVConfig(t *testing.T, url, user, pass string) {
	t.Helper()
	origURL, origUser, origPass := configs.WEB_DAV_URL, configs.WEB_DAV_USER, configs.WEB_DAV_PASSWORD
	configs.WEB_DAV_URL = url
	configs.WEB_DAV_USER = user
	configs.WEB_DAV_PASSWORD = pass
	t.Cleanup(func() {
		configs.WEB_DAV_URL = origURL
		configs.WEB_DAV_USER = origUser
		configs.WEB_DAV_PASSWORD = origPass
	})
}

// 在临时目录里写一个内容已知的文件，返回完整路径。
func writeTempFile(t *testing.T, name, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return path
}

// 用一个简易 WebDAV stub 录下 PUT 请求并校验拼装出的访问 URL 与上传内容。
func TestUploadToNasUploadsFileAndReturnsURL(t *testing.T) {
	const (
		user    = "dev"
		pass    = "secret"
		newName = "abc123.xlsx"
		content = "hello-nas"
	)

	var (
		mu          sync.Mutex
		putBody     []byte
		putPath     string
		gotAuthUser string
		gotAuthOK   bool
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			u, p, ok := r.BasicAuth()
			if !ok {
				// gowebdav 默认是 AutoAuth：第一发不带凭据，看到 401+WWW-Authenticate
				// 才升级成 Basic 重试。这里走真实的挑战流程。
				w.Header().Set("WWW-Authenticate", `Basic realm="test"`)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			body, _ := io.ReadAll(r.Body)
			mu.Lock()
			putPath = r.URL.Path
			putBody = body
			gotAuthUser = u
			gotAuthOK = p == pass
			mu.Unlock()
			w.WriteHeader(http.StatusCreated)
		default:
			// gowebdav 在某些路径下会发 PROPFIND/MKCOL 之类，宽松放行。
			w.WriteHeader(http.StatusOK)
		}
	}))
	t.Cleanup(srv.Close)

	withWebDAVConfig(t, srv.URL, user, pass)
	filePath := writeTempFile(t, "local-source.xlsx", content)

	gotURL, err := utils.UploadToNas(filePath, newName)
	if err != nil {
		t.Fatalf("UploadToNas returned error: %v", err)
	}

	wantURL := srv.URL + "/" + newName
	if gotURL != wantURL {
		t.Fatalf("UploadToNas() URL = %q, want %q", gotURL, wantURL)
	}

	mu.Lock()
	defer mu.Unlock()
	if !strings.HasSuffix(putPath, "/local-source.xlsx") {
		t.Fatalf("PUT path = %q, want suffix /local-source.xlsx (gowebdav 使用本地文件名作为路径)", putPath)
	}
	if string(putBody) != content {
		t.Fatalf("uploaded body = %q, want %q", string(putBody), content)
	}
	if !gotAuthOK || gotAuthUser != user {
		t.Fatalf("auth header user=%q ok=%v, want user=%q ok=true", gotAuthUser, gotAuthOK, user)
	}
}

// WebDAV 返回 5xx 时应当回传错误，且不返回拼装好的 URL。
func TestUploadToNasReturnsErrorOnServerFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			http.Error(w, "boom", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	withWebDAVConfig(t, srv.URL, "dev", "secret")
	filePath := writeTempFile(t, "fail.xlsx", "payload")

	gotURL, err := utils.UploadToNas(filePath, "fail.xlsx")
	if err == nil {
		t.Fatal("UploadToNas should return error when server responds 500")
	}
	if gotURL != "" {
		t.Fatalf("UploadToNas() URL = %q on failure, want empty string", gotURL)
	}
}

// 本地文件不存在时必须直接报错，并且不应触发任何 HTTP 请求。
func TestUploadToNasReturnsErrorWhenLocalFileMissing(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	withWebDAVConfig(t, srv.URL, "dev", "secret")
	missing := filepath.Join(t.TempDir(), "does-not-exist.xlsx")

	gotURL, err := utils.UploadToNas(missing, "does-not-exist.xlsx")
	if err == nil {
		t.Fatal("UploadToNas should return error when local file does not exist")
	}
	if gotURL != "" {
		t.Fatalf("UploadToNas() URL = %q on missing-file error, want empty string", gotURL)
	}
	if hits != 0 {
		t.Fatalf("server received %d requests on missing-file path, want 0", hits)
	}
}
