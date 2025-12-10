package usecase

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"
	"webscraper-v2/internal/domain/entity"
	"webscraper-v2/internal/domain/repository"
	"webscraper-v2/internal/infrastructure/config"
	pkgerrors "webscraper-v2/pkg/errors"
	"webscraper-v2/pkg/validator"

	"github.com/robfig/cron/v3"
)

type ScheduleUseCase struct {
	scheduleRepo repository.ScheduleRepository
	scrapingUC   *ScrapingUseCase
	config       *config.Config
	cron         *cron.Cron
	cronParser   cron.Parser
	activeJobs   map[int64]cron.EntryID
	mu           sync.RWMutex
	isStarted    bool
	validator    *validator.Validator
}

func NewScheduleUseCase(scheduleRepo repository.ScheduleRepository, scrapingUC *ScrapingUseCase, cfg *config.Config) *ScheduleUseCase {
	c := cron.New(cron.WithSeconds())
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	return &ScheduleUseCase{
		scheduleRepo: scheduleRepo,
		scrapingUC:   scrapingUC,
		config:       cfg,
		cron:         c,
		cronParser:   parser,
		activeJobs:   make(map[int64]cron.EntryID),
		isStarted:    false,
		validator:    validator.NewValidator(),
	}
}

func (uc *ScheduleUseCase) CreateSchedule(req *entity.CreateScheduleRequest, userID int64) (*entity.Schedule, error) {

	if err := uc.validateScheduleRequest(req); err != nil {
		return nil, err
	}

	if err := uc.validator.ValidateCronExpression(req.CronExpr); err != nil {
		return nil, pkgerrors.ValidationError(err.Error())
	}

	nextRun, err := uc.calculateNextRun(req.CronExpr)
	if err != nil {
		return nil, pkgerrors.InternalError("failed to calculate next run", err)
	}

	schedule := &entity.Schedule{
		UserID:   userID,
		Name:     strings.TrimSpace(req.Name),
		URL:      strings.TrimSpace(req.URL),
		CronExpr: strings.TrimSpace(req.CronExpr),
		Active:   true,
		NextRun:  &nextRun,
		RunCount: 0,
	}

	if err := uc.scheduleRepo.Create(schedule); err != nil {
		return nil, pkgerrors.DatabaseError("create schedule", err)
	}

	uc.mu.Lock()
	if uc.isStarted && schedule.Active {
		if err := uc.addJobToCronUnsafe(schedule); err != nil {
			log.Printf("Warning: Could not add job to cron for schedule %d: %v", schedule.ID, err)
		}
	}
	uc.mu.Unlock()

	log.Printf("✅ Schedule created: %s (ID: %d) - Next run: %v", schedule.Name, schedule.ID, nextRun)
	return schedule, nil
}

func (uc *ScheduleUseCase) GetSchedulesByUser(userID int64) ([]*entity.Schedule, error) {
	schedules, err := uc.scheduleRepo.FindByUserID(userID)
	if err != nil {
		return nil, pkgerrors.DatabaseError("get user schedules", err)
	}
	return schedules, nil
}

func (uc *ScheduleUseCase) GetSchedule(id int64) (*entity.Schedule, error) {
	schedule, err := uc.scheduleRepo.FindByID(id)
	if err != nil {
		return nil, pkgerrors.DatabaseError("get schedule", err)
	}
	return schedule, nil
}

