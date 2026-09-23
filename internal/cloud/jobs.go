package cloud

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"

    "github.com/sepanta/schemagit/internal/store"
)

func CreateJobAPI(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "method not allowed", 405)
        return
    }

    var body struct {
        Type    string            `json:"type"`
        Payload map[string]string `json:"payload"`
    }

    json.NewDecoder(r.Body).Decode(&body)

    job := store.Job{
        Type:    body.Type,
        Payload: body.Payload,
    }

    _ = projStore.CreateJob(job)

    w.Write([]byte(`{"ok":true}`))
}

func FetchNextJobAPI(w http.ResponseWriter, r *http.Request) {
    job, err := projStore.FetchNextPendingJob()
    if err != nil {
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

    id, _ := strconv.ParseInt(idStr, 10, 64)

    projStore.UpdateJobStatus(id, status, errMsg)

    w.Write([]byte(`{"ok":true}`))
}

func IncrementJobAttemptsAPI(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, _ := strconv.ParseInt(idStr, 10, 64)

    projStore.IncrementJobAttempts(id)

    w.Write([]byte(`{"ok":true}`))
}
