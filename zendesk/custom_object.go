package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type CustomObjectRecord struct {
	Url                string                 `json:"url,omitempty"`
	Name               string                 `json:"name,omitempty"`
	ID                 string                 `json:"id,omitempty"`
	CustomObjectKey    string                 `json:"custom_object_key"`
	CustomObjectFields map[string]interface{} `json:"custom_object_fields" binding:"required"`
	CreatedByUserID    string                 `json:"created_by_user_id,omitempty"`
	UpdatedByUserID    string                 `json:"updated_by_user_id,omitempty"`
	CreatedAt          time.Time              `json:"created_at,omitempty"`
	UpdatedAt          time.Time              `json:"updated_at,omitempty"`
	ExternalID         string                 `json:"external_id,omitempty"`
}

type CustomObjectFieldOptions struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	RawName string `json:"raw_name"`
	Value   string `json:"value"`
}

type CustomObjectField struct {
	Active              bool                       `json:"active"`
	CreatedAt           time.Time                  `json:"created_at"`
	Description         string                     `json:"description"`
	ID                  int64                      `json:"id"`
	Key                 string                     `json:"key"`
	Position            int                        `json:"position"`
	RawDescription      string                     `json:"raw_description"`
	RawTitle            string                     `json:"raw_title"`
	RegexpForValidation interface{}                `json:"regexp_for_validation"`
	System              bool                       `json:"system"`
	Title               string                     `json:"title"`
	Type                string                     `json:"type"`
	UpdatedAt           time.Time                  `json:"updated_at"`
	URL                 string                     `json:"url"`
	CustomFieldOptions  []CustomObjectFieldOptions `json:"custom_field_options,omitempty"`
}

// CustomObjectAPI an interface containing all custom object related methods
type CustomObjectAPI interface {
	CreateCustomObjectRecord(
		ctx context.Context, record CustomObjectRecord, customObjectKey string) (CustomObjectRecord, error)
	AutocompleteSearchCustomObjectRecords(
		ctx context.Context,
		customObjectKey string,
		opts *AutocompleteSearchCustomObjectRecordsOptions,
	) ([]CustomObjectRecord, CursorPaginationMeta, error)
	SearchCustomObjectRecords(
		ctx context.Context, customObjectKey string, opts *SearchCustomObjectRecordsOptions,
	) ([]CustomObjectRecord, CursorPaginationMeta, int64, error)
	ListCustomObjectRecords(
		ctx context.Context, customObjectKey string, opts *CustomObjectListOptions) ([]CustomObjectRecord, CursorPaginationMeta, error)
	ShowCustomObjectRecord(
		ctx context.Context, customObjectKey string, customObjectRecordID string,
	) (*CustomObjectRecord, error)
	UpdateCustomObjectRecord(
		ctx context.Context, customObjectKey string, customObjectRecordID string, record CustomObjectRecord,
	) (*CustomObjectRecord, error)
	GetSourcesByTarget(
		ctx context.Context,
		fieldID string,
		sourceType string,
		targetID string,
		targetType string,
		opts *PageOptions,
	) (GetSourceByTargetResult, error)
	DeleteCustomObjectRecord(
		ctx context.Context,
		record CustomObjectRecord,
	) error
	ListCustomObjectFields(
		ctx context.Context,
		customObjectKey string,
	) ([]CustomObjectField, error)
	ShowCustomObjectField(
		ctx context.Context, customObjectKey string, customObjectFieldKeyOrId string,
	) (CustomObjectField, error)
	FilterCustomObjectRecords(
		ctx context.Context, customObjectKey string, filterBody interface{}, opts *SearchCustomObjectRecordsOptions,
	) ([]CustomObjectRecord, CursorPaginationMeta, int64, error)
}

// CreateCustomObjectRecord CreateCustomObject create a custom object record
func (z *Client) CreateCustomObjectRecord(
	ctx context.Context, record CustomObjectRecord, customObjectKey string,
) (CustomObjectRecord, error) {

	var data, result struct {
		CustomObjectRecord CustomObjectRecord `json:"custom_object_record"`
	}
	data.CustomObjectRecord = record

	body, err := z.post(ctx, fmt.Sprintf("/custom_objects/%s/records.json", customObjectKey), data)
	if err != nil {
		return CustomObjectRecord{}, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return CustomObjectRecord{}, err
	}
	return result.CustomObjectRecord, nil
}

// CustomObjectListOptions custom object list options
type CustomObjectListOptions struct {
	CursorPagination
	Ids         string `url:"filter[ids],omitempty"`
	ExternalIds string `url:"filter[external_ids],omitempty"`
}

// AutocompleteSearchCustomObjectRecordsOptions custom object search
type AutocompleteSearchCustomObjectRecordsOptions struct {
	Name string `url:"name,omitempty"`
	CursorPagination
}

