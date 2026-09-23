package main

import (
    "fmt"
    "schemagit/internal/schema/diff"
    "schemagit/internal/schema/generator"
    "schemagit/internal/schema/planner"
)

func main() {
    d := diff.SchemaDiff{
        AddedTables: []diff.Table{
            {Name: "posts", Columns: []diff.Column{{Name: "id", Type: "INTEGER"}}},
        },
    }

    mig := generator.Generate(d)
    plan := planner.BuildPlan(d, mig)

    for _, step := range plan.Steps {
        fmt.Printf("Step: %s\nRollback: %s\nType: %s\n\n", step.Statement, step.Rollback, step.Type)
    }
}
