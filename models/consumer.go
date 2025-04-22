package models

import (
	"case-study-kredit-plus/library/types"
	"time"
)

type ConsumerBulk struct {
	ID           string    `json:"ID" db:"id"`
	NIK          string    `json:"NIK" db:"NIK"`
	FullName     string    `json:"FullName" db:"full_name"`
	LegalName    string    `json:"LegalName" db:"legal_name"`
	PlaceOfBirth string    `json:"PlaceOfBirth" db:"place_of_birth"`
	DateOfBirth  time.Time `json:"DateOfBirth" db:"date_of_birth"`
	Salary       float64   `json:"Salary" db:"salary"`
	KTPImgURL    string    `json:"KTPImgURL" db:"ktp_img_url"`
	SelfieImgURL string    `json:"SelfieImgURL" db:"selfie_img_url"`

	StatusID   string `json:"StatusID" db:"status_id"`
	StatusName string `json:"StatusName" db:"status_name"`
}

type Consumer struct {
	ID           string    `json:"ID" db:"id"`
	NIK          string    `json:"NIK" db:"NIK"`
	FullName     string    `json:"FullName" db:"full_name"`
	LegalName    string    `json:"LegalName" db:"legal_name"`
	PlaceOfBirth string    `json:"PlaceOfBirth" db:"place_of_birth"`
	DateOfBirth  time.Time `json:"DateOfBirth" db:"date_of_birth"`
	Salary       float64   `json:"Salary" db:"salary"`
	KTPImgURL    string    `json:"KTPImgURL" db:"ktp_img_url"`
	SelfieImgURL string    `json:"SelfieImgURL" db:"selfie_img_url"`

	StatusID string `json:"StatusID" db:"status_id"`
	Status   Status `json:"Status"`
}

type FindAllConsumerParams struct {
	FindAllParams  types.FindAllParams
	NIK            string
	FullName       string
	LegalName      string
	PlaceOfBirth   string
	MinDateOfBirth string
	MaxDateOfBirth string
	MinSalary      float64
	MaxSalary      float64
}
