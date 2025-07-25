package entity

import "time"

type Schedule struct {
	ID        uint64     `json:"id"`
	UserID    uint64     `json:"user_id"`
	Name      string     `json:"name"`
	URL       string     `json:"url"`
	CronExpr  string     `json:"cron_expression"`
	Active    bool       `json:"active"`
	LastRun   *time.Time `json:"last_run,omitempty"`
	NextRun   *time.Time `json:"next_run,omitempty"`
	RunCount  uint64     `json:"run_count"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CreatecheduleRequest struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	CronExpr string `json:"cron_expression"`
}

type UpdateScheduleRequest struct {
	Name     *string `json:"name,omitempty"`
	URL      *string `json:"url,omitempty"`
	CronExpr *string `json:"cron_expression,omitempty"`
	Active   *bool   `json:"active,omitempty"`
}
