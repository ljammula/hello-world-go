package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
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
