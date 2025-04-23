package models

import (
	"case-study-kredit-plus/library/types"
	"time"
)

type ConsumerBulk struct {
	ID           string    `json:"ID" db:"id" validate:"omitempty,uuid4"`
	NIK          string    `json:"NIK" db:"NIK" validate:"len=16,numeric"`
	FullName     string    `json:"FullName" db:"full_name"`
	LegalName    string    `json:"LegalName" db:"legal_name"`
	PlaceOfBirth string    `json:"PlaceOfBirth" db:"place_of_birth"`
	DateOfBirth  time.Time `json:"DateOfBirth" db:"date_of_birth"`
	Salary       float64   `json:"Salary" db:"salary" validate:"max=99999999999"`
	KTPImgURL    string    `json:"KTPImgURL" db:"ktp_img_url" validate:"url"`
	SelfieImgURL string    `json:"SelfieImgURL" db:"selfie_img_url" validate:"url"`

	StatusID   string `json:"StatusID" db:"status_id"`
	StatusName string `json:"StatusName" db:"status_name"`
}

type Consumer struct {
	ID           string    `json:"ID" db:"id" validate:"omitempty,uuid4"`
	NIK          string    `json:"NIK" db:"NIK" validate:"len=16,numeric"`
	FullName     string    `json:"FullName" db:"full_name"`
	LegalName    string    `json:"LegalName" db:"legal_name"`
	PlaceOfBirth string    `json:"PlaceOfBirth" db:"place_of_birth"`
	DateOfBirth  time.Time `json:"DateOfBirth" db:"date_of_birth"`
	Salary       float64   `json:"Salary" db:"salary" validate:"max=99999999999"`
	KTPImgURL    string    `json:"KTPImgURL" db:"ktp_img_url" validate:"url"`
	SelfieImgURL string    `json:"SelfieImgURL" db:"selfie_img_url" validate:"url"`

	StatusID string `json:"StatusID" db:"status_id"`
	Status   Status `json:"Status"`
}

type FindAllConsumerParams struct {
	FindAllParams  types.FindAllParams
	NIK            string `validate:"omitempty,len=16,numeric"`
	FullName       string
	LegalName      string
	PlaceOfBirth   string
	MinDateOfBirth string
	MaxDateOfBirth string
	MinSalary      float64 `validate:"numeric"`
	MaxSalary      float64 `validate:"numeric"`
}
