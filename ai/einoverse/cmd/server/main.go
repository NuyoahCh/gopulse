package main

import (
	"fmt"

	"github.com/Nuyoahch/einoverse/internal/config"
	"github.com/Nuyoahch/einoverse/internal/handler"
	kbHandler "github.com/Nuyoahch/einoverse/internal/handler/knowledgebase"
	leaveHandler "github.com/Nuyoahch/einoverse/internal/handler/leave"
	kbRepo "github.com/Nuyoahch/einoverse/internal/repository/knowledgebase"
	leaveRepo "github.com/Nuyoahch/einoverse/internal/repository/leave"
	kbService "github.com/Nuyoahch/einoverse/internal/service/knowledgebase"
	leaveService "github.com/Nuyoahch/einoverse/internal/service/leave"
	"github.com/Nuyoahch/einoverse/pkg/eino"
	"github.com/Nuyoahch/einoverse/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化日志
	if err := logger.Init(cfg.Log.Level); err != nil {
		panic(fmt.Sprintf("failed to init logger: %v", err))
	}
	defer logger.Logger.Sync()

	logger.Logger.Info("starting application", zap.Any("config", cfg))

	// 初始化 Eino 客户端
	einoClient := eino.NewClient()
	if !einoClient.IsAvailable() {
		logger.Logger.Warn("Eino API key not configured, some features may not work")
	}

	// 初始化仓储
	kbRepo := kbRepo.NewInMemoryRepository()
	leaveRepo := leaveRepo.NewInMemoryRepository()

	// 初始化服务
	kbService := kbService.NewService(kbRepo, einoClient, logger.Logger)
	leaveService := leaveService.NewService(leaveRepo, einoClient, logger.Logger)

	// 初始化处理器
	kbHandler := kbHandler.NewHandler(kbService, logger.Logger)
	leaveHandler := leaveHandler.NewHandler(leaveService, logger.Logger)

	// 设置 Gin
	if cfg.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建 Gin 路由
	router := gin.New()
	// 使用日志中间件记录请求信息
	router.Use(handler.LoggerMiddleware(logger.Logger))
	// 使用恢复中间件捕获panic并记录错误
	router.Use(handler.RecoveryMiddleware(logger.Logger))

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		// 返回健康检查响应
		c.JSON(200, gin.H{
			"status":        "ok",
			"model_backend": getModelBackend(einoClient),
		})
	})

	// 知识库路由
	kb := router.Group("/api/v1/knowledgebase")
	{
		kb.POST("/documents", kbHandler.CreateDocument)
		kb.GET("/documents/:id", kbHandler.GetDocument)
		kb.GET("/documents/search", kbHandler.SearchDocuments)
		kb.POST("/ask", kbHandler.AskQuestion)
	}

	// 请假路由
	leave := router.Group("/api/v1/leave")
	{
		leave.POST("/applications", leaveHandler.CreateApplication)
		leave.GET("/applications/:id", leaveHandler.GetApplication)
		leave.POST("/applications/:id/approve", leaveHandler.ApproveApplication)
	}

	// 启动服务器
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	logger.Logger.Info("server starting", zap.String("address", addr))

	if err := router.Run(addr); err != nil {
		logger.Logger.Fatal("server start failed", zap.Error(err))
	}
}

// getModelBackend 获取模型后端
func getModelBackend(client *eino.Client) string {
	if client.IsAvailable() {
		return "eino"
	}
	return "mock"
}
