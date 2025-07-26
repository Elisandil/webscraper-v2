package repository

import (
	"time"

	"webscraper-v2/internal/domain/entity"
)

type ScheduleRepository interface {
	Create(schedule *entity.Schedule) error
	FindByID(id int64) (*entity.Schedule, error)
	FindByUserID(userID int64) ([]*entity.Schedule, error)
	FindActiveSchedules() ([]*entity.Schedule, error)
	Update(schedule *entity.Schedule) error
	Delete(id int64) error
	UpdateLastRun(id int64, lastRun time.Time, runCount int) error
	UpdateNextRun(id int64, nextRun time.Time) error
}
