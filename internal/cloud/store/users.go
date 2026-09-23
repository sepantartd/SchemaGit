package store

import (
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "github.com/google/uuid"
)

type User struct {
    Email    string
    Password string
}

func hashPassword(pw string) string {
    h := sha256.Sum256([]byte(pw))
    return hex.EncodeToString(h[:])
}

func CreateUser(email, password string) error {
    _, err := DB.Exec(`
        INSERT INTO users (email, password)
        VALUES (?, ?)
    `, email, hashPassword(password))

    return err
}

func Authenticate(email, password string) (string, error) {
    row := DB.QueryRow(`
        SELECT password FROM users WHERE email = ?
    `, email)

    var stored string
    err := row.Scan(&stored)
    if err != nil {
        return "", err
    }

    if stored != hashPassword(password) {
        return "", errors.New("invalid password")
    }

    token := uuid.New().String()

    DB.Exec(`
        INSERT INTO sessions (token, email)
        VALUES (?, ?)
    `, token, email)

    return token, nil
}

func ValidateToken(token string) (string, error) {
    row := DB.QueryRow(`
        SELECT email FROM sessions WHERE token = ?
    `, token)

    var email string
    err := row.Scan(&email)
    if err != nil {
        return "", err
    }

    return email, nil
}
