package main

import (
    "fmt"
    "os"

    "github.com/sepanta/schemagit/internal/cli"
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
        if err := cli.RunDiff(); err != nil {
            log.Error(err.Error())
        }
    case "apply":
        if err := cli.RunApply(); err != nil {
            log.Error(err.Error())
        }
    default:
        log.Error("Unknown command: " + cmd)
    }
}
