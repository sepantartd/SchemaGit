package runner

import (
    "encoding/json"

    "github.com/google/uuid"
    "github.com/sepanta/schemagit/internal/cloud/store"
    "github.com/sepanta/schemagit/internal/cloud/webhook_delivery"
)

type Runner struct {
    Queue          *Queue
    Executor       *Executor
    CompletionHook func(*Job)
}

func NewRunner(agentURL string) *Runner {
    return &Runner{
        Queue:    NewQueue(),
        Executor: NewExecutor(agentURL),
        CompletionHook: func(job *Job) {
            if job.OrgID == "" {
                return
            }
            payload, _ := json.Marshal(map[string]string{"project": job.ProjectID, "status": string(job.Status), "job": job.ID})
            owner, _ := store.OrgOwner(job.OrgID)
            go store.DispatchNotification(job.OrgID, "migration.completed", string(payload), owner)
            store.AddAudit(job.OrgID, "migration.completed", string(job.Status))
            go webhook_delivery.Dispatch(job.OrgID, "migration.completed", map[string]interface{}{
                "project": job.ProjectID,
                "status":  job.Status,
                "job":     job.ID,
            })
        },
    }
}

func (r *Runner) RunMigration(projectID string, plan []string) (*Job, error) {
    return r.RunMigrationForOrg(projectID, "", plan)
}

func (r *Runner) RunMigrationForOrg(projectID, orgID string, plan []string) (*Job, error) {
    job := Job{ID: uuid.New().String(), ProjectID: projectID, OrgID: orgID, Plan: plan, Status: Pending}
    r.Queue.Add(job)
    next := r.Queue.Next()
    if next == nil {
        return nil, nil
    }
    if err := r.Executor.Execute(next); err != nil {
        return next, err
    }
    if next.Status == Done && r.CompletionHook != nil {
        r.CompletionHook(next)
    }
    return next, nil
}
