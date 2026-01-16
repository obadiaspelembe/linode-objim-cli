package commons


import (
    "os"
    "path/filepath"
    "testing"
)

func writeFile(t *testing.T, path, content string) {
    t.Helper()
    if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
        t.Fatalf("failed to create directories for %s: %v", path, err)
    }
    if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
        t.Fatalf("failed to write file %s: %v", path, err)
    }
}

func TestLoadConfig(t *testing.T) {
    t.Run("success - loads token and region from profile", func(t *testing.T) {
        // Arrange
        td := t.TempDir()
        // Ensure os.UserHomeDir() resolves to our temp dir on all platforms
        t.Setenv("HOME", td)
        t.Setenv("USERPROFILE", td)

        iniPath := filepath.Join(td, ".linodeobjim", "credentials.ini")
        iniContent := `
[default]
token = abc123
region = us-east

[another]
token = xyz789
region = eu-west
`
        writeFile(t, iniPath, iniContent)

        // Act
        cfg, err := LoadConfig("default")

        // Assert
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if cfg.Token != "abc123" {
            t.Errorf("Token mismatch: got %q, want %q", cfg.Token, "abc123")
        }
        if cfg.Region != "us-east" {
            t.Errorf("Region mismatch: got %q, want %q", cfg.Region, "us-east")
        }
    })

    t.Run("error - configuration file missing", func(t *testing.T) {
        td := t.TempDir()
        t.Setenv("HOME", td)
        t.Setenv("USERPROFILE", td)

        // Do NOT create the file; LoadConfig should error.
        _, err := LoadConfig("default")
        if err == nil {
            t.Fatalf("expected error but got nil")
        }
        if got := err.Error(); got == "" || !containsAll(got, []string{
            "failed to load configuration file",
        }) {
            t.Errorf("unexpected error message: %q", got)
        }
    })

    t.Run("error - profile missing", func(t *testing.T) {
        td := t.TempDir()
        t.Setenv("HOME", td)
        t.Setenv("USERPROFILE", td)

        iniPath := filepath.Join(td, ".linodeobjim", "credentials.ini")
        iniContent := `
[default]
token = abc123
region = us-east
`
        writeFile(t, iniPath, iniContent)

        _, err := LoadConfig("nonexistent")
        if err == nil {
            t.Fatalf("expected error but got nil")
        }
        if got := err.Error(); got == "" || !containsAll(got, []string{
            "profile", "not found",
        }) {
            t.Errorf("unexpected error message: %q", got)
        }
    })

    t.Run("error - token missing in profile", func(t *testing.T) {
        td := t.TempDir()
        t.Setenv("HOME", td)
        t.Setenv("USERPROFILE", td)

        iniPath := filepath.Join(td, ".linodeobjim", "credentials.ini")
        iniContent := `
[default]
region = us-east
`
        writeFile(t, iniPath, iniContent)

        _, err := LoadConfig("default")
        if err == nil {
            t.Fatalf("expected error but got nil")
        }
        if got := err.Error(); got == "" || !containsAll(got, []string{
            "token not found", "default",
        }) {
            t.Errorf("unexpected error message: %q", got)
        }
    })

    t.Run("error - region missing in profile", func(t *testing.T) {
        td := t.TempDir()
        t.Setenv("HOME", td)
        t.Setenv("USERPROFILE", td)

        iniPath := filepath.Join(td, ".linodeobjim", "credentials.ini")
        iniContent := `
[default]
token = abc123
`
        writeFile(t, iniPath, iniContent)

        _, err := LoadConfig("default")
        if err == nil {
            t.Fatalf("expected error but got nil")
        }
        if got := err.Error(); got == "" || !containsAll(got, []string{
            "region not found", "default",
        }) {
            t.Errorf("unexpected error message: %q", got)
        }
    })

    t.Run("error - invalid INI format", func(t *testing.T) {
        td := t.TempDir()
        t.Setenv("HOME", td)
        t.Setenv("USERPROFILE", td)

        iniPath := filepath.Join(td, ".linodeobjim", "credentials.ini")
        // Malformed INI to force parse error
        iniContent := `
[default
token = abc123
region = us-east
`
        writeFile(t, iniPath, iniContent)

        _, err := LoadConfig("default")
        if err == nil {
            t.Fatalf("expected error but got nil")
        }
        if got := err.Error(); got == "" || !containsAll(got, []string{
            "failed to load configuration file",
        }) {
            t.Errorf("unexpected error message: %q", got)
        }
    })
}

// containsAll checks if s contains all substrings in parts.
func containsAll(s string, parts []string) bool {
    for _, p := range parts {
        if !contains(s, p) {
            return false
        }
    }
    return true
}

// contains is a tiny helper to avoid importing strings for brevity.
func contains(s, sub string) bool {
    return len(sub) == 0 || (len(s) >= len(sub) && indexOf(s, sub) >= 0)
}

// indexOf naive substring search to avoid extra imports in this snippet.
func indexOf(s, sub string) int {
outer:
    for i := 0; i+len(sub) <= len(s); i++ {
        for j := 0; j < len(sub); j++ {
            if s[i+j] != sub[j] {
                continue outer
            }
        }
        return i
    }
    return -1
}
