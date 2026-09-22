package migration

import (
    "fmt"

    "github.com/sepanta/schemagit/internal/diff"
    "github.com/sepanta/schemagit/internal/schema"
)

func GenerateMySQLSQL(d *diff.SchemaDiff) ([]string, error) {
    stmts := []string{}

    for _, change := range d.Changes {
        switch change.Type {

        case diff.AddTable:
            stmts = append(stmts, fmt.Sprintf("CREATE TABLE %s ();", change.TableName))

        case diff.DropTable:
            stmts = append(stmts, fmt.Sprintf("DROP TABLE IF EXISTS %s;", change.TableName))

        case diff.AddColumn:
            stmts = append(stmts, fmt.Sprintf(
                "ALTER TABLE %s ADD COLUMN %s %s;",
                change.TableName, change.Column.Name, change.Column.Type,
            ))

        case diff.DropColumn:
            stmts = append(stmts, fmt.Sprintf(
                "ALTER TABLE %s DROP COLUMN %s;",
                change.TableName, change.Column.Name,
            ))

        case diff.ModifyColumn:
            stmts = append(stmts, fmt.Sprintf(
                "ALTER TABLE %s MODIFY COLUMN %s %s;",
                change.TableName, change.Column.Name, change.Column.Type,
            ))
        }
    }

    return stmts, nil
}
