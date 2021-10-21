package constant

type LendingStatus string

const (
	LendingDraft    LendingStatus = "DRAFT"
	LendingCanceled LendingStatus = "CANCELED"
	LendingActive   LendingStatus = "ACTIVE"
	LendingInactive LendingStatus = "INACTIVE"
)
