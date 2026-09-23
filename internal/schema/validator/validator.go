package validator

import (
    "fmt"
    "schemagit/internal/schema/diff"
)

type ValidationIssue struct {
    Severity string // "warning" or "error"
    Message  string
}

type ValidationResult struct {
    Issues []ValidationIssue
}

func Validate(d diff.SchemaDiff) ValidationResult {
    result := ValidationResult{}

    // حذف جدول
    for _, t := range d.RemovedTables {
        result.Issues = append(result.Issues, ValidationIssue{
            Severity: "warning",
            Message:  fmt.Sprintf("Dropping table %s may cause data loss", t.Name),
        })
    }

    // تغییرات در جدول‌های موجود
    for _, mt := range d.ModifiedTables {

        // ستون‌های حذف‌شده
        for _, c := range mt.RemovedColumns {
            result.Issues = append(result.Issues, ValidationIssue{
                Severity: "warning",
                Message:  fmt.Sprintf("Dropping column %s.%s may cause data loss", mt.Name, c.Name),
            })
        }

        // ستون‌های جدید
        for _, c := range mt.AddedColumns {
            if !c.Nullable && c.Default == "" {
                result.Issues = append(result.Issues, ValidationIssue{
                    Severity: "error",
                    Message:  fmt.Sprintf("Adding NOT NULL column %s.%s without default will fail", mt.Name, c.Name),
                })
            }
        }

        // ستون‌های تغییرکرده
        for _, c := range mt.ModifiedColumns {
            result.Issues = append(result.Issues, ValidationIssue{
                Severity: "warning",
                Message:  fmt.Sprintf("Modifying column %s.%s may cause data inconsistency", mt.Name, c.Name),
            })
        }
    }

    return result
}
