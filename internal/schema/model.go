package schema

type Column struct {
    Name string
    Type string
    NotNull bool
    DefaultValue *string
}

type Table struct {
    Name    string
    Columns []Column
}

type DatabaseSchema struct {
    Tables []Table
}
