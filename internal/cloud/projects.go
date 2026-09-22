package cloud

import "github.com/sepanta/schemagit/internal/store"

func GetProjectByID(id int64) *store.Project {
    projects := GetProjects()
    for _, p := range projects {
        if p.ID == id {
            return &p
        }
    }
    return nil
}

func UpdateProject(p *store.Project) error {
    return projStore.UpdateProject(*p)
}

func DeleteProject(id int64) error {
    _ = projStore.DeleteProject(id)
    _ = logsStore.DeleteLogsByProject(id)
    _ = projStore.DeleteMembersByProject(id)
    _ = projStore.DeleteWebhooksByProject(id)
    return nil
}

func GetMembersByProject(id int64) []store.ProjectMember {
    return projStore.LoadMembersByProject(id)
}

func AddMemberToProject(id int64, email, role string) error {
    m := store.ProjectMember{
        ProjectID: id,
        Email:     email,
        Role:      role,
    }
    return projStore.AddProjectMember(m)
}

func DeleteMember(memberID int64) error {
    return projStore.DeleteMember(memberID)
}

func GetWebhooksByProject(id int64) []store.ProjectWebhook {
    return projStore.LoadWebhooksByProject(id)
}

func AddWebhookToProject(id int64, url, secret, wtype string) error {
    w := store.ProjectWebhook{
        ProjectID: id,
        URL:       url,
        Secret:    secret,
        Type:      wtype,
    }
    return projStore.AddWebhook(w)
}

func DeleteWebhook(id int64) error {
    return projStore.DeleteWebhook(id)
}
