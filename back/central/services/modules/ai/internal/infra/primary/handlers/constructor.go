package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IHandlers interface {
	GetRecommendation(c *gin.Context)
	AssistantChat(c *gin.Context)
	GetAssistantState(c *gin.Context)
	MarkAssistantIntroSeen(c *gin.Context)
	SubmitAssistantFeedback(c *gin.Context)
	MarkAssistantClick(c *gin.Context)
	ListAssistantAlerts(c *gin.Context)
	GetAssistantAlertsUnread(c *gin.Context)
	MarkAssistantAlertsSeen(c *gin.Context)
	ListAssistantReviewMessages(c *gin.Context)
	GetAssistantReviewSummary(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

type Handlers struct {
	uc  app.IUseCase
	log log.ILogger
}

func New(uc app.IUseCase, logger log.ILogger) IHandlers {
	return &Handlers{uc: uc, log: logger}
}

func accessScope(c *gin.Context) (dtos.AccessScope, bool) {
	userID, ok := middleware.GetUserID(c)
	if !ok || userID == 0 {
		return dtos.AccessScope{}, false
	}
	tokenBusinessID, _ := middleware.GetBusinessIDFromContext(c)

	scope := dtos.AccessScope{UserID: userID, TokenBusinessID: tokenBusinessID}
	if tokenBusinessID == 0 {
		if raw := c.Query("business_id"); raw != "" {
			if parsed, err := strconv.ParseUint(raw, 10, 64); err == nil {
				scope.RequestedBusinessID = uint(parsed)
			}
		}
	}
	return scope, true
}
