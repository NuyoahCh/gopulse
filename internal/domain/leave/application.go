package leave

import "time"

// ApplicationStatus 申请状态
type ApplicationStatus string

// 申请状态
const (
	StatusPending  ApplicationStatus = "pending"  // 待审批
	StatusApproved ApplicationStatus = "approved" // 已批准
	StatusRejected ApplicationStatus = "rejected" // 已拒绝
)

// LeaveType 请假类型
type LeaveType string

// 请假类型
const (
	LeaveTypeSick     LeaveType = "sick"     // 病假
	LeaveTypePersonal LeaveType = "personal" // 事假
	LeaveTypeAnnual   LeaveType = "annual"   // 年假
	LeaveTypeAdjust   LeaveType = "adjust"   // 调休
)

// Application 请假申请
type Application struct {
	ID           string            `json:"id"`
	EmployeeID   string            `json:"employee_id"`
	EmployeeName string            `json:"employee_name"`
	Supervisor   string            `json:"supervisor"`
	LeaveType    LeaveType         `json:"leave_type"`
	StartDate    time.Time         `json:"start_date"`
	EndDate      time.Time         `json:"end_date"`
	Days         float64           `json:"days"`
	Reason       string            `json:"reason"`
	WorkHandover string            `json:"work_handover,omitempty"`
	Status       ApplicationStatus `json:"status"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// CreateRequest 创建申请请求
type CreateRequest struct {
	EmployeeID   string `json:"employee_id" binding:"required"`
	EmployeeName string `json:"employee_name" binding:"required"`
	Supervisor   string `json:"supervisor" binding:"required"`
	Text         string `json:"text" binding:"required"`
}
