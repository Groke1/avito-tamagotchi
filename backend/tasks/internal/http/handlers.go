package http

import (
	"errors"
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/controller"
	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/entity"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	getTaskHandler       *controller.GetTaskHandler
	getTodayTasksHandler *controller.GetTodayTasksHandler
	completeTaskHandler  *controller.CompleteTaskHandler
}

func NewHandlers(getTaskHandler *controller.GetTaskHandler, getTodayTasksHandler *controller.GetTodayTasksHandler, completeTaskHandler *controller.CompleteTaskHandler) *Handlers {
	return &Handlers{
		getTaskHandler:       getTaskHandler,
		getTodayTasksHandler: getTodayTasksHandler,
		completeTaskHandler:  completeTaskHandler,
	}
}

func (h *Handlers) GetTask(con *gin.Context) {
	userID, _ := con.Get("userID")
	uid, _ := userID.(string)
	if uid == "" {
		con.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	taskId := con.Param("task_id")
	query := controller.GetTaskQuery{TaskId: taskId, UserId: uid}
	result, err := h.getTaskHandler.Handle(con.Request.Context(), query)
	if err != nil {
		con.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	con.JSON(http.StatusOK, result)

}

func (h *Handlers) GetTodayTasks(con *gin.Context) {
	userID, _ := con.Get("userID")
	uid, _ := userID.(string)
	if uid == "" {
		con.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	result, err := h.getTodayTasksHandler.Handle(con.Request.Context(), controller.GetTodayTasksQuery{UserID: uid})
	if err != nil {
		con.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	con.JSON(http.StatusOK, result)
}

func (h *Handlers) CompleteTask(con *gin.Context) {
	userID, _ := con.Get("userID")
	uid, _ := userID.(string)
	if uid == "" {
		con.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	taskID := con.Param("task_id")
	result, err := h.completeTaskHandler.Handle(con.Request.Context(), controller.CompleteTaskQuery{TaskID: taskID, UserID: uid})
	if err != nil {
		if errors.Is(err, entity.ErrTaskNotFound) {
			con.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		if errors.Is(err, entity.ErrTaskAlreadyCompleted) {
			con.JSON(http.StatusConflict, gin.H{"error": "task already completed"})
			return
		}
		con.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	con.JSON(http.StatusOK, result)
}
