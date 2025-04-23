package models

import (
	"case-study-kredit-plus/library/types"
)

type ConsumerCreditLimitBulk struct {
	ID         string  `json:"ID" db:"id" validate:"omitempty,uuid4"`
	ConsumerID string  `json:"ConsumerID" db:"consumer_id" validate:"required,uuid4"`
	Month1     float64 `json:"Month1" db:"1_month" validate:"numeric"`
	Month2     float64 `json:"Month2" db:"2_month" validate:"numeric"`
	Month3     float64 `json:"Month3" db:"3_month" validate:"numeric"`
	Month6     float64 `json:"Month6" db:"6_month" validate:"numeric"`

	StatusID   string `json:"StatusID" db:"status_id"`
	StatusName string `json:"StatusName" db:"status_name"`

	ConsumerName string `json:"ConsumerName" db:"consumer_name"`
}

type ConsumerCreditLimit struct {
	ID         string  `json:"ID" db:"id" validate:"omitempty,uuid4"`
	ConsumerID string  `json:"ConsumerID" db:"consumer_id" validate:"required,uuid4"`
	Month1     float64 `json:"Month1" db:"1_month" validate:"numeric"`
	Month2     float64 `json:"Month2" db:"2_month" validate:"numeric"`
	Month3     float64 `json:"Month3" db:"3_month" validate:"numeric"`
	Month6     float64 `json:"Month6" db:"6_month" validate:"numeric"`

	StatusID string `json:"StatusID" db:"status_id"`
	Status   Status `json:"Status"`

	Consumer *IDNameTemplate `json:"Consumer"`
}

type FindAllConsumerCreditLimitParams struct {
	FindAllParams types.FindAllParams
	ConsumerID    string `validate:"omitempty,uuid4"`
}

// Credit Limit Avalability
type ConsumerCreditLimitAvailability struct {
	ConsumerID     string  `json:"ConsumerID" db:"consumer_id" validate:"required,uuid4"`
	RemainingLimit float64 `json:"RemainingLimit" db:"remaining_limit" validate:"numeric"`
}
