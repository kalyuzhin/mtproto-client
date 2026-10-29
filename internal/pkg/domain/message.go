package domain

import "time"

type Message struct {
	messageID int64
	message   string
	// chat where this message have been sent
	peerID int64
	date   time.Time
}
