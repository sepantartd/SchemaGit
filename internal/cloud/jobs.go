package cloud

import (
    "encoding/json"
    "time"
)

type Job struct {
    ID      int64             `json:"id"`
    Type    string            `json:"type"`
    Payload map[string]string `json:"payload"`
    Created int64             `json:"created"`
}

func FetchNextJob(apiKey string) *Job {
    // Cloud API call
    res := apiCall("/jobs/next", apiKey)
    if res == nil {
        return nil
    }

    var job Job
    json.Unmarshal(res.Body, &job)
    return &job
}

func ReportJobSuccess(apiKey string, jobID int64) {
    apiCall("/jobs/success?id="+fmtInt(jobID), apiKey)
}

func ReportJobFailure(apiKey string, jobID int64, err string) {
    apiCall("/jobs/failure?id="+fmtInt(jobID)+"&error="+err, apiKey)
}
