package app

import (
	"context"
	"errors"
	"github.com/Nuyoahch/gopulse/architecture /ddd/internal/task/domain"
	"time"
)

// 这里你可以自己换成你项目的错误处理方式
var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrTaskAlreadyExists = errors.New("task already exists")
)

// TaskService 应用服务：组合多个领域对象 / 仓储，完成一个用例
type TaskService struct {
	repo    domain.TaskRepository
	clockFn func() time.Time // 方便测试注入 fake time
}

func NewTaskService(repo domain.TaskRepository) *TaskService {
	return &TaskService{
		repo:    repo,
		clockFn: time.Now,
	}
}

// CreateTaskCmd 创建任务用例的输入参数
type CreateTaskCmd struct {
	ID          string
	Name        string
	Description string
}

// CreateTask 创建一个新任务
func (s *TaskService) CreateTask(ctx context.Context, cmd CreateTaskCmd) (*domain.Task, error) {
	// 简单示例：先检查是否已存在同 ID 任务
	existing, _ := s.repo.FindByID(domain.TaskID(cmd.ID))
	if existing != nil {
		return nil, ErrTaskAlreadyExists
	}

	task, err := domain.NewTask(
		domain.TaskID(cmd.ID),
		cmd.Name,
		cmd.Description,
		s.clockFn(),
	)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(task); err != nil {
		return nil, err
	}

	// 这里将来可以在任务创建后，调用你的 GMP 调度器，把 Task 投递进去
	// 比如：s.scheduler.Submit(task.ToG())

	return task, nil
}

// StartTask 开始执行任务
func (s *TaskService) StartTask(ctx context.Context, id string) (*domain.Task, error) {
	task, err := s.repo.FindByID(domain.TaskID(id))
	if err != nil || task == nil {
		return nil, ErrTaskNotFound
	}

	if err := task.Start(s.clockFn()); err != nil {
		return nil, err
	}
	if err := s.repo.Save(task); err != nil {
		return nil, err
	}
	return task, nil
}

// CompleteTask 完成任务
func (s *TaskService) CompleteTask(ctx context.Context, id string) (*domain.Task, error) {
	task, err := s.repo.FindByID(domain.TaskID(id))
	if err != nil || task == nil {
		return nil, ErrTaskNotFound
	}

	if err := task.Complete(s.clockFn()); err != nil {
		return nil, err
	}
	if err := s.repo.Save(task); err != nil {
		return nil, err
	}
	return task, nil
}

// ListTasks 查询任务列表
func (s *TaskService) ListTasks(ctx context.Context) ([]*domain.Task, error) {
	return s.repo.ListAll()
}
