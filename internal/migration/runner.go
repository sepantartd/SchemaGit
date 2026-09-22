package migration

import (
    "database/sql"
    "fmt"

    "github.com/sepanta/schemagit/internal/log"
)

func ApplySQL(db *sql.DB, stmts []string) error {
    for _, stmt := range stmts {
        log.Info("Executing: " + stmt)

        _, err := db.Exec(stmt)
        if err != nil {
            // اگر دستور کامنت بود، اجرا نمی‌شود
            if stmt[:2] == "--" {
                log.Info("Skipped (comment): " + stmt)
                continue
            }

            return fmt.Errorf("error executing '%s': %w", stmt, err)
        }
    }

    return nil
}
