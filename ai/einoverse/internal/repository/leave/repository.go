package leave

import (
	"sync"
	"time"

	"github.com/Nuyoahch/einoverse/internal/domain/leave"
	"github.com/google/uuid"
)

// Repository 请假申请仓储接口
type Repository interface {
	Create(app *leave.Application) error
	GetByID(id string) (*leave.Application, error)
	GetByEmployeeID(employeeID string, offset, limit int) ([]leave.Application, int, error)
	Update(id string, app *leave.Application) error
	UpdateStatus(id string, status leave.ApplicationStatus) error
}

// InMemoryRepository 内存实现的仓储
type InMemoryRepository struct {
	apps  map[string]*leave.Application
	mutex sync.RWMutex
}

// NewInMemoryRepository 创建内存仓储
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		apps: make(map[string]*leave.Application),
	}
}

// Create 创建请假申请
func (r *InMemoryRepository) Create(app *leave.Application) error {
	// 控制并发操作
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 如果申请ID为空，则生成一个UUID
	if app.ID == "" {
		app.ID = uuid.New().String()
	}
	// 如果创建时间为空，则设置为当前时间
	now := time.Now()
	if app.CreatedAt.IsZero() {
		app.CreatedAt = now
	}
	// 如果更新时间为空，则设置为当前时间
	app.UpdatedAt = now

	r.apps[app.ID] = app
	return nil
}

// GetByID 根据ID获取请假申请
func (r *InMemoryRepository) GetByID(id string) (*leave.Application, error) {
	// 控制并发操作
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	app, exists := r.apps[id]
	if !exists {
		// 如果申请不存在，则返回空
		return nil, nil
	}

	// 返回副本
	appCopy := *app
	return &appCopy, nil
}

// GetByEmployeeID 根据EmployeeID获取请假申请
func (r *InMemoryRepository) GetByEmployeeID(employeeID string, offset, limit int) ([]leave.Application, int, error) {
	// 控制并发操作
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// 匹配ID
	var matched []leave.Application
	for _, app := range r.apps {
		if app.EmployeeID == employeeID {
			matched = append(matched, *app)
		}
	}

	total := len(matched)
	if offset >= total {
		return []leave.Application{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return matched[offset:end], total, nil
}

// Update 更新操作
func (r *InMemoryRepository) Update(id string, app *leave.Application) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.apps[id]; !exists {
		return nil
	}

	app.UpdatedAt = time.Now()
	r.apps[id] = app
	return nil
}

// UpdateStatus 更新状态
func (r *InMemoryRepository) UpdateStatus(id string, status leave.ApplicationStatus) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	app, exists := r.apps[id]
	if !exists {
		return nil
	}

	app.Status = status
	app.UpdatedAt = time.Now()
	return nil
}
