package repository

import "webscraper/internal/domain/entity"

type UserRepository interface {
	Create(user *entity.User) error
	FindByUsername(username string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindByID(id int64) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id int64) error
	ExistsUsername(username string) (bool, error)
	ExistsEmail(email string) (bool, error)
}
