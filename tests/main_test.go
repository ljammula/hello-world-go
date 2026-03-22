package main_test

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

var binaryPath string

func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "helloworld")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create temp dir: %v\n", err)
		os.Exit(1)
	}
	binaryPath = filepath.Join(tmp, "helloworld")
	if err := exec.Command("go", "build", "-o", binaryPath, "../").Run(); err != nil {
		fmt.Fprintf(os.Stderr, "build failed: %v\n", err)
		os.Exit(1)
	}
	code := m.Run()
	os.RemoveAll(tmp)
	os.Exit(code)
}

func TestHelloWorldEnglish(t *testing.T) {
	out, err := exec.Command(binaryPath, "-lang=en").CombinedOutput()
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	output := strings.TrimSpace(string(out))
	want := "Hello, World!"
	if output != want {
		t.Errorf("Got %q, want %q", output, want)
	}
}

func TestHelloWorldTelugu(t *testing.T) {
	out, err := exec.Command(binaryPath, "-lang=te").CombinedOutput()
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	output := strings.TrimSpace(string(out))
	want := "హలో వరల్డ్!"
	if output != want {
		t.Errorf("Got %q, want %q", output, want)
	}
}

func TestHelloWorldHindi(t *testing.T) {
	out, err := exec.Command(binaryPath, "-lang=hi").CombinedOutput()
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	output := strings.TrimSpace(string(out))
	want := "नमस्ते दुनिया!"
	if output != want {
		t.Errorf("Got %q, want %q", output, want)
	}
}

func TestHelloWorldDefault(t *testing.T) {
	out, err := exec.Command(binaryPath).CombinedOutput()
	if err != nil {
		t.Fatalf("Run failed: %v\nOutput: %s", err, out)
	}
	output := strings.TrimSpace(string(out))
	want := "Hello, World!"
	if output != want {
		t.Errorf("Got %q, want %q", output, want)
	}
}

func TestHelloWorldExplicitCLIMode(t *testing.T) {
	out, err := exec.Command(binaryPath, "-mode=cli", "-lang=en").CombinedOutput()
	if err != nil {
		t.Fatalf("Run failed: %v\nOutput: %s", err, out)
	}
	output := strings.TrimSpace(string(out))
	want := "Hello, World!"
	if output != want {
		t.Errorf("Got %q, want %q", output, want)
	}
}

func TestHelloWorldInvalidLang(t *testing.T) {
	out, err := exec.Command(binaryPath, "-lang=fr").CombinedOutput()
	if err == nil {
		t.Fatalf("Run should fail for unsupported language")
	}
	output := strings.TrimSpace(string(out))
	re := regexp.MustCompile(`Unsupported language\. Use one of: [a-z, ]+\.`)
	if !re.MatchString(output) {
		t.Errorf("Got %q, want pattern %q", output, re.String())
	}
}

func TestInvalidModeFails(t *testing.T) {
	out, err := exec.Command(binaryPath, "-mode=banana").CombinedOutput()
	if err == nil {
		t.Fatalf("Run should fail for unsupported mode")
	}
	output := strings.TrimSpace(string(out))
	want := `Unsupported mode "banana". Use one of: cli, http.`
	if output != want {
		t.Fatalf("Got %q, want %q", output, want)
	}
}

func findFreeAddress(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to get free port: %v", err)
	}
	defer l.Close()
	return l.Addr().String()
}

func startHTTPServer(t *testing.T, dir string, addr string) {
	t.Helper()
	cmd := exec.Command(binaryPath, "-mode=http", "-addr", addr)
	cmd.Dir = dir

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start HTTP server: %v", err)
	}

	t.Cleanup(func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		_, _ = cmd.Process.Wait()
	})

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := http.Get("http://" + addr + "/")
		if err == nil {
			resp.Body.Close()
			return
		}
		time.Sleep(50 * time.Millisecond)
	}

	t.Fatalf("server did not start in time; output: %s", out.String())
}

func TestHTTPServesIndexAtRoot(t *testing.T) {
	projectDir := filepath.Clean("..")
	addr := findFreeAddress(t)
	startHTTPServer(t, projectDir, addr)

	resp, err := http.Get("http://" + addr + "/")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("got status %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed reading response body: %v", err)
	}
	expected, err := os.ReadFile(filepath.Join(projectDir, "index.html"))
	if err != nil {
		t.Fatalf("failed reading expected index.html: %v", err)
	}

	gotBody := strings.TrimSpace(string(body))
	wantBody := strings.TrimSpace(string(expected))
	if gotBody != wantBody {
		t.Fatalf("unexpected body\n got: %q\nwant: %q", gotBody, wantBody)
	}
}

func TestHTTPMissingIndexReturnsClear404(t *testing.T) {
	addr := findFreeAddress(t)
	tempDir := t.TempDir()
	startHTTPServer(t, tempDir, addr)

	resp, err := http.Get("http://" + addr + "/")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("got status %d, want %d", resp.StatusCode, http.StatusNotFound)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed reading response body: %v", err)
	}
	got := string(body)
	want := "index file not found: index.html"
	if !strings.Contains(got, want) {
		t.Fatalf("got body %q, want it to contain %q", got, want)
	}
}
