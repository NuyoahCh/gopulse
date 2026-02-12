package knowledgebase

import (
	"fmt"
	"strings"

	"github.com/Nuyoahch/einoverse/internal/domain/knowledgebase"
	kbRepo "github.com/Nuyoahch/einoverse/internal/repository/knowledgebase"
	"github.com/Nuyoahch/einoverse/pkg/eino"
	"github.com/Nuyoahch/einoverse/pkg/errors"
	"go.uber.org/zap"
)

// Service 知识库服务
type Service struct {
	repo   kbRepo.Repository
	eino   *eino.Client
	logger *zap.Logger
}

// NewService 创建知识库服务
func NewService(repo kbRepo.Repository, einoClient *eino.Client, logger *zap.Logger) *Service {
	return &Service{
		repo:   repo,
		eino:   einoClient,
		logger: logger,
	}
}

// CreateDocument 创建文档
func (s *Service) CreateDocument(doc *knowledgebase.Document) error {
	if doc.Title == "" || doc.Content == "" {
		return errors.ErrInvalidInput
	}

	if err := s.repo.Create(doc); err != nil {
		s.logger.Error("failed to create document", zap.Error(err))
		return fmt.Errorf("create document failed: %w", err)
	}

	s.logger.Info("document created", zap.String("id", doc.ID), zap.String("title", doc.Title))
	return nil
}

// SearchDocuments 搜索文档
func (s *Service) SearchDocuments(query string, limit int) ([]knowledgebase.SearchResult, error) {
	if query == "" {
		return []knowledgebase.SearchResult{}, nil
	}

	if limit <= 0 {
		limit = 5
	}

	results, err := s.repo.Search(query, limit)
	if err != nil {
		s.logger.Error("failed to search documents", zap.Error(err))
		return nil, fmt.Errorf("search documents failed: %w", err)
	}

	return results, nil
}

// AskQuestion 问答
func (s *Service) AskQuestion(question string) (*knowledgebase.QAResponse, error) {
	if question == "" {
		return nil, errors.ErrInvalidInput
	}

	// 提取问题的关键词用于搜索
	searchKeywords := extractSearchKeywords(question)
	s.logger.Info("extracted search keywords for question",
		zap.String("question", question),
		zap.String("keywords", searchKeywords))

	// 使用关键词检索相关文档
	results, err := s.SearchDocuments(searchKeywords, 5)
	if err != nil {
		return nil, err
	}

	// 记录搜索结果
	s.logger.Info("search results for question",
		zap.String("keywords", searchKeywords),
		zap.Int("result_count", len(results)))

	// 如果关键词搜索没有结果，尝试使用原问题进行搜索
	if len(results) == 0 {
		s.logger.Info("no results with keywords, trying original question",
			zap.String("original_question", question))
		results, err = s.SearchDocuments(question, 5)
		if err != nil {
			return nil, err
		}
		s.logger.Info("search results with original question",
			zap.Int("result_count", len(results)))
	}

	// 构建上下文
	context := buildContext(results)

	// 调用 LLM 生成回答
	prompt := buildQuestionPrompt(question, context)
	messages := []eino.Message{
		{
			Role:    "system",
			Content: "你是一个专业的企业知识库助手，擅长从文档中提取信息并准确回答员工的问题。回答要简洁、准确、专业。",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	answer, err := s.eino.ChatCompletion(messages)
	if err != nil {
		s.logger.Error("LLM chat completion failed", zap.Error(err))
		return nil, errors.ErrLLMServiceError
	}

	sourceDocs := extractSourceDocs(results)

	return &knowledgebase.QAResponse{
		Answer:     strings.TrimSpace(answer),
		SourceDocs: sourceDocs,
		Confidence: calculateConfidence(results),
	}, nil
}

// GetDocument 获取文档
func (s *Service) GetDocument(id string) (*knowledgebase.Document, error) {
	doc, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("failed to get document", zap.String("id", id), zap.Error(err))
		return nil, fmt.Errorf("get document failed: %w", err)
	}

	if doc == nil {
		return nil, errors.ErrDocumentNotFound
	}

	return doc, nil
}

// extractSearchKeywords 从问题中提取搜索关键词
// 移除常见的疑问词和标点，提取核心关键词用于搜索
func extractSearchKeywords(question string) string {
	// 移除常见的疑问词和助词，但保留实体词（如"公司"、"年假"等）
	stopWords := []string{
		"请问", "请问一下", "我想问", "我想知道", "能不能告诉我",
		"什么是", "有哪些", "怎么", "如何", "为什么",
		"？", "?", "。", ".", "，", ",",
		"的", "是", "什么", "吗", "呢", "啊",
	}

	result := question
	for _, word := range stopWords {
		result = strings.ReplaceAll(result, word, " ")
	}

	// 移除多余空格并分割成词
	fields := strings.Fields(result)
	if len(fields) == 0 {
		// 如果提取失败，返回原问题（去除标点）
		result = strings.Trim(question, "？?。.，, ")
		return result
	}

	// 保留所有有意义的词（长度大于等于2的词）
	meaningfulWords := make([]string, 0, len(fields))
	for _, word := range fields {
		if len([]rune(word)) >= 2 { // 使用rune计算中文字符长度
			meaningfulWords = append(meaningfulWords, word)
		}
	}

	if len(meaningfulWords) == 0 {
		// 如果没有有意义的词，返回原问题（去除标点）
		result = strings.Trim(question, "？?。.，, ")
		return result
	}

	// 如果词太多，只取前4个最重要的词
	if len(meaningfulWords) > 4 {
		meaningfulWords = meaningfulWords[:4]
	}

	// 重新组合，用空格分隔
	result = strings.Join(meaningfulWords, " ")

	return strings.TrimSpace(result)
}

// 辅助函数
func buildContext(results []knowledgebase.SearchResult) string {
	if len(results) == 0 {
		return "知识库中没有找到相关内容。"
	}

	var builder strings.Builder
	builder.WriteString("相关知识库内容：\n\n")
	for i, result := range results {
		builder.WriteString(fmt.Sprintf("[文档 %d] %s\n", i+1, result.Document.Title))
		builder.WriteString(result.Document.Content)
		builder.WriteString("\n\n")
	}

	return builder.String()
}

func buildQuestionPrompt(question, context string) string {
	return fmt.Sprintf(`基于以下知识库内容回答问题：

%s

问题：%s

请用简洁明了的方式回答问题。如果知识库中没有相关信息，请如实告知。`, context, question)
}

func extractSourceDocs(results []knowledgebase.SearchResult) []string {
	docs := make([]string, 0, len(results))
	for _, result := range results {
		docs = append(docs, result.Document.ID)
	}
	return docs
}

func calculateConfidence(results []knowledgebase.SearchResult) float64 {
	if len(results) == 0 {
		return 0.0
	}

	totalScore := 0.0
	for _, result := range results {
		totalScore += result.Score
	}

	avgScore := totalScore / float64(len(results))
	// 归一化到 0-1 范围
	if avgScore > 10.0 {
		return 1.0
	}
	return avgScore / 10.0
}
