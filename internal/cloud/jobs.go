package cloud

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/sepanta/schemagit/internal/store"
)

var projStore *store.Store

// این تابع رو جایی که Cloud رو بالا میاری صدا می‌کنی
func InitStore(s *store.Store) {
    projStore = s
}

//
// APIهای HTTP برای مدیریت Jobها
//

func CreateJobAPI(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "method not allowed", 405)
        return
    }

    var body struct {
        Type    string            `json:"type"`
        Payload map[string]string `json:"payload"`
    }

    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, "invalid json", 400)
        return
    }

    job := store.Job{
        Type:    body.Type,
        Payload: body.Payload,
    }

    if err := projStore.CreateJob(job); err != nil {
        http.Error(w, "cannot create job", 500)
        return
    }

    w.Write([]byte(`{"ok":true}`))
}

func FetchNextJobAPI(w http.ResponseWriter, r *http.Request) {
    job, err := projStore.FetchNextPendingJob()
    if err != nil || job == nil {
        w.Write([]byte(`{"job":null}`))
        return
    }

    raw, _ := json.Marshal(job)
    w.Write(raw)
}

func UpdateJobStatusAPI(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    status := r.URL.Query().Get("status")
    errMsg := r.URL.Query().Get("error")

    if idStr == "" || status == "" {
        http.Error(w, "missing params", 400)
        return
    }

    id, _ := strconv.ParseInt(idStr, 10, 64)

    projStore.UpdateJobStatus(id, status, errMsg)

    w.Write([]byte(`{"ok":true}`))
}

func IncrementJobAttemptsAPI(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    if idStr == "" {
        http.Error(w, "missing id", 400)
        return
    }

    id, _ := strconv.ParseInt(idStr, 10, 64)

    projStore.IncrementJobAttempts(id)

    w.Write([]byte(`{"ok":true}`))
}

//
// توابعی که Agent استفاده می‌کند
//

func CreateJob(apiKey, jobType string, payload map[string]string) {
    // اینجا می‌تونی بعداً validate APIKey هم اضافه کنی
    j := store.Job{
        Type:    jobType,
        Payload: payload,
    }
    _ = projStore.CreateJob(j)
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