// ListCustomObjectRecords list objects
// https://developer.zendesk.com/api-reference/custom-data/custom-objects/custom_object_records/#list-custom-object-records
func (z *Client) ListCustomObjectRecords(
	ctx context.Context, customObjectKey string, opts *CustomObjectListOptions) ([]CustomObjectRecord, CursorPaginationMeta, error) {
	var result struct {
		CustomObjectRecords []CustomObjectRecord `json:"custom_object_records"`
		Meta                CursorPaginationMeta `json:"meta"`
	}
	tmp := opts
	if tmp == nil {
		tmp = &CustomObjectListOptions{}
	}
	url := fmt.Sprintf("/custom_objects/%s/records", customObjectKey)
	urlWithOptions, err := addOptions(url, tmp)
	body, err := z.get(ctx, urlWithOptions)

	if err != nil {
		return nil, result.Meta, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, result.Meta, err
	}
	return result.CustomObjectRecords, result.Meta, nil
}

// AutocompleteSearchCustomObjectRecords search for a custom object record by the name field
// https://developer.zendesk.com/api-reference/custom-objects/custom_object_records/#autocomplete-custom-object-record-search
func (z *Client) AutocompleteSearchCustomObjectRecords(
	ctx context.Context, customObjectKey string, opts *AutocompleteSearchCustomObjectRecordsOptions,
) ([]CustomObjectRecord, CursorPaginationMeta, error) {
	var result struct {
		CustomObjectRecords []CustomObjectRecord `json:"custom_object_records"`
		Meta                CursorPaginationMeta `json:"meta"`
	}
	tmp := opts
	if tmp == nil {
		tmp = &AutocompleteSearchCustomObjectRecordsOptions{}
	}
	url := fmt.Sprintf("/custom_objects/%s/records/autocomplete", customObjectKey)
	urlWithOptions, err := addOptions(url, tmp)
	body, err := z.get(ctx, urlWithOptions)

	if err != nil {
		return nil, CursorPaginationMeta{}, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, CursorPaginationMeta{}, err
	}
	return result.CustomObjectRecords, result.Meta, nil
}

type SearchCustomObjectRecordsOptions struct {
	CursorPagination

	// One of name, created_at, updated_at, -name, -created_at, or -updated_at.
	// The - denotes the sort will be descending. Defaults to sorting by relevance.
	Sort string `url:"sort,omitempty"`

	// Query string
	Query string `url:"query,omitempty"`

	// ExternalID string
	ExternalID string `url:"external_id,omitempty"`
}

// SearchCustomObjectRecords search for a custom object record by the name field
// https://developer.zendesk.com/api-reference/custom-objects/custom_object_records/#search-custom-object-records
func (z *Client) SearchCustomObjectRecords(
	ctx context.Context, customObjectKey string, opts *SearchCustomObjectRecordsOptions,
) ([]CustomObjectRecord, CursorPaginationMeta, int64, error) {
	var result struct {
		Count               int64                `json:"count"`
		CustomObjectRecords []CustomObjectRecord `json:"custom_object_records"`
		Meta                CursorPaginationMeta `json:"meta"`
	}
	tmp := opts
	if tmp == nil {
		tmp = &SearchCustomObjectRecordsOptions{}
	}
	url := fmt.Sprintf("/custom_objects/%s/records/search", customObjectKey)
	urlWithOptions, err := addOptions(url, tmp)
	body, err := z.get(ctx, urlWithOptions)

	if err != nil {
		return nil, result.Meta, 0, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, result.Meta, 0, err
	}
	return result.CustomObjectRecords, result.Meta, result.Count, nil
}

// ShowCustomObjectRecord returns a custom record for a specific object using a provided id.
// https://developer.zendesk.com/api-reference/custom-objects/custom_object_records/#show-custom-object-record
func (z *Client) ShowCustomObjectRecord(
	ctx context.Context, customObjectKey string, customObjectRecordID string,
) (*CustomObjectRecord, error) {
	var result struct {
		CustomObjectRecord CustomObjectRecord `json:"custom_object_record"`
	}

	url := fmt.Sprintf("/custom_objects/%s/records/%s", customObjectKey, customObjectRecordID)
	body, err := z.get(ctx, url)

	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &result)

	if err != nil {
		return nil, err
	}
	return &result.CustomObjectRecord, nil
}

// UpdateCustomObjectRecord Updates an individual custom object record
// https://developer.zendesk.com/api-reference/custom-objects/custom_object_records/#update-custom-object-record
func (z *Client) UpdateCustomObjectRecord(
	ctx context.Context, customObjectKey string, customObjectRecordID string, record CustomObjectRecord,
) (*CustomObjectRecord, error) {
	var data, result struct {
		CustomObjectRecord CustomObjectRecord `json:"custom_object_record"`
	}
	data.CustomObjectRecord = record

	url := fmt.Sprintf("/custom_objects/%s/records/%s", customObjectKey, customObjectRecordID)
	body, err := z.patch(ctx, url, data)

	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &result)

	if err != nil {
		return nil, err
	}
	return &result.CustomObjectRecord, nil
}

