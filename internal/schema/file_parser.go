package schema

import (
    "os"
    "regexp"
    "strings"
)

func LoadDesiredSchemaFromFile(path string) (*DatabaseSchema, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    content := string(data)

    // regex برای پیدا کردن CREATE TABLE
    re := regexp.MustCompile(`CREATE TABLE\s+(\w+)\s*\(([^;]+)\);`)
    matches := re.FindAllStringSubmatch(content, -1)

    schema := &DatabaseSchema{}

    for _, m := range matches {
        tableName := m[1]
        colsRaw := m[2]

        table := Table{Name: tableName}

        cols := strings.Split(colsRaw, ",")
        for _, col := range cols {
            col = strings.TrimSpace(col)

            parts := strings.Split(col, " ")
            if len(parts) < 2 {
                continue
            }

            name := parts[0]
            ctype := parts[1]

            table.Columns = append(table.Columns, Column{
                Name: name,
                Type: ctype,
            })
        }

        schema.Tables = append(schema.Tables, table)
    }

    return schema, nil
}
