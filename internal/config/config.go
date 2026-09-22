package config

type Config struct {
    DBPath     string
    SchemaPath string
}

func Default() *Config {
    return &Config{
        DBPath:     "local.db",
        SchemaPath: "schema.sql",
    }
}