func (uc *ScheduleUseCase) UpdateSchedule(id int64, req *entity.UpdateScheduleRequest, userID int64) (*entity.Schedule, error) {
	schedule, err := uc.scheduleRepo.FindByID(id)
	if err != nil {
		return nil, pkgerrors.DatabaseError("find schedule", err)
	}
	if schedule == nil {
		return nil, pkgerrors.NotFoundError("schedule")
	}
	if schedule.UserID != userID {
		return nil, pkgerrors.New(pkgerrors.CodeAuthorization, "unauthorized access to schedule", pkgerrors.ErrUnauthorized)
	}

	uc.removeJobFromCron(schedule.ID)

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if err := uc.validator.ValidateRequired(name, "name"); err != nil {
			return nil, pkgerrors.ValidationError(err.Error())
		}
		schedule.Name = name
	}

	if req.URL != nil {
		url := strings.TrimSpace(*req.URL)
		if err := uc.validator.ValidateURL(url); err != nil {
			return nil, pkgerrors.ValidationError(err.Error())
		}
		schedule.URL = url
	}

	if req.CronExpr != nil {
		newCronExpr := strings.TrimSpace(*req.CronExpr)
		if err := uc.validator.ValidateCronExpression(newCronExpr); err != nil {
			return nil, pkgerrors.ValidationError(err.Error())
		}
		schedule.CronExpr = newCronExpr

		nextRun, err := uc.calculateNextRun(schedule.CronExpr)
		if err != nil {
			return nil, pkgerrors.InternalError("failed to calculate next run", err)
		}
		schedule.NextRun = &nextRun
	}

	if req.Active != nil {
		schedule.Active = *req.Active
	}

	if err := uc.scheduleRepo.Update(schedule); err != nil {
		return nil, pkgerrors.DatabaseError("update schedule", err)
	}
	if uc.isStarted && schedule.Active {
		if err := uc.addJobToCron(schedule); err != nil {
			log.Printf("Warning: Could not add updated job to cron for schedule %d: %v", schedule.ID, err)
		}
	}

	log.Printf("✅ Schedule updated: %s (ID: %d)", schedule.Name, schedule.ID)
	return schedule, nil
}

func (uc *ScheduleUseCase) DeleteSchedule(id int64, userID int64) error {
	schedule, err := uc.scheduleRepo.FindByID(id)
	if err != nil {
		return pkgerrors.DatabaseError("find schedule", err)
	}
	if schedule == nil {
		return pkgerrors.NotFoundError("schedule")
	}
	if schedule.UserID != userID {
		return pkgerrors.New(pkgerrors.CodeAuthorization, "unauthorized access to schedule", pkgerrors.ErrUnauthorized)
	}

	uc.removeJobFromCron(id)
	if err := uc.scheduleRepo.Delete(id); err != nil {
		return pkgerrors.DatabaseError("delete schedule", err)
	}

	log.Printf("✅ Schedule deleted: %s (ID: %d)", schedule.Name, schedule.ID)
	return nil
}

func (uc *ScheduleUseCase) StartScheduler() {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.isStarted {
		log.Println("⚠️  Scheduler is already started")
		return
	}

	log.Println("🚀 Starting scheduler...")

	if err := uc.loadActiveSchedulesUnsafe(); err != nil {
		log.Printf("❌ Error loading active schedules: %v", err)
	}

	uc.cron.Start()
	uc.isStarted = true
	log.Printf("✅ Scheduler started successfully with %d active jobs", len(uc.activeJobs))
}

func (uc *ScheduleUseCase) StopScheduler() {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if !uc.isStarted {
		log.Println("⚠️  Scheduler is not running")
		return
	}

	log.Println("🛑 Stopping scheduler...")
	uc.cron.Stop()
	uc.activeJobs = make(map[int64]cron.EntryID)
	uc.isStarted = false
	log.Println("✅ Scheduler stopped successfully")
}

func (uc *ScheduleUseCase) GetSchedulerStatus() map[string]interface{} {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	return map[string]interface{}{
		"is_running":   uc.isStarted,
		"active_jobs":  len(uc.activeJobs),
		"cron_entries": len(uc.cron.Entries()),
	}
}

func (uc *ScheduleUseCase) loadActiveSchedulesUnsafe() error {
	schedules, err := uc.scheduleRepo.FindActiveSchedules()
	if err != nil {
		return pkgerrors.Wrap(err, "error loading active schedules")
	}

	log.Printf("📋 Loading %d active schedules...", len(schedules))
	successCount := 0

	for _, schedule := range schedules {
		if err := uc.addJobToCronUnsafe(schedule); err != nil {
			log.Printf("❌ Failed to add schedule %d to cron: %v", schedule.ID, err)
		} else {
			successCount++
			log.Printf("✅ Added schedule: %s (ID: %d, Cron: %s)",
				schedule.Name, schedule.ID, schedule.CronExpr)
		}
	}

	log.Printf("📊 Successfully loaded %d/%d schedules", successCount, len(schedules))
	return nil
}

