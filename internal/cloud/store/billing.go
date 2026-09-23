package store

import (
    "errors"
)

type BillingPlan struct {
    Email         string
    Plan          string
    MaxProjects   int
    MaxMigrations int
}

func GetBillingPlan(email string) (*BillingPlan, error) {
    row := DB.QueryRow(`
        SELECT plan FROM billing WHERE email = ?
    `, email)

    var plan string
    err := row.Scan(&plan)
    if err != nil {
        return nil, err
    }

    if plan == "pro" {
        return &BillingPlan{
            Email:         email,
            Plan:          "pro",
            MaxProjects:   50,
            MaxMigrations: 10000,
        }, nil
    }

    return &BillingPlan{
        Email:         email,
        Plan:          "free",
        MaxProjects:   3,
        MaxMigrations: 100,
    }, nil
}

func SetBillingPlan(email, plan string) error {
    if plan != "free" && plan != "pro" {
        return errors.New("invalid plan")
    }

    _, err := DB.Exec(`
        INSERT OR REPLACE INTO billing (email, plan)
        VALUES (?, ?)
    `, email, plan)

    return err
}
