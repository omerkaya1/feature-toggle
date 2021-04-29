package internal

import "time"

// Feature is the structure that holds feature-related data
// NOTE: I really hate to expose the entities to the higher levels,
//  but the time is limited, so I keep it this way.
// NOTE: the following is an assumption
//  - customer id == user
type Feature struct {
	Inverted bool `json:"inverted"`
	Active   bool `json:"active"`
	// In a real app, I'd prefer managing features in the DB by its primary key
	ID            int64     `json:"-"`
	DisplayName   string    `json:"displayName,omitempty"`
	TechnicalName string    `json:"technicalName"`
	Description   string    `json:"description"`
	CustomerIDs   []string  `json:"customerIDs"`
	Expires       time.Time `json:"expiresOn"`
}
