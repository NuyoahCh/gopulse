package domain

// TaskRepository 仓储接口：领域层只关心接口，不关心具体用 MySQL、Redis 还是内存
type TaskRepository interface {
	Save(task *Task) error
	FindByID(id TaskID) (*Task, error)
	ListAll() ([]*Task, error)
}
