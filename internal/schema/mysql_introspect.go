package schema

import (
    "database/sql"
)

func LoadMySQLSchema(db *sql.DB) (*DatabaseSchema, error) {
    schema := &DatabaseSchema{}

    tables, err := db.Query(`
        SELECT table_name 
        FROM information_schema.tables 
        WHERE table_schema = DATABASE();
    `)
    if err != nil {
        return nil, err
    }
    defer tables.Close()

    for tables.Next() {
        var tableName string
        tables.Scan(&tableName)

        table := Table{Name: tableName}

        cols, err := db.Query(`
            SELECT column_name, column_type, is_nullable, column_default
            FROM information_schema.columns
            WHERE table_name = ? AND table_schema = DATABASE();
        `, tableName)
        if err != nil {
            return nil, err
        }

        for cols.Next() {
            var name, dtype, nullable string
            var defaultValue *string

            cols.Scan(&name, &dtype, &nullable, &defaultValue)

            table.Columns = append(table.Columns, Column{
                Name:         name,
                Type:         dtype,
                NotNull:      nullable == "NO",
                DefaultValue: defaultValue,
            })
        }

        cols.Close()
        schema.Tables = append(schema.Tables, table)
    }

    return schema, nil
}
