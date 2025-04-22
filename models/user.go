package models

import (
	"case-study-kredit-plus/library/types"
)

type UserBulk struct {
	ID                 string `json:"ID" db:"id"`
	Name               string `json:"Name" db:"name"`
	Email              string `json:"Email" db:"email"`
	Username           string `json:"Username" db:"username"`
	CountryCallingCode string `json:"CountryCallingCode" db:"country_calling_code"`
	PhoneNumber        string `json:"PhoneNumber" db:"phone_number"`
	Password           string `json:"Password" db:"password"`

	StatusID   string `json:"StatusID" db:"status_id"`
	StatusName string `json:"StatusName" db:"status_name"`
}

type User struct {
	ID                 string `json:"ID" db:"id"`
	Name               string `json:"Name" db:"name"`
	Email              string `json:"Email" db:"email"`
	Username           string `json:"Username" db:"username"`
	CountryCallingCode string `json:"CountryCallingCode" db:"country_calling_code"`
	PhoneNumber        string `json:"PhoneNumber" db:"phone_number"`
	Password           string `json:"Password" db:"password"`

	StatusID string `json:"StatusID" db:"status_id"`
	Status   Status `json:"Status"`
}

type UserJWTContent struct {
	ID    string `json:"ID" db:"id"`
	Name  string `json:"Name" db:"name" validate:"required"`
	Token string `json:"Token"`
	Email string `json:"Email" db:"email" validate:"required"`

	StatusID string `json:"StatusID" db:"status_id"`
	Status   Status `json:"Status"`
}

type UserUpdatePassword struct {
	ID                 string `json:"ID"`
	OldPassword        string `json:"OldPassword" validate:"required"`
	NewPassword        string `json:"NewPassword" validate:"required"`
	NewPasswordConfirm string `json:"NewPasswordConfirm" validate:"required"`
}

type FindAllUserParams struct {
	FindAllParams      types.FindAllParams
	UserID             string
	Name               string
	Username           string
	Email              string
	CountryCallingCode string
	PhoneNumber        string
	Password           string
}
