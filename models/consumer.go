package models

import (
	"case-study-kredit-plus/library/types"
	"time"
)

type ConsumerBulk struct {
	ID           string    `json:"ID" db:"id" validate:"omitempty,uuid4"`
	NIK          string    `json:"NIK" db:"NIK" validate:"len=16"`
	FullName     string    `json:"FullName" db:"full_name" validate:"alphanum"`
	LegalName    string    `json:"LegalName" db:"legal_name" validate:"alphanum"`
	PlaceOfBirth string    `json:"PlaceOfBirth" db:"place_of_birth" validate:"alphanum"`
	DateOfBirth  time.Time `json:"DateOfBirth" db:"date_of_birth"`
	Salary       float64   `json:"Salary" db:"salary" validate:"max=11"`
	KTPImgURL    string    `json:"KTPImgURL" db:"ktp_img_url" validate:"url"`
	SelfieImgURL string    `json:"SelfieImgURL" db:"selfie_img_url" validate:"url"`

	StatusID   string `json:"StatusID" db:"status_id"`
	StatusName string `json:"StatusName" db:"status_name"`
}

type Consumer struct {
	ID           string    `json:"ID" db:"id" validate:"omitempty,uuid4"`
	NIK          string    `json:"NIK" db:"NIK" validate:"len=16"`
	FullName     string    `json:"FullName" db:"full_name" validate:"alphanum"`
	LegalName    string    `json:"LegalName" db:"legal_name" validate:"alphanum"`
	PlaceOfBirth string    `json:"PlaceOfBirth" db:"place_of_birth" validate:"alphanum"`
	DateOfBirth  time.Time `json:"DateOfBirth" db:"date_of_birth"`
	Salary       float64   `json:"Salary" db:"salary" validate:"max=11"`
	KTPImgURL    string    `json:"KTPImgURL" db:"ktp_img_url" validate:"url"`
	SelfieImgURL string    `json:"SelfieImgURL" db:"selfie_img_url" validate:"url"`

	StatusID string `json:"StatusID" db:"status_id"`
	Status   Status `json:"Status"`
}

type FindAllConsumerParams struct {
	FindAllParams  types.FindAllParams
	NIK            string `validate:"len=16"`
	FullName       string `validate:"alphanum"`
	LegalName      string `validate:"alphanum"`
	PlaceOfBirth   string
	MinDateOfBirth string
	MaxDateOfBirth string
	MinSalary      float64 `validate:"numeric"`
	MaxSalary      float64 `validate:"numeric"`
}
