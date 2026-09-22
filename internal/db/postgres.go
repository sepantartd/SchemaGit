package db

import (
    "database/sql"
    "fmt"

    _ "github.com/lib/pq"
)

type PostgresDB struct {
    Conn *sql.DB
}

func OpenPostgres(host string, port int, user, pass, dbname string) (*PostgresDB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, pass, dbname,
    )

    conn, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }

    return &PostgresDB{Conn: conn}, nil
}
