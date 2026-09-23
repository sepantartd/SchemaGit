package cloud

import (
    "encoding/json"
    "net/http"
    "strconv"
)

func AgentLogsAPI(w http.ResponseWriter, r *http.Request) {
    pidStr := r.URL.Query().Get("project_id")
    pid, _ := strconv.ParseInt(pidStr, 10, 64)

    logs := projStore.LoadAgentLogsByProject(pid)

    raw, _ := json.Marshal(logs)
    w.Write(raw)
}
