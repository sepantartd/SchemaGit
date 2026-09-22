package git

import (
    "os/exec"
    "strings"
)

func GetCommitHash() (string, error) {
    out, err := exec.Command("git", "rev-parse", "HEAD").Output()
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(out)), nil
}
