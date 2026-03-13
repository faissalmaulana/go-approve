package contract

import "time"

type CreateApprovalRoom struct {
	Title       string
	DueAt       time.Time
	Filepaths   []string
	SubmitterId string
}
