package log

import (
    "fmt"
)

func Info(msg string) {
    fmt.Println("[INFO]", msg)
}

func Error(msg string) {
    fmt.Println("[ERROR]", msg)
}
