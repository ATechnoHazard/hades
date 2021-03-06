package entities

import (
	"time"
)

type Event struct {
	ID                    uint           `json:"event_id" gorm:"primary_key;AUTO_INCREMENT"`
	Days                  uint           `json:"days"`
	OrganizationID        uint           `json:"org_id"`
	Name                  string         `json:"name" gorm:"unique"`
	Budget                string         `json:"budget"`
	Description           string         `json:"description"`
	Category              string         `json:"category"`
	Venue                 string         `json:"venue"`
	Attendance            string         `json:"attendance"`
	ExpectedParticipants  string         `json:"expected_participants"`
	PROrequest            string         `json:"pro_request"`
	CampusEngineerRequest string         `json:"campus_engineer_request"`
	Duration              string         `json:"duration"`
	Status                string         `json:"status"`
	ToDate                time.Time      `json:"to_date"`
	FromDate              time.Time      `json:"from_date"`
	ToTime                time.Time      `json:"to_time"`
	FromTime              time.Time      `json:"from_time"`
	DeletedAt             *time.Time     `json:"-" sql:"index"`
	Attendees             []Participant  `json:"attendees" gorm:"many2many:participant_events;"`
	Guests                []Guest        `json:"guests" gorm:"many2many:guest_event;"`
	Segments              []EventSegment `json:"-"`
}
