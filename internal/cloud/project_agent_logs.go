package cloud

import (
	"html/template"
	"net/http"
	"strconv"
)

var projectAgentLogsTpl = template.Must(template.ParseFiles("dashboard/project_logs_agent.html"))

// ProjectAgentLogs renders the agent logs for a project.
func ProjectAgentLogs(w http.ResponseWriter, r *http.Request) {
	pid := parseID(r)

	logs := projStore.LoadAgentLogsByProject(pid)

	render(w, "project_logs_agent.html", map[string]interface{}{
		"AgentLogs": logs,
	})
}

func parseID(r *http.Request) int64 {
	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	return id
}

func render(w http.ResponseWriter, name string, data interface{}) {
	if name == "project_logs_agent.html" {
		if err := projectAgentLogsTpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "template not found", http.StatusNotFound)
}

func init() {
	http.HandleFunc("/project/logs/agent", ProjectAgentLogs)
}
