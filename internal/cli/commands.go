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

//
// -------------------------
//        DIFF COMMAND
// -------------------------
//

func RunDiff() error {
    cfg := config.Default()

    if !git.IsGitRepo() {
        return fmt.Errorf("not inside a git repository")
    }

    log.Info("Determining database type: " + cfg.DBType)

    var currentSchema *schema.DatabaseSchema

    //
    // Load current schema depending on DB type
    //
    if cfg.DBType == "sqlite" {
        log.Info("Opening SQLite database: " + cfg.DBPath)
        database, err := db.Open(cfg.DBPath)
        if err != nil {
            return err
        }

        log.Info("Loading current SQLite schema...")
        currentSchema, err = schema.LoadCurrentSchemaFromDB(database)
        if err != nil {
            return err
        }

    } else if cfg.DBType == "postgres" {
        log.Info("Connecting to PostgreSQL...")
        pg, err := db.OpenPostgres(cfg.PGHost, cfg.PGPort, cfg.PGUser, cfg.PGPass, cfg.PGName)
        if err != nil {
            return err
        }

        log.Info("Loading current PostgreSQL schema...")
        currentSchema, err = schema.LoadPostgresSchema(pg.Conn)
        if err != nil {
            return err
        }
    }

    //
    // Load desired schema (from store or schema.sql)
    //
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
        log.Info("No stored schema for this commit → using schema.sql")
        desiredSchema, err = schema.LoadDesiredSchemaFromFile(cfg.SchemaPath)
        if err != nil {
            return err
        }
    }

    //
    // Compute diff
    //
    log.Info("Computing diff...")
    d, err := diff.ComputeDiff(currentSchema, desiredSchema)
    if err != nil {
        return err
    }

    //
    // Print diff
    //
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

//
// -------------------------
//        APPLY COMMAND
// -------------------------
//

func RunApply() error {
    cfg := config.Default()

    log.Info("Determining database type: " + cfg.DBType)

    var currentSchema *schema.DatabaseSchema

    //
    // Load current schema depending on DB type
    //
    if cfg.DBType == "sqlite" {
        log.Info("Opening SQLite database: " + cfg.DBPath)
        database, err := db.Open(cfg.DBPath)
        if err != nil {
            return err
        }

        log.Info("Loading current SQLite schema...")
        currentSchema, err = schema.LoadCurrentSchemaFromDB(database)
        if err != nil {
            return err
        }

    } else if cfg.DBType == "postgres" {
        log.Info("Connecting to PostgreSQL...")
        pg, err := db.OpenPostgres(cfg.PGHost, cfg.PGPort, cfg.PGUser, cfg.PGPass, cfg.PGName)
        if err != nil {
            return err
        }

        log.Info("Loading current PostgreSQL schema...")
        currentSchema, err = schema.LoadPostgresSchema(pg.Conn)
        if err != nil {
            return err
        }
    }

    //
    // Load desired schema (from store or schema.sql)
    //
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
        log.Info("No stored schema for this commit → using schema.sql")
        desiredSchema, err = schema.LoadDesiredSchemaFromFile(cfg.SchemaPath)
        if err != nil {
            return err
        }
    }

    //
    // Compute diff
    //
    log.Info("Computing diff...")
    d, err := diff.ComputeDiff(currentSchema, desiredSchema)
    if err != nil {
        return err
    }

    //
    // Generate SQL depending on DB type
    //
    log.Info("Generating SQL...")
    var stmts []string

    if cfg.DBType == "sqlite" {
        stmts, err = migration.GenerateSQL(d)
    } else if cfg.DBType == "postgres" {
        stmts, err = migration.GeneratePostgresSQL(d)
    }

    if err != nil {
        return err
    }

    //
    // Apply SQL
    //
    log.Info("Applying changes...")

    if cfg.DBType == "sqlite" {
        database, err := db.Open(cfg.DBPath)
        if err != nil {
            return err
        }
        if err := migration.ApplySQL(database, stmts); err != nil {
            return err
        }

    } else if cfg.DBType == "postgres" {
        pg, err := db.OpenPostgres(cfg.PGHost, cfg.PGPort, cfg.PGUser, cfg.PGPass, cfg.PGName)
        if err != nil {
            return err
        }
        if err := migration.ApplySQL(pg.Conn, stmts); err != nil {
            return err
        }
    }

    //
    // Save desired schema for this commit
    //
    log.Info("Saving desired schema for commit: " + commit)
    if err := st.SaveSchema(commit, desiredSchema); err != nil {
        return err
    }

    log.Info("Database synced successfully.")
    return nil
}
