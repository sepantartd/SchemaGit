package runner

import "fmt"

type Executor struct {
    AgentURL string
}

func NewExecutor(agentURL string) *Executor {
    return &Executor{AgentURL: agentURL}
}

func (e *Executor) Execute(job *Job) error {
    // TODO: ارسال به Agent
    fmt.Println("Sending migration to agent:", e.AgentURL)

    // شبیه‌سازی نتیجه
    job.Status = Done
    job.Result = "Migration executed successfully"

    return nil
}
