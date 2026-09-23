package store

import (
    "encoding/json"
    "schemagit/internal/schema/planner"
)

func SavePlan(repo string, prNumber int, plan planner.MigrationPlan) error {
    data, _ := json.Marshal(plan)

    _, err := DB.Exec(`
        INSERT INTO plans (id, repo, pr_number, plan_json)
        VALUES (?, ?, ?, ?)
    `, repo+"#"+string(prNumber), repo, prNumber, string(data))

    return err
}

func LoadPlan(repo string, prNumber int) (*planner.MigrationPlan, error) {
    row := DB.QueryRow(`
        SELECT plan_json FROM plans WHERE repo = ? AND pr_number = ?
    `, repo, prNumber)

    var jsonStr string
    err := row.Scan(&jsonStr)
    if err != nil {
        return nil, err
    }

    var plan planner.MigrationPlan
    json.Unmarshal([]byte(jsonStr), &plan)
    return &plan, nil
}
