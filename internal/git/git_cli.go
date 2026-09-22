package git

import (
    "os/exec"
    "strings"
)

func GetCurrentBranch() (string, error) {
    out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(out)), nil
}

func FindSchemaFilePath() (string, error) {
    // در MVP فرض می‌کنیم schema.sql در ریشه مخزن است
    return "schema.sql", nil
}

func IsGitRepo() bool {
    _, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
    return err == nil
}
