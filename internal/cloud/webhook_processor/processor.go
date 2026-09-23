package webhook_processor

import (
    "fmt"
    "schemagit/internal/schema/diff"
    "schemagit/internal/schema/generator"
    "schemagit/internal/schema/validator"
    "schemagit/internal/schema/planner"
    "schemagit/internal/cloud/runner"
    "schemagit/internal/cloud/github"
)

type Processor struct {
    Runner *runner.Runner
    GitHub *github.GitHubClient
}

func NewProcessor(r *runner.Runner, gh *github.GitHubClient) *Processor {
    return &Processor{
        Runner: r,
        GitHub: gh,
    }
}

func (p *Processor) HandlePullRequestEvent(repo string, prNumber int, commitSHA string) error {
    fmt.Printf("Processing PR #%d from %s (commit: %s)\n", prNumber, repo, commitSHA)

    // مرحله ۱: گرفتن اسکیماها از GitHub
    oldSchema, newSchema := p.GitHub.FetchSchemas(repo, prNumber, commitSHA)

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

    // مرحله ۶: ذخیرهٔ نتیجه در Cloud
    storePlan(repo, prNumber, plan)

    fmt.Println("Webhook processing completed.")
    return nil
}

func storePlan(repo string, prNumber int, plan planner.MigrationPlan) {
    // TODO: ذخیره در دیتابیس Cloud
    fmt.Printf("Storing plan for PR #%d in repo %s\n", prNumber, repo)
}
