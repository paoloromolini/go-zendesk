package zendesk

import (
	"context"
	"encoding/json"
	"time"
)

// Requests is struct for requests payload
type Requests struct {
	URL              string          `json:"url,omitempty"`
	ID               int             `json:"id,omitempty"`
	Status           string          `json:"status,omitempty"`
	Priority         string          `json:"priority,omitempty"`
	Type             string          `json:"type,omitempty"`
	Subject          string          `json:"subject,omitempty"`
	Description      string          `json:"description,omitempty"`
	OrganizationID   int64           `json:"organization_id,omitempty"`
	Via              RequestsVia     `json:"via,omitempty"`
	CustomFields     []CustomField   `json:"custom_fields,omitempty"`
	RequesterID      int64           `json:"requester_id,omitempty"`
	CollaboratorIds  []int64         `json:"collaborator_ids,omitempty"`
	EmailCcIds       []int64         `json:"email_cc_ids,omitempty"`
	IsPublic         bool            `json:"is_public,omitempty"`
	DueAt            *time.Time      `json:"due_at,omitempty"`
	CanBeSolvedByMe  bool            `json:"can_be_solved_by_me,omitempty"`
	CreatedAt        *time.Time      `json:"created_at,omitempty"`
	UpdatedAt        *time.Time      `json:"updated_at,omitempty"`
	Recipient        string          `json:"recipient,omitempty"`
	FollowupSourceID int64           `json:"followup_source_id,omitempty"`
	AssigneeID       int64           `json:"assignee_id,omitempty"`
	TicketFormID     int64           `json:"ticket_form_id,omitempty"`
	CustomStatusID   int64           `json:"custom_status_id,omitempty"`
	Fields           []RequestsField `json:"fields,omitempty"`
}

type RequestsVia struct {
	Channel string `json:"channel"`
	Source  struct {
		From interface{} `json:"from"`
		To   interface{} `json:"to"`
		Rel  string      `json:"rel"`
	} `json:"source"`
}

type RequestsField struct {
	ID    int64 `json:"id"`
	Value any   `json:"value"`
}

// RequestsOptions are the options that can be provided to the requests search API
//
// ref: https://developer.zendesk.com/api-reference/ticketing/tickets/ticket-requests/#parameters-1
type RequestsOptions struct {
	SearchOptions  `json:",inline"`
	OrganizationID int64 `url:"organization_id"`
}

type SearchRequestsAPI interface {
	SearchRequests(ctx context.Context, opts *RequestsOptions) ([]Requests, Page, error)
}

// Search requests allows end users to query zendesk requests search api
//
// ref: https://developer.zendesk.com/api-reference/ticketing/tickets/ticket-requests/
func (z *Client) SearchRequests(ctx context.Context, opts *RequestsOptions) ([]Requests, Page, error) {
	var data struct {
		Requests []Requests `json:"requests"`
		Page
	}

	if opts == nil {
		return []Requests{}, Page{}, &OptionsError{opts}
	}

	u, err := addOptions("/requests/search.json", opts)
	if err != nil {
		return []Requests{}, Page{}, &OptionsError{opts}
	}

	body, err := z.get(ctx, u)
	if err != nil {
		return []Requests{}, Page{}, &OptionsError{opts}
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return []Requests{}, Page{}, &OptionsError{opts}
	}

	return data.Requests, data.Page, nil
}
