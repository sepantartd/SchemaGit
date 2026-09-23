package generator

import (
    "fmt"
    "schemagit/internal/schema/diff"
)

type Migration struct {
    Statements []string
}

func Generate(d diff.SchemaDiff) Migration {
    mig := Migration{}

    // جدول‌های جدید
    for _, t := range d.AddedTables {
        stmt := fmt.Sprintf("CREATE TABLE %s (", t.Name)
        for i, c := range t.Columns {
            stmt += fmt.Sprintf("%s %s", c.Name, c.Type)
            if !c.Nullable {
                stmt += " NOT NULL"
            }
            if c.Default != "" {
                stmt += fmt.Sprintf(" DEFAULT %s", c.Default)
            }
            if i < len(t.Columns)-1 {
                stmt += ", "
            }
        }
        stmt += ");"
        mig.Statements = append(mig.Statements, stmt)
    }

    // جدول‌های حذف‌شده
    for _, t := range d.RemovedTables {
        stmt := fmt.Sprintf("DROP TABLE %s;", t.Name)
        mig.Statements = append(mig.Statements, stmt)
    }

    // جدول‌های تغییرکرده
    for _, mt := range d.ModifiedTables {
        // ستون‌های جدید
        for _, c := range mt.AddedColumns {
            stmt := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", mt.Name, c.Name, c.Type)
            if !c.Nullable {
                stmt += " NOT NULL"
            }
            if c.Default != "" {
                stmt += fmt.Sprintf(" DEFAULT %s", c.Default)
            }
            stmt += ";"
            mig.Statements = append(mig.Statements, stmt)
        }

        // ستون‌های حذف‌شده
        for _, c := range mt.RemovedColumns {
            stmt := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;", mt.Name, c.Name)
            mig.Statements = append(mig.Statements, stmt)
        }

        // ستون‌های تغییرکرده
        for _, c := range mt.ModifiedColumns {
            stmt := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s;", mt.Name, c.Name, c.Type)
            mig.Statements = append(mig.Statements, stmt)
        }
    }

    return mig
}
