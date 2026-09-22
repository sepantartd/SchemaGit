package cli

import (
    "fmt"

    "github.com/sepanta/schemagit/internal/config"
    "github.com/sepanta/schemagit/internal/db"
    "github.com/sepanta/schemagit/internal/diff"
    "github.com/sepanta/schemagit/internal/git"
    "github.com/sepanta/schemagit/internal/log"
    "github.com/sepanta/schemagit/internal/migration"
    "github.com/sepanta/schemagit/internal/schema"
)

func RunDiff() error {
    cfg := config.Default()

    if !git.IsGitRepo() {
        return fmt.Errorf("not inside a git repository")
    }

    log.Info("Loading current DB schema...")
    database, err := db.Open(cfg.DBPath)
    if err != nil {
        return err
    }
    currentSchema, err := schema.LoadCurrentSchemaFromDB(database)
    if err != nil {
        return err
    }

    log.Info("Loading desired schema from file...")
    desiredSchema, err := schema.LoadDesiredSchemaFromFile(cfg.SchemaPath)
    if err != nil {
        return err
    }

    log.Info("Computing diff...")
    d, err := diff.ComputeDiff(currentSchema, desiredSchema)
    if err != nil {
        return err
    }

    log.Info("Changes:")
    for _, c := range d.Changes {
        fmt.Printf("- %s on table %s", c.Type, c.TableName)
        if c.Column != nil {
            fmt.Printf(" (column: %s)", c.Column.Name)
        }
        fmt.Println()
    }

    return nil
}

func RunApply() error {
    cfg := config.Default()

    log.Info("Opening database...")
    database, err := db.Open(cfg.DBPath)
    if err != nil {
        return err
    }

    log.Info("Loading current schema...")
    currentSchema, err := schema.LoadCurrentSchemaFromDB(database)
    if err != nil {
        return err
    }

    log.Info("Loading desired schema...")
    desiredSchema, err := schema.LoadDesiredSchemaFromFile(cfg.SchemaPath)
    if err != nil {
        return err
    }

    log.Info("Computing diff...")
    d, err := diff.ComputeDiff(currentSchema, desiredSchema)
    if err != nil {
        return err
    }

    log.Info("Generating SQL...")
    stmts, err := migration.GenerateSQL(d)
    if err != nil {
        return err
    }

    log.Info("Applying changes...")
    err = migration.ApplySQL(database, stmts)
    if err != nil {
        return err
    }

    log.Info("Database synced successfully.")
    return nil
}
