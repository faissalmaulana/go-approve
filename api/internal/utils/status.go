package utils

type Status int

const (
	Approved Status = iota
	Rejected
	Pending
	Accepted
)

func (s Status) String() string {
	return []string{"approved", "rejected", "pending", "accepted"}[s]
}
