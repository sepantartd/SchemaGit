package cloud

import (
    "os"
    "os/exec"
    "strings"
)

func GitCloneOrPull(repo, branch, token, workDir string) error {
    // Private repo → replace https://github.com/... with token auth
    if token != "" {
        repo = strings.Replace(repo, "https://", "https://"+token+"@", 1)
    }

    if _, err := os.Stat(workDir); os.IsNotExist(err) {
        cmd := exec.Command("git", "clone", "-b", branch, repo, workDir)
        return cmd.Run()
    }

    cmd := exec.Command("git", "-C", workDir, "pull")
    return cmd.Run()
}
