package easy_replication

import "github.com/go-mysql-org/go-mysql/mysql"

type CheckpointProvider interface {
	Exists() bool
	Get() mysql.Position
	Save(pos mysql.Position) error
	Flush()
}
