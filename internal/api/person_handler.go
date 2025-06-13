package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/zthiagovalle/rinha-de-backend-2023-q3/internal/store"
	"github.com/zthiagovalle/rinha-de-backend-2023-q3/internal/utils"
)

type createPersonRequest struct {
	Username  *string   `json:"apelido"`
	Name      *string   `json:"nome"`
	BirthDate *string   `json:"nascimento"`
	Stack     *[]string `json:"stack"`
}

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

func (ph *PersonHandler) HandleGetPerson(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "id is required"})
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "invalid uuid"})
		return
	}

	person, err := ph.personStore.GetPersonByID(id)
	if err != nil {
		ph.logger.Printf("ERROR: getPersonByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if person == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"message": "person not found"})
		return
	}

	utils.WriteEntity(w, http.StatusOK, person)
}

func (ph *PersonHandler) HandleGetPersonsByTerm(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("t")
	if term == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "term is required"})
		return
	}

	limit := 50
	persons, err := ph.personStore.GetPersonsByTerm(term, limit)
	if err != nil {
		ph.logger.Printf("ERROR: getPersonsByTerm: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if len(persons) == 0 {
		utils.WriteEntity(w, http.StatusOK, []any{})
		return
	}

	utils.WriteEntity(w, http.StatusOK, persons)
}

func (ph *PersonHandler) HandleCreatePerson(w http.ResponseWriter, r *http.Request) {
	var req createPersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request body"})
		return
	}

	if err := ph.validateCreatePerson(req); err != nil {
		utils.WriteJSON(w, http.StatusUnprocessableEntity, utils.Envelope{"error": err.Error()})
		return
	}

	birthDate, err := time.Parse("2006-01-02", *req.BirthDate)
	if err != nil {
		ph.logger.Printf("ERROR: parse birthDate: %v", err)
		utils.WriteJSON(w, http.StatusUnprocessableEntity, utils.Envelope{"error": "nascimento must be in YYYY-MM-DD format"})
		return
	}
	person := &store.Person{
		Username:  *req.Username,
		Name:      *req.Name,
		BirthDate: store.DateOnly{Time: birthDate},
		Stack:     req.Stack,
	}

	id, err := ph.personStore.CreatePerson(person)
	if err != nil {
		if err.Error() == store.ErrPersonUsernameAlreadyExists {
			utils.WriteJSON(w, http.StatusUnprocessableEntity, utils.Envelope{"error": "apelido already exists"})
			return
		}

		ph.logger.Printf("ERROR: createPerson: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/pessoas/%s", id.String()))
	w.WriteHeader(http.StatusCreated)
}

func (ph *PersonHandler) validateCreatePerson(req createPersonRequest) error {
	if req.Username == nil || *req.Username == "" {
		return fmt.Errorf("apelido is required")
	}
	if len(*req.Username) > 32 {
		return fmt.Errorf("apelido must be at most 32 characters long")
	}

	if req.Name == nil || *req.Name == "" {
		return fmt.Errorf("nome is required")
	}
	if len(*req.Name) > 100 {
		return fmt.Errorf("nome must be at most 100 characters long")
	}

	if req.BirthDate == nil || *req.BirthDate == "" {
		return fmt.Errorf("nascimento is required")
	}
	birthDateReger := regexp.MustCompile(`^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])$`)
	if !birthDateReger.MatchString(*req.BirthDate) {
		return fmt.Errorf("nascimento must be in YYYY-MM-DD format")
	}

	if req.Stack != nil && len(*req.Stack) > 0 {
		for _, stack := range *req.Stack {
			if len(stack) > 32 {
				return fmt.Errorf("each stack item must be at most 32 characters long")
			}
		}
	}

	return nil
}
