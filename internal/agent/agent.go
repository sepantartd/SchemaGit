package agent

import (
    "fmt"
    "time"

    "github.com/sepanta/schemagit/internal/cloud"
    "github.com/sepanta/schemagit/internal/log"
    "github.com/sepanta/schemagit/internal/schema"
    "github.com/sepanta/schemagit/internal/migrate"
)

type Agent struct {
    APIKey string
}

func NewAgent(apiKey string) *Agent {
    return &Agent{APIKey: apiKey}
}

func (a *Agent) Start() {
    log.Info("Agent started")

    for {
        job := cloud.FetchNextJob(a.APIKey)
        if job == nil {
            time.Sleep(2 * time.Second)
            continue
        }

        log.Info("Running job: " + job.Type)

        cloud.MarkJobRunning(a.APIKey, job.ID)

        err := a.runJob(job)
        if err != nil {
            cloud.MarkJobFailed(a.APIKey, job.ID, err.Error())
        } else {
            cloud.MarkJobDone(a.APIKey, job.ID)
        }
    }
}

func (a *Agent) runJob(job *cloud.Job) error {
    switch job.Type {

    case "github_push":
        return a.handleGitHubPush(job)

    case "manual_apply":
        return a.handleManualApply(job)

    default:
        return fmt.Errorf("unknown job type: %s", job.Type)
    }
}

func (a *Agent) handleGitHubPush(job *cloud.Job) error {
    repo := job.Payload["repo"]
    branch := job.Payload["branch"]
    projectID := job.PayloadInt("project_id")

    // Clone/Pull Repo
    err := cloud.CloneOrPullRepo(repo, branch)
    if err != nil {
        return err
    }

    desired, err := schema.LoadDesiredSchemaFromFile("schema.sql")
    if err != nil {
        return err
    }

    current, err := schema.LoadCurrentSchema(projectID)
    if err != nil {
        return err
    }

    diff := schema.Diff(current, desired)
    if diff.IsEmpty() {
        return nil
    }

    mig, err := migrate.Generate(diff)
    if err != nil {
        return err
    }

    return migrate.Apply(projectID, mig)
}

func (a *Agent) handleManualApply(job *cloud.Job) error {
    projectID := job.PayloadInt("project_id")
    schemaPath := job.Payload["schema_path"]

    desired, err := schema.LoadDesiredSchemaFromFile(schemaPath)
    if err != nil {
        return err
    }

    current, err := schema.LoadCurrentSchema(projectID)
    if err != nil {
        return err
    }

    diff := schema.Diff(current, desired)
    if diff.IsEmpty() {
        return nil
    }

    mig, err := migrate.Generate(diff)
    if err != nil {
        return err
    }

    return migrate.Apply(projectID, mig)
}        log.Info("Processing job: " + job.Type)

        err := a.processJob(job)
        if err != nil {
            cloud.ReportJobFailure(a.APIKey, job.ID, err.Error())
        } else {
            cloud.ReportJobSuccess(a.APIKey, job.ID)
        }
    }
}

func (a *Agent) processJob(job *cloud.Job) error {
    switch job.Type {

    case "github_push":
        return a.handleGitHubPush(job)

    case "manual_apply":
        return a.handleManualApply(job)

    default:
        return fmt.Errorf("unknown job type: %s", job.Type)
    }
}

func (a *Agent) handleGitHubPush(job *cloud.Job) error {
    repo := job.Payload["repo"]
    branch := job.Payload["branch"]
    projectID := int64(job.PayloadInt("project_id"))

    // Clone or pull
    err := a.cloneOrPullRepo(repo, branch)
    if err != nil {
        return err
    }

    // Load desired schema
    desired, err := schema.LoadDesiredSchemaFromFile(a.WorkDir + "/schema.sql")
    if err != nil {
        return err
    }

    // Load current schema
    current, err := schema.LoadCurrentSchema(projectID)
    if err != nil {
        return err
    }

    // Diff
    diff := schema.Diff(current, desired)
    if diff.IsEmpty() {
        log.Info("No changes detected")
        return nil
    }

    // Generate migration
    mig, err := migrate.Generate(diff)
    if err != nil {
        return err
    }

    // Apply migration
    err = migrate.Apply(projectID, mig)
    if err != nil {
        return err
    }

    return nil
}

func (a *Agent) handleManualApply(job *cloud.Job) error {
    projectID := int64(job.PayloadInt("project_id"))

    desired, err := schema.LoadDesiredSchemaFromFile(job.Payload["schema_path"])
    if err != nil {
        return err
    }

    current, err := schema.LoadCurrentSchema(projectID)
    if err != nil {
        return err
    }

    diff := schema.Diff(current, desired)
    if diff.IsEmpty() {
        return nil
    }

    mig, err := migrate.Generate(diff)
    if err != nil {
        return err
    }

    return migrate.Apply(projectID, mig)
}

func (a *Agent) cloneOrPullRepo(repo, branch string) error {
    if _, err := os.Stat(a.WorkDir); os.IsNotExist(err) {
        cmd := exec.Command("git", "clone", "-b", branch, repo, a.WorkDir)
        return cmd.Run()
    }

    cmd := exec.Command("git", "-C", a.WorkDir, "pull")
    return cmd.Run()
}
