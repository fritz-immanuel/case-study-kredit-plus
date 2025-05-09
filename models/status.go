package models

var (
	STATUS_ACTIVE   = "1"
	STATUS_INACTIVE = "0"

	DEFAULT_STATUS_ID = STATUS_ACTIVE
)

type Status struct {
	ID   string `db:"id" json:"ID"`
	Name string `db:"name" json:"Name"`
}
