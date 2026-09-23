package agent

import (
    "fmt"
    "time"

    "github.com/sepanta/schemagit/internal/cloud"
    "github.com/sepanta/schemagit/internal/log"
    "github.com/sepanta/schemagit/internal/migrate"
    "github.com/sepanta/schemagit/internal/schema"
)

type Agent struct {
    APIKey  string
    WorkDir string
    Token   string
}

func NewAgent(apiKey, workDir, token string) *Agent {
    return &Agent{APIKey: apiKey, WorkDir: workDir, Token: token}
}

func (a *Agent) Start() {
    log.Info("Agent started")

    for {
        job := cloud.FetchNextJob(a.APIKey)
        if job == nil {
            time.Sleep(2 * time.Second)
            continue
        }

        cloud.MarkJobRunning(a.APIKey, job.ID)

        err := a.runJob(job)
        if err != nil {
            cloud.MarkJobFailed(a.APIKey, job.ID, err.Error())
        } else {
            cloud.MarkJobDone(a.APIKey, job.ID)
        }
    }
}

func (a *Agent) runJob(job *store.Job) error {
    switch job.Type {

    case "github_push":
        return a.handlePush(job)

    case "github_pr_diff":
        return a.handlePRDiff(job)

    case "github_pr_apply":
        return a.handlePRApply(job)

    default:
        return fmt.Errorf("unknown job type: %s", job.Type)
    }
}

func (a *Agent) handlePush(job *store.Job) error {
    repo := job.Payload["repo"]
    branch := job.Payload["branch"]
    projectID := job.PayloadInt("project_id")

    err := cloud.GitCloneOrPull(repo, branch, a.Token, a.WorkDir)
    if err != nil {
        return err
    }

    desired, err := schema.LoadDesiredSchemaFromFile(a.WorkDir + "/schema.sql")
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

func (a *Agent) handlePRDiff(job *store.Job) error {
    repo := job.Payload["repo"]
    branch := job.Payload["branch"]
    projectID := job.PayloadInt("project_id")

    err := cloud.GitCloneOrPull(repo, branch, a.Token, a.WorkDir)
    if err != nil {
        return err
    }

    desired, err := schema.LoadDesiredSchemaFromFile(a.WorkDir + "/schema.sql")
    if err != nil {
        return err
    }

    current, err := schema.LoadCurrentSchema(projectID)
    if err != nil {
        return err
    }

    diff := schema.Diff(current, desired)

    cloud.SaveCloudLog(store.CloudLog{
        Type:      "pr_diff",
        Success:   true,
        Error:     "",
        Timestamp: time.Now().Unix(),
        ProjectID: projectID,
    })

    return nil
}

func (a *Agent) handlePRApply(job *store.Job) error {
    repo := job.Payload["repo"]
    branch := job.Payload["branch"]
    projectID := job.PayloadInt("project_id")

    err := cloud.GitCloneOrPull(repo, branch, a.Token, a.WorkDir)
    if err != nil {
        return err
    }

    desired, err := schema.LoadDesiredSchemaFromFile(a.WorkDir + "/schema.sql")
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
