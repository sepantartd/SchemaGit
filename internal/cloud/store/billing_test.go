package store

import (
    "os"
    "testing"
)

func TestBillingLimitsDefaultsAndOverrides(t *testing.T) {
    t.Setenv("BILLING_FREE_MAX_PROJECTS", "")
    t.Setenv("BILLING_FREE_MAX_MIGRATIONS", "")
    t.Setenv("BILLING_PRO_MAX_PROJECTS", "")
    t.Setenv("BILLING_PRO_MAX_MIGRATIONS", "")

    freeProjects, freeMigrations := billingLimits("free")
    if freeProjects != defaultFreeMaxProjects || freeMigrations != defaultFreeMaxMigrations {
        t.Fatalf("unexpected free defaults: %d/%d", freeProjects, freeMigrations)
    }

    t.Setenv("BILLING_PRO_MAX_PROJECTS", "7")
    t.Setenv("BILLING_PRO_MAX_MIGRATIONS", "42")
    proProjects, proMigrations := billingLimits("pro")
    if proProjects != 7 || proMigrations != 42 {
        t.Fatalf("unexpected pro overrides: %d/%d", proProjects, proMigrations)
    }

    t.Setenv("BILLING_PRO_MAX_PROJECTS", "-1")
    if projects, _ := billingLimits("pro"); projects != defaultProMaxProjects {
        t.Fatalf("negative limit should use the default: %d", projects)
    }

    _ = os.Getenv("BILLING_FREE_MAX_PROJECTS")
}
