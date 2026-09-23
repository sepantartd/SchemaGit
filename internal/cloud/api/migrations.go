package api

import (
    "encoding/json"
    "net/http"

    "schemagit/internal/cloud/runner"
    "schemagit/internal/cloud/store"
    "schemagit/internal/schema/diff"
    "schemagit/internal/schema/generator"
    "schemagit/internal/schema/planner"
    "schemagit/internal/schema/validator"
)

type RunMigrationRequest struct {
    OldSchema []diff.Table     `json:"old_schema"`
    NewSchema []diff.Table     `json:"new_schema"`
    Diff      *diff.SchemaDiff `json:"diff"`
}

type RunMigrationResponse struct {
    JobID  string `json:"job_id"`
    Status string `json:"status"`
    Result string `json:"result,omitempty"`
}

func RunMigrationHandler(r *runner.Runner) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        token := req.Header.Get("Authorization")
        email, err := store.ValidateToken(token)
        if err != nil {
            http.Error(w, "invalid token", http.StatusUnauthorized)
            return
        }
        if !store.CanCreateMigration(email) {
            http.Error(w, "migration limit reached for your plan", http.StatusForbidden)
            return
        }

        var body RunMigrationRequest
        if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
            http.Error(w, "invalid request", http.StatusBadRequest)
            return
        }

        var d diff.SchemaDiff
        if body.Diff != nil {
            d = *body.Diff
        } else {
            d = diff.Diff(body.OldSchema, body.NewSchema)
        }

        validation := validator.Validate(d)
        for _, issue := range validation.Issues {
            if issue.Severity == "error" {
                http.Error(w, issue.Message, http.StatusBadRequest)
                return
            }
        }

        mig := generator.Generate(d)
        plan := planner.BuildPlan(d, mig)
        projectID := "project-id"
        job, err := r.RunMigration(projectID, extractStatements(plan))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        if err := store.RecordMigration(email, projectID); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(RunMigrationResponse{
            JobID: job.ID, Status: string(job.Status), Result: job.Result,
        })
    }
}

func extractStatements(plan planner.MigrationPlan) []string {
    stmts := []string{}
    for _, step := range plan.Steps {
        stmts = append(stmts, step.Statement)
    }
    return stmts
}
