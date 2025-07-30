package entity

import "time"

type Schedule struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	Name      string     `json:"name"`
	URL       string     `json:"url"`
	CronExpr  string     `json:"cron_expression"`
	Active    bool       `json:"active"`
	LastRun   *time.Time `json:"last_run,omitempty"`
	NextRun   *time.Time `json:"next_run,omitempty"`
	RunCount  int        `json:"run_count"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CreateScheduleRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	URL      string `json:"url" validate:"required,url"`
	CronExpr string `json:"cron_expression" validate:"required"`
}

type UpdateScheduleRequest struct {
	Name     *string `json:"name,omitempty"`
	URL      *string `json:"url,omitempty"`
	CronExpr *string `json:"cron_expression,omitempty"`
	Active   *bool   `json:"active,omitempty"`
}
