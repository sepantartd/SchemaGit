package main

import (
    "fmt"
    "schemagit/internal/schema/diff"
)

func main() {
    oldSchema := []diff.Table{
        {Name: "users", Columns: []diff.Column{{Name: "id", Type: "INTEGER"}}},
    }
    newSchema := []diff.Table{
        {Name: "users", Columns: []diff.Column{{Name: "id", Type: "INTEGER"}, {Name: "email", Type: "TEXT"}}},
        {Name: "posts", Columns: []diff.Column{{Name: "id", Type: "INTEGER"}}},
    }

    result := diff.Diff(oldSchema, newSchema)
    fmt.Printf("Added tables: %v\n", result.AddedTables)
    fmt.Printf("Modified tables: %v\n", result.ModifiedTables)
}
