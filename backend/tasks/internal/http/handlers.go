package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/controller"
	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/entity"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	getTaskHandler           *controller.GetTaskHandler
	getTodayTasksHandler     *controller.GetTodayTasksHandler
	completeTaskHandler      *controller.CompleteTaskHandler
	getCompletedTasksHandler *controller.GetCompletedTasksHandler
}

func NewHandlers(
	getTaskHandler *controller.GetTaskHandler,
	getTodayTasksHandler *controller.GetTodayTasksHandler,
	completeTaskHandler *controller.CompleteTaskHandler,
	getCompletedTasksHandler *controller.GetCompletedTasksHandler,
) *Handlers {
	return &Handlers{
		getTaskHandler:           getTaskHandler,
		getTodayTasksHandler:     getTodayTasksHandler,
		completeTaskHandler:      completeTaskHandler,
		getCompletedTasksHandler: getCompletedTasksHandler,
	}
}

func (h *Handlers) GetTask(con *gin.Context) {
	userID, _ := con.Get("userID")
	uid, _ := userID.(string)
	if uid == "" {
		SendError(con, controller.ErrUnauthorized)
		return
	}
	taskID := con.Param("task_id")
	query := controller.GetTaskQuery{TaskID: taskID, UserID: uid}
	result, err := h.getTaskHandler.Handle(con.Request.Context(), query)
	if err != nil {
		if errors.Is(err, entity.ErrTaskNotFound) {
			SendError(con, entity.ErrTaskNotFound)
			return
		}
		SendError(con, err)
		return
	}
	con.JSON(http.StatusOK, result)
}

func (h *Handlers) GetTodayTasks(con *gin.Context) {
	userID, _ := con.Get("userID")
	uid, _ := userID.(string)
	if uid == "" {
		SendError(con, controller.ErrUnauthorized)
		return
	}

	result, err := h.getTodayTasksHandler.Handle(con.Request.Context(), controller.GetTodayTasksQuery{UserID: uid})
	if err != nil {
		SendError(con, err)
		return
	}

	con.JSON(http.StatusOK, result)
}

func (h *Handlers) CompleteTask(con *gin.Context) {
	userID, _ := con.Get("userID")
	uid, _ := userID.(string)
	if uid == "" {
		SendError(con, controller.ErrUnauthorized)
		return
	}

	taskID := con.Param("task_id")
	result, err := h.completeTaskHandler.Handle(con.Request.Context(), controller.CompleteTaskQuery{TaskID: taskID, UserID: uid})
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, entity.ErrTaskNotFound) {
			SendError(con, entity.ErrTaskNotFound)
			return
		}
		if errors.Is(err, entity.ErrTaskAlreadyCompleted) {
			SendError(con, entity.ErrTaskAlreadyCompleted)
			return
		}
		SendError(con, err)
		return
	}

	con.JSON(http.StatusOK, result)
}

func (h *Handlers) GetTodayTasksForInternal(con *gin.Context) {
	uid := con.Param("user_id")
	if uid == "" {
		SendError(con, controller.ErrInvalidRequest)
		return
	}

	result, err := h.getCompletedTasksHandler.Handle(
		con.Request.Context(),
		controller.GetCompletedTasksQuery{UserID: uid},
	)
	if err != nil {
		SendError(con, err)
		return
	}

	con.JSON(http.StatusOK, result)
}
