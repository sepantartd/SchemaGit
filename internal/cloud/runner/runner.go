package runner

import "github.com/google/uuid"

import "schemagit/internal/cloud/store"

type Runner struct {
    Queue    *Queue
    Executor *Executor
}

func NewRunner(agentURL string) *Runner {
    return &Runner{
        Queue:    NewQueue(),
        Executor: NewExecutor(agentURL),
    }
}

func (r *Runner) RunMigration(projectID string, plan []string) (*Job, error) {
    job := Job{
        ID:        uuid.New().String(),
        ProjectID: projectID,
        Plan:      plan,
        Status:    Pending,
    }

    store.SaveJob(&job)

    r.Queue.Add(job)

    next := r.Queue.Next()
    if next == nil {
        return nil, nil
    }

    err := r.Executor.Execute(next)

    store.SaveJob(next)

    return next, err
}
