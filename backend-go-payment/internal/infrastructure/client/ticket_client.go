package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/harri88/eticketing-challenge/internal/domain"
)

type ticketClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewTicketClient(baseURL string) domain.TicketClient {
	return &ticketClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *ticketClient) ConfirmPayment(ctx context.Context, orderID string) error {
	payload := map[string]string{"order_id": orderID}
	jsonBody, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/v1/checkout/confirm-payment", c.baseURL), bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ticket service responded with status: %d", resp.StatusCode)
	}

	return nil
}
