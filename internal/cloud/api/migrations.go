package api

import (
	"encoding/json"
	"net/http"

	"schemagit/internal/cloud/runner"
	"schemagit/internal/cloud/store"
)

// Existing migration API now carries the organization context when supplied.
type RunMigrationRequest struct {
	OldSchema []diff.Table     `json:"old_schema"`
	NewSchema []diff.Table     `json:"new_schema"`
	Diff      *diff.SchemaDiff `json:"diff"`
	OrgID     string           `json:"org_id"`
}
