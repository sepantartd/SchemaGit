package diff

type Column struct {
    Name     string
    Type     string
    Nullable bool
    Default  string
}

type Table struct {
    Name    string
    Columns []Column
}

type ModifiedTable struct {
    Name            string
    AddedColumns    []Column
    RemovedColumns  []Column
    ModifiedColumns []Column
}

type SchemaDiff struct {
    AddedTables    []Table
    RemovedTables  []Table
    ModifiedTables []ModifiedTable
}

func Diff(oldSchema, newSchema []Table) SchemaDiff {
    diff := SchemaDiff{}

    oldMap := map[string]Table{}
    newMap := map[string]Table{}

    for _, t := range oldSchema {
        oldMap[t.Name] = t
    }
    for _, t := range newSchema {
        newMap[t.Name] = t
    }

    // جدول‌های جدید
    for name, newTable := range newMap {
        if _, ok := oldMap[name]; !ok {
            diff.AddedTables = append(diff.AddedTables, newTable)
        }
    }

    // جدول‌های حذف‌شده
    for name, oldTable := range oldMap {
        if _, ok := newMap[name]; !ok {
            diff.RemovedTables = append(diff.RemovedTables, oldTable)
        }
    }

    // جدول‌های تغییرکرده
    for name, oldTable := range oldMap {
        newTable, ok := newMap[name]
        if !ok {
            continue
        }

        mt := ModifiedTable{Name: name}

        oldCols := map[string]Column{}
        newCols := map[string]Column{}

        for _, c := range oldTable.Columns {
            oldCols[c.Name] = c
        }
        for _, c := range newTable.Columns {
            newCols[c.Name] = c
        }

        // ستون‌های جدید
        for cname, nc := range newCols {
            if _, ok := oldCols[cname]; !ok {
                mt.AddedColumns = append(mt.AddedColumns, nc)
            }
        }

        // ستون‌های حذف‌شده
        for cname, oc := range oldCols {
            if _, ok := newCols[cname]; !ok {
                mt.RemovedColumns = append(mt.RemovedColumns, oc)
            }
        }

        // ستون‌های تغییرکرده
        for cname, oc := range oldCols {
            nc, ok := newCols[cname]
            if !ok {
                continue
            }

            if oc.Type != nc.Type ||
                oc.Nullable != nc.Nullable ||
                oc.Default != nc.Default {
                mt.ModifiedColumns = append(mt.ModifiedColumns, nc)
            }
        }

        if len(mt.AddedColumns) > 0 ||
            len(mt.RemovedColumns) > 0 ||
            len(mt.ModifiedColumns) > 0 {
            diff.ModifiedTables = append(diff.ModifiedTables, mt)
        }
    }

    return diff
}
