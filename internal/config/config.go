package config

//
// Config ساختار کامل تنظیمات پروژه SchemaGit
// این ساختار هم SQLite و هم PostgreSQL را پشتیبانی می‌کند.
// در آینده می‌توان MySQL و سایر دیتابیس‌ها را نیز اضافه کرد.
//
type Config struct {
    // نوع دیتابیس: sqlite یا postgres
    DBType string

    // -------------------------
    // SQLite
    // -------------------------
    DBPath string

    // -------------------------
    // PostgreSQL
    // -------------------------
    PGHost string
    PGPort int
    PGUser string
    PGPass string
    PGName string

    // -------------------------
    // مسیر فایل اسکیما
    // -------------------------
    SchemaPath string

    // -------------------------
    // مسیر دیتابیس داخلی SchemaGit
    // برای ذخیره اسکیما هر کامیت
    // -------------------------
    StorePath string
}

//
// Default مقداردهی پیش‌فرض تنظیمات
//
func Default() *Config {
    return &Config{
        // دیتابیس پیش‌فرض
        DBType: "sqlite",

        // SQLite
        DBPath: "local.db",

        // PostgreSQL
        PGHost: "localhost",
        PGPort: 5432,
        PGUser: "postgres",
        PGPass: "postgres",
        PGName: "postgres",

        // فایل اسکیما
        SchemaPath: "schema.sql",

        // دیتابیس داخلی SchemaGit
        StorePath: ".schemagit.db",
    }
}
