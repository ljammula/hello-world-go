package main_test

import (
    "os/exec"
    "strings"
    "testing"
    "regexp"
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

func TestHelloWorldHindi(t *testing.T) {
    out, err := exec.Command("go", "run", "../main.go", "-lang=hi").CombinedOutput()
    if err != nil {
        t.Fatalf("Run failed: %v", err)
    }
    output := strings.TrimSpace(string(out))
    want := "नमस्ते दुनिया!"
    if output != want {
        t.Errorf("Got %q, want %q", output, want)
    }
}

func TestHelloWorldInvalidLang(t *testing.T) {
    out, err := exec.Command("go", "run", "../main.go", "-lang=fr").CombinedOutput()
    if err == nil {
        t.Fatalf("Run should fail for unsupported language")
    }
    output := strings.TrimSpace(string(out))
    re := regexp.MustCompile(`Unsupported language. Use one of: [a-z, ]+\.`)
    if !re.MatchString(output) {
        t.Errorf("Got %q, want pattern %q", output, re.String())
    }
}
