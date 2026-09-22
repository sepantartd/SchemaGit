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
    "github.com/sepanta/schemagit/internal/store"
)

func RunDiff() error {
    cfg := config.Default()

    if !git.IsGitRepo() {
        return fmt.Errorf("not inside a git repository")
    }

    log.Info("Opening database...")
    database, err := db.Open(cfg.DBPath)
    if err != nil {
        return err
    }

    log.Info("Loading current DB schema...")
    currentSchema, err := schema.LoadCurrentSchemaFromDB(database)
    if err != nil {
        return err
    }

    log.Info("Determining desired schema...")
    commit, err := git.GetCommitHash()
    if err != nil {
        return err
    }

    st, err := store.Open(".schemagit.db")
    if err != nil {
        return err
    }

    desiredSchema, err := st.LoadSchema(commit)
    if err != nil {
        log.Info("No stored schema for this commit, falling back to schema.sql")
        desiredSchema, err = schema.LoadDesiredSchemaFromFile(cfg.SchemaPath)
        if err != nil {
            return err
        }
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

    log.Info("Determining desired schema...")
    commit, err := git.GetCommitHash()
    if err != nil {
        return err
    }

    st, err := store.Open(".schemagit.db")
    if err != nil {
        return err
    }

    desiredSchema, err := st.LoadSchema(commit)
    if err != nil {
        log.Info("No stored schema for this commit, using schema.sql")
        desiredSchema, err = schema.LoadDesiredSchemaFromFile(cfg.SchemaPath)
        if err != nil {
            return err
        }
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
    if err := migration.ApplySQL(database, stmts); err != nil {
        return err
    }

    log.Info("Saving desired schema for this commit...")
    if err := st.SaveSchema(commit, desiredSchema); err != nil {
        return err
    }

    log.Info("Database synced successfully for commit: " + commit)
    return nil
}
