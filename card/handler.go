package card

import (
	"encoding/json"
	"log"
	"net/http"
)

type Handler struct {
}

const (
	InvalidCard = 1 << iota
	InvalidDate = 1 << iota
)

func (h *Handler) ValidateCard(res http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var cardReq ValidateCardRequest

	if err := json.NewDecoder(r.Body).Decode(&cardReq); err != nil {
		return
	}

	c := cardReq.Card()

	var response ValidateCardResponse
	valid := luhnAlgorithm(c)
	if response.Valid = valid; !valid {
		response.SetError(InvalidCard, ErrInvalidCardNumber)
	}

	err := validateDate(c)
	if err != nil {
		response.Valid = false
		response.SetError(InvalidDate, err)
	}

	switch {
	case !response.Valid:
		res.WriteHeader(http.StatusBadRequest)
	default:
		res.WriteHeader(http.StatusOK)
	}

	err = json.NewEncoder(res).Encode(response)
	if err != nil {
		log.Println(err)
		return
	}
}

type ValidateCardRequest struct {
	CardNumber      string `json:"cardNumber"`
	ExpirationMonth int    `json:"expirationMonth"`
	ExpirationYear  int    `json:"expirationYear"`
}

func (r *ValidateCardRequest) Card() *Card {
	return &Card{
		Number:          r.CardNumber,
		ExpirationMonth: r.ExpirationMonth,
		ExpirationYear:  r.ExpirationYear,
	}
}

type ValidateCardResponse struct {
	Valid bool         `json:"valid"`
	Error errorMessage `json:"error"`
}

func (e *ValidateCardResponse) SetError(code int, err error) {
	e.Error = errorMessage{Message: err.Error(), Code: code}
}

type errorMessage struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
