package config

type Config struct {
    DBType string

    // SQLite
    DBPath string

    // PostgreSQL
    PGHost string
    PGPort int
    PGUser string
    PGPass string
    PGName string

    // MySQL
    MYHost string
    MYPort int
    MYUser string
    MYPass string
    MYName string

    SchemaPath string
    StorePath  string
}

func Default() *Config {
    return &Config{
        DBType: "sqlite",

        DBPath: "local.db",

        PGHost: "localhost",
        PGPort: 5432,
        PGUser: "postgres",
        PGPass: "postgres",
        PGName: "postgres",

        MYHost: "localhost",
        MYPort: 3306,
        MYUser: "root",
        MYPass: "",
        MYName: "test",

        SchemaPath: "schema.sql",
        StorePath:  ".schemagit.db",
    }
}
