package zendesk

import (
	"context"
	"encoding/json"
	"time"
)

// Article is a struct for articles payload
type Article struct {
	AuthorID          int64      `json:"author_id,omitempty"`
	Body              string     `json:"body,omitempty"`
	CommentsDisabled  bool       `json:"comments_disabled,omitempty"`
	ContentTagIDs     []string   `json:"content_tag_ids,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
	Draft             bool       `json:"draft,omitempty"`
	EditedAt          *time.Time `json:"edited_at,omitempty"`
	HtmlUrl           string     `json:"html_url,omitempty"`
	ID                int64      `json:"id,omitempty"`
	LabelNames        []string   `json:"label_names,omitempty"`
	Locale            string     `json:"string,omitempty"`
	Outdated          bool       `json:"outdated,omitempty"`
	OutdatedLocales   []string   `json:"outdated_locales,omitempty"`
	PermissionGroupID int        `json:"permission_id,omitempty"`
	Position          int        `json:"position,omitempty"`
	Promoted          bool       `json:"promoted,omitempty"`
	SectionID         int        `json:"section_id,omitempty"`
	SourceLocale      string     `json:"source_locale,omitempty"`
	Title             string     `json:"title,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
	Url               string     `json:"url,omitempty"`
	UserSegmentID     int        `json:"user_segment_id,omitempty"`
	VoteCount         int        `json:"vote_count,omitempty"`
	VoteSum           int        `json:"vote_sum,omitempty"`
}

// ListArticles list all articles in Help Center
//
// https://developer.zendesk.com/api-reference/help_center/help-center-api/articles/#list-articles
func (z *Client) ListArticles(ctx context.Context, opts *TicketListOptions) ([]Article, Page, error) {
	var data struct {
		Articles []Article `json:"articles"`
		Page
	}

	tmp := opts
	if tmp == nil {
		tmp = &TicketListOptions{}
	}
	
	u, err := addOptions("/help_center/articles.json", tmp)
	if err != nil {
		return nil, Page{}, err
	}
	
	body, err := z.get(ctx, u)
	if err != nil {
		return nil, Page{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, Page{}, err
	}
	return data.Articles, data.Page, nil
}
