package cloud

import "github.com/sepanta/schemagit/internal/store"

var projStore *store.Store

func InitProjects(path string) {
    st, err := store.Open(path)
    if err != nil {
        return
    }
    projStore = st
}

func GetProjects() []store.Project {
    if projStore == nil {
        return []store.Project{}
    }
    return projStore.LoadProjects()
}

func AddProject(name string) {
    if projStore == nil {
        return
    }
    projStore.SaveProject(store.Project{
        Name: name,
    })
}
