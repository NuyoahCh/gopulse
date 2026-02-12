package leave

// ApprovalDecision 审批决策
type ApprovalDecision string

// 审批决策
const (
	DecisionApproved    ApprovalDecision = "approved"    // 已批准
	DecisionRejected    ApprovalDecision = "rejected"    // 已拒绝
	DecisionConditional ApprovalDecision = "conditional" // 待审批
)

// ApprovalRecommendation 审批建议
type ApprovalRecommendation struct {
	Decision    ApprovalDecision `json:"decision"`              // 审批决策
	Confidence  float64          `json:"confidence"`            // 置信度
	Reason      string           `json:"reason"`                // 原因
	Suggestions []string         `json:"suggestions,omitempty"` // 建议
}

// ApprovalRequest 审批请求
type ApprovalRequest struct {
	ApplicationID string `json:"application_id,omitempty"` // 申请ID（从URL路径获取，此处可选）
	Approver      string `json:"approver" binding:"required"`       // 审批人
	Comments      string `json:"comments,omitempty"`                // 备注
}
