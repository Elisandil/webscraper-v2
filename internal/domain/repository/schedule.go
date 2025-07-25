package repository

import (
	"time"

	"webscraper-v2/internal/domain/entity"
)

type ScheduleRepository interface {
	Create(schedule *entity.Schedule) error
	FindByID(id uint64) (*entity.Schedule, error)
	FindByUserID(userID uint64) (*[]entity.Schedule, error)
	FindActiveSchedules() (*[]entity.Schedule, error)
	Update(schedule *entity.Schedule) error
	Delete(id uint64) error
	UpdateLastRun(id uint64, lastRun time.Time, runCount uint64) error
	UpdateNextRun(id uint64, nextRun time.Time) error
}
