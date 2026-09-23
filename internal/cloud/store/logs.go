package store

func SaveLog(jobID string, log string) error {
    _, err := DB.Exec(`
        INSERT INTO logs (id, job_id, log)
        VALUES (hex(randomblob(16)), ?, ?)
    `, jobID, log)

    return err
}

func LoadLogs(jobID string) ([]string, error) {
    rows, err := DB.Query(`
        SELECT log FROM logs WHERE job_id = ?
    `, jobID)
    if err != nil {
        return nil, err
    }

    logs := []string{}
    for rows.Next() {
        var l string
        rows.Scan(&l)
        logs = append(logs, l)
    }

    return logs, nil
}
