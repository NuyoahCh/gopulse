package leave

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Nuyoahch/einoverse/internal/domain/leave"
	leaveRepo "github.com/Nuyoahch/einoverse/internal/repository/leave"
	"github.com/Nuyoahch/einoverse/pkg/eino"
	"github.com/Nuyoahch/einoverse/pkg/errors"
	"go.uber.org/zap"
)

// Service 请假服务
type Service struct {
	repo   leaveRepo.Repository
	eino   *eino.Client
	logger *zap.Logger
}

// NewService 创建请假服务
func NewService(repo leaveRepo.Repository, einoClient *eino.Client, logger *zap.Logger) *Service {
	return &Service{
		repo:   repo,
		eino:   einoClient,
		logger: logger,
	}
}


// CreateApplication 从文本创建请假申请
func (s *Service) CreateApplication(req *leave.CreateRequest) (*leave.Application, error) {
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// 调用 LLM 从文本提取结构化信息
	app, err := s.extractApplicationFromText(req)
	if err != nil {
		s.logger.Error("failed to extract application from text", zap.Error(err))
		return nil, fmt.Errorf("extract application failed: %w", err)
	}

	// 设置状态
	app.Status = leave.StatusPending

	// 保存到仓储
	if err := s.repo.Create(app); err != nil {
		s.logger.Error("failed to create application", zap.Error(err))
		return nil, fmt.Errorf("create application failed: %w", err)
	}

	s.logger.Info("application created",
		zap.String("id", app.ID),
		zap.String("employee_id", app.EmployeeID))

	return app, nil
}

// GetApplication 获取申请
func (s *Service) GetApplication(id string) (*leave.Application, error) {
	app, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("failed to get application", zap.String("id", id), zap.Error(err))
		return nil, fmt.Errorf("get application failed: %w", err)
	}

	if app == nil {
		return nil, errors.ErrApplicationNotFound
	}

	return app, nil
}

// ApproveApplication 审批申请并给出建议
func (s *Service) ApproveApplication(appID string, approver string) (*leave.ApprovalRecommendation, error) {
	app, err := s.GetApplication(appID)
	if err != nil {
		return nil, err
	}

	// 调用 LLM 生成审批建议
	recommendation, err := s.generateApprovalRecommendation(app, approver)
	if err != nil {
		s.logger.Error("failed to generate approval recommendation", zap.Error(err))
		return nil, fmt.Errorf("generate approval recommendation failed: %w", err)
	}

	return recommendation, nil
}

// 辅助方法
func (s *Service) validateCreateRequest(req *leave.CreateRequest) error {
	if req.EmployeeID == "" || req.EmployeeName == "" || req.Supervisor == "" || req.Text == "" {
		return errors.ErrInvalidInput
	}
	return nil
}

