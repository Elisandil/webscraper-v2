// internal/usecase/schedule.go
package usecase

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
	"webscraper-v2/internal/config"
	"webscraper-v2/internal/domain/entity"
	"webscraper-v2/internal/domain/repository"

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
	}
}

func (uc *ScheduleUseCase) CreateSchedule(req *entity.CreateScheduleRequest, userID int64) (*entity.Schedule, error) {

	if err := uc.validateCronExpression(req.CronExpr); err != nil {
		return nil, fmt.Errorf("invalid cron expression: %w", err)
	}
	nextRun, err := uc.calculateNextRun(req.CronExpr)

	if err != nil {
		return nil, fmt.Errorf("error calculating next run: %w", err)
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
		return nil, fmt.Errorf("error creating schedule: %w", err)
	}

	// IMPORTANTE: Registrar el job en el cron si el scheduler está iniciado
	if uc.isStarted && schedule.Active {

		if err := uc.addJobToCron(schedule); err != nil {
			log.Printf("Warning: Could not add job to cron for schedule %d: %v", schedule.ID, err)
		}
	}
	log.Printf("✅ Schedule created: %s (ID: %d) - Next run: %v", schedule.Name, schedule.ID, nextRun)
	return schedule, nil
}

func (uc *ScheduleUseCase) GetSchedulesByUser(userID int64) ([]*entity.Schedule, error) {
	return uc.scheduleRepo.FindByUserID(userID)
}

func (uc *ScheduleUseCase) GetSchedule(id int64) (*entity.Schedule, error) {
	return uc.scheduleRepo.FindByID(id)
}

func (uc *ScheduleUseCase) UpdateSchedule(id int64, req *entity.UpdateScheduleRequest, userID int64) (*entity.Schedule, error) {
	schedule, err := uc.scheduleRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("error finding schedule: %w", err)
	}
	if schedule == nil {
		return nil, fmt.Errorf("schedule not found")
	}
	if schedule.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to schedule")
	}
	// Remover job anterior del cron si existe
	uc.removeJobFromCron(schedule.ID)

	// Actualizar campos si se proporcionaron
	if req.Name != nil {
		schedule.Name = strings.TrimSpace(*req.Name)
	}
	if req.URL != nil {
		schedule.URL = strings.TrimSpace(*req.URL)
	}
	if req.CronExpr != nil {
		newCronExpr := strings.TrimSpace(*req.CronExpr)

		if err := uc.validateCronExpression(newCronExpr); err != nil {
			return nil, fmt.Errorf("invalid cron expression: %w", err)
		}
		schedule.CronExpr = newCronExpr
		// Recalcular próxima ejecución
		nextRun, err := uc.calculateNextRun(schedule.CronExpr)

		if err != nil {
			return nil, fmt.Errorf("error calculating next run: %w", err)
		}
		schedule.NextRun = &nextRun
	}

	if req.Active != nil {
		schedule.Active = *req.Active
	}
	if err := uc.scheduleRepo.Update(schedule); err != nil {
		return nil, fmt.Errorf("error updating schedule: %w", err)
	}
	// Reagregar job al cron si está activo
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
		return fmt.Errorf("error finding schedule: %w", err)
	}
	if schedule == nil {
		return fmt.Errorf("schedule not found")
	}
	if schedule.UserID != userID {
		return fmt.Errorf("unauthorized access to schedule")
	}
	uc.removeJobFromCron(id)

	if err := uc.scheduleRepo.Delete(id); err != nil {
		return fmt.Errorf("error deleting schedule: %w", err)
	}
	log.Printf("✅ Schedule deleted: %s (ID: %d)", schedule.Name, schedule.ID)
	return nil
}

// Métodos del scheduler
func (uc *ScheduleUseCase) StartScheduler() {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.isStarted {
		log.Println("⚠️  Scheduler is already started")
		return
	}
	log.Println("🚀 Starting scheduler...")

	// Cargar todos los schedules activos y registrarlos en el cron
	if err := uc.loadActiveSchedules(); err != nil {
		log.Printf("❌ Error loading active schedules: %v", err)
	}
	// Iniciar el cron
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
	uc.activeJobs = make(map[int64]cron.EntryID) // Limpiar jobs activos
	uc.isStarted = false
	log.Println("✅ Scheduler stopped successfully")
}

func (uc *ScheduleUseCase) loadActiveSchedules() error {
	schedules, err := uc.scheduleRepo.FindActiveSchedules()

	if err != nil {
		return fmt.Errorf("error loading active schedules: %w", err)
	}
	log.Printf("📋 Loading %d active schedules...", len(schedules))
	successCount := 0

	for _, schedule := range schedules {

		if err := uc.addJobToCron(schedule); err != nil {
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

	if _, exists := uc.activeJobs[schedule.ID]; exists {
		log.Printf("⚠️  Job already exists for schedule %d, skipping", schedule.ID)
		return nil
	}
	jobFunc := func() {
		uc.executeSchedule(schedule)
	}
	// Agregar job al cron
	entryID, err := uc.cron.AddFunc(schedule.CronExpr, jobFunc)

	if err != nil {
		return fmt.Errorf("error adding job to cron: %w", err)
	}
	// Guardar referencia del job
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

func (uc *ScheduleUseCase) executeSchedule(schedule *entity.Schedule) {
	log.Printf("🔄 Executing scheduled scraping: %s (ID: %d, URL: %s)",
		schedule.Name, schedule.ID, schedule.URL)

	now := time.Now()
	// Ejecutar scraping
	_, err := uc.scrapingUC.ScrapeURL(schedule.URL, schedule.UserID)

	if err != nil {
		log.Printf("❌ Error executing scheduled scraping %d: %v", schedule.ID, err)
	} else {
		log.Printf("✅ Scheduled scraping completed successfully: %s", schedule.Name)
	}
	// Actualizar última ejecución y contador
	newRunCount := schedule.RunCount + 1

	if err := uc.scheduleRepo.UpdateLastRun(schedule.ID, now, newRunCount); err != nil {
		log.Printf("❌ Error updating last run for schedule %d: %v", schedule.ID, err)
	}
	// Calcular y actualizar próxima ejecución
	nextRun, err := uc.calculateNextRun(schedule.CronExpr)

	if err != nil {
		log.Printf("❌ Error calculating next run for schedule %d: %v", schedule.ID, err)
		return
	}
	if err := uc.scheduleRepo.UpdateNextRun(schedule.ID, nextRun); err != nil {
		log.Printf("❌ Error updating next run for schedule %d: %v", schedule.ID, err)
	} else {
		log.Printf("📅 Next run for schedule %d: %v", schedule.ID, nextRun)
	}
}

func (uc *ScheduleUseCase) validateCronExpression(cronExpr string) error {
	_, err := uc.cronParser.Parse(cronExpr)
	return err
}

func (uc *ScheduleUseCase) calculateNextRun(cronExpr string) (time.Time, error) {
	schedule, err := uc.cronParser.Parse(cronExpr)

	if err != nil {
		return time.Time{}, err
	}
	return schedule.Next(time.Now()), nil
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
