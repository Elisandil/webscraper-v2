package repository

import "webscraper/domain/entity"

type ScrapingRepository interface {
	Save(result *entity.ScrapingResult) error
	FindAll() ([]*entity.ScrapingResult, error)
	FindByID(id int64) (*entity.ScrapingResult, error)
	Delete(id int64) error
}
