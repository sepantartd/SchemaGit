package main

import (
    "fmt"
    "schemagit/internal/schema/diff"
    "schemagit/internal/schema/generator"
)

func main() {
    d := diff.SchemaDiff{
        AddedTables: []diff.Table{
            {Name: "posts", Columns: []diff.Column{{Name: "id", Type: "INTEGER"}}},
        },
    }

    mig := generator.Generate(d)
    for _, stmt := range mig.Statements {
        fmt.Println(stmt)
    }
}
