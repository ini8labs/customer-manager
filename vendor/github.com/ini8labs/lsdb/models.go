package lsdb

type EventParticipantInfo struct {
	EventID string
	ParticipantInfo
}

type ParticipantInfo struct {
	UserID     string
	BetNumbers []int
	Amount     int
}

type WinnerInfo struct {
	EventID   string
	UserID    string
	WinType   string
	AmountWon int
}
