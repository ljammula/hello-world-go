package tests

import (
    "os/exec"
    "strings"
    "testing"
)

func TestHelloWorldEnglish(t *testing.T) {
    out, err := exec.Command("go", "run", "../main.go", "-lang=en").CombinedOutput()
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
    out, err := exec.Command("go", "run", "../main.go", "-lang=te").CombinedOutput()
    if err != nil {
        t.Fatalf("Run failed: %v", err)
    }
    output := strings.TrimSpace(string(out))
    want := "హలో వరల్డ్!"
    if output != want {
        t.Errorf("Got %q, want %q", output, want)
    }
}

func TestHelloWorldInvalidLang(t *testing.T) {
    out, err := exec.Command("go", "run", "../main.go", "-lang=fr").CombinedOutput()
    if err != nil {
        t.Fatalf("Run failed: %v", err)
    }
    output := strings.TrimSpace(string(out))
    want := "Unsupported language. Use 'en' or 'te'."
    if output != want {
        t.Errorf("Got %q, want %q", output, want)
    }
}
