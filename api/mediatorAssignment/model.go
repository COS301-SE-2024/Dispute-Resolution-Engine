package mediatorassignment

import "gorm.io/gorm"

type DBModel interface {
}

type DBModelReal struct {
	DB *gorm.DB
}