func (uc *ScheduleUseCase) addJobToCron(schedule *entity.Schedule) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	return uc.addJobToCronUnsafe(schedule)
}

func (uc *ScheduleUseCase) addJobToCronUnsafe(schedule *entity.Schedule) error {

	if _, exists := uc.activeJobs[schedule.ID]; exists {
		log.Printf("⚠️  Job already exists for schedule %d, skipping", schedule.ID)
		return nil
	}

	scheduleID := schedule.ID
	scheduleName := schedule.Name
	scheduleURL := schedule.URL
	jobFunc := func() {
		uc.executeScheduleByID(scheduleID, scheduleName, scheduleURL)
	}

	entryID, err := uc.cron.AddFunc(schedule.CronExpr, jobFunc)
	if err != nil {
		return pkgerrors.Wrap(err, "error adding job to cron")
	}

	uc.activeJobs[schedule.ID] = entryID
	log.Printf("🕒 Job registered for schedule %d with cron expression: %s", schedule.ID, schedule.CronExpr)
	return nil
}

func (uc *ScheduleUseCase) removeJobFromCron(scheduleID int64) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if entryID, exists := uc.activeJobs[scheduleID]; exists {
		uc.cron.Remove(entryID)
		delete(uc.activeJobs, scheduleID)
		log.Printf("🗑️  Removed job for schedule %d", scheduleID)
	}
}

func (uc *ScheduleUseCase) executeScheduleByID(scheduleID int64, scheduleName, scheduleURL string) {
	log.Printf("🔄 Executing scheduled scraping: %s (ID: %d, URL: %s)",
		scheduleName, scheduleID, scheduleURL)

	schedule, err := uc.scheduleRepo.FindByID(scheduleID)
	if err != nil {
		log.Printf("❌ Error fetching schedule %d: %v", scheduleID, err)
		return
	}
	if schedule == nil {
		log.Printf("❌ Schedule %d not found, removing from cron", scheduleID)
		uc.removeJobFromCron(scheduleID)
		return
	}
	if !schedule.Active {
		log.Printf("⚠️  Schedule %d is no longer active, removing from cron", scheduleID)
		uc.removeJobFromCron(scheduleID)
		return
	}

	now := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	_, err = uc.scrapingUC.ScrapeURL(ctx, schedule.URL, schedule.UserID)
	if err != nil {
		log.Printf("❌ Error executing scheduled scraping %d: %v", scheduleID, err)
	} else {
		log.Printf("✅ Scheduled scraping completed successfully: %s", schedule.Name)
	}

	newRunCount := schedule.RunCount + 1
	if err := uc.scheduleRepo.UpdateLastRun(scheduleID, now, newRunCount); err != nil {
		log.Printf("❌ Error updating last run for schedule %d: %v", scheduleID, err)
	}

	nextRun, err := uc.calculateNextRun(schedule.CronExpr)
	if err != nil {
		log.Printf("❌ Error calculating next run for schedule %d: %v", scheduleID, err)
		return
	}
	if err := uc.scheduleRepo.UpdateNextRun(scheduleID, nextRun); err != nil {
		log.Printf("❌ Error updating next run for schedule %d: %v", scheduleID, err)
	} else {
		log.Printf("📅 Next run for schedule %d: %v", scheduleID, nextRun)
	}
}

func (uc *ScheduleUseCase) calculateNextRun(cronExpr string) (time.Time, error) {
	schedule, err := uc.cronParser.Parse(cronExpr)
	if err != nil {
		return time.Time{}, err
	}
	return schedule.Next(time.Now()), nil
}

func (uc *ScheduleUseCase) validateScheduleRequest(req *entity.CreateScheduleRequest) error {

	if err := uc.validator.ValidateStruct(req, "schedule request"); err != nil {
		return pkgerrors.ValidationError(err.Error())
	}
	if err := uc.validator.ValidateRequired(req.Name, "name"); err != nil {
		return pkgerrors.ValidationError(err.Error())
	}
	if err := uc.validator.ValidateURL(req.URL); err != nil {
		return pkgerrors.ValidationError(err.Error())
	}
	if err := uc.validator.ValidateRequired(req.CronExpr, "cron expression"); err != nil {
		return pkgerrors.ValidationError(err.Error())
	}
	return nil
}
