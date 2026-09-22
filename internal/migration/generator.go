package migration

import (
    "fmt"

    "github.com/sepanta/schemagit/internal/diff"
    "github.com/sepanta/schemagit/internal/schema"
)

func GenerateSQL(d *diff.SchemaDiff) ([]string, error) {
    stmts := []string{}

    for _, change := range d.Changes {
        switch change.Type {

        case diff.AddTable:
            stmts = append(stmts, generateCreateTable(change.TableName))

        case diff.DropTable:
            stmts = append(stmts, fmt.Sprintf("DROP TABLE IF EXISTS %s;", change.TableName))

        case diff.AddColumn:
            stmts = append(stmts, generateAddColumn(change.TableName, change.Column))

        case diff.DropColumn:
            stmts = append(stmts, generateDropColumn(change.TableName, change.Column))

        case diff.ModifyColumn:
            stmts = append(stmts, generateModifyColumn(change.TableName, change.Column))
        }
    }

    return stmts, nil
}

func generateCreateTable(name string) string {
    return fmt.Sprintf("CREATE TABLE %s ();", name)
}

func generateAddColumn(table string, col *schema.Column) string {
    return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table, col.Name, col.Type)
}

func generateDropColumn(table string, col *schema.Column) string {
    // SQLite از DROP COLUMN پشتیبانی نمی‌کند
    // در نسخه MVP فقط لاگ می‌کنیم
    return fmt.Sprintf("-- DROP COLUMN %s FROM %s (SQLite does not support DROP COLUMN)", col.Name, table)
}

func generateModifyColumn(table string, col *schema.Column) string {
    // SQLite از MODIFY COLUMN پشتیبانی نمی‌کند
    // در نسخه MVP فقط لاگ می‌کنیم
    return fmt.Sprintf("-- MODIFY COLUMN %s IN %s TO TYPE %s (SQLite limitation)", col.Name, table, col.Type)
}
