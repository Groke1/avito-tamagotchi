package reward

import (
	"errors"
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/httpx"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (c *controller) GetDefinition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	definition, err := c.service.GetDefinition(r.Context(), code)

	if err != nil {
		if errors.Is(err, entity.ErrRewardDefinitionNotFound) {
			httpx.WriteError(w, http.StatusNotFound, httpx.ErrNotFoundDefinition)
			return
		}

		c.logger.Error("failed to get definition", zap.String("code", code), zap.Error(err))
		httpx.WriteError(w, http.StatusInternalServerError, httpx.ErrInternal)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, getDefinitionResponse{
		Code:        definition.Code,
		Name:        definition.Name,
		Description: definition.Description,
	})
}
