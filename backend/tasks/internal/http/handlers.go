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

func NewHandlers(
	getTaskHandler *controller.GetTaskHandler,
	getTodayTasksHandler *controller.GetTodayTasksHandler,
	completeTaskHandler *controller.CompleteTaskHandler,
) *Handlers {
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
		sendError(con, http.StatusUnauthorized, "UNAUTHORIZED", "Требуется повторная авторизация")
		return
	}
	taskId := con.Param("task_id")
	query := controller.GetTaskQuery{TaskId: taskId, UserId: uid}
	result, err := h.getTaskHandler.Handle(con.Request.Context(), query)
	if err != nil {
		if errors.Is(err, entity.ErrTaskNotFound) {
			sendError(con, http.StatusNotFound, "TASK_NOT_FOUND", "Задание не найдено")
			return
		}
		sendError(con, http.StatusInternalServerError, "INTERNAL_ERROR", "Внутренняя ошибка сервера")
		return
	}
	con.JSON(http.StatusOK, result)

}

func (h *Handlers) GetTodayTasks(con *gin.Context) {
	userID, _ := con.Get("userID")
	uid, _ := userID.(string)
	if uid == "" {
		sendError(con, http.StatusUnauthorized, "UNAUTHORIZED", "Требуется повторная авторизация")
		return
	}

	result, err := h.getTodayTasksHandler.Handle(con.Request.Context(), controller.GetTodayTasksQuery{UserID: uid})
	if err != nil {
		sendError(con, http.StatusInternalServerError, "INTERNAL_ERROR", "Внутренняя ошибка сервера")
		return
	}

	con.JSON(http.StatusOK, result)
}

func (h *Handlers) CompleteTask(con *gin.Context) {
	userID, _ := con.Get("userID")
	uid, _ := userID.(string)
	if uid == "" {
		sendError(con, http.StatusUnauthorized, "UNAUTHORIZED", "Требуется повторная авторизация")
		return
	}

	taskID := con.Param("task_id")
	result, err := h.completeTaskHandler.Handle(con.Request.Context(), controller.CompleteTaskQuery{TaskID: taskID, UserID: uid})
	if err != nil {
		if errors.Is(err, entity.ErrTaskNotFound) {
			sendError(con, http.StatusNotFound, "TASK_NOT_FOUND", "Задание не найдено")
			return
		}
		if errors.Is(err, entity.ErrTaskAlreadyCompleted) {
			sendError(con, http.StatusConflict, "TASK_ALREADY_COMPLETED", "Награда за это задание уже получена")
			return
		}
		sendError(con, http.StatusInternalServerError, "INTERNAL_ERROR", "Внутренняя ошибка сервера")
		return
	}

	con.JSON(http.StatusOK, result)
}
