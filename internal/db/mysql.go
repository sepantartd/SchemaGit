package db

import (
    "database/sql"
    "fmt"

    _ "github.com/go-sql-driver/mysql"
)

type MySQLDB struct {
    Conn *sql.DB
}

func OpenMySQL(host string, port int, user, pass, dbname string) (*MySQLDB, error) {
    dsn := fmt.Sprintf(
        "%s:%s@tcp(%s:%d)/%s",
        user, pass, host, port, dbname,
    )

    conn, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }

    return &MySQLDB{Conn: conn}, nil
}
