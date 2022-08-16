package flashcards

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/hascorp/hasuniversity/dto"
)

type FlashcardHandler struct {
	Session *gocql.Session
}

func (h *FlashcardHandler) GetAllFlashcards(w http.ResponseWriter, r *http.Request) {
	var err error
	query := r.URL.Query()
	limit := -1
	if query.Has("limit") {
		limit, err = strconv.Atoi(query.Get("limit"))
		if err != nil {
			log.Println("failed to parse limit as an int", query.Get("limit"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	cards, err := dto.GetAllCardSets(r.Context(), h.Session, limit)
	if err != nil {
		log.Println("uh oh in API handler", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payload, err := json.Marshal(cards)
	if err != nil {
		log.Println("failed to marshal", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(payload)
	if err != nil {
		log.Println("failed to write payload", err)
	}
}

func (h *FlashcardHandler) GetFlashcardSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	fc, err := dto.GetCardSet(r.Context(), h.Session, vars["uuid"])
	if err != nil {
		log.Println("uh oh in API handler", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payload, err := json.Marshal(fc)
	if err != nil {
		log.Println("failed to marshal", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(payload)
	if err != nil {
		log.Println("failed to write payload", err)
	}
}

func (h *FlashcardHandler) AddFlashcardSet(w http.ResponseWriter, r *http.Request) {
	var c dto.CardSet
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("failed to read request body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(b, &c)
	if err != nil {
		log.Println("failed to unmarshal json", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	returned, err := dto.AddCardSet(r.Context(), h.Session, &c)
	if err != nil {
		log.Println("failed to add flashcard", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(returned)
	if err != nil {
		log.Println("failed to marshal flashcard", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Println("failed to write response body", err)
	}
}
