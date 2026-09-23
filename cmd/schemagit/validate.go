package main

import (
    "fmt"
    "schemagit/internal/schema/diff"
    "schemagit/internal/schema/validator"
)

func main() {
    d := diff.SchemaDiff{
        RemovedTables: []diff.Table{{Name: "old_table"}},
        ModifiedTables: []diff.ModifiedTable{
            {Name: "users", RemovedColumns: []diff.Column{{Name: "email", Type: "TEXT"}}},
        },
    }

    result := validator.Validate(d)
    for _, issue := range result.Issues {
        fmt.Printf("[%s] %s\n", issue.Severity, issue.Message)
    }
}
