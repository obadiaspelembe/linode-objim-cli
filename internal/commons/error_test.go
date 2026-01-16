
package commons

import (
    "bytes"
    "errors"
    "os"
    "os/exec"
    "regexp"
    "strings"
    "testing"
)

// stripANSI removes ANSI color codes so we can assert on plain text
var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(s string) string {
    return ansiRegexp.ReplaceAllString(s, "")
}

func TestErrorCheck_NoError_NoOutput_NoExit(t *testing.T) {
    // Capture stdout
    origStdout := os.Stdout
    r, w, err := os.Pipe()
    if err != nil {
        t.Fatalf("failed to create pipe: %v", err)
    }
    os.Stdout = w

    // Call with nil error — should not print and not exit
    ErrorCheck(nil, "should not print")

    // Restore stdout
    _ = w.Close()
    os.Stdout = origStdout

    var buf bytes.Buffer
    if _, err := buf.ReadFrom(r); err != nil {
        t.Fatalf("failed to read from pipe: %v", err)
    }
    _ = r.Close()

    if got := buf.String(); got != "" {
        t.Fatalf("expected no output, got %q", got)
    }
}

// We use a helper subprocess to validate os.Exit(1) behavior safely.
func TestErrorCheck_WithError_ExitsAndPrints(t *testing.T) {
    const (
        wantMsg   = "Error getting signed URL"
        wantError = "boom: something went wrong"
    )

    // Build a command to re-run this test binary and execute TestHelperProcess
    cmd := exec.Command(os.Args[0], "-test.run", "^TestHelperProcess$")
    // Pass env vars to trigger and parameterize the helper
    cmd.Env = append(os.Environ(),
        "RUN_ERRORCHECK=1",
        "ERR_TEXT="+wantError,
        "ERR_MSG="+wantMsg,
    )

    out, err := cmd.CombinedOutput()

    // We expect the process to exit with code 1
    if err == nil {
        t.Fatalf("expected process to exit with non-nil error (status 1), got nil")
    }
    exitErr, ok := err.(*exec.ExitError)
    if !ok {
        t.Fatalf("expected *exec.ExitError, got %T: %v", err, err)
    }
    if status := exitErr.ExitCode(); status != 1 {
        t.Fatalf("expected exit code 1, got %d; output:\n%s", status, string(out))
    }

    // The output contains colored text; strip ANSI to assert content
    clean := stripANSI(string(out))

    // ErrorCheck prints: "%s : %s" where both parts are colored
    // Verify both the message and error are present
    if !strings.Contains(clean, wantMsg) {
        t.Fatalf("output missing message %q; got:\n%s", wantMsg, clean)
    }
    if !strings.Contains(clean, wantError) {
        t.Fatalf("output missing error %q; got:\n%s", wantError, clean)
    }

    // Optional: verify the separator " : " is there as emitted by fmt.Printf
    if !strings.Contains(clean, " : ") {
        t.Fatalf("output missing separator ' : '; got:\n%s", clean)
    }
}

// TestHelperProcess is executed in a child process to exercise os.Exit(1)
// It is selected via -test.run and gated by RUN_ERRORCHECK.
func TestHelperProcess(t *testing.T) {
    if os.Getenv("RUN_ERRORCHECK") != "1" {
        return // do nothing when running regular tests
    }
    errText := os.Getenv("ERR_TEXT")
    msg := os.Getenv("ERR_MSG")

    // Call the function under test; this should print and exit(1)
    ErrorCheck(errors.New(errText), msg)

    // If ErrorCheck did not exit, fail to signal unexpected behavior
    t.Fatalf("ErrorCheck did not call os.Exit(1) as expected")
}
