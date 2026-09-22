package main

import (
    "fmt"
    "os"

    "github.com/sepanta/schemagit/internal/log"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("SchemaGit v0.1")
        fmt.Println("Usage:")
        fmt.Println("  schemagit diff")
        fmt.Println("  schemagit apply")
        return
    }

    cmd := os.Args[1]

    switch cmd {
    case "diff":
        log.Info("Diff command will be implemented in next phase")
    case "apply":
        log.Info("Apply command will be implemented in next phase")
    default:
        log.Error("Unknown command: " + cmd)
    }
}