// GetSourceByTargetResult result of the Get Source By Target
type GetSourceByTargetResult struct {
	CustomObjectRecords []CustomObjectRecord `json:"custom_object_records,omitempty"`
	Users               []User               `json:"users,omitempty"`
	Organizations       []Organization       `json:"organizations,omitempty"`
	Tickets             []Ticket             `json:"tickets,omitempty"`
	Page
}

// GetSourcesByTarget Returns a list of source objects whose values are populated with the id of a related target object
// https://developer.zendesk.com/api-reference/ticketing/lookup_relationships/lookup_relationships/#get-sources-by-target
func (z *Client) GetSourcesByTarget(
	ctx context.Context,
	targetType string,
	targetID string,
	fieldID string,
	sourceType string,
	opts *PageOptions,
) (GetSourceByTargetResult, error) {
	result := GetSourceByTargetResult{}
	tmp := opts
	if tmp == nil {
		tmp = &PageOptions{}
	}
	url := fmt.Sprintf("/%s/%s/relationship_fields/%s/%s", targetType, targetID, fieldID, sourceType)
	urlWithOptions, err := addOptions(url, tmp)
	body, err := z.get(ctx, urlWithOptions)

	if err != nil {
		return GetSourceByTargetResult{}, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return GetSourceByTargetResult{}, err
	}
	return result, nil
}

// DeleteCustomObjectRecord Delete a custom object record
// https://developer.zendesk.com/api-reference/custom-data/custom-objects/custom_object_records/#delete-custom-object-record
func (z *Client) DeleteCustomObjectRecord(
	ctx context.Context,
	record CustomObjectRecord,
) error {
	endpointURL := fmt.Sprintf("/custom_objects/%s/records/%s", record.CustomObjectKey, record.ID)
	err := z.delete(ctx, endpointURL)
	if err != nil {
		return err
	}
	return nil
}

// ListCustomObjectFields Lists all undeleted custom fields for the specified object.
// https://developer.zendesk.com/api-reference/custom-data/custom-objects/custom_object_fields/#list-custom-object-fields
func (z *Client) ListCustomObjectFields(
	ctx context.Context,
	customObjectKey string) ([]CustomObjectField, error) {

	var result struct {
		CustomObjectFields []CustomObjectField `json:"custom_object_fields"`
	}

	url := fmt.Sprintf("/custom_objects/%s/fields", customObjectKey)
	body, err := z.get(ctx, url)

	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &result)

	if err != nil {
		return nil, err
	}
	return result.CustomObjectFields, nil
}

// FilterCustomObjectRecords
// https://developer.zendesk.com/api-reference/custom-data/custom-objects/
// custom_object_records/#filtered-search-of-custom-object-records
func (z *Client) FilterCustomObjectRecords(
	ctx context.Context, customObjectKey string, filterBody interface{}, opts *SearchCustomObjectRecordsOptions,
) ([]CustomObjectRecord, CursorPaginationMeta, int64, error) {
	var data struct {
		FilterBody interface{} `json:"filter"`
	}
	data.FilterBody = filterBody
	var result struct {
		Count               int64                `json:"count"`
		CustomObjectRecords []CustomObjectRecord `json:"custom_object_records"`
		Meta                CursorPaginationMeta `json:"meta"`
	}
	tmp := opts
	if tmp == nil {
		tmp = &SearchCustomObjectRecordsOptions{}
	}
	url := fmt.Sprintf("/custom_objects/%s/records/search", customObjectKey)
	urlWithOptions, err := addOptions(url, tmp)
	body, err := z.post(ctx, urlWithOptions, data)
	if err != nil {
		return nil, CursorPaginationMeta{}, 0, err
	}
	err = json.Unmarshal(body, &result)
	return result.CustomObjectRecords, result.Meta, result.Count, nil
}

// ShowCustomObjectField Returns a custom field for a specific object using a provided key or id of the field.
func (z *Client) ShowCustomObjectField(
	ctx context.Context, customObjectKey string, customObjectFieldKeyOrId string,
) (CustomObjectField, error) {
	var result struct {
		CustomObjectField CustomObjectField `json:"custom_object_field"`
	}

	url := fmt.Sprintf("/custom_objects/%s/fields/%s", customObjectKey, customObjectFieldKeyOrId)
	body, err := z.get(ctx, url)

	if err != nil {
		return CustomObjectField{}, err
	}
	err = json.Unmarshal(body, &result)

	if err != nil {
		return CustomObjectField{}, err
	}
	return result.CustomObjectField, nil
}
