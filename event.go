package easy_replication

import "github.com/go-mysql-org/go-mysql/canal"

type EventAction int

const (
	EventActionInsert = 1
	EventActionUpdate = 2
	EventActionDelete = 3
)

type Event struct {
	Action   string
	Schema   string
	Table    string
	RawEvent canal.RowsEvent
}
