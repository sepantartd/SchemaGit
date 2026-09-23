package planner

import (
	"schemagit/internal/schema/diff"
	"schemagit/internal/schema/generator"
)

type PlanStep struct {
	Statement string
	Rollback  string
	Type      string // "create_table", "drop_table", "add_column", "drop_column", "modify_column"
}

type MigrationPlan struct {
	Steps []PlanStep
}

func BuildPlan(d diff.SchemaDiff, mig generator.Migration) MigrationPlan {
	plan := MigrationPlan{}

	for _, stmt := range mig.Statements {
		step := PlanStep{
			Statement: stmt,
			Rollback:  buildRollback(stmt),
			Type:      classify(stmt),
		}
		plan.Steps = append(plan.Steps, step)
	}

	return plan
}

func classify(stmt string) string {
	if startsWith(stmt, "CREATE TABLE") {
		return "create_table"
	}
	if startsWith(stmt, "DROP TABLE") {
		return "drop_table"
	}
	if startsWith(stmt, "ALTER TABLE") && contains(stmt, "ADD COLUMN") {
		return "add_column"
	}
	if startsWith(stmt, "ALTER TABLE") && contains(stmt, "DROP COLUMN") {
		return "drop_column"
	}
	return "modify_column"
}

func buildRollback(stmt string) string {
	if startsWith(stmt, "CREATE TABLE") {
		return convert(stmt, "CREATE TABLE", "DROP TABLE")
	}
	if startsWith(stmt, "DROP TABLE") {
		return "-- rollback not possible for DROP TABLE"
	}
	if contains(stmt, "ADD COLUMN") {
		return convert(stmt, "ADD COLUMN", "DROP COLUMN")
	}
	if contains(stmt, "DROP COLUMN") {
		return "-- rollback not possible for DROP COLUMN"
	}
	return "-- rollback for modify_column not implemented"
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (indexOf(s, sub) != -1)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func convert(stmt, from, to string) string {
	return replace(stmt, from, to)
}

func replace(s, from, to string) string {
	idx := indexOf(s, from)
	if idx == -1 {
		return s
	}
	return s[:idx] + to + s[idx+len(from):]
}
