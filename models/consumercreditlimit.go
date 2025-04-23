package models

import (
	"case-study-kredit-plus/library/types"
)

type ConsumerCreditLimitBulk struct {
	ID         string  `json:"ID" db:"id"`
	ConsumerID string  `json:"ConsumerID" db:"consumer_id"`
	Month1     float64 `json:"Month1" db:"1_month"`
	Month2     float64 `json:"Month2" db:"2_month"`
	Month3     float64 `json:"Month3" db:"3_month"`
	Month6     float64 `json:"Month6" db:"6_month"`

	StatusID   string `json:"StatusID" db:"status_id"`
	StatusName string `json:"StatusName" db:"status_name"`

	ConsumerName string `json:"ConsumerName" db:"consumer_name"`
}

type ConsumerCreditLimit struct {
	ID         string  `json:"ID" db:"id"`
	ConsumerID string  `json:"ConsumerID" db:"consumer_id"`
	Month1     float64 `json:"Month1" db:"1_month"`
	Month2     float64 `json:"Month2" db:"2_month"`
	Month3     float64 `json:"Month3" db:"3_month"`
	Month6     float64 `json:"Month6" db:"6_month"`

	StatusID string `json:"StatusID" db:"status_id"`
	Status   Status `json:"Status"`

	Consumer *IDNameTemplate `json:"Consumer"`
}

type FindAllConsumerCreditLimitParams struct {
	FindAllParams types.FindAllParams
	ConsumerID    string
}

// Credit Limit Avalability
type ConsumerCreditLimitAvailability struct {
	ConsumerID     string  `json:"ConsumerID" db:"consumer_id"`
	RemainingLimit float64 `json:"RemainingLimit" db:"remaining_limit"`
}
