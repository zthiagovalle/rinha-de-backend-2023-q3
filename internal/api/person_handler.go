package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zthiagovalle/rinha-de-backend-2023-q3/internal/store"
	"github.com/zthiagovalle/rinha-de-backend-2023-q3/internal/utils"
)

type PersonHandler struct {
	logger      *log.Logger
	personStore store.PersonStore
}

func NewPersonHandler(logger *log.Logger, personStore store.PersonStore) *PersonHandler {
	return &PersonHandler{
		logger:      logger,
		personStore: personStore,
	}
}

func (ph *PersonHandler) HandleCountPersons(w http.ResponseWriter, r *http.Request) {
	totalPersons, err := ph.personStore.CountPersons()
	if err != nil {
		ph.logger.Printf("ERROR: countPersons: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "%d", totalPersons)
}
