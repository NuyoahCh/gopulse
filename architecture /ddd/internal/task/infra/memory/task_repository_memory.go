package memory

import (
	"github.com/Nuyoahch/gopulse/architecture /ddd/internal/task/domain"
	"sync"
)

// InMemoryTaskRepository 一个简单的内存仓储实现，适合 demo / 单元测试
type InMemoryTaskRepository struct {
	mu    sync.RWMutex
	store map[domain.TaskID]*domain.Task
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		store: make(map[domain.TaskID]*domain.Task),
	}
}

func (r *InMemoryTaskRepository) Save(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[task.ID()] = task
	return nil
}

func (r *InMemoryTaskRepository) FindByID(id domain.TaskID) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if t, ok := r.store[id]; ok {
		return t, nil
	}
	return nil, nil
}

func (r *InMemoryTaskRepository) ListAll() ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.Task, 0, len(r.store))
	for _, t := range r.store {
		result = append(result, t)
	}
	return result, nil
}
