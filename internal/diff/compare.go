package diff

import (
    "github.com/sepanta/schemagit/internal/schema"
)

func ComputeDiff(current, desired *schema.DatabaseSchema) (*SchemaDiff, error) {
    d := NewDiff()

    currentTables := map[string]schema.Table{}
    desiredTables := map[string]schema.Table{}

    for _, t := range current.Tables {
        currentTables[t.Name] = t
    }
    for _, t := range desired.Tables {
        desiredTables[t.Name] = t
    }

    // 1. جدول‌های جدید
    for name, desiredTable := range desiredTables {
        if _, exists := currentTables[name]; !exists {
            d.Add(TableChange{
                Type:      AddTable,
                TableName: name,
            })
        } else {
            // بررسی ستون‌ها
            currTable := currentTables[name]

            currCols := map[string]schema.Column{}
            desCols := map[string]schema.Column{}

            for _, c := range currTable.Columns {
                currCols[c.Name] = c
            }
            for _, c := range desiredTable.Columns {
                desCols[c.Name] = c
            }

            // ستون جدید
            for colName, col := range desCols {
                if _, exists := currCols[colName]; !exists {
                    d.Add(TableChange{
                        Type:      AddColumn,
                        TableName: name,
                        Column:    &col,
                    })
                }
            }

            // ستون حذف‌شده
            for colName, col := range currCols {
                if _, exists := desCols[colName]; !exists {
                    d.Add(TableChange{
                        Type:      DropColumn,
                        TableName: name,
                        Column:    &col,
                    })
                }
            }

            // تغییر نوع ستون
            for colName, desCol := range desCols {
                if currCol, exists := currCols[colName]; exists {
                    if currCol.Type != desCol.Type {
                        d.Add(TableChange{
                            Type:      ModifyColumn,
                            TableName: name,
                            Column:    &desCol,
                        })
                    }
                }
            }
        }
    }

    // 2. جدول حذف‌شده
    for name := range currentTables {
        if _, exists := desiredTables[name]; !exists {
            d.Add(TableChange{
                Type:      DropTable,
                TableName: name,
            })
        }
    }

    return d, nil
}
