package diff

import "github.com/sepanta/schemagit/internal/schema"

type ChangeType string

const (
    AddTable    ChangeType = "ADD_TABLE"
    DropTable   ChangeType = "DROP_TABLE"
    AddColumn   ChangeType = "ADD_COLUMN"
    DropColumn  ChangeType = "DROP_COLUMN"
    ModifyColumn ChangeType = "MODIFY_COLUMN"
)

type TableChange struct {
    Type      ChangeType
    TableName string
    Column    *schema.Column
}

type SchemaDiff struct {
    Changes []TableChange
}

func NewDiff() *SchemaDiff {
    return &SchemaDiff{
        Changes: []TableChange{},
    }
}

func (d *SchemaDiff) Add(change TableChange) {
    d.Changes = append(d.Changes, change)
}
