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
