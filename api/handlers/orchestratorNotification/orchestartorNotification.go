package orchestratornotification

import (
	"api/handlers/notifications"
	"api/models"
	"api/utilities"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrchestratorInterface interface {
	NotifyEvent(c *gin.Context)
}

type OrchestratorNotification struct {
	notifications.EmailSystem
}

func SetupNotificationRoutes(group *gin.RouterGroup, h OrchestratorInterface) {
	group.POST("", h.NotifyEvent)
	/*
		group.Handle("/reset-password", middleware.RoleMiddleware(http.AuthFunc(h.ResetPassword), 0)).Methods(http.MethodPost)
		// router.Handle("/verify", middleware.RoleMiddleware(http.AuthFunc(h.Verify), 0)).Methods(http.MethodPost)
	*/
}

func NewOrchestratorNotification(db *gorm.DB) OrchestratorInterface {
	return &OrchestratorNotification{
		notifications.NewHandler(db),
	}
}

func (e *OrchestratorNotification) NotifyEvent(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()

	//bind the request body to the struct
	var req models.NotifyEventOrchestrator
	if err := c.BindJSON(&req); err != nil {
		logger.WithError(err).Error("Invalid request from Orchestrator")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request"})
		return
	}

	if req.ActiveWorkflowID == nil || req.CurrentState == nil {
		logger.Error("Invalid request from Orchestrator")
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid request"})
		return
	}

	//get the dispute details using ID from request body
	e.EmailSystem.NotifyDisputeStateChanged(c, *req.ActiveWorkflowID, *req.CurrentState)
	logger.Info("Email notification sent successfully, via Orchestrator")
}