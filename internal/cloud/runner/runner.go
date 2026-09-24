package runner

import "github.com/google/uuid"

type Runner struct {
	Queue          *Queue
	Executor       *Executor
	CompletionHook func(*Job)
}

func NewRunner(agentURL string) *Runner {
	return &Runner{Queue: NewQueue(), Executor: NewExecutor(agentURL)}
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
