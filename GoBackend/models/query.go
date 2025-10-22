package models

import "time"

type Query struct {
	Id          int16     `json:"id"`
	CustomerId  string    `json:"customer_id"`
	CreatedAt   time.Time `json:"created_at"`
	QueryText   string    `json:"query_text"`
	Answer      string    `json:"answer"`
	QueryStatus Status    `json:"query_status"`
}

type QueryStatus struct {
	Answer      string `json:"answer"`
	QueryStatus Status `json:"query_status"`
}

type FAQ struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type Status int

const (
	PENDING Status = iota
	RESOLVED
	UNRESOLVED
)
