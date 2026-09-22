package cloud

import (
    "time"

    "github.com/sepanta/schemagit/internal/log"
    "github.com/sepanta/schemagit/internal/store"
)

type CloudLog struct {
    ID        int64
    Type      string // diff / apply / webhook
    Success   bool
    Error     string
    Timestamp int64
}

var logsStore *store.Store

func InitLogs(path string) {
    st, err := store.Open(path)
    if err != nil {
        log.Error("cannot open cloud logs store: " + err.Error())
        return
    }
    logsStore = st
}

func AddLog(t string, success bool, errMsg string) {
    if logsStore == nil {
        return
    }

    entry := CloudLog{
        Type:      t,
        Success:   success,
        Error:     errMsg,
        Timestamp: time.Now().Unix(),
    }

    logsStore.SaveCloudLog(entry)
}

func GetLogs() []CloudLog {
    if logsStore == nil {
        return []CloudLog{}
    }
    return logsStore.LoadCloudLogs()
}
