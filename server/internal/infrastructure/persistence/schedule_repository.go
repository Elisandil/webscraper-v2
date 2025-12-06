package persistence

import (
	"database/sql"
	"fmt"
	"time"
	"webscraper-v2/internal/domain/entity"
	"webscraper-v2/internal/domain/repository"
	"webscraper-v2/internal/infrastructure/database"
)

type scheduleRepository struct {
	db *database.SQLiteDB
}

func NewScheduleRepository(db *database.SQLiteDB) repository.ScheduleRepository {
	return &scheduleRepository{db: db}
}

func (r *scheduleRepository) Create(schedule *entity.Schedule) error {
	query := `INSERT INTO schedules (user_id, name, url, cron_expression, active, last_run, next_run, run_count, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	schedule.CreatedAt = now
	schedule.UpdatedAt = now

	res, err := r.db.Exec(query,
		schedule.UserID, schedule.Name, schedule.URL, schedule.CronExpr,
		schedule.Active, schedule.LastRun, schedule.NextRun, schedule.RunCount,
		schedule.CreatedAt, schedule.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating schedule: %w", err)
	}
	id, err := res.LastInsertId()

	if err != nil {
		return fmt.Errorf("error getting last insert id: %w", err)
	}
	schedule.ID = id
	return nil
}

func (r *scheduleRepository) FindByID(id int64) (*entity.Schedule, error) {
	query := `SELECT id, user_id, name, url, cron_expression, active, last_run, next_run, run_count, created_at, updated_at 
			  FROM schedules WHERE id = ?`

	schedule := &entity.Schedule{}
	var lastRun, nextRun, createdAt, updatedAt sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&schedule.ID, &schedule.UserID, &schedule.Name, &schedule.URL,
		&schedule.CronExpr, &schedule.Active, &lastRun, &nextRun,
		&schedule.RunCount, &createdAt, &updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding schedule by id: %w", err)
	}

	if err := r.parseTimestamps(schedule, lastRun, nextRun, createdAt, updatedAt); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (r *scheduleRepository) FindByUserID(userID int64) ([]*entity.Schedule, error) {
	query := `SELECT id, user_id, name, url, cron_expression, active, last_run, next_run, run_count, created_at, updated_at 
			  FROM schedules WHERE user_id = ? ORDER BY created_at DESC`

	return r.findSchedules(query, userID)
}

func (r *scheduleRepository) FindActiveSchedules() ([]*entity.Schedule, error) {
	query := `SELECT id, user_id, name, url, cron_expression, active, last_run, next_run, run_count, created_at, updated_at 
			  FROM schedules WHERE active = true ORDER BY next_run ASC`

	return r.findSchedules(query)
}

func (r *scheduleRepository) Update(schedule *entity.Schedule) error {
	query := `UPDATE schedules SET name = ?, url = ?, cron_expression = ?, active = ?, updated_at = ? 
			  WHERE id = ?`

	schedule.UpdatedAt = time.Now()
	_, err := r.db.Exec(query,
		schedule.Name, schedule.URL, schedule.CronExpr, schedule.Active,
		schedule.UpdatedAt, schedule.ID)

	if err != nil {
		return fmt.Errorf("error updating schedule: %w", err)
	}
	return nil
}

func (r *scheduleRepository) Delete(id int64) error {
	query := `DELETE FROM schedules WHERE id = ?`
	_, err := r.db.Exec(query, id)

	if err != nil {
		return fmt.Errorf("error deleting schedule: %w", err)
	}
	return nil
}

func (r *scheduleRepository) UpdateLastRun(id int64, lastRun time.Time, runCount int) error {
	query := `UPDATE schedules SET last_run = ?, run_count = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, lastRun, runCount, time.Now(), id)

	if err != nil {
		return fmt.Errorf("error updating last run: %w", err)
	}
	return nil
}

func (r *scheduleRepository) UpdateNextRun(id int64, nextRun time.Time) error {
	query := `UPDATE schedules SET next_run = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, nextRun, time.Now(), id)

	if err != nil {
		return fmt.Errorf("error updating next run: %w", err)
	}
	return nil
}

func (r *scheduleRepository) findSchedules(query string, args ...interface{}) ([]*entity.Schedule, error) {
	rows, err := r.db.Query(query, args...)

	if err != nil {
		return nil, fmt.Errorf("error querying schedules: %w", err)
	}
	defer rows.Close()
	var schedules []*entity.Schedule

	for rows.Next() {
		schedule := &entity.Schedule{}
		var lastRun, nextRun, createdAt, updatedAt sql.NullString

		err := rows.Scan(
			&schedule.ID,
			&schedule.UserID,
			&schedule.Name,
			&schedule.URL,
			&schedule.CronExpr,
			&schedule.Active,
			&lastRun,
			&nextRun,
			&schedule.RunCount,
			&createdAt,
			&updatedAt)

		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		if err := r.parseTimestamps(schedule, lastRun, nextRun, createdAt, updatedAt); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	return schedules, nil
}

func (r *scheduleRepository) parseTimestamps(schedule *entity.Schedule, lastRun, nextRun, createdAt, updatedAt sql.NullString) error {
	var err error

	if lastRun.Valid {
		t, err := r.parseDateTime(lastRun.String)

		if err != nil {
			return fmt.Errorf("error parsing last_run: %w", err)
		}
		schedule.LastRun = &t
	}

	if nextRun.Valid {
		t, err := r.parseDateTime(nextRun.String)

		if err != nil {
			return fmt.Errorf("error parsing next_run: %w", err)
		}
		schedule.NextRun = &t
	}

	if createdAt.Valid {
		schedule.CreatedAt, err = r.parseDateTime(createdAt.String)
		if err != nil {
			return fmt.Errorf("error parsing created_at: %w", err)
		}
	}

	if updatedAt.Valid {
		schedule.UpdatedAt, err = r.parseDateTime(updatedAt.String)
		if err != nil {
			return fmt.Errorf("error parsing updated_at: %w", err)
		}
	}

	return nil
}

func (r *scheduleRepository) parseDateTime(dateStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		time.DateTime,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse datetime: %s", dateStr)
}
