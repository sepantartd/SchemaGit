package webhook_processor

import (
    "fmt"
    "schemagit/internal/schema/diff"
    "schemagit/internal/schema/generator"
    "schemagit/internal/schema/validator"
    "schemagit/internal/schema/planner"
    "schemagit/internal/cloud/runner"
)

type Processor struct {
    Runner *runner.Runner
}

func NewProcessor(r *runner.Runner) *Processor {
    return &Processor{Runner: r}
}

func (p *Processor) HandlePullRequestEvent(repo string, prNumber int, commitSHA string) error {
    fmt.Printf("Processing PR #%d from %s (commit: %s)\n", prNumber, repo, commitSHA)

    // مرحله ۱: گرفتن اسکیماها از ریپو
    oldSchema, newSchema := fetchSchemasFromRepo(repo, prNumber, commitSHA)

    // مرحله ۲: Diff
    d := diff.Diff(oldSchema, newSchema)

    // مرحله ۳: Validate
    validation := validator.Validate(d)
    for _, issue := range validation.Issues {
        if issue.Severity == "error" {
            fmt.Println("Migration unsafe:", issue.Message)
            return fmt.Errorf("unsafe migration: %s", issue.Message)
        }
    }

    // مرحله ۴: Generate Migration
    mig := generator.Generate(d)

    // مرحله ۵: Build Plan
    plan := planner.BuildPlan(d, mig)

    // مرحله ۶: ذخیرهٔ نتیجه در Cloud (TODO)
    storePlan(repo, prNumber, plan)

    // مرحله ۷: اجرای Migration (اختیاری)
    // job, err := p.Runner.RunMigration(repo, extractStatements(plan))
    // if err != nil {
    //     return err
    // }

    fmt.Println("Webhook processing completed.")
    return nil
}

func fetchSchemasFromRepo(repo string, prNumber int, commitSHA string) ([]diff.Table, []diff.Table) {
    // TODO: اتصال به GitHub API
    fmt.Println("Fetching schemas from repo:", repo)
    return []diff.Table{}, []diff.Table{}
}

func storePlan(repo string, prNumber int, plan planner.MigrationPlan) {
    // TODO: ذخیره در دیتابیس Cloud
    fmt.Printf("Storing plan for PR #%d in repo %s\n", prNumber, repo)
}

func extractStatements(plan planner.MigrationPlan) []string {
    stmts := []string{}
    for _, step := range plan.Steps {
        stmts = append(stmts, step.Statement)
    }
    return stmts
}
