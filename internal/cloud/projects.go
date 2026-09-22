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
