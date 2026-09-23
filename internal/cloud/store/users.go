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
    row := DB.QueryRow(`SELECT password FROM users WHERE email = ?`, email)

    var stored string
    if err := row.Scan(&stored); err != nil {
        return "", err
    }
    if stored != hashPassword(password) {
        return "", errors.New("invalid password")
    }

    token := uuid.New().String()
    if _, err := DB.Exec(`INSERT INTO sessions (token, email) VALUES (?, ?)`, token, email); err != nil {
        return "", err
    }
    return token, nil
}

func ValidateToken(token string) (string, error) {
    row := DB.QueryRow(`SELECT email FROM sessions WHERE token = ?`, token)
    var email string
    if err := row.Scan(&email); err != nil {
        return "", err
    }
    return email, nil
}

func ChangePassword(email, old, new string) error {
    if new == "" {
        return errors.New("new password is required")
    }

    row := DB.QueryRow(`SELECT password FROM users WHERE email = ?`, email)
    var stored string
    if err := row.Scan(&stored); err != nil {
        return err
    }
    if stored != hashPassword(old) {
        return errors.New("invalid old password")
    }

    _, err := DB.Exec(`UPDATE users SET password = ? WHERE email = ?`, hashPassword(new), email)
    return err
}

func DeleteUser(email string) {
    // Delete dependent records explicitly so this also works for existing databases
    // that were created without foreign-key cascading enabled.
    DB.Exec(`DELETE FROM sessions WHERE email = ?`, email)
    DB.Exec(`DELETE FROM billing WHERE email = ?`, email)
    DB.Exec(`DELETE FROM migrations WHERE owner = ?`, email)
    DB.Exec(`DELETE FROM projects WHERE owner = ?`, email)
    DB.Exec(`DELETE FROM users WHERE email = ?`, email)
}
