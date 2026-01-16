
package commons

import (
    "bytes"
    "os" 
    "testing"
)


// captureStdout runs f while capturing anything written to os.Stdout.
func captureStdout(t *testing.T, f func()) string {
    t.Helper()
    orig := os.Stdout
    r, w, err := os.Pipe()
    if err != nil {
        t.Fatalf("failed to create pipe: %v", err)
    }
    os.Stdout = w

    // Execute target function
    f()

    // Restore stdout
    _ = w.Close()
    os.Stdout = orig

    var buf bytes.Buffer
    if _, err := buf.ReadFrom(r); err != nil {
        t.Fatalf("failed to read from pipe: %v", err)
    }
    _ = r.Close()

    return buf.String()
}

func TestNewPrinter(t *testing.T) {
    p := NewPrinter()
    if p == nil {
        t.Fatalf("NewPrinter returned nil")
    }
}

func TestPrinter_Success(t *testing.T) {
    p := NewPrinter()
    got := captureStdout(t, func() {
        p.Success("done")
    })
    clean := stripANSI(got)
    if clean != "done" {
        t.Fatalf("unexpected output: got %q, want %q", clean, "done")
    }
}

func TestPrinter_SuccessEx(t *testing.T) {
    p := NewPrinter()
    got := captureStdout(t, func() {
        p.SuccessEx("complete", "[OK]")
    })
    clean := stripANSI(got)
    // SuccessEx prints: "<exclude> <green(content)>"
    want := "[OK] complete"
    if clean != want {
        t.Fatalf("unexpected output: got %q, want %q", clean, want)
    }
}

func TestPrinter_Info(t *testing.T) {
    p := NewPrinter()
    got := captureStdout(t, func() {
        p.Info("info msg")
    })
    clean := stripANSI(got)
    if clean != "info msg" {
        t.Fatalf("unexpected output: got %q, want %q", clean, "info msg")
    }
}

func TestPrinter_InfoEx(t *testing.T) {
    p := NewPrinter()
    got := captureStdout(t, func() {
        p.InfoEx("starting", ">>")
    })
    clean := stripANSI(got)
    want := ">> starting"
    if clean != want {
        t.Fatalf("unexpected output: got %q, want %q", clean, want)
    }
}

func TestPrinter_Error(t *testing.T) {
    p := NewPrinter()
    got := captureStdout(t, func() {
        p.Error("failed")
    })
    clean := stripANSI(got)
    if clean != "failed" {
        t.Fatalf("unexpected output: got %q, want %q", clean, "failed")
    }
}

// Optional: edge-case tests to ensure no panics with empty strings
func TestPrinter_EmptyInputs(t *testing.T) {
    p := NewPrinter()

    _ = captureStdout(t, func() { p.Success("") })
    _ = captureStdout(t, func() { p.SuccessEx("", "") })
    _ = captureStdout(t, func() { p.Info("") })
    _ = captureStdout(t, func() { p.InfoEx("", "") })
    _ = captureStdout(t, func() { p.Error("") })
    // If any panic occurred, the test would fail automatically.
}
