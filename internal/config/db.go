package config

import "github.com/sepanta/schemagit/internal/store"

func SetDBForProject(p *store.Project) {
    cfg := Default()

    cfg.DBType = p.DBType

    if p.DBType == "sqlite" {
        cfg.SQLitePath = p.SQLitePath
    }

    if p.DBType == "postgres" {
        cfg.PGHost = p.PGHost
        cfg.PGPort = p.PGPort
        cfg.PGUser = p.PGUser
        cfg.PGPass = p.PGPass
        cfg.PGName = p.PGName
    }

    if p.DBType == "mysql" {
        cfg.MYHost = p.MYHost
        cfg.MYPort = p.MYPort
        cfg.MYUser = p.MYUser
        cfg.MYPass = p.MYPass
        cfg.MYName = p.MYName
    }
}
