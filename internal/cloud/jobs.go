package cloud

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/sepanta/schemagit/internal/store"
)

func CreateJob(apiKey, jobType string, payload map[string]string) {
    j := store.Job{
        Type:    jobType,
        Payload: payload,
    }
    projStore.CreateJob(j)
}

func FetchNextJob(apiKey string) *store.Job {
    job, err := projStore.FetchNextPendingJob()
    if err != nil {
        return nil
    }
    return job
}

func MarkJobRunning(apiKey string, id int64) {
    projStore.UpdateJobStatus(id, "running", "")
}

func MarkJobDone(apiKey string, id int64) {
    projStore.UpdateJobStatus(id, "done", "")
}

func MarkJobFailed(apiKey string, id int64, err string) {
    projStore.UpdateJobStatus(id, "failed", err)
    projStore.IncrementJobAttempts(id)
}
