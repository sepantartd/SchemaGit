package runner

type JobStatus string

const (
    Pending JobStatus = "pending"
    Running JobStatus = "running"
    Done    JobStatus = "done"
    Failed  JobStatus = "failed"
)

type Job struct {
    ID        string
    ProjectID string
    Plan      []string
    Status    JobStatus
    Result    string
}

type Queue struct {
    jobs []Job
}

func NewQueue() *Queue {
    return &Queue{jobs: []Job{}}
}

func (q *Queue) Add(job Job) {
    q.jobs = append(q.jobs, job)
}

func (q *Queue) Next() *Job {
    for i := range q.jobs {
        if q.jobs[i].Status == Pending {
            q.jobs[i].Status = Running
            return &q.jobs[i]
        }
    }
    return nil
}
