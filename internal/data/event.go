package data

import (
	"github.com/google/uuid"
	"time"
)

type EventRecord struct {
	ID        uuid.UUID `json:"ID"`
	Type      EventType `json:"Type"`
	TimeStamp time.Time `json:"TimeStamp"`
}
type EventType int

const (
	CompanyCreated = iota
	CompanyUpdated
	CompanyDeleted
)

func (e EventType) String() string {
	return [...]string{"CompanyCreated", "CompanyUpdated", "CompanyDeleted"}[e]
}
