package webhook_processor

import (
    "fmt"
    "schemagit/internal/schema/diff"
    "schemagit/internal/schema/generator"
    "schemagit/internal/schema/validator"
    "schemagit/internal/schema/planner"
    "schemagit/internal/cloud/runner"
    "schemagit/internal/cloud/github"
    "schemagit/internal/cloud/store"
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

    oldSchema, newSchema := p.GitHub.FetchSchemas(repo, prNumber, commitSHA)

    d := diff.Diff(oldSchema, newSchema)

    validation := validator.Validate(d)
    for _, issue := range validation.Issues {
        if issue.Severity == "error" {
            store.SavePlanError(repo, prNumber, issue.Message)
            return fmt.Errorf("unsafe migration: %s", issue.Message)
        }
    }

    mig := generator.Generate(d)
    plan := planner.BuildPlan(d, mig)

    store.SavePlan(repo, prNumber, plan)

    fmt.Println("Webhook processing completed.")
    return nil
}
