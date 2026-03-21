package tests

import (
    "os/exec"
    "strings"
    "testing"
)

func TestHelloWorld(t *testing.T) {
    out, err := exec.Command("go", "run", "../main.go").CombinedOutput()
    if err != nil {
        t.Fatalf("Run failed: %v", err)
    }
    output := strings.TrimSpace(string(out))
    want := "Hello, World!"
    if output != want {
        t.Errorf("Got %q, want %q", output, want)
    }
}
