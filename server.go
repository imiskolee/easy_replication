package easy_replication

import (
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

type ReplicationServer struct {
	checkpointProvider CheckpointProvider
	config             canal.Config
	canal              *canal.Canal
	tasks              []Task
}

func (s *ReplicationServer) OnRotate(roateEvent *replication.RotateEvent) error {
	return nil
}

func (s *ReplicationServer) OnTableChanged(schema string, table string) error {
	return nil
}

func (s *ReplicationServer) OnDDL(nextPos mysql.Position, queryEvent *replication.QueryEvent) error {
	return nil
}

func (s *ReplicationServer) OnRow(e *canal.RowsEvent) error {
	for _, task := range s.tasks {
		s.doTask(e, &task)
	}
	return nil
}

func (s *ReplicationServer) OnXID(nextPos mysql.Position) error {
	return nil
}

func (s *ReplicationServer) OnGTID(gtid mysql.GTIDSet) error {
	return nil
}

func (s *ReplicationServer) OnPosSynced(pos mysql.Position, set mysql.GTIDSet, force bool) error {
	return s.checkpointProvider.Save(pos)
}

func (s *ReplicationServer) String() string {
	return ""
}

func (s *ReplicationServer) AddTask(task Task) *ReplicationServer {
	s.tasks = append(s.tasks, task)
	return s
}

func (s *ReplicationServer) doTask(e *canal.RowsEvent, task *Task) error {
	event := &Event{
		Action:   e.Action,
		Schema:   e.Table.Schema,
		Table:    e.Table.Name,
		RawEvent: *e,
	}
	if task.Ack {
		for i := 0; i < task.Retry; i++ {
			if err := task.Handler(event); err == nil {
				break
			}
		}
	}
	return nil
}

func (s *ReplicationServer) Init() {
	instance, err := canal.NewCanal(&s.config)
	if err != nil {
		panic("Can not Init Canal Server: " + err.Error())
	}
	s.canal = instance
	s.canal.SetEventHandler(s)
}

func (s *ReplicationServer) Run() {
	if !s.checkpointProvider.Exists() {
		s.canal.Run()
		return
	}

	pos := s.checkpointProvider.Get()
	s.canal.RunFrom(pos)
}
