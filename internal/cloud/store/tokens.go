package store

import (
    "github.com/google/uuid"
)

func CreateToken(email string) (string, error) {
    token := uuid.New().String()

    _, err := DB.Exec(`
        INSERT INTO api_tokens (token, email)
        VALUES (?, ?)
    `, token, email)

    return token, err
}

func ListTokens(email string) ([]string, error) {
    rows, err := DB.Query(`
        SELECT token FROM api_tokens WHERE email = ?
    `, email)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var out []string
    for rows.Next() {
        var t string
        if err := rows.Scan(&t); err != nil {
            return nil, err
        }
        out = append(out, t)
    }

    return out, nil
}

func DeleteToken(token string) error {
    _, err := DB.Exec(`
        DELETE FROM api_tokens WHERE token = ?
    `, token)
    return err
}

func ValidateAPIToken(token string) (string, error) {
    row := DB.QueryRow(`
        SELECT email FROM api_tokens WHERE token = ?
    `, token)

    var email string
    err := row.Scan(&email)
    if err != nil {
        return "", err
    }

    return email, nil
}
