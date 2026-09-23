package store

import "schemagit/internal/cloud/runner"

func SaveJob(job *runner.Job) error {
    _, err := DB.Exec(`
        INSERT OR REPLACE INTO jobs (id, project_id, status, result)
        VALUES (?, ?, ?, ?)
    `, job.ID, job.ProjectID, job.Status, job.Result)

    return err
}

func LoadJob(id string) (*runner.Job, error) {
    row := DB.QueryRow(`
        SELECT project_id, status, result FROM jobs WHERE id = ?
    `, id)

    var projectID, status, result string
    err := row.Scan(&projectID, &status, &result)
    if err != nil {
        return nil, err
    }

    return &runner.Job{
        ID:        id,
        ProjectID: projectID,
        Status:    runner.JobStatus(status),
        Result:    result,
    }, nil
}
