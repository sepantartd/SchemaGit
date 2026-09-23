package cloud

import (
    "net/http"
    "strconv"
)

func ProjectAgentLogs(w http.ResponseWriter, r *http.Request) {
    pidStr := r.URL.Query().Get("project_id")
    pid, _ := strconv.ParseInt(pidStr, 10, 64)

    logs := projStore.LoadAgentLogsByProject(pid)

    render(w, "project_logs_agent.html", map[string]interface{}{
        "AgentLogs": logs,
    })
}
