package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type CustomStatus struct {
	Active                bool      `json:"active"`
	AgentLabel            string    `json:"agent_label"`
	CreatedAt             time.Time `json:"created_at"`
	Default               bool      `json:"default"`
	Description           string    `json:"description"`
	EndUserDescription    string    `json:"end_user_description"`
	EndUserLabel          string    `json:"end_user_label"`
	ID                    int64     `json:"id"`
	RawAgentLabel         string    `json:"raw_agent_label"`
	RawDescription        string    `json:"raw_description"`
	RawEndUserDescription string    `json:"raw_end_user_description"`
	RawEndUserLabel       string    `json:"raw_end_user_label"`
	StatusCategory        string    `json:"status_category"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// CustomStatusAPI an interface containing all custom status related methods
type CustomStatusAPI interface {
	ListCustomTicketStatuses(ctx context.Context, opts *CustomTicketStatusOptions) ([]CustomStatus, error)
	ShowCustomTicketStatus(ctx context.Context, customStatusID int64) (*CustomStatus, error)
}

// CustomTicketStatusOptions list custom status options
type CustomTicketStatusOptions struct {
	Active           bool   `url:"active,omitempty"`
	Default          bool   `url:"default,omitempty"`
	StatusCategories string `url:"status_categories,omitempty"`
}

// ListCustomTicketStatuses mocks base method.
// https://developer.zendesk.com/api-reference/ticketing/tickets/custom_ticket_statuses/#list-custom-ticket-statuses
func (z *Client) ListCustomTicketStatuses(
	ctx context.Context, opts *CustomTicketStatusOptions,
) ([]CustomStatus, error) {
	var data struct {
		CustomStatuses []CustomStatus `json:"custom_statuses"`
	}
	tmp := opts
	if tmp == nil {
		tmp = &CustomTicketStatusOptions{}
	}
	u, err := addOptions("/custom_statuses.json", tmp)
	body, err := z.get(ctx, u)
	if err != nil {
		return []CustomStatus{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return []CustomStatus{}, err
	}
	return data.CustomStatuses, nil
}

// ShowCustomTicketStatus returns the custom ticket status object
func (z *Client) ShowCustomTicketStatus(ctx context.Context, customStatusID int64) (*CustomStatus, error) {
	var result struct {
		CustomStatus CustomStatus `json:"custom_status"`
	}
	url := fmt.Sprintf("/custom_statuses/%d", customStatusID)
	body, err := z.get(ctx, url)

	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &result)

	if err != nil {
		return nil, err
	}
	return &result.CustomStatus, nil
}
