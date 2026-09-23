package api

import (
    "encoding/json"
    "net/http"

    "schemagit/internal/schema/diff"
    "schemagit/internal/schema/generator"
    "schemagit/internal/schema/validator"
    "schemagit/internal/schema/planner"
    "schemagit/internal/cloud/runner"
    "schemagit/internal/cloud/store"
)

type RunMigrationRequest struct {
    OldSchema []diff.Table       `json:"old_schema"`
    NewSchema []diff.Table       `json:"new_schema"`
    Diff      *diff.SchemaDiff   `json:"diff"`
}

type RunMigrationResponse struct {
    JobID  string `json:"job_id"`
    Status string `json:"status"`
    Result string `json:"result,omitempty"`
}

func RunMigrationHandler(r *runner.Runner) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {

        var body RunMigrationRequest
        json.NewDecoder(req.Body).Decode(&body)

        // مرحله ۱: Diff
        var d diff.SchemaDiff
        if body.Diff != nil {
            d = *body.Diff
        } else {
            d = diff.Diff(body.OldSchema, body.NewSchema)
        }

        // مرحله ۲: Validate
        validation := validator.Validate(d)
        for _, issue := range validation.Issues {
            if issue.Severity == "error" {
                http.Error(w, issue.Message, http.StatusBadRequest)
                return
            }
        }

        // مرحله ۳: Generate Migration
        mig := generator.Generate(d)

        // مرحله ۴: Build Plan
        plan := planner.BuildPlan(d, mig)

        // مرحله ۵: Run Migration
        job, err := r.RunMigration("project-id", extractStatements(plan))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // مرحله ۶: پاسخ
        resp := RunMigrationResponse{
            JobID:  job.ID,
            Status: string(job.Status),
            Result: job.Result,
        }

        // Save job logs
        store.SaveLog(job.ID, job.Result)

        json.NewEncoder(w).Encode(resp)
    }
}

func extractStatements(plan planner.MigrationPlan) []string {
    stmts := []string{}
    for _, step := range plan.Steps {
        stmts = append(stmts, step.Statement)
    }
    return stmts
}
