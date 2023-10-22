package network

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/wirepair/netcode"
	"io"
	"log"
	"net/http"
)

type NetcodeToken struct {
	ClientID     uint64 `json:"client_id"`
	ConnectToken string `json:"connect_token"`
}

func (m *Manager) getConnectToken(ctx context.Context) (uint64, *netcode.ConnectToken, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.tokenUrl, bytes.NewReader([]byte("{}")))
	if err != nil {
		return 0, nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	var response NetcodeToken
	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	if err := json.Unmarshal(bts, &response); err != nil {
		return 0, nil, fmt.Errorf("unable to decode json: %v.\nReceived message: %s", err, string(bts))
	}
	log.Printf("%d: %s", response.ClientID, response.ConnectToken)

	tokenBts, err := base64.StdEncoding.DecodeString(response.ConnectToken)
	if err != nil {
		return 0, nil, err
	}
	connToken, err := netcode.ReadConnectToken(tokenBts)
	if err != nil {
		return 0, nil, err
	}
	return response.ClientID, connToken, nil
}
