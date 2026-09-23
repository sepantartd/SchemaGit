package store

import (
    "errors"
    "os"
    "strconv"
)

type BillingPlan struct {
    Email         string
    Plan          string
    MaxProjects   int
    MaxMigrations int
}

const (
    defaultFreeMaxProjects   = 3
    defaultFreeMaxMigrations = 100
    defaultProMaxProjects    = 50
    defaultProMaxMigrations  = 10000
)

func envLimit(name string, fallback int) int {
    value, err := strconv.Atoi(os.Getenv(name))
    if err != nil || value < 0 {
        return fallback
    }
    return value
}

func billingLimits(plan string) (int, int) {
    if plan == "pro" {
        return envLimit("BILLING_PRO_MAX_PROJECTS", defaultProMaxProjects), envLimit("BILLING_PRO_MAX_MIGRATIONS", defaultProMaxMigrations)
    }
    return envLimit("BILLING_FREE_MAX_PROJECTS", defaultFreeMaxProjects), envLimit("BILLING_FREE_MAX_MIGRATIONS", defaultFreeMaxMigrations)
}

func GetBillingPlan(email string) (*BillingPlan, error) {
    row := DB.QueryRow(`
        SELECT plan FROM billing WHERE email = ?
    `, email)

    var plan string
    err := row.Scan(&plan)
    if err != nil {
        // Existing users without a billing row remain on the free plan.
        if plan == "" {
            plan = "free"
        } else {
            return nil, err
        }
    }

    maxProjects, maxMigrations := billingLimits(plan)
    return &BillingPlan{
        Email:         email,
        Plan:          plan,
        MaxProjects:   maxProjects,
        MaxMigrations: maxMigrations,
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

func CanCreateProject(email string) bool {
    plan, err := GetBillingPlan(email)
    if err != nil {
        return false
    }
    return CountProjects(email) < plan.MaxProjects
}

func CanCreateMigration(email string) bool {
    plan, err := GetBillingPlan(email)
    if err != nil {
        return false
    }
    return CountMigrations(email) < plan.MaxMigrations
}
