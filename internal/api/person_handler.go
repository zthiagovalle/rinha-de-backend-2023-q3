package api

import (
	"log"
	"net/http"

	"github.com/zthiagovalle/rinha-de-backend-2023-q3/internal/utils"
)

type PersonHandler struct {
	logger *log.Logger
}

func NewPersonHandler(logger *log.Logger) *PersonHandler {
	return &PersonHandler{
		logger: logger,
	}
}

func (ph *PersonHandler) HandleCountPersons(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"contagem": 10})
}
