package main

import (
    "fmt"
    "os"

    "github.com/sepanta/schemagit/internal/cli"
    "github.com/sepanta/schemagit/internal/log"
    "github.com/sepanta/schemagit/internal/ui"
    "github.com/sepanta/schemagit/internal/cloud"
)

func printHelp() {
    fmt.Println("SchemaGit v0.4")
    fmt.Println("")
    fmt.Println("Usage:")
    fmt.Println("  schemagit diff        Show differences between DB and desired schema")
    fmt.Println("  schemagit apply       Apply schema changes to database")
    fmt.Println("  schemagit ui          Start Web Dashboard")
    fmt.Println("  schemagit cloud       Start Cloud API Server")
    fmt.Println("")
    fmt.Println("Examples:")
    fmt.Println("  schemagit diff")
    fmt.Println("  schemagit apply")
    fmt.Println("  schemagit ui")
    fmt.Println("  schemagit cloud")
}

func main() {
    if len(os.Args) < 2 {
        printHelp()
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

    case "ui":
        ui.StartServer()

    case "cloud":
        cloud.StartAPIServer()

    case "help":
        printHelp()

    default:
        log.Error("Unknown command: " + cmd)
        printHelp()
    }
}
