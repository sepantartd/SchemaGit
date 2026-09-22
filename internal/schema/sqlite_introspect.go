package schema

import (
    "database/sql"
)

func LoadCurrentSchemaFromDB(db *sql.DB) (*DatabaseSchema, error) {
    rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    schema := &DatabaseSchema{}

    for rows.Next() {
        var tableName string
        rows.Scan(&tableName)

        table := Table{Name: tableName}

        colRows, err := db.Query("PRAGMA table_info(" + tableName + ")")
        if err != nil {
            return nil, err
        }

        for colRows.Next() {
            var cid int
            var name, ctype string
            var notnull int
            var dfltValue sql.NullString
            var pk int

            colRows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk)

            var def *string
            if dfltValue.Valid {
                def = &dfltValue.String
            }

            table.Columns = append(table.Columns, Column{
                Name:        name,
                Type:        ctype,
                NotNull:     notnull == 1,
                DefaultValue: def,
            })
        }

        colRows.Close()
        schema.Tables = append(schema.Tables, table)
    }

    return schema, nil
}
