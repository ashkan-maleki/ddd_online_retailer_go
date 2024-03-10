package repository

import "gorm.io/gorm"

type GormRepository struct {
	db *gorm.DB
}
