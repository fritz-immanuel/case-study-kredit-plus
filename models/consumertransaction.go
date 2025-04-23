package models

import (
	"case-study-kredit-plus/library/types"
)

type ConsumerTransactionBulk struct {
	ID                string  `json:"ID" db:"id"`
	ConsumerID        string  `json:"ConsumerID" db:"consumer_id"`
	ContractNumber    string  `json:"ContractNumber" db:"contract_number"`
	OTR               float64 `json:"OTR" db:"OTR"`
	AdminFee          float64 `json:"AdminFee" db:"admin_fee"`
	InstallmentAmount float64 `json:"InstallmentAmount" db:"installment_amount"`
	LoanTerm          int     `json:"LoanTerm" db:"loan_term"`
	InterestAmount    float64 `json:"InterestAmount" db:"interest_amount"`
	TotalAmount       float64 `json:"TotalAmount" db:"total_amount"`
	AssetName         string  `json:"AssetName" db:"asset_name"`

	StatusID   string `json:"StatusID" db:"status_id"`
	StatusName string `json:"StatusName" db:"status_name"`

	ConsumerName string `json:"ConsumerName" db:"consumer_name"`
}

type ConsumerTransaction struct {
	ID                string  `json:"ID" db:"id"`
	ConsumerID        string  `json:"ConsumerID" db:"consumer_id"`
	ContractNumber    string  `json:"ContractNumber" db:"contract_number"`
	OTR               float64 `json:"OTR" db:"OTR"`
	AdminFee          float64 `json:"AdminFee" db:"admin_fee"`
	InstallmentAmount float64 `json:"InstallmentAmount" db:"installment_amount"`
	LoanTerm          int     `json:"LoanTerm" db:"loan_term"`
	InterestAmount    float64 `json:"InterestAmount" db:"interest_amount"`
	TotalAmount       float64 `json:"TotalAmount" db:"total_amount"`
	AssetName         string  `json:"AssetName" db:"asset_name"`

	StatusID string `json:"StatusID" db:"status_id"`
	Status   Status `json:"Status"`

	Consumer *IDNameTemplate `json:"Consumer"`
}

type FindAllConsumerTransactionParams struct {
	FindAllParams  types.FindAllParams
	ConsumerID     string
	ContractNumber string
	LoanTerm       int
}