func (s *Service) extractApplicationFromText(req *leave.CreateRequest) (*leave.Application, error) {
	prompt := buildExtractionPrompt(req)
	messages := []eino.Message{
		{
			Role:    "system",
			Content: "你是一个专业的请假单生成助手，擅长从自然语言文本中提取结构化信息。输出必须是有效的JSON格式。",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := s.eino.ChatCompletion(messages)
	if err != nil {
		return nil, fmt.Errorf("LLM chat completion failed: %w", err)
	}

	// 清理响应，提取JSON
	response = cleanJSONResponse(response)
	
	// 记录原始响应以便调试
	s.logger.Info("LLM response for application extraction", 
		zap.String("response", response))

	// 使用自定义JSON解析以支持YYYY-MM-DD日期格式
	var appData map[string]interface{}
	if err := json.Unmarshal([]byte(response), &appData); err != nil {
		return nil, fmt.Errorf("parse LLM response failed: %w", err)
	}

	// 手动解析日期字段
	var app leave.Application
	if startDateStr, ok := appData["start_date"].(string); ok && startDateStr != "" {
		parsedDate, err := parseDate(startDateStr)
		if err != nil {
			s.logger.Warn("failed to parse start_date", 
				zap.String("date", startDateStr),
				zap.Error(err))
		} else {
			app.StartDate = parsedDate
		}
	}
	
	if endDateStr, ok := appData["end_date"].(string); ok && endDateStr != "" {
		parsedDate, err := parseDate(endDateStr)
		if err != nil {
			s.logger.Warn("failed to parse end_date", 
				zap.String("date", endDateStr),
				zap.Error(err))
		} else {
			app.EndDate = parsedDate
		}
	}
	
	// 解析其他字段
	if leaveType, ok := appData["leave_type"].(string); ok {
		app.LeaveType = leave.LeaveType(leaveType)
	}
	if reason, ok := appData["reason"].(string); ok {
		app.Reason = reason
	}
	if workHandover, ok := appData["work_handover"].(string); ok {
		app.WorkHandover = workHandover
	}

	// 设置基本信息
	app.EmployeeID = req.EmployeeID
	app.EmployeeName = req.EmployeeName
	app.Supervisor = req.Supervisor

	// 计算天数（包含开始和结束日期）
	if !app.StartDate.IsZero() && !app.EndDate.IsZero() {
		// 将日期归一化到当天00:00:00，确保计算准确
		start := time.Date(app.StartDate.Year(), app.StartDate.Month(), app.StartDate.Day(), 0, 0, 0, 0, app.StartDate.Location())
		end := time.Date(app.EndDate.Year(), app.EndDate.Month(), app.EndDate.Day(), 0, 0, 0, 0, app.EndDate.Location())
		
		// 计算天数差，并加1（因为包含开始和结束日期）
		days := int(end.Sub(start).Hours() / 24) + 1
		app.Days = float64(days)
	}

	return &app, nil
}

func (s *Service) generateApprovalRecommendation(app *leave.Application, approver string) (*leave.ApprovalRecommendation, error) {
	appJSON, _ := json.MarshalIndent(app, "", "  ")

	prompt := fmt.Sprintf(`作为审批人（%s），请审核以下请假申请，并给出审批建议。

请假申请：
%s

请综合考虑以下因素：
1. 请假类型是否合理
2. 请假时间是否合适
3. 工作交接是否到位
4. 是否符合公司请假政策

输出JSON格式：
{
  "decision": "approved/rejected/conditional",
  "confidence": 0.0-1.0之间的数值,
  "reason": "审批理由和建议",
  "suggestions": ["建议1", "建议2"]
}

只输出JSON，不要其他文字。`, approver, string(appJSON))

	messages := []eino.Message{
		{
			Role:    "system",
			Content: "你是一个专业的HR审批助手，擅长根据公司政策和实际情况给出合理的审批建议。输出必须是有效的JSON格式。",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := s.eino.ChatCompletion(messages)
	if err != nil {
		return nil, fmt.Errorf("LLM chat completion failed: %w", err)
	}

	response = cleanJSONResponse(response)

	var rec leave.ApprovalRecommendation
	if err := json.Unmarshal([]byte(response), &rec); err != nil {
		return nil, fmt.Errorf("parse LLM response failed: %w", err)
	}

	return &rec, nil
}

func buildExtractionPrompt(req *leave.CreateRequest) string {
	return fmt.Sprintf(`请从以下员工的请假文本中提取信息，生成结构化的请假申请单。输出格式为JSON。

员工信息：
- 员工ID：%s
- 员工姓名：%s
- 直属主管：%s

请假文本：%s

如果文本中没有明确提到开始日期、结束日期、请假类型、原因，请根据文本内容合理推断。
输出JSON格式：
{
  "leave_type": "sick/personal/annual/adjust",
  "start_date": "YYYY-MM-DD",
  "end_date": "YYYY-MM-DD",
  "reason": "请假原因",
  "work_handover": "工作交接情况（如有）"
}

只输出JSON，不要其他文字。`, req.EmployeeID, req.EmployeeName, req.Supervisor, req.Text)
}

// parseDate 解析日期字符串，支持多种格式
// 返回的时间设置为当天的00:00:00，时区为本地时区
func parseDate(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("empty date string")
	}
	
	var parsedTime time.Time
	var err error
	
	// 尝试不同的日期格式
	formats := []string{
		"2006-01-02T15:04:05Z07:00", // ISO 8601 完整格式
		"2006-01-02T15:04:05Z",      // ISO 8601 UTC格式
		"2006-01-02 15:04:05",       // 标准格式
		"2006-01-02",                // 仅日期格式 YYYY-MM-DD
		"2006/01/02",                // 斜杠格式
	}
	
	for _, format := range formats {
		if parsedTime, err = time.Parse(format, dateStr); err == nil {
			// 如果只包含日期部分（YYYY-MM-DD），将时间设置为00:00:00，使用本地时区
			if format == "2006-01-02" || format == "2006/01/02" {
				// 获取本地时区
				loc := time.Now().Location()
				parsedTime = time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(), 0, 0, 0, 0, loc)
			}
			return parsedTime, nil
		}
	}
	
	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

func cleanJSONResponse(response string) string {
	response = strings.TrimSpace(response)

	// 移除代码块标记
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimSpace(response)
	} else if strings.HasPrefix(response, "```") {
		response = strings.TrimPrefix(response, "```")
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimSpace(response)
	}

	return response
}
