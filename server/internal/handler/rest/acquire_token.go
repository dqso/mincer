package rest

import (
	"encoding/base64"
	"encoding/json"
	"github.com/dqso/mincer/server/internal/log"
	"net/http"
)

func (h Handler) AcquireToken(w http.ResponseWriter, r *http.Request) {
	var request acquireTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := request.Validate(); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	clientID, connectToken, err := h.usecase.AcquireToken(r.Context())
	if err != nil {
		h.logger.Error("unable to acquire a token", log.Err(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var response acquireTokenResponse
	response.DTO(clientID, connectToken)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

type acquireTokenRequest struct{}

func (acquireTokenRequest) Validate() error {
	return nil
}

type acquireTokenResponse struct {
	ClientID     uint64 `json:"client_id"`
	ConnectToken string `json:"connect_token"`
}

func (r *acquireTokenResponse) DTO(clientID uint64, connectToken []byte) {
	r.ClientID = clientID
	r.ConnectToken = base64.StdEncoding.EncodeToString(connectToken)
}
