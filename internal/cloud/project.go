package cloud

type Project struct {
    ID        int64
    Name      string
    DBType    string
    DBPath    string
    PGHost    string
    PGPort    int
    PGUser    string
    PGPass    string
    PGName    string
    MYHost    string
    MYPort    int
    MYUser    string
    MYPass    string
    MYName    string
}
