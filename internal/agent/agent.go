package agent

import (
    "fmt"
    "time"

    "github.com/sepanta/schemagit/internal/cloud"
    "github.com/sepanta/schemagit/internal/log"
    "github.com/sepanta/schemagit/internal/migrate"
    "github.com/sepanta/schemagit/internal/schema"
    "github.com/sepanta/schemagit/internal/store"
)

type Agent struct {
    cfg *AgentConfig
}

func NewAgent(cfg *AgentConfig) *Agent {
    return &Agent{cfg: cfg}
}

func (a *Agent) Start() {
    log.Info("Agent started")

    for {
        job := cloud.FetchNextJob(a.cfg.APIKey)
        if job == nil {
            time.Sleep(2 * time.Second)
            continue
        }

        // SECURITY: Check project access
        projectID := job.PayloadInt("project_id")
        if !a.isProjectAllowed(projectID) {
            cloud.MarkJobFailed(a.cfg.APIKey, job.ID, "agent not allowed for this project")
            continue
        }

        cloud.MarkJobRunning(a.cfg.APIKey, job.ID)

        err := a.runJob(job)
        if err != nil {
            cloud.MarkJobFailed(a.cfg.APIKey, job.ID, err.Error())
        } else {
            cloud.MarkJobDone(a.cfg.APIKey, job.ID)
        }
    }
}

func (a *Agent) isProjectAllowed(pid int64) bool {
    for _, allowed := range a.cfg.Projects {
        if allowed == pid {
            return true
        }
    }
    return false
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
