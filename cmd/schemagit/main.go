package main

import (
    "fmt"
    "os"

    "github.com/sepanta/schemagit/internal/cli"
    "github.com/sepanta/schemagit/internal/log"
)

func printHelp() {
    fmt.Println("SchemaGit v0.2")
    fmt.Println("")
    fmt.Println("Usage:")
    fmt.Println("  schemagit diff      Show differences between DB and desired schema")
    fmt.Println("  schemagit apply     Apply schema changes to database")
    fmt.Println("")
    fmt.Println("Examples:")
    fmt.Println("  schemagit diff")
    fmt.Println("  schemagit apply")
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

    case "help":
        printHelp()

    default:
        log.Error("Unknown command: " + cmd)
        printHelp()
    }
}
