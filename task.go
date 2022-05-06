package easy_replication

type RowHandler func(raw *Event) error

type Task struct {
	Name          string
	Ack           bool
	Retry         int
	IncludeTables []string //*  db.*  db.table_name  db.table*
	ExcludeTables []string
	Handler       RowHandler
}
