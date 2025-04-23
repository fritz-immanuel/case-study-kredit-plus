package models

import (
	"case-study-kredit-plus/library/types"
)

type ConsumerTransactionBulk struct {
	ID                string  `json:"ID" db:"id" validate:"omitempty,uuid4"`
	ConsumerID        string  `json:"ConsumerID" db:"consumer_id" validate:"required,uuid4"`
	ContractNumber    string  `json:"ContractNumber" db:"contract_number"`
	OTR               float64 `json:"OTR" db:"OTR" validate:"numeric"`
	AdminFee          float64 `json:"AdminFee" db:"admin_fee" validate:"numeric"`
	InstallmentAmount float64 `json:"InstallmentAmount" db:"installment_amount" validate:"numeric"`
	LoanTerm          int     `json:"LoanTerm" db:"loan_term" validate:"oneof=1 2 3 6"`
	InterestAmount    float64 `json:"InterestAmount" db:"interest_amount" validate:"numeric"`
	TotalAmount       float64 `json:"TotalAmount" db:"total_amount" validate:"numeric"`
	AssetName         string  `json:"AssetName" db:"asset_name"`

	StatusID   string `json:"StatusID" db:"status_id"`
	StatusName string `json:"StatusName" db:"status_name"`

	ConsumerName string `json:"ConsumerName" db:"consumer_name"`
}

type ConsumerTransaction struct {
	ID                string  `json:"ID" db:"id" validate:"omitempty,uuid4"`
	ConsumerID        string  `json:"ConsumerID" db:"consumer_id" validate:"required,uuid4"`
	ContractNumber    string  `json:"ContractNumber" db:"contract_number"`
	OTR               float64 `json:"OTR" db:"OTR" validate:"numeric"`
	AdminFee          float64 `json:"AdminFee" db:"admin_fee" validate:"numeric"`
	InstallmentAmount float64 `json:"InstallmentAmount" db:"installment_amount" validate:"numeric"`
	LoanTerm          int     `json:"LoanTerm" db:"loan_term" validate:"oneof=1 2 3 6"`
	InterestAmount    float64 `json:"InterestAmount" db:"interest_amount" validate:"numeric"`
	TotalAmount       float64 `json:"TotalAmount" db:"total_amount" validate:"numeric"`
	AssetName         string  `json:"AssetName" db:"asset_name"`

	StatusID string `json:"StatusID" db:"status_id"`
	Status   Status `json:"Status"`

	Consumer *IDNameTemplate `json:"Consumer"`
}

type FindAllConsumerTransactionParams struct {
	FindAllParams  types.FindAllParams
	ConsumerID     string `validate:"omitempty,uuid4"`
	ContractNumber string
	LoanTerm       int `validate:"omitempty,oneof=1 2 3 6"`
}
