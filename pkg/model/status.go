package model

import "time"

var (
	// StatusOK ok
	StatusOK = "STATUS_OK_eduid_ladok_"
	// StatusFail not ok
	StatusFail = "STATUS_FAIL_eduid_ladok_"
)

// Status type
type Status struct {
	Name       string    `json:"name,omitempty"`
	SchoolName string    `json:"school_name,omitempty"`
	Healthy    bool      `json:"healthy,omitempty"`
	Status     string    `json:"status,omitempty"`
	Message    string    `json:"message,omitempty"`
	Timestamp  time.Time `json:"timestamp,omitempty"`
}

// ManyStatus contains many status objects
type ManyStatus []*Status

// Check checks the status of each status, return the first that does not pass.
func (s ManyStatus) Check() *Status {
	for _, status := range s {
		if s == nil {
			continue
		}
		if !status.Healthy {
			return status
		}
	}
	status := &Status{
		Healthy:   true,
		Status:    StatusOK,
		Timestamp: time.Now(),
	}
	return status
}
