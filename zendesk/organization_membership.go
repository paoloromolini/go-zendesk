package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type (
	// OrganizationMembership is struct for organization membership payload
	// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_memberships/
	OrganizationMembership struct {
		ID             int64     `json:"id,omitempty"`
		URL            string    `json:"url,omitempty"`
		UserID         int64     `json:"user_id"`
		OrganizationID int64     `json:"organization_id"`
		Default        bool      `json:"default"`
		Name           string    `json:"organization_name"`
		CreatedAt      time.Time `json:"created_at,omitempty"`
		UpdatedAt      time.Time `json:"updated_at,omitempty"`
	}

	// OrganizationMembershipListOptions is a struct for options for organization membership list
	// ref: https://developer.zendesk.com/api-reference/ticketing/organizations/organization_memberships/
	OrganizationMembershipListOptions struct {
		PageOptions
		OrganizationID int64 `json:"organization_id,omitempty" url:"organization_id,omitempty"`
		UserID         int64 `json:"user_id,omitempty" url:"user_id,omitempty"`
	}

	// OrganizationMembershipOptions is a struct for options for organization membership
	// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_memberships/
	OrganizationMembershipOptions struct {
		OrganizationID int64 `json:"organization_id,omitempty"`
		UserID         int64 `json:"user_id,omitempty"`
	}

	// OrganizationMembershipsList is a struct for organization memberships payload
	OrganizationMembershipsList struct {
		OrganizationMemberships []OrganizationMembershipOptions `json:"organization_memberships"`
	}

	JobStatus struct {
		ID       string      `json:"id"`
		URL      string      `json:"url"`
		Total    interface{} `json:"total"`
		Progress interface{} `json:"progress"`
		Status   string      `json:"status"`
		Message  string      `json:"message"`
		Results  struct {
			Success bool `json:"success"`
		} `json:"results"`
	}

	// OrganizationMembershipAPI is an interface containing organization membership related methods
	OrganizationMembershipAPI interface {
		GetOrganizationMemberships(context.Context, *OrganizationMembershipListOptions) ([]OrganizationMembership, Page, error)
		CreateOrganizationMembership(context.Context, OrganizationMembershipOptions) (OrganizationMembership, error)
		CreateManyOrganizationMemberships(context.Context, OrganizationMembershipsList) (JobStatus, error)
		DeleteManyOrganizationMemberships(context.Context, []string) error
		SetDefaultOrganization(context.Context, OrganizationMembershipOptions) (OrganizationMembership, error)
		GetUserMemberships(context.Context, *OrganizationMembershipListOptions) ([]OrganizationMembership, Page, error)
	}
)

// GetUserMemberships gets the memberships of the specified users
// ref: https://developer.zendesk.com/api-reference/ticketing/organizations/organization_memberships/#list-memberships
func (z *Client) GetUserMemberships(ctx context.Context, opts *OrganizationMembershipListOptions) ([]OrganizationMembership, Page, error) {
	var result struct {
		OrganizationMemberships []OrganizationMembership `json:"organization_memberships"`
		Page
	}

	u, err := addOptions(fmt.Sprintf("/users/%d/organization_memberships.json", opts.UserID), opts.PageOptions)
	if err != nil {
		return nil, Page{}, err
	}

	body, err := z.get(ctx, u)
	if err != nil {
		return nil, Page{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, Page{}, err
	}

	return result.OrganizationMemberships, result.Page, nil
}

// GetOrganizationMemberships gets the memberships of the specified organization
// ref: https://developer.zendesk.com/api-reference/ticketing/organizations/organization_memberships/
func (z *Client) GetOrganizationMemberships(ctx context.Context, opts *OrganizationMembershipListOptions) ([]OrganizationMembership, Page, error) {
	var result struct {
		OrganizationMemberships []OrganizationMembership `json:"organization_memberships"`
		Page
	}

	tmp := opts
	if tmp == nil {
		tmp = new(OrganizationMembershipListOptions)
	}

	u, err := addOptions("/organization_memberships.json", tmp)
	if err != nil {
		return nil, Page{}, err
	}

	body, err := z.get(ctx, u)
	if err != nil {
		return nil, Page{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, Page{}, err
	}

	return result.OrganizationMemberships, result.Page, nil
}

// CreateOrganizationMembership creates an organization membership for an existing user and org
// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_memberships/#create-membership
func (z *Client) CreateOrganizationMembership(ctx context.Context, opts OrganizationMembershipOptions) (OrganizationMembership, error) {
	var data, result struct {
		OrganizationMembership OrganizationMembership `json:"organization_membership"`
	}

	data.OrganizationMembership = OrganizationMembership{
		UserID:         opts.UserID,
		OrganizationID: opts.OrganizationID,
	}

	body, err := z.post(ctx, "/organization_memberships.json", data)

	if err != nil {
		return OrganizationMembership{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return OrganizationMembership{}, err
	}

	return result.OrganizationMembership, err
}

// SetDefaultOrganization sets the default organization for a user that has a membership in that org
// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_memberships/#set-organization-as-default
func (z *Client) SetDefaultOrganization(ctx context.Context, opts OrganizationMembershipOptions) (OrganizationMembership, error) {
	var result struct {
		OrganizationMembership OrganizationMembership `json:"organization_membership"`
	}

	body, err := z.put(ctx, fmt.Sprintf("/users/%d/organizations/%d/make_default.json", opts.UserID, opts.OrganizationID), nil)
	if err != nil {
		return OrganizationMembership{}, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return OrganizationMembership{}, err
	}

	return result.OrganizationMembership, nil
}

// CreateManyOrganizationMemberships creates many organization membership for existing users and organizations
// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_memberships/#create-many-memberships
func (z *Client) CreateManyOrganizationMemberships(ctx context.Context, data OrganizationMembershipsList) (JobStatus, error) {
	var result struct {
		JobStatus JobStatus `json:"job_status"`
	}

	body, err := z.post(ctx, "/organization_memberships/create_many", data)

	if err != nil {
		return JobStatus{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return JobStatus{}, err
	}

	return result.JobStatus, err
}

// DeleteManyOrganizationMemberships deletes many organization membership passing organization memberships IDs
// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_memberships/#bulk-delete-memberships
func (z *Client) DeleteManyOrganizationMemberships(ctx context.Context, IDs []string) error {
	ids := strings.Join(IDs, ",")
	err := z.delete(ctx, fmt.Sprintf("/organization_memberships/destroy_many?ids=%v", ids))

	if err != nil {
		return err
	}

	return nil
}